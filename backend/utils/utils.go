package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/xhermitx/gitpulse-01/backend/config"
)

const TOKEN_EXPIRATION_TIME = time.Hour * 24 * 30

func ResponseWriter(w http.ResponseWriter, status int, msg any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(msg)
}

func ErrResponseWriter(w http.ResponseWriter, status int, err error) {
	ResponseWriter(w, status, map[string]string{"error": err.Error()})
}

func ParseRequestBody(r *http.Request, payload any) error {
	if r.Body == nil {
		return errors.New("empty request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func GenerateToken(id string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(TOKEN_EXPIRATION_TIME).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(config.Envs.AuthSecret))
	if err != nil {
		log.Fatal(err)
	}

	return tokenString
}

func GetToken(r *http.Request) (string, error) {
	// Extract the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Println("missing authorization header")
		return "", fmt.Errorf("authorization header missing")
	}

	// Split the header to get the token part
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		log.Println("incorrect authorization header")
		return "", fmt.Errorf("authorization header format must be Bearer {token}")
	}

	return headerParts[1], nil
}
