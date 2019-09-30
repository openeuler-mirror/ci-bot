package cibot

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
)

type CLAHandler struct {
	Context context.Context
}

type CLARequest struct {
	IsSuccess   bool    `json:"isSuccess,omitempty"`
	Description *string `json:"description,omitempty"`
}

type CLAResult struct {
	Type        int     `json:"type,omitempty"`
	Name        *string `json:"name,omitempty"`
	Title       *string `json:"title,omitempty"`
	Corporation *string `json:"corporation,omitempty"`
	Address     *string `json:"address,omitempty"`
	Date        *string `json:"date,omitempty"`
	Email       *string `json:"email,omitempty"`
	Telephone   *string `json:"telephone,omitempty"`
	Fax         *string `json:"fax,omitempty"`
}

// ServeHTTP validates an incoming cla request.
func (s *CLAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Info("received a cla request")
	if r.Method == "POST" {
		// read body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			glog.Errorf("read body error: %v", err)
			return
		}

		// unmarshal request
		var clarequest CLARequest
		err = json.Unmarshal(body, &clarequest)
		if err != nil {
			glog.Errorf("unmarshal body error: %v", err)
			return
		}

		glog.Infof("cla request content: %v", clarequest)

		// constuct result
		claresult := CLARequest{
			IsSuccess: true,
		}
		result, err := json.Marshal(claresult)
		if err != nil {
			glog.Errorf("marshal result error: %v", err)
			return
		}

		// CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Content type
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))
	} else {
		glog.Infof("unsupport request method: %s", r.Method)
	}
}
