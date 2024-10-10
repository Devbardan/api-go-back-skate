package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"skate/db"
	"skate/models"

	"github.com/gorilla/mux"
)

func GetUserFollowersHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user models.User
	var userFollowers []models.UserFollower

	db.DB.Where("nik_name = ?", params["username"]).First(&user)
	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	db.DB.Preload("User").Preload("Follower").Where("user_id = ?", user.ID).Find(&userFollowers)
	if len(userFollowers) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No followers found for this user"))
		return
	}
	json.NewEncoder(w).Encode(&userFollowers)
}

func FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	var newFollower models.UserFollower
	params := mux.Vars(r)

	var user models.User
	db.DB.Where("nik_name = ?", params["username"]).First(&user)
	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}
	newFollower.UserID = user.ID

	var follower models.User
	db.DB.Where("nik_name = ?", params["followerId"]).First(&follower)
	if follower.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Follower not found"))
		return
	}
	newFollower.FollowerID = follower.ID

	// Comprobar si el usuario ya está siguiendo al otro usuario
	var existingFollower models.UserFollower
	if err := db.DB.Where("user_id = ? AND follower_id = ?", newFollower.UserID, newFollower.FollowerID).First(&existingFollower).Error; err == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ya estás siguiendo a este usuario"))
		return
	}

	createFollower := db.DB.Create(&newFollower)
	if createFollower.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(createFollower.Error.Error()))
		return
	}
	json.NewEncoder(w).Encode(&newFollower)
}

func UnfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	username := params["username"]
	followerNikName := params["followerNikName"]

	log.Printf("username: %s, followerNikName: %s", username, followerNikName)

	if username == "" || followerNikName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Faltan parámetros en la solicitud"))
		return
	}

	var user models.User
	var follower models.User

	// Buscar el ID del usuario usando el username
	if err := db.DB.Where("nik_name = ?", username).First(&user).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Usuario no encontrado"))
		return
	}

	// Buscar el ID del seguidor usando el followerNikName
	if err := db.DB.Where("nik_name = ?", followerNikName).First(&follower).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Seguidor no encontrado"))
		return
	}

	var existingFollower models.UserFollower
	// Verificar si la relación existe en la tabla follower_user
	if err := db.DB.Where("user_id = ? AND follower_id = ?", user.ID, follower.ID).First(&existingFollower).Error; err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No estás siguiendo a este usuario"))
		return
	}

	// Eliminar la relación en la tabla follower_user
	db.DB.Delete(&existingFollower)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Has dejado de seguir a este usuario"))
}