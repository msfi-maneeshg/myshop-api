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
	var fileName = vars["file-name"]
	var folder1 = vars["f1"]
	var folder2 = vars["f2"]
	var folder3 = vars["f3"]

	var fullPath string
	if folder1 != "" {
		fullPath = folder1
	}
	if folder2 != "" {
		fullPath = folder1 + "/" + folder2
	}
	if folder2 != "" {
		fullPath = folder1 + "/" + folder2 + "/" + folder3
	}
	data, _ := ioutil.ReadFile(fullPath + "/" + fileName)
	w.Write(data)
	r.Body.Close()
}
