package lib

import (
	"bytes"
	"html/template"
	"math/rand"
	user_entity "miniproject-alterra/module/user/entity"
)

const DATE_WITH_DAY_FORMAT = "2006-01-02 15:04:05 Monday"

func Contains(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

func RandomString(n int) string {

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)

}

func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func CheckIsAdmin(user user_entity.User) bool {

	if user.Role == "admin" {
		return true
	}

	return false

}
