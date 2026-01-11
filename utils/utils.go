package utils

import (
	"encoding/json"
	"net/http"
	"subscribly/customerrors"

	"golang.org/x/crypto/bcrypt"
)

func PasswordHasher(password string) (string, error) {

	hashedPassInBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassInBytes), nil

}

func ComparePassword(password string, hashedPassword string) (error) {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil{
		return err
	}

	return nil
}

func ErrorWriter(w http.ResponseWriter, err error, errorCode int) {

	response := customerrors.MyErrors{Message: err.Error() }

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(errorCode)
	json.NewEncoder(w).Encode(response)

}

func ResponseWriter(w http.ResponseWriter, response any, successCode int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(successCode)
	json.NewEncoder(w).Encode(response)
}