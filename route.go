package main

import (
	"encoding/json"
	"myshop-api/api/admin"
	"myshop-api/api/common"
	"myshop-api/api/front"
	"net/http"

	"github.com/gorilla/mux"
)

func addRouters(router *mux.Router) {
	router.HandleFunc("/health", health)
	router.HandleFunc("/image/{f1}/{f2}/{f3}/{file-name}", common.GetCDNImagePath).Methods("GET")

	//------------ Admin APIs
	router.HandleFunc("/admin/add-product", admin.AddProduct).Methods("POST")
	router.HandleFunc("/admin/product-list", admin.GetProductList).Methods("GET")
	router.HandleFunc("/admin/product-detail/{productID}", admin.GetProductList).Methods("GET")
	router.HandleFunc("/admin/product-detail/{productID}", admin.EditProduct).Methods("UPDATE")

	router.HandleFunc("/admin/order-details/{orderID}", admin.GetOrderDetails).Methods("GET")
	router.HandleFunc("/admin/order-list/{limit}/{offset}", admin.GetOrderList).Methods("GET")
	router.HandleFunc("/admin/order-list/{orderType}/{limit}/{offset}", admin.GetOrderList).Methods("GET")
	router.HandleFunc("/admin/order-status", admin.UpdateOrderStatus).Methods("UPDATE")

	router.HandleFunc("/admin/check-login", admin.CheckUserLogin).Methods("POST")
	router.HandleFunc("/admin/change-password", admin.ChangePassword).Methods("UPDATE")

	//----------Front APIs
	router.HandleFunc("/check-login", front.CheckUserLogin).Methods("POST")
	router.HandleFunc("/register", front.Register).Methods("POST")
	router.HandleFunc("/check-email", front.CheckUserEmail).Methods("POST")
	router.HandleFunc("/product-list", front.GetProductList).Methods("GET")
	router.HandleFunc("/product-list/{productIDs}", front.GetProductList).Methods("GET")
	router.HandleFunc("/product-detail/{productID}", front.GetProductList).Methods("GET")
	router.HandleFunc("/place-order/{userID}", front.PlaceOrder).Methods("POST")
	router.HandleFunc("/order-list/{userID}/{limit}/{offset}", front.GetOrderList).Methods("GET")
	router.HandleFunc("/order-list/{userID}/{orderType}/{limit}/{offset}", front.GetOrderList).Methods("GET")
	router.HandleFunc("/order-details/{orderID}", front.GetOrderDetails).Methods("GET")

}

func health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
