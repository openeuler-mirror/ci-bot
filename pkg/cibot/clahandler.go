package cibot

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitee.com/openeuler/ci-bot/pkg/cibot/database"
	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/golang/glog"
	"golang.org/x/oauth2"
)

type CLAHandler struct {
	Context context.Context
}

type CLARequest struct {
	Type        int     `json:"type,omitempty"`
	Name        *string `json:"name,omitempty"`
	Title       *string `json:"title,omitempty"`
	Corporation *string `json:"corporation,omitempty"`
	Address     *string `json:"address,omitempty"`
	Date        *string `json:"date,omitempty"`
	Email       *string `json:"email,omitempty"`
	Telephone   *string `json:"telephone,omitempty"`
	Fax         *string `json:"fax,omitempty"`
	Code        *string `json:"code,omitempty"`
	Lang        *string `json:"lang,omitempty"`
	Client      *string `json:"client,omitempty"`
	AccessKey   string  `json:"-"`
}

type CLAResult struct {
	IsSuccess   bool   `json:"isSuccess"`
	ErrorCode   int    `json:"errorCode"`
	Description string `json:"description,omitempty"`
}

const (
	ErrorCode_OK = iota
	ErrorCode_ServerHandleError
	ErrorCode_EmailError
	ErrorCode_TelephoneError
	ErrorCode_EmailNotTheSameError
)

const COOKIE_KEY string = "cla-info"

// ServeHTTP validates an incoming cla request.
func (s *CLAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Info("received a cla request")
	if r.Method == "POST" {
		// add logs
		glog.Infof("CLA Request Header: %v", r.Header)
		glog.Infof("CLA Request Host: %v", r.Host)
		glog.Infof("CLA Request RemoteAddr: %v", r.RemoteAddr)

		// read body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.HandleResult(w, CLAResult{
				IsSuccess:   false,
				Description: fmt.Sprintf("read body error: %v", err),
				ErrorCode:   ErrorCode_ServerHandleError,
			})
			return
		}

		// unmarshal request
		var clarequest CLARequest
		err = json.Unmarshal(body, &clarequest)
		if err != nil {
			s.HandleResult(w, CLAResult{
				IsSuccess:   false,
				Description: fmt.Sprintf("unmarshal json error: %v", err),
				ErrorCode:   ErrorCode_ServerHandleError,
			})
			return
		}

		cookie, err := r.Cookie(COOKIE_KEY)
		if err == nil {
			clarequest.AccessKey = cookie.Value
		} else {
			glog.Infof("Get cookie err: %v", err)
		}

		glog.Infof("cla request content: %v", clarequest)
		s.HandleRequest(w, clarequest)
	} else if r.Method == "OPTIONS" {
		// CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		glog.Infof("finish request method: %s", r.Method)
	} else {
		glog.Infof("unsupport request method: %s", r.Method)
	}
}

