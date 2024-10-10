package routes

import (
	"encoding/json"
	"net/http"
	"skate/db"
	"skate/models"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(userID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID // Aquí incluimos el ID del usuario en los claims del token
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte("SkatePostgresDeveloperZamus"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials models.User
	json.NewDecoder(r.Body).Decode(&credentials)
	var userInDB models.User
	db.DB.Where("nik_name = ?", credentials.NikName).First(&userInDB)
	if userInDB.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Usuario no encontrado"))
		return
	}
	if !CheckPasswordHash(credentials.Password, userInDB.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Contraseña incorrecta"))
		return
	}

	tokenString, err := GenerateJWT(userInDB.ID) // Pasamos el ID del usuario en lugar del nikname
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error al generar el token"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":   tokenString,
		"NikName": credentials.NikName,
		"userId":  userInDB.ID, // Incluimos el ID del usuario en la respuesta
	})
}
