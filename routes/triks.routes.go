package routes

import (
	"encoding/json"
	"net/http"
	"skate/db"
	"skate/models"

	"github.com/gorilla/mux"
)

func GetTrikHandler(w http.ResponseWriter, r *http.Request) {
	var trik models.Trik
	params := mux.Vars(r)
	db.DB.First(&trik, params["id"])
	if trik.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}
	json.NewEncoder(w).Encode(&trik)
}
func GetTriksHandler(w http.ResponseWriter, r *http.Request) {
	var triks []models.Trik
	db.DB.Find(&triks)
	json.NewEncoder(w).Encode(&triks)
}

func CreateTrik(w http.ResponseWriter, r *http.Request) {
	var newTrik models.Trik
	json.NewDecoder(r.Body).Decode(&newTrik)
	createTrik := db.DB.Create(&newTrik)
	err := createTrik.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	json.NewEncoder(w).Encode(&newTrik)

}
func DeleteTrik(w http.ResponseWriter, r *http.Request) {
	var trik models.Trik
	params := mux.Vars(r)
	db.DB.First(&trik, params["id"])
	if trik.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}
	db.DB.Unscoped().Delete(&trik) /* unscoped para eliminar permanente, si lo quitas queda con fecha de eliminacion recuerda*/
	w.WriteHeader(http.StatusOK)
}
