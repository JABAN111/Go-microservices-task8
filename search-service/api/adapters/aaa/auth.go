package aaa

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	hashCost = 14
)

var (
	errInvalidAccount      = errors.New("invalid username or password")
	errFailedCreatingToken = errors.New("failed to create token")
	errTokenIsInvalid      = errors.New("token is invalid")
)

// Authentication, Authorization, Accounting
type AAA struct {
	users     map[string]string
	tokenTTL  time.Duration
	log       *slog.Logger
	secretKey []byte
	salt      string
}

func New(tokenTTL time.Duration, log *slog.Logger) (AAA, error) {
	adminUser, ok := os.LookupEnv("ADMIN_USER")
	if !ok {
		return AAA{}, fmt.Errorf("could not get admin user from env")
	}
	adminPassword, ok := os.LookupEnv("ADMIN_PASSWORD")
	if !ok {
		return AAA{}, fmt.Errorf("could not get admin password from env")
	}
	secretKeyStr, ok := os.LookupEnv("SECRET_KEY")
	if !ok {
		return AAA{}, fmt.Errorf("could not get secret key from env")
	}
	salt, ok := os.LookupEnv("SALT")
	if !ok {
		return AAA{}, fmt.Errorf("could not get salt from env")
	}

	hashPwd, err := hashPassword(adminPassword, salt)
	if err != nil {
		return AAA{}, fmt.Errorf("can't hash default admin: %w", err)
	}
	return AAA{
		users:     map[string]string{adminUser: hashPwd},
		tokenTTL:  tokenTTL,
		log:       log,
		secretKey: []byte(secretKeyStr),
		salt:      salt,
	}, nil
}

func (a AAA) Login(name, password string) (string, error) {
	storedHash, ok := a.users[name]
	if !ok {
		return "", errInvalidAccount
	}

	if !verifyPassword(password, storedHash, a.salt) {
		return "", errInvalidAccount
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	now := time.Now()
	claims["exp"] = now.Add(a.tokenTTL).Unix()
	claims["authorized"] = true
	claims["iss"] = "searchServiceApi"
	claims["sub"] = name
	claims["aud"] = "api"
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := token.SignedString(a.secretKey)
	if err != nil {
		return "", errFailedCreatingToken
	}

	return tokenString, nil
}

func (a AAA) Verify(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.secretKey, nil
	})
	if err != nil {
		a.log.Error("failed to parse token", slog.String("error", err.Error()))
		return errTokenIsInvalid
	}

	if !token.Valid {
		a.log.Warn("invalid token")
		return errTokenIsInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		a.log.Warn("invalid token claims")
		return errTokenIsInvalid
	}

	name, ok := claims["sub"]
	if !ok {
		a.log.Error("user doesn't contain sub", "name", name)
		return errTokenIsInvalid
	}
	if _, ok := a.users[name.(string)]; !ok {
		a.log.Error("user does not exist in the app", "name", name)
		return errTokenIsInvalid
	}

	fmt.Printf("Got token %v", token)
	return nil
}

func hashPassword(password, salt string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password+salt), hashCost)
	return string(bytes), err
}

func verifyPassword(password, hash, salt string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+salt))
	return err == nil
}
