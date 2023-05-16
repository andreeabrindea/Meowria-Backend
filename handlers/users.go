package handlers

import (
	"Meowria-backend/db"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID       int       `json:"user_id"`
	SessionToken string    `json:"session_token"`
	Expiry       time.Time `json:"expiry"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")

}

func ParseIDFromPath(r *http.Request) (int, error) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		return 0, errors.New("invalid path")
	}

	idStr := pathParts[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("id should be an integer")
	}
	return id, nil
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	users, err := db.GetAllUsers("postgres://exupvkwi:FQOURrIUoc19JWoXYZ6ywiC5PRTER4N-@balarama.db.elephantsql.com/exupvkwi")
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	usersJS, _ := json.MarshalIndent(users, "", "  ")
	_, err = w.Write(usersJS)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetUsersById(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := ParseIDFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	users, err := db.GetUsersByIdFromDB("postgres://exupvkwi:FQOURrIUoc19JWoXYZ6ywiC5PRTER4N-@balarama.db.elephantsql.com/exupvkwi", id)
	if err != nil {
		return
	}
	if len(users) == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", "application/json")
		output, _ := json.MarshalIndent(users, "", "  ")
		_, err = w.Write(output)
		if err != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Set headers for the main request
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Parse the JSON request body into a User struct
	var user db.Users
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the user into the database
	err = db.InsertUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour).Unix(), // Set expiration time
	})

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the JWT token in the response
	response := LoginResponse{
		UserID:       user.ID,
		SessionToken: tokenString,
		Expiry:       time.Now().Add(time.Hour),
	}

	// Encode the response struct as JSON and write to the response body
	jsonResp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResp)
	if err != nil {
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Handle preflight request
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Set headers for the main request
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	// Parse the login request from the request body
	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	username := loginReq.Username
	password := loginReq.Password

	// Validate the credentials against a database
	user, err := db.GetUserByUsername(username)
	if err != nil {
		http.Error(w, "something off", http.StatusUnauthorized)
		return
	}
	if db.CheckPassword(user.Password, password) == false {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour).Unix(), // Set expiration time
	})

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the JWT token in the response
	response := LoginResponse{
		UserID:       user.ID,
		SessionToken: tokenString,
		Expiry:       time.Now().Add(time.Hour),
	}

	// Encode the response struct as JSON and write to the response body
	jsonResp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResp)
	if err != nil {
		return
	}
}
