package handlers

import (
	"Meowria-backend/db"
	"encoding/json"
	"net/http"
)

func GetAllRecords(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	users, err := db.GetAllRecords("postgres://exupvkwi:FQOURrIUoc19JWoXYZ6ywiC5PRTER4N-@balarama.db.elephantsql.com/exupvkwi")
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

func GetRecordsByUserId(w http.ResponseWriter, r *http.Request) {
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
	records, err := db.GetRecordsByUserId(id)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	output, _ := json.MarshalIndent(records, "", "  ")
	_, err = w.Write(output)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CreateRecord(w http.ResponseWriter, r *http.Request) {
	// Cors
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Parse the JSON request body into a User struct
	var record db.SymptomsRecord
	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the record into the database
	err = db.InsertRecord(record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
