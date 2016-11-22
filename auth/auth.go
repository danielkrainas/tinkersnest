package auth

import (
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"time"

	"golang.org/x/crypto/pbkdf2"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/danielkrainas/tinkersnest/api/v1"
)

var (
	sharedKey = []byte("secret foo")
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

func BearerToken(u *v1.User) (string, error) {
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: sharedKey}, nil)
	if err != nil {
		return "", err
	}

	c := jwt.Claims{
		Subject: u.Name,
		Issuer:  "tinkernest",
		Expiry:  jwt.NewNumericDate(time.Now().Add(5 * time.Hour)),
	}

	return jwt.Signed(sig).Claims(c).CompactSerialize()
}

func VerifyBearerToken(rawToken string) (string, error) {
	token, err := jwt.ParseSigned(rawToken)
	if err != nil {
		return "", err
	}

	c := &jwt.Claims{}
	if err := token.Claims(sharedKey, c); err != nil {
		return "", err
	}

	err = c.Validate(jwt.Expected{
		Issuer: "tinkersnest",
	})

	if err != nil {
		return "", errors.New("token invalid")
	}

	return c.Subject, nil
}
