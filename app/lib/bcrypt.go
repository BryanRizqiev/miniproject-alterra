package lib

import "golang.org/x/crypto/bcrypt"

func BcryptHashPassword(password string) (string, error) {

	var passwordBytes = []byte(password)

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)

	return string(hashedPasswordBytes), err

}

func BcryptMatchingPassword(hashedPassword, currPassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currPassword))

	return err == nil

}
