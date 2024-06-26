package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"secure/pkg/util"
	"strings"

	"github.com/gorilla/mux"
)

type Permission int

const (
	Read  = "read"
	Write = "write"
)

type User struct {
	Name        string   `json:"name"`
	Token       string   `json:"token"`
	Permissions []string `json:"permissions"`
}

var PermissionShorthand map[string]string = map[string]string{
	Read:  "r",
	Write: "w",
}

type Auth struct {
	Users []User `json:"users"`
}

func (auth *Auth) Middleware(permissions []string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			splits := strings.SplitN(token, " ", 2)
			if len(splits) != 2 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Bad authorization format"))
				return
			}
			token = splits[1]
			var u *User
			for _, user := range auth.Users {
				if user.Token == token {
					u = &user
					break
				}
			}
			if u == nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unknown authorization"))
				return
			}
			for _, p := range permissions {
				if !util.Contains(u.Permissions, p) {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(fmt.Sprintf("Missing %s permissions", p)))
					return
				}
			}
			r.Header.Set("Authorization", u.Name)
			next.ServeHTTP(w, r)
		})
	}
}

func (auth *Auth) Populate(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to load users:\n%s", err)
	}
	err = json.Unmarshal(file, &auth.Users)
	if err != nil {
		log.Fatalf("Failed to read users.json:\n%s", err)
	}

	for _, u := range auth.Users {
		for _, p := range u.Permissions {
			if _, ok := PermissionShorthand[p]; !ok {
				log.Printf(`Unknown permission: "%s"`, p)
			}
		}
	}
}
