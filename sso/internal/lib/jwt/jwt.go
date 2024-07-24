package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wlcmtunknwndth/messagio_test/sso/internal/domain/models"
	"os"
	"time"
)

const (
	secretKeyEnv  = "secret"
	errNoKeyFound = "no secret key found"
)

func NewToken(user *models.User, duration time.Duration) (string, error) {
	const op = "sso.lib.jwt.NewToken"
	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)

	claims["uid"] = user.ID
	claims["user"] = user.Username
	claims["exp"] = time.Now().Add(duration).Unix()

	secret, ok := os.LookupEnv(secretKeyEnv)
	if !ok {
		return "", fmt.Errorf("%s: %s", op, errNoKeyFound)
	}

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}
