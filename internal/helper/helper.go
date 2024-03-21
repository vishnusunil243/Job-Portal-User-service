package helper

import (
	"fmt"
	"time"
	"unicode"

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
func GetNumberInString(s string) int {
	num := 0
	if len(s) == 0 {
		return 0
	}
	for _, sr := range s {
		if unicode.IsNumber(sr) {
			num = int(sr - '0')
		}
	}
	return num
}
func ConvertStringToDate(str string) (time.Time, error) {
	layout := "02-01-2006"
	date, err := time.Parse(layout, str)
	if err != nil {
		return time.Time{}, fmt.Errorf("please provide a valid start date")
	}
	return date, nil
}
