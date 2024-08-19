package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Note struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

var notes = []Note{}
var nextID = 1
var mu sync.Mutex

func getNotesHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}
func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var newNote Note
	err := json.NewDecoder(r.Body).Decode(&newNote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newNote.ID = nextID
	nextID++
	notes = append(notes, newNote)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newNote)
}
func main() {
	http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getNotesHandler(w, r)
		case "POST":
			createNoteHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	fmt.Println("Server is listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
