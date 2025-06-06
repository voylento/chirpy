package auth

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

const MaxBcryptPasswordBytes = 72

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:				"chirpy",
		IssuedAt:			jwt.NewNumericDate(time.Now()),
		ExpiresAt:		jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:			userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signedString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected jwt signing method: %v", token.Header["alg"])
		}

		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("Invalid token")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("Error extracting jwt claims")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

var ErrorInvalidAuthHeader= errors.New("Invalid Authorization Header")

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")	
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", ErrorInvalidAuthHeader
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", ErrorInvalidAuthHeader
	}

	return token, nil
}


