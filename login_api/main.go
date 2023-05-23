package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = map[string]string{
	"username1": "password1",
	"username2": "password2",
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	password, exists := users[user.Username]
	if !exists || password != user.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenString := createToken(user.Username)
	w.Write([]byte(tokenString))
}

func createToken(username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString([]byte("your-secret-key"))
	return tokenString
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the home page!"))
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/home", homeHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
