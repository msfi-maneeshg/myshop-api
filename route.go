package main

import (
	"encoding/json"
	"myshop-api/api/admin"
	"myshop-api/api/common"
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
	router.HandleFunc("/admin/product-detail/{productID}", admin.EditProduct).Methods("update")

	router.HandleFunc("/admin/order-list/{limit}/{offset}", admin.GetOrderList).Methods("GET")
	router.HandleFunc("/admin/order-list/{orderType}/{limit}/{offset}", admin.GetOrderList).Methods("GET")

}

func health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
