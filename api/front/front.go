package front

import (
	"myshop-api/api/common"
	"net/http"
)

//CheckUserLogin :
func CheckUserLogin(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userID")
	password := r.FormValue("password")

	userName, err := CheckUserLoginDetails(userID, password)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Getting error while checking user login details. Error:"+err.Error())
		return
	}

	if userName == "" {
		common.APIResponse(w, http.StatusNotFound, "Invalid Credentials")
		return
	}

	common.APIResponse(w, http.StatusOK, UserDetails{Username: userName, UserID: userID})
}
