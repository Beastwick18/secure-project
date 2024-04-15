package auth

import (
	"log"
	"net/http"
	"strings"
)

type User struct {
	name   string
	apiKey string
	read   bool
	write  bool
}

type Auth struct {
	users []User
}

func (auth *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		found := false
		token := r.Header.Get("Authorization")
		splits := strings.SplitN(token, " ", 2)
		if len(splits) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad authorization format"))
			return
		}
		token = splits[1]
		for _, user := range auth.users {
			if user.apiKey == token {
				found = true
				break
			}
		}
		if !found {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unknown authorization"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (auth *Auth) Populate() {
	auth.users = []User{
		{name: "normal", apiKey: "key1", read: true, write: false},
		{name: "sudo", apiKey: "key2", read: true, write: true},
	}
}

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"strings"
//
// 	"github.com/golang-jwt/jwt/v5"
// )
//
// type User []struct {
// 	Username    string   `json:"name"`
// 	Permissions []string `json:"permissions"`
// }
//
// func Auth(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		token := r.Header.Get("Authorization")
// 		if token == "" {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			return
// 		}
//
// 		token = strings.Replace(token, "Bearer ", "", 1)
//
// 		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("Unexpected signing method")
// 			}
// 			return []byte("secret"), nil
// 		})
//
// 		if err != nil || !parsedToken.Valid {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			return
// 		}
//
// 		claims, ok := parsedToken.Claims.(jwt.MapClaims)
// 		if !ok {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			return
// 		}
//
// 		userID, ok := claims["user_id"].(string)
// 		if !ok {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			return
// 		}
//
// 		ctx := context.WithValue(r.Context(), "user_id", userID)
// 		h.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }
//
// func GetUsers() {
//
// }
