package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/love-sunshine30/logReader/handlers"
	"github.com/love-sunshine30/logReader/models"
)

func main() {
	// initializing the db
	err := models.InitDB()
	if err != nil {
		fmt.Printf("Database connection failed: %v", err)
		return
	}

	// router
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/health", handlers.Health)
	router.HandlerFunc(http.MethodPost, "/upload", handlers.Upload)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Printf("Starting the server on port %s", srv.Addr)
	srv.ListenAndServe()
}
