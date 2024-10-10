package routes

import (
	"encoding/json"
	"net/http"
	"skate/db"
	"skate/models"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	json.NewDecoder(r.Body).Decode(&newUser)
	hashedPassword, err := HashPassword(newUser.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error al hashear la contraseña"))
		return
	}
	newUser.Password = hashedPassword
	createTrik := db.DB.Create(&newUser)
	err = createTrik.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	tokenString, err := GenerateJWT(newUser.NikName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error al generar el token"))
		return
	}

	response := map[string]interface{}{
		"user":  &newUser,
		"token": tokenString,
	}

	json.NewEncoder(w).Encode(response)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)
	db.DB.Where("nik_name = ?", params["username"]).First(&user)
	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}
	db.DB.Model(&user).Association("Triks").Find(&user.Triks)
	db.DB.Model(&user).Association("Followers").Find(&user.Followers)
	db.DB.Model(&user).Association("Nemesis").Find(&user.Nemesis)
	json.NewEncoder(w).Encode(&user)
}

func GetNiknameHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)

	db.DB.Where("nik_name = ?", params["nik_name"]).First(&user)

	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}

	db.DB.Model(&user).Association("Triks").Find(&user.Triks)
	db.DB.Model(&user).Association("Followers").Find(&user.Followers)
	db.DB.Model(&user).Association("Nemesis").Find(&user.Nemesis)
	json.NewEncoder(w).Encode(&user)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	db.DB.Preload("Triks").Preload("Nemesis").Preload("Followers").Find(&users)
	json.NewEncoder(w).Encode(&users)
}

func EditUserHandler(w http.ResponseWriter, r *http.Request) {
	var updatedUser models.User
	json.NewDecoder(r.Body).Decode(&updatedUser)
	params := mux.Vars(r)

	var user models.User
	db.DB.Where("nik_name = ?", params["username"]).First(&user)
	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}

	if updatedUser.Password != "" {
		hashedPassword, err := HashPassword(updatedUser.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error al hashear la contraseña"))
			return
		}
		updatedUser.Password = hashedPassword
	} else {
		updatedUser.Password = user.Password
	}

	db.DB.Model(&user).Updates(updatedUser)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&user)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)
	db.DB.Where("nik_name = ?", params["username"]).First(&user)
	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}
	db.DB.Unscoped().Delete(&user)
	w.WriteHeader(http.StatusOK)
}
