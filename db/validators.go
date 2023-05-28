package db

import "regexp"

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9.'_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
	if emailRegex.MatchString(email) == false {
		return false
	}
	return true
}

func IsValidUsername(username string) bool {
	pattern := "^[a-zA-Z0-9!\\-_.]+$"
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(username)
}
