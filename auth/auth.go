package auth

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	SALT_SIZE                = 32
	PASSWORD_HASH_ITERATIONS = 4096
	PASSWORD_KEY_LENGTH      = 32
)

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SALT_SIZE)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

func HashPassword(password string, salt []byte) string {
	dk := pbkdf2.Key([]byte(password), salt, PASSWORD_HASH_ITERATIONS, PASSWORD_KEY_LENGTH, sha512.New)
	return fmt.Sprintf("%x", dk)
}
