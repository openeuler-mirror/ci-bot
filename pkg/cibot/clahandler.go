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

type CLAResult struct {
	IsSuccess   bool    `json:"isSuccess,omitempty"`
	Description *string `json:"description,omitempty"`
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
		claresult := CLAResult{
			IsSuccess: true,
		}
		result, err := json.Marshal(claresult)
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
	} else if r.Method == "OPTIONS" {
		// CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		glog.Infof("finish request method: %s", r.Method)
	} else {
		glog.Infof("unsupport request method: %s", r.Method)
	}
}
