package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Username string

var UserName Username = "Username"

var secretKey = []byte("internal_API_key")

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error decoding the body")
		return
	}
	if user.Username == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Username or password missing")
		return
	}
	//add dummy username and password to db so that when we come here we can check if the user exists
	// and has entered the correct password
	expTime := time.Now().Add(time.Minute * 5)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("error generating the JWT token")
		return
	}
	w.Header().Set("Auth", token)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Auth")

		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil || !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(context.Background(), UserName, claims.Username)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
