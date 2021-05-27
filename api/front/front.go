package front

import (
	"myshop-api/api/common"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
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

//GetProductList:
func GetProductList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var productID = vars["productID"]
	var productIDs = vars["productIDs"]

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	lowStockOrder := strings.ToLower(r.URL.Query().Get("lowStock"))
	if lowStockOrder != "asc" && lowStockOrder != "desc" {
		lowStockOrder = ""
	}

	newStock := strings.ToLower(r.URL.Query().Get("newStock"))
	if newStock != "asc" && newStock != "desc" {
		newStock = ""
	}

	if strings.Contains(r.URL.RequestURI(), "product-detail") {
		productIDs = ""
	} else if strings.Contains(r.URL.RequestURI(), "product-list") {
		productID = ""
	}

	r.URL.RequestURI()

	productList, err := GetProductsDetail(productID, productIDs, lowStockOrder, newStock, limit)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Getting error while getting product list. Error:"+err.Error())
		return
	}

	if productID != "" {
		if len(productList) == 0 {
			common.APIResponse(w, http.StatusNotFound, "No product found!")
			return
		}
		common.APIResponse(w, http.StatusOK, productList[0])
		return
	}

	common.APIResponse(w, http.StatusOK, productList)
}
