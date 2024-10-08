package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/xhermitx/gitpulse-01/backend/config"
	"github.com/xhermitx/gitpulse-01/backend/types"
	"github.com/xhermitx/gitpulse-01/backend/utils"
)

func GenerateToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(config.Envs.JWTExpiration).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(config.Envs.AuthSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
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

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Envs.AuthSecret), nil
	})
}

func AuthMiddleware(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := GetToken(r)
		if err != nil {
			log.Println("Error fetching token from header")
			utils.ErrResponseWriter(w, http.StatusUnauthorized, err)
			return
		}

		token, err := ValidateToken(tokenString)
		if err != nil {
			log.Println("1 : Error validating token")
			utils.ErrResponseWriter(w, http.StatusUnauthorized, err)
			return
		}

		if !token.Valid {
			log.Println("2 : Error validating token")
			utils.ErrResponseWriter(w, http.StatusUnauthorized, err)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userId := claims["user_id"].(string)

		log.Println("user_id", userId)

		user, err := store.FindUserById(userId)
		if err != nil || user == nil {
			log.Println("Error finding user")
			utils.ErrResponseWriter(w, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}

		ctx := context.WithValue(context.Background(), types.UserContext("user_id"), userId)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}
