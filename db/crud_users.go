package db

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func InsertUser(user Users) (error, int) {
	// Open a database connection
	db, err := sql.Open("postgres", "postgres://exupvkwi:FQOURrIUoc19JWoXYZ6ywiC5PRTER4N-@balarama.db.elephantsql.com/exupvkwi")
	if err != nil {
		return err, 0
	}
	defer db.Close()

	// Check if the username already exists
	_, err = GetUserByUsername(user.UserName)
	if err == nil {
		return fmt.Errorf("username taken"), 0
	} else if err != nil && err.Error() != "user not found" {
		return err, 0
	}

	// Check if the email already exists
	_, err = GetUserByEmail(user.Email)
	if err == nil {
		return fmt.Errorf("email already exists"), 0
	} else if err != nil && err.Error() != "user not found" {
		return err, 0
	}

	if isValidEmail(user.Email) == false {
		return fmt.Errorf("invalid email"), 0
	}

	if len(user.UserName) < 3 {
		return fmt.Errorf("username should contain at least 3 characters (letters and digits)"), 0
	}
	if IsValidUsername(user.UserName) == false {
		return fmt.Errorf("invalid username"), 0
	}

	// Get the next unique ID from the sequence generator
	var id int
	err = db.QueryRow("SELECT nextval('users_id_seq')").Scan(&id)
	if err != nil {
		return err, id
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err, id
	}

	// Execute the SQL INSERT statement
	stmt, err := db.Prepare("INSERT INTO Users(id, firstname, surname, username, email, password) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		return err, id
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, user.FirstName, user.SurName, user.UserName, user.Email, hashedPassword)
	if err != nil {
		return err, id
	}

	return nil, id
}
