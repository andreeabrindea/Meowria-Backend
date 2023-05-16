package main

import (
	"Meowria-backend/handlers"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc(
		"/api/users",
		handlers.GetAllUsers,
	)
	http.HandleFunc(
		"/api/users/",
		handlers.GetUsersById,
	)
	http.HandleFunc(
		"/register",
		handlers.CreateUser,
	)
	http.HandleFunc(
		"/login",
		handlers.LoginHandler,
	)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
