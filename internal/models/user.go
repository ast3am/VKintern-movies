package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (u *User) CheckCreds(email, password string) bool {
	if u.Email != email || u.Password != password {
		return false
	}
	return true
}

func checkPassword(userPwdHash, requestPwdHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPwdHash), []byte(requestPwdHash))
	return err == nil
}
