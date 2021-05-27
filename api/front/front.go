package front

import (
	"encoding/json"
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

	if userID == "" || password == "" {
		common.APIResponse(w, http.StatusNotFound, "Invalid Credentials")
		return
	}

	objUserInfo, err := CheckUserLoginDetails(userID, password)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Getting error while checking user login details. Error:"+err.Error())
		return
	}

	if objUserInfo.Username == "" {
		common.APIResponse(w, http.StatusNotFound, "Invalid Credentials")
		return
	}

	common.APIResponse(w, http.StatusOK, UserDetails{Username: objUserInfo.Username, UserID: userID})
}

//CheckUserEmail :
func CheckUserEmail(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userID")
	if userID == "" {
		common.APIResponse(w, http.StatusNotFound, "Invalid Credentials")
		return
	}
	objUserInfo, err := CheckUserLoginDetails(userID, "")
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Getting error while checking user email. Error:"+err.Error())
		return
	}

	if objUserInfo.Username == "" {
		common.APIResponse(w, http.StatusNotFound, "Invalid Credentials")
		return
	}

	common.APIResponse(w, http.StatusOK, UserDetails{Username: objUserInfo.Username, UserID: userID})
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

//PlaceOrder :
func PlaceOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var userID = vars["userID"]
	var objUserOrderInput UserOrderInput
	//---------- check user details
	objUserInfo, err := CheckUserLoginDetails(userID, "")
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Getting error while checking user email. Error:"+err.Error())
		return
	}

	if objUserInfo.Username == "" {
		common.APIResponse(w, http.StatusNotFound, "Invalid Credentials")
		return
	}

	//---------- check user input details
	//------check body request
	if r.Body == nil {
		common.APIResponse(w, http.StatusBadRequest, "Request body can not be blank")
		return
	}
	err = json.NewDecoder(r.Body).Decode(&objUserOrderInput)
	if err != nil {
		common.APIResponse(w, http.StatusBadRequest, "Error:"+err.Error())
		return
	}

	if len(objUserOrderInput.UserOrderInfo) == 0 {
		common.APIResponse(w, http.StatusBadRequest, "Request body can not be blank")
		return
	}

	if objUserOrderInput.Address == "" || objUserOrderInput.Name == "" || objUserOrderInput.MobileNumber == 0 {
		common.APIResponse(w, http.StatusBadRequest, "Request body can not be blank")
		return
	}

	productIDs := []string{}
	productQuantityMap := map[int64]int64{}
	//---------- validate data
	for _, product := range objUserOrderInput.UserOrderInfo {
		if _, okay := productQuantityMap[product.ProductID]; !okay {
			productIDs = append(productIDs, strconv.Itoa(int(product.ProductID)))
			productQuantityMap[product.ProductID] = product.ProductQuantity
		}
	}

	productsDetails, err := GetProductsDetail("", strings.Join(productIDs, ","), "", "", 0)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Getting error while getting product list. Error:"+err.Error())
		return
	}
	if len(productsDetails) != len(productIDs) {
		common.APIResponse(w, http.StatusBadRequest, "Something wrong with products! Please try again")
		return
	}

	var totalPayment float64
	for index, productDetails := range productsDetails {

		if productDetails.ProductQuantity < productQuantityMap[productDetails.ProductID] {
			common.APIResponse(w, http.StatusBadRequest, "Something wrong with products! Please try again")
			return
		}

		productDetails.ProductQuantity = productQuantityMap[productDetails.ProductID]
		discountedPrice := productDetails.ProductPrize * (100 - productDetails.ProductDiscount) / 100

		totalPayment = totalPayment + discountedPrice
		productsDetails[index] = productDetails

	}

	orderID, err := InsertUserOrder(objUserInfo, objUserOrderInput, totalPayment)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Getting error while placing user order. Error:"+err.Error())
		return
	}

	err = InsertUserOrderDetails(orderID, objUserInfo, productsDetails, totalPayment)
	if err != nil {
		common.APIResponse(w, http.StatusInternalServerError, "Getting error while inserting order details. Error:"+err.Error())
		return
	}

	for _, productDetails := range productsDetails {
		err = UpdateProductStock(productDetails)
		if err != nil {
			common.APIResponse(w, http.StatusInternalServerError, "Getting error while inserting order details. Error:"+err.Error())
			return
		}
	}
	common.APIResponse(w, http.StatusOK, "Order has been placed.")
}
