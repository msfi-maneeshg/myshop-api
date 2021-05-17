package common

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// Response :
type Response struct {
	Error   interface{} `json:"error,omitempty"`
	Content interface{} `json:"content,omitempty"`
}

//APIResponse : to send response in request
func APIResponse(w http.ResponseWriter, status int, output interface{}) {
	var objResponce Response
	if status == http.StatusOK {
		objResponce.Content = output
	} else {
		objResponce.Error = output
	}
	finalOutput, _ := json.Marshal(objResponce)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(finalOutput)
	return
}

//GetImage :
func GetCDNImagePath(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var filePath = vars["file-path"]
	var fileName = vars["file-name"]
	data, _ := ioutil.ReadFile("images/" + filePath + "/" + fileName)
	w.Write(data)
	r.Body.Close()
}
