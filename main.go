package main

import (
	"log"
	"net/http"
	"skate/db"
	"skate/models"
	"skate/routes"

	"github.com/gorilla/mux"
)

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func OptionsHandler(w http.ResponseWriter, r *http.Request) {
    enableCors(w)
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }
}
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)
		next.ServeHTTP(w, r)
	})
}

func main() {
	db.DBConnection()
	db.DB.AutoMigrate(models.User{})
	db.DB.AutoMigrate(models.UserFollower{})
	db.DB.AutoMigrate(models.Trik{})
	db.DB.AutoMigrate(models.Nemesis{})
	router := mux.NewRouter()

	// Manejo global de solicitudes OPTIONS
	router.Methods("OPTIONS").HandlerFunc(OptionsHandler)

	router.HandleFunc("/", routes.HomeHandler)
	router.HandleFunc("/user", routes.CreateUserHandler).Methods("POST")
	router.HandleFunc("/users", routes.GetUsersHandler).Methods("GET")
	router.HandleFunc("/user/{username}", routes.EditUserHandler).Methods("PUT")
	router.HandleFunc("/user/{username}", routes.DeleteUserHandler).Methods("DELETE")
	router.HandleFunc("/user/{username}/followers", routes.GetUserFollowersHandler).Methods("GET")
	router.HandleFunc("/user/{username}", routes.GetUserHandler).Methods("GET")
	router.HandleFunc("/users/{username}/followers/{followerId}", routes.FollowUserHandler).Methods("POST")
	router.HandleFunc("/user/{username}/follower/{followerNikName}", routes.UnfollowUserHandler).Methods("DELETE")
	router.HandleFunc("/trik", routes.CreateTrik).Methods("POST")
	router.HandleFunc("/triks", routes.GetTriksHandler).Methods("GET")
	router.HandleFunc("/trik/{id}", routes.GetTrikHandler).Methods("GET")
	router.HandleFunc("/trik/{id}", routes.DeleteTrik).Methods("DELETE")
	router.HandleFunc("/nemeses", routes.GetNemesesHandler).Methods("GET")
	router.HandleFunc("/nemese", routes.CreateNemesisHandler).Methods("POST")
	router.HandleFunc("/login", routes.LoginHandler).Methods("POST")

	router.Use(corsMiddleware)

	log.Fatal(http.ListenAndServe(":5000", router))
}
