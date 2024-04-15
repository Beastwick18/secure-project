package main

import (
	"encoding/json"
	"log"
	"net/http"
	"secure/auth"
	"secure/name"
	"secure/phone"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type PhoneBookList struct {
	Username    string `json:"name"`
	PhoneNumber string `json:"phone"`
}

var phoneBook []PhoneBookList

func main() {
	router := mux.NewRouter()

	a := auth.Auth{}
	a.Populate()
	router.Use(a.Middleware)
	router.HandleFunc("/PhoneBook/list", retreiveAllEntries).Methods("GET")
	router.HandleFunc("/PhoneBook/add", insertNewPhonebook).Methods("POST")
	router.HandleFunc("/PhoneBook/deleteByName", deletePhonebookEntryByName).Methods("PUT").Queries("name", "{name}")
	router.HandleFunc("/PhoneBook/deleteByNumber", deletePhonebookEntryByNumber).Methods("PUT").Queries("number", "{number}")

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func retreiveAllEntries(w http.ResponseWriter, _ *http.Request) {
	entries, err := json.Marshal(phoneBook)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(entries)
}

func insertNewPhonebook(w http.ResponseWriter, r *http.Request) {
	var entry PhoneBookList
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON format"))
		return
	}

	if !name.ValidName(entry.Username) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid name"))
		return
	}

	if !phone.ValidPhone(entry.PhoneNumber) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid phone number"))
		return
	}

	phoneBook = append(phoneBook, entry)
	w.WriteHeader(http.StatusOK)
}

func deletePhonebookEntryByName(w http.ResponseWriter, r *http.Request) {
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

	// Delete entry from in-memory phonebook
	deleted := false
	for i, entry := range phoneBook {
		if entry.Username == name_str {
			phoneBook = append(phoneBook[:i], phoneBook[i+1:]...)
			deleted = true
			break
		}
	}
	if !deleted {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Name not found"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deletePhonebookEntryByNumber(w http.ResponseWriter, r *http.Request) {
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

	// Delete entry from in-memory phonebook
	deleted := false
	for i, entry := range phoneBook {
		if entry.PhoneNumber == number {
			phoneBook = append(phoneBook[:i], phoneBook[i+1:]...)
			deleted = true
			break
		}
	}
	if !deleted {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Entered number not found"))
		return
	}

	w.WriteHeader(http.StatusOK)
}
