package entity

import "regexp"

type User struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

var loginRegex = regexp.MustCompile("^([a-zA-Z0-9@]).{6,}$")
var passwordRegex = regexp.MustCompile(`^[\S\s]{8,}$`)

// Лучше всего валидацию производить на клиенте
func IsLoginValid(login string) bool {
	return loginRegex.MatchString(login)
}

func IsPasswordValid(password string) bool {
	return loginRegex.MatchString(password)
}
