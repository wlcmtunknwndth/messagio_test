package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

const (
	secretKeyEnv  = "jwtkey"
	idKey         = "id"
	usernameKey   = "user"
	adminKey      = "isadmin"
	expirationKey = "exp"

	scope = "common.jwt.NewToken."
)

type Info struct {
	ID       int64
	Username string
	IsAdmin  bool
}

var (
	ErrNoKeyFound    = errors.New("no secret key found")
	ErrInvalidToken  = errors.New("invalid token")
	errMapAssertion  = errors.New("couldn't assert to map claims")
	errAssertion     = errors.New("wrong assertion")
	errValueNotFound = errors.New("value not found")
)

func NewToken(id int64, username string, duration time.Duration, isAdmin bool) (string, error) {
	const op = scope + "NewToken"
	//token := jwt.New(jwt.SigningMethodHS512)
	//
	//claims, ok := token.Claims.(jwt.MapClaims)
	//if !ok {
	//	return "", fmt.Errorf("%s: %w", op, errMapAssertion)
	//}
	claims := jwt.MapClaims{
		idKey:         id,
		usernameKey:   username,
		expirationKey: time.Now().Add(duration).Unix(),
		adminKey:      isAdmin,
	}
	//claims[idKey] = id
	//claims[usernameKey] = username
	//claims[expirationKey] = time.Now().Add(duration).Unix()
	//claims[adminKey] = isAdmin

	secret, ok := os.LookupEnv(secretKeyEnv)
	if !ok {
		return "", fmt.Errorf("%s: %w", op, ErrNoKeyFound)
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))

	//tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}

func GetInfo(token string) (*Info, error) {
	const op = scope + "GetInfo"

	claims, err := castToMapClaims(token)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var inf Info
	var ok bool
	inf.Username, ok = claims[usernameKey].(string)
	if inf.Username == "" {
		return nil, fmt.Errorf("%s: %w", op, errAssertion)
	} else if !ok {
		return nil, fmt.Errorf("%s: %w", op, ErrNoKeyFound)
	}

	inf.ID, ok = claims[idKey].(int64)
	if inf.ID == 0 {
		return nil, fmt.Errorf("%s: %w", op, errAssertion)
	} else if !ok {
		return nil, fmt.Errorf("%s: %w", op, ErrNoKeyFound)
	}

	inf.IsAdmin, ok = claims[adminKey].(bool)
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, ErrNoKeyFound)
	}

	return &inf, nil
}

func Access(token string) (bool, error) {
	const op = scope + "Access"

	claims, err := jwt.Parse(token, getKey)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if !claims.Valid {
		return false, fmt.Errorf("%s: %w", op, ErrInvalidToken)
	}

	return claims.Valid, nil
}

func IsAdmin(token string) (bool, error) {
	const op = scope + "IsAdmin"

	claims, err := castToMapClaims(token)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	res, ok := claims[adminKey].(bool)
	if !ok {
		return false, fmt.Errorf("%s: %w", op, errAssertion)
	}

	return res, nil
}

func GetUsername(token string) (string, error) {
	const op = scope + "GetUsername"

	claims, err := castToMapClaims(token)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	username, ok := claims[usernameKey].(string)
	if !ok {
		return "", fmt.Errorf("%s: %w", op, errAssertion)
	} else if username == "" {
		return "", fmt.Errorf("%s: %w", op, ErrNoKeyFound)
	}
	return username, nil
}

func GetID(token string) (int64, error) {
	const op = scope + "GetID"

	claims, err := castToMapClaims(token)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	idInterface, ok := claims[idKey]
	if !ok {
		return 0, fmt.Errorf("%s: %w", op, errValueNotFound)
	}

	id, ok := idInterface.(int64)
	if !ok {
		return 0, fmt.Errorf("%s: %w", op, errAssertion)
	} else if id == 0 {
		return 0, fmt.Errorf("%s: %w", op, ErrNoKeyFound)
	}

	return id, nil
}

func castToMapClaims(token string) (jwt.MapClaims, error) {
	const op = scope + "castToMapClaims"

	tkn, err := jwt.Parse(token, getKey)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, errAssertion)
	}

	return claims, nil
}

func getKey(token *jwt.Token) (interface{}, error) {
	const op = scope + "getKey"
	key, ok := os.LookupEnv(secretKeyEnv)
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, ErrNoKeyFound)
	}
	return []byte(key), nil
}
