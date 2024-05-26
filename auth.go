package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func WithJWTAuth(handlerFunc http.HandlerFunc, store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the token from req... from (Auth Header)
		tokenString := GetTokenFromRequest(r)

		// validate the token
		token, err := validateJWT(tokenString)
		if err != nil {
			log.Println("error validating token: ", err)
			WriteJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "permission denied"})
			return
		}

		if !token.Valid {
			log.Println("token is invalid")
			WriteJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "invalid jwt token"})
			return
		}

		// get the user_id from the token
		claims := token.Claims.(jwt.MapClaims)
		id := claims["user_id"].(string)

		u, err := store.GetUserByID(id)
		if err != nil {
			log.Println("error getting user by id: ", err)
			WriteJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "permission denied"})
			return
		}

		// case where jwt is valid, but credentials are not found in the db
		if u == nil {
			log.Println("user not found")
			WriteJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "invalid jwt token"})
			return
		}

		// then, call the HandlerFunc and continue with the endpoint...
		handlerFunc(w, r)
	}
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}

func validateJWT(t string) (*jwt.Token, error) {
	secret := Envs.JWTSecret

	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func HashPassword(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CreateJWT(secret []byte, id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    strconv.Itoa(int(id)),
		"expires_at": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
