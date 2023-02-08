package hash

import (
	"bytes"
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

const saltLenght = 8

func GenerateSalt() []byte {
	salt := make([]byte, saltLenght)
	rand.Read(salt)

	return salt
}

func HashPassword(salt []byte, password string) []byte {
	result := append([]byte{}, salt...)
	passwordHash := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)

	return append(result, passwordHash...)
}

func ComparePassword(hashedPassword []byte, password string) bool {
	salt := hashedPassword[0:saltLenght]
	userPasswordHash := HashPassword(salt, password)

	return bytes.Equal(userPasswordHash, hashedPassword)
}
