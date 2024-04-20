package auth

import (
	"context"
	"fmt"
	"net/http"
	"secure/pkg/util"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestPopulate(t *testing.T) {
	var a Auth
	a.Populate("../../users.json")
	if len(a.Users) != 2 {
		t.Fatalf("There should only be 2 users, read and read-write")
	}
	var ro *User
	var rw *User
	if a.Users[0].Name == "read-only" {
		ro = &a.Users[0]
		rw = &a.Users[1]
	} else {
		ro = &a.Users[1]
		rw = &a.Users[0]
	}
	if ro == nil {
		t.Fatalf("Could not find read-only user")
	}
	if rw == nil {
		t.Fatalf("Could not find read-write user")
	}
	if !util.Contains(ro.Permissions, Read) {
		t.Fatalf("Read only permissions does not include %s", Read)
	}
	if !util.Contains(rw.Permissions, Read) {
		t.Fatalf("Read write permissions does not include %s", Read)
	}
	if !util.Contains(rw.Permissions, Write) {
		t.Fatalf("Read write permissions does not include %s", Write)
	}
	if len(ro.Token) <= 0 {
		t.Fatalf("Read only token is empty")
	}
	if len(rw.Token) <= 0 {
		t.Fatalf("Read only token is empty")
	}
}

func TestReadOnlyAllowed(t *testing.T) {
	s, readOnlyKey, _ := setupTest([]string{Read})

	if status := getResponse(t, readOnlyKey); status != http.StatusOK {
		t.Fatalf("Response code was %d, but should have been %d", status, http.StatusOK)
	}

	stopServer(t, s)
}

func TestReadOnlyDenied(t *testing.T) {
	s, readOnlyKey, _ := setupTest([]string{Write})

	if status := getResponse(t, readOnlyKey); status != http.StatusUnauthorized {
		t.Fatalf("Response code was %d, but should have been %d", status, http.StatusUnauthorized)
	}

	stopServer(t, s)
}

func TestReadWriteAllowed(t *testing.T) {
	s, _, readWriteKey := setupTest([]string{Read, Write})

	if status := getResponse(t, readWriteKey); status != http.StatusOK {
		t.Fatalf("Response code was %d, but should have been %d", status, http.StatusOK)
	}

	stopServer(t, s)
}

func TestReadWriteDenied(t *testing.T) {
	s, readOnlyKey, _ := setupTest([]string{Read, Write})

	if status := getResponse(t, readOnlyKey); status != http.StatusUnauthorized {
		t.Fatalf("Response code was %d, but should have been %d", status, http.StatusUnauthorized)
	}

	stopServer(t, s)
}

func TestUnknownUser(t *testing.T) {
	s, _, _ := setupTest([]string{Write})

	if status := getResponse(t, "unknown_key"); status != http.StatusUnauthorized {
		t.Fatalf("Response code was %d, but should have been %d", status, http.StatusUnauthorized)
	}

	stopServer(t, s)
}

func setupTest(permissions []string) (*http.Server, string, string) {
	readOnlyKey := "read_only_key"
	readWriteKey := "read_write_key"
	readOnly := User{Name: "read-only", Token: readOnlyKey, Permissions: []string{Read}}
	readWrite := User{Name: "read-write", Token: readWriteKey, Permissions: []string{Read, Write}}
	a := &Auth{Users: []User{readOnly, readWrite}}

	mw := a.Middleware(permissions)
	router := mux.NewRouter()
	router.Use(mw)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	s := &http.Server{
		Addr:         ":8081",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		if err := s.ListenAndServe(); err != nil {
			return
		}
	}()

	return s, readOnlyKey, readWriteKey
}

func getResponse(t *testing.T, key string) int {
	req, err := http.NewRequest("GET", "http://127.0.0.1:8081/", nil)
	if err != nil {
		t.Fatalf("Failed to build request:\n%s", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", key))
	client := &http.Client{}

	attempts := 0

	for attempts < 10 {
		resp, err := client.Do(req)

		if err != nil {
			client.Timeout = time.Second
			continue
		}
		return resp.StatusCode
	}
	return -1
}

func stopServer(t *testing.T, s *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		t.Fatalf("Failed to stop server:\n%s", err)
	}
}
