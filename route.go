package main

import (
	"encoding/json"
	"myshop-api/api/admin"
	"net/http"

	"github.com/gorilla/mux"
)

func addRouters(router *mux.Router) {
	router.HandleFunc("/health", health)

	//------------ Admin APIs

	router.HandleFunc("/admin/add-product", admin.AddProduct).Methods("POST")
	router.HandleFunc("/admin/product-list", admin.GetProductList).Methods("GET")
}

func health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
