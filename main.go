package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"secure/pkg/auth"
	"secure/pkg/database"
	"secure/pkg/name"
	"secure/pkg/phone"

	"github.com/gorilla/mux"
)

type Context struct {
	pb    *database.PhoneBook
	audit *log.Logger
}

func main() {
	router := mux.NewRouter()
	read := router.NewRoute().Subrouter()
	readwrite := router.NewRoute().Subrouter()

	phonebookPath := "phonebook.db"
	if len(os.Args) > 1 {
		phonebookPath = os.Args[1]
	}

	auditPath := "audit.log"
	if len(os.Args) > 2 {
		auditPath = os.Args[2]
	}

	os.MkdirAll(filepath.Dir(auditPath), os.ModePerm)
	auditFile, err := os.OpenFile(auditPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Failed to create audit.log:\n%s", err)
	}
	w := io.MultiWriter(os.Stdout, auditFile)

	pb, err := database.CreateTable(phonebookPath)
	if err != nil {
		log.Fatalf("Failed to open database:\n%s", err)
	}
	log.Printf("Successfully loaded %s", phonebookPath)
	ctx := Context{
		pb: pb,
	}
	defer ctx.pb.Close()

	ctx.audit = log.New(w, "<Audit> ", log.LstdFlags)

	a := auth.Auth{}
	a.Populate()
	read.Use(a.Middleware([]string{auth.Read}))
	readwrite.Use(a.Middleware([]string{auth.Read, auth.Write}))
	router.Use(LogMiddleware)
	read.HandleFunc("/PhoneBook/list", ctx.retreiveAllEntries).Methods("GET")
	readwrite.HandleFunc("/PhoneBook/add", ctx.insertNewPhonebook).Methods("POST")
	readwrite.HandleFunc("/PhoneBook/deleteByName", ctx.deletePhonebookEntryByName).Methods("PUT").Queries("name", "{name}")
	readwrite.HandleFunc("/PhoneBook/deleteByNumber", ctx.deletePhonebookEntryByNumber).Methods("PUT").Queries("number", "{number}")

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("%s %s", r.Method, r.RequestURI))
		next.ServeHTTP(w, r)
	})
}

func (ctx *Context) retreiveAllEntries(w http.ResponseWriter, r *http.Request) {
	entries, err := ctx.pb.ListAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json_string, err := json.Marshal(entries)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx.audit.Printf(`[%s] List all users`, r.Header.Get("Authorization"))
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_string)
}

func (ctx *Context) insertNewPhonebook(w http.ResponseWriter, r *http.Request) {
	var entry database.UserEntry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON format"))
		return
	}

	if !name.ValidName(entry.Name) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid name"))
		return
	}

	if !phone.ValidPhone(entry.Phone) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid phone number"))
		return
	}

	err = ctx.pb.Append(&entry)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Failed to insert new entry into database")
		return
	}

	ctx.audit.Printf(`[%s] Added user "%s"`, r.Header.Get("Authorization"), entry.Name)
	w.WriteHeader(http.StatusOK)
}

func (ctx *Context) deletePhonebookEntryByName(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	params := mux.Vars(r)
	name_str := params["name"]
	if name_str == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Name is a mandatory field"))
		return
	}

	if !name.ValidName(name_str) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid phone number"))
		return
	}

	// Delete entry from database
	deleted, err := ctx.pb.DeleteByName(name_str)
	if err != nil {
		log.Printf("Error while trying to delete by name:\n%s", err)
	}
	if !deleted {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Name not found"))
		return
	}

	ctx.audit.Printf(`[%s] Removed user "%s"`, r.Header.Get("Authorization"), name_str)
	w.WriteHeader(http.StatusOK)
}

func (ctx *Context) deletePhonebookEntryByNumber(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	params := mux.Vars(r)
	number := params["number"]
	if number == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Number not provided"))
		return
	}

	if !phone.ValidPhone(number) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid phone number"))
		return
	}

	name, deleted, err := ctx.pb.DeleteByPhone(number)
	if err != nil {
		log.Printf("Error while trying to delete by phone number:\n%s", err)
	}
	if !deleted {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Entered number not found"))
		return
	}

	ctx.audit.Printf(`[%s] Removed user "%s"`, r.Header.Get("Authorization"), name)
	w.WriteHeader(http.StatusOK)
}
