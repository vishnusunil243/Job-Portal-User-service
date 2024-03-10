package helper

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
func CompareHashedPassword(hashedPass, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
