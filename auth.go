package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func WithJWTAuth(handlerFunc http.HandlerFunc, store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the request header "Authorization" key
		tokenString := GetTokenFromRequest(r)

		// validate the token
		token, err := validateJWT(tokenString)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
			return
		}

		if !token.Valid {
			WriteJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
			return
		}

		// get the userID from the token
		claims := token.Claims.(jwt.MapClaims)
		userID := claims["user_id"].(string)

		_, err = store.GetUserByID(userID)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "Failed to get User"})
			return
		}

		// call the handlerFunc and continue the request
		handlerFunc(w, r)
	}
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")

	if tokenAuth == "" {
		return ""
	}

	return tokenAuth
}

func validateJWT(t string) (token *jwt.Token, err error) {
	secret := Envs.JWTSecret
	return jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}

		return []byte(secret), nil
	})
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CreateJWT(secret []byte, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   userID,
		"expiresAt": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
