package main

import (
	"fmt"
	"log"
	"myshop-api/api/data"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	//-------connection database
	err := data.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	//-------setting up route
	router := mux.NewRouter()
	addRouters(router)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "UPDATE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	fmt.Println("Server is started...")
	log.Fatal(http.ListenAndServe(":8000", c.Handler(router)))
}
