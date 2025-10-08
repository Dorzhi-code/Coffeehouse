package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthRequest struct{
	Login string `json:"login"`
	Password string `json:"password"`
}

var secretKey = []byte("my_secret_key")
// const userCtxKey = context("user_id")

func LogingHandler(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		print(err.Error())
		return
	}

	// Провека логина и пароля
	if req.Login != "admin" || req.Password != "admin"{
		http.Error(w, "Invalid data", http.StatusUnauthorized)
		return
	}

	token, err := GenerateToken(1)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string] string{
		"token":token,
	})
}

func GenerateToken(userId int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "token do not provide", http.StatusUnauthorized)
			return 
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := ParseToken(tokenString)
		if err != nil || !token.Valid{
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return 
		}
		claims := token.Claims.(jwt.MapClaims)
		user_id := int(claims["user_id"].(float64))
		ctx := context.WithValue(r.Context(), "user_id", user_id)
		next(w, r.WithContext(ctx))
	}
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Get access")
}
