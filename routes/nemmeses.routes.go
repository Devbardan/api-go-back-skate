package routes

import (
	"encoding/json"
	"net/http"
	"skate/db"
	"skate/models"
)

func CreateNemesisHandler(w http.ResponseWriter, r *http.Request) {
	var newNemesis models.Nemesis
	json.NewDecoder(r.Body).Decode(&newNemesis)
	createNemesis := db.DB.Create(&newNemesis)
	err := createNemesis.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	json.NewEncoder(w).Encode(&newNemesis)
}

func GetNemesesHandler(w http.ResponseWriter, r *http.Request) {
	var nemeses []models.Nemesis
	db.DB.Find(&nemeses)
	json.NewEncoder(w).Encode(&nemeses)
}
