package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type PhoneBookList struct {
	Username    string `json:"name"`
	PhoneNumber string `json:"phone"`
}

var phoneBook []PhoneBookList

func main() {

	// router := mux.NewRouter()
	//
	// router.HandleFunc("/PhoneBook/list", retreiveAllEntries).Methods("GET")
	// router.HandleFunc("/PhoneBook/add", insertNewPhonebook).Methods("POST")
	// router.HandleFunc("/PhoneBook/deleteByName", deletePhonebookEntryByName).Methods("PUT").Queries("name", "{name}")
	// router.HandleFunc("/PhoneBook/deleteByNumber", deletePhonebookEntryByNumber).Methods("PUT").Queries("number", "{number}")
	//
	log.Println("Starting server...")
	// log.Fatal(http.ListenAndServe(":8080", router))
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

	phoneBook = append(phoneBook, entry)
	w.WriteHeader(http.StatusOK)
}

func deletePhonebookEntryByName(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	params := mux.Vars(r)
	name := params["name"]
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Name is a mandatory field"))
		return
	}

	// Delete entry from in-memory phonebook
	deleted := false
	for i, entry := range phoneBook {
		if entry.Username == name {
			phoneBook = append(phoneBook[:i], phoneBook[i+1:]...)
			deleted = true
			break
		}
	}
	if !deleted {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Name  %s not found", name)))
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
		w.Write([]byte(fmt.Sprintf("Entered number %s not found", number)))
		return
	}

	w.WriteHeader(http.StatusOK)
}
