package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"log"
)

type Users struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstname"`
	SurName   string `json:"surname"`
	UserName  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func GetAllUsers(connection string) ([]Users, error) {
	conn, err := pgx.Connect(context.Background(), connection)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(conn, context.Background())

	rows, err := conn.Query(context.Background(), "SELECT * FROM Users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []Users
	for rows.Next() {
		user := Users{}
		err = rows.Scan(&user.ID, &user.FirstName, &user.SurName, &user.UserName, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUsersByIdFromDB(connection string, id int) ([]Users, error) {
	conn, err := pgx.Connect(context.Background(), connection)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(conn, context.Background())

	rows, err := conn.Query(context.Background(), "SELECT * FROM Users WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []Users
	for rows.Next() {
		user := Users{}
		err = rows.Scan(&user.ID, &user.FirstName, &user.SurName, &user.UserName, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}
func InsertUser(user Users) error {
	// Open a database connection
	db, err := sql.Open("postgres", "postgres://exupvkwi:FQOURrIUoc19JWoXYZ6ywiC5PRTER4N-@balarama.db.elephantsql.com/exupvkwi")
	if err != nil {
		return err
	}
	defer db.Close()

	// Execute the SQL INSERT statement
	stmt, err := db.Prepare("INSERT INTO Users(id, firstname, surname, username, email, password) VALUES ($1, $2, $3, $4, $5, $6 )")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID, user.FirstName, user.SurName, user.UserName, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}
func GetUserByUsername(username string) (Users, error) {
	conn, err := pgx.Connect(context.Background(), "postgres://exupvkwi:FQOURrIUoc19JWoXYZ6ywiC5PRTER4N-@balarama.db.elephantsql.com/exupvkwi")
	if err != nil {
		log.Fatal(err)
		return Users{}, err
	}
	defer conn.Close(context.Background())

	row, err := conn.Query(context.Background(), "SELECT id, firstname, surname, username, email, password FROM Users WHERE username=$1", username)
	if err != nil {
		return Users{}, err
	}
	defer row.Close()

	user := Users{}
	if row.Next() {
		err = row.Scan(&user.ID, &user.FirstName, &user.SurName, &user.UserName, &user.Email, &user.Password)
		if err != nil {
			return Users{}, err
		}
	} else {
		return Users{}, fmt.Errorf("user not found")
	}

	return user, nil
}

func CheckPassword(userPassword string, password string) bool {
	return userPassword == password
}