// HandleResult output result to client
func (s *CLAHandler) HandleResult(w http.ResponseWriter, r CLAResult) {
	// log error code
	if !r.IsSuccess {
		glog.Errorf("handle result error code: %v description: %v",
			r.ErrorCode, r.Description)
	}

	// constuct result
	result, err := json.Marshal(r)
	if err != nil {
		glog.Errorf("marshal result error: %v", err)
		return
	}

	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// Content type
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

// HandleRequest handles the cla request
func (s *CLAHandler) HandleRequest(w http.ResponseWriter, request CLARequest) {
	// build model object
	if *request.Code == "" || *request.Client == "" || *request.Lang == "" {
		s.HandleResult(w, CLAResult{
			IsSuccess:   false,
			Description: fmt.Sprintf("request parameter error"),
			ErrorCode:   ErrorCode_ServerHandleError,
		})
		return
	}

	accesskey := request.AccessKey

	if accesskey == "" {
		token, err := GetToken(*request.Code, *request.Client, *request.Lang)

		if err != nil {
			s.HandleResult(w, CLAResult{
				IsSuccess:   false,
				Description: fmt.Sprintf("request gitee user error: %v", err),
				ErrorCode:   ErrorCode_ServerHandleError,
			})
			return
		}
		accesskey = token.AccessToken
		glog.Infof("access key get successfully.")

	}

	cookie := http.Cookie{Name: COOKIE_KEY, Value: accesskey, Path: "/", MaxAge: 86400}
	http.SetCookie(w, &cookie)

	emails, err := GetEmails(accesskey)
	if err != nil {
		s.HandleResult(w, CLAResult{
			IsSuccess:   false,
			Description: fmt.Sprintf("request parameter error: %v", err),
			ErrorCode:   ErrorCode_ServerHandleError,
		})
		return
	}


        primaryEmail := ""

        for _, email := range emails {
               if email.State == "confirmed"  {
                       for _, t := range email.Scope {
                              if t == "primary" {
                                      primaryEmail = email.Email
                                      break
                              } 
                       } 
               } 
        } 

	if primaryEmail == "" || primaryEmail != *request.Email {
		s.HandleResult(w, CLAResult{
			IsSuccess:   false,
			Description: "The email is not the same as gitee account email.",
			ErrorCode:   ErrorCode_EmailNotTheSameError,
		})
		return
	}

	cds := database.CLADetails{
		Type: request.Type,
	}
	if request.Name != nil {
		cds.Name = *request.Name
	}
	if request.Title != nil {
		cds.Title = *request.Title
	}
	if request.Corporation != nil {
		cds.Corporation = *request.Corporation
	}
	if request.Address != nil {
		cds.Address = *request.Address
	}
	if request.Date != nil {
		cds.Date = *request.Date
	}
	if request.Email != nil {
		cds.Email = *request.Email
	}
	if request.Telephone != nil {
		cds.Telephone = *request.Telephone
	}
	if request.Fax != nil {
		cds.Fax = *request.Fax
	}

	// tostring
	data, err := cds.ToString()
	if err != nil {
		s.HandleResult(w, CLAResult{
			IsSuccess:   false,
			Description: fmt.Sprintf("request tostring error: %v", err),
			ErrorCode:   ErrorCode_ServerHandleError,
		})
		return
	}
	glog.Infof("add cla details data: %s", data)

	// Check email in database
	var lenEmail int
	err = database.DBConnection.Model(&database.CLADetails{}).
		Where("email = ?", cds.Email).Count(&lenEmail).Error
	if err != nil {
		s.HandleResult(w, CLAResult{
			IsSuccess:   false,
			Description: fmt.Sprintf("check email exitency error: %v", err),
			ErrorCode:   ErrorCode_ServerHandleError,
		})
		return
	}
	if lenEmail > 0 {
		s.HandleResult(w, CLAResult{
			IsSuccess:   false,
			Description: "email is already registered",
			ErrorCode:   ErrorCode_EmailError,
		})
		return
	}

	/* Check telephone in database
	var lenTelephone int
	err = database.DBConnection.Model(&database.CLADetails{}).
		Where("telephone = ?", cds.Telephone).Count(&lenTelephone).Error
	if err != nil {
		s.HandleResult(w, CLAResult{
			IsSuccess:   false,
			Description: fmt.Sprintf("check telephone exitency error: %v", err),
			ErrorCode:   ErrorCode_ServerHandleError,
		})
		return
	}
	if lenTelephone > 0 {
		s.HandleResult(w, CLAResult{
			IsSuccess:   false,
			Description: "telephone is already registered",
			ErrorCode:   ErrorCode_TelephoneError,
		})
		return
	}*/

	// add cla in database
	err = database.DBConnection.Create(&cds).Error
	if err != nil {
		s.HandleResult(w, CLAResult{
			IsSuccess:   false,
			Description: fmt.Sprintf("add cla details error: %v", err),
			ErrorCode:   ErrorCode_ServerHandleError,
		})
		return
	}

	// constuct result
	s.HandleResult(w, CLAResult{
		IsSuccess: true,
		ErrorCode: ErrorCode_OK,
	})
}

func GetUser(ak string) (gitee.User, error) {

	ctx2 := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ak},
	)

	// configuration
	giteeConf := gitee.NewConfiguration()
	giteeConf.HTTPClient = oauth2.NewClient(ctx2, ts)

	// git client
	giteeClient := gitee.NewAPIClient(giteeConf)

	user, _, err := giteeClient.UsersApi.GetV5User(ctx2, nil)
	return user, err
}

func GetEmails(ak string) ([]gitee.Email, error) {
        ctx2 := context.Background()
        ts := oauth2.StaticTokenSource(
                &oauth2.Token{AccessToken: ak},
        )

        // configuration
        giteeConf := gitee.NewConfiguration()
        giteeConf.HTTPClient = oauth2.NewClient(ctx2, ts)

        // git client
        giteeClient := gitee.NewAPIClient(giteeConf)

        emails, _, err := giteeClient.EmailsApi.GetV5Emails(ctx2, nil)
        return emails, err


}
