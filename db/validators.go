package db

import "regexp"

func isValidEmail(email string) bool {
	pattern := "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(email)
}

func IsValidUsername(username string) bool {
	pattern := "^[a-zA-Z0-9!\\-_.]+$"
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(username)
}
