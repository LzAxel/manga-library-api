package hash

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password, salt string) (string, error) {
	passwordByte := []byte(password + salt)
	passwordHash, err := bcrypt.GenerateFromPassword(passwordByte, 10)
	if err != nil {
		return "", err
	}

	return string(passwordHash), nil
}

func ComparePassword(hashedPassword, password string) error {

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
