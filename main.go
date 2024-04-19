package main

import (
	"log"
	"net/http"
	"skate/db"
	"skate/models"
	"skate/routes"

	"github.com/gorilla/mux"
)

func main() {
	/*
	 */
	db.DBConnection()
	db.DB.AutoMigrate(models.Trik{})
	db.DB.AutoMigrate(models.User{})

	router := mux.NewRouter()

	router.HandleFunc("/", routes.HomeHandler)

	router.HandleFunc("/trik", routes.CreateTrik).Methods("POST")
	router.HandleFunc("/user", routes.CreateUserHandler).Methods("POST")
	router.HandleFunc("/triks", routes.GetTriksHandler).Methods("GET")
	router.HandleFunc("/users", routes.GetUsersHandler).Methods("GET")
	router.HandleFunc("/trik/{id}", routes.GetTrikHandler).Methods("GET")
	router.HandleFunc("/user/{id}", routes.GetUserHandler).Methods("GET")
	router.HandleFunc("/trik/{id}", routes.DeleteTrik).Methods("DELETE")
	router.HandleFunc("/user/{id}", routes.DeleteUserHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", router))
}
