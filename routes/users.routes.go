package routes

import (
	"encoding/json"
	"net/http"
	"skate/db"
	"skate/models"

	"github.com/gorilla/mux"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)
	db.DB.First(&user, params["id"])
	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}
	db.DB.Model(&user).Association("Triks").Find(&user.Triks)
	json.NewEncoder(w).Encode(&user)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	db.DB.Find(&users)
	json.NewEncoder(w).Encode(&users)

}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	json.NewDecoder(r.Body).Decode(&newUser)
	createTrik := db.DB.Create(&newUser)
	err := createTrik.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	json.NewEncoder(w).Encode(&newUser)

}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)
	db.DB.First(&user, params["id"])
	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}
	db.DB.Unscoped().Delete(&user) /* unscoped para eliminar permanente, si lo quitas queda con fecha de eliminacion recuerda*/
	w.WriteHeader(http.StatusOK)
}
