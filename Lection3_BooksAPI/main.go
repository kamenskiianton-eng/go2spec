package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
)

const port = ":3000"

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book

func init() {
	books = []Book{
		{ID: "1", Title: "Мастер и Маргарита", Author: "Михаил Булгаков"},
		{ID: "2", Title: "1984", Author: "Джордж Оруэлл"},
		{ID: "3", Title: "Маленький принц", Author: "Антуан де Сент-Экзюпери"},
	}
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	idx := slices.IndexFunc(books, func(b Book) bool {
		return b.ID == id
	})
	if idx != -1 {
		json.NewEncoder(w).Encode(books[idx])
		return
	}

	http.Error(w, fmt.Sprintf("Book with id=%s is not found", id), http.StatusNotFound)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var newBook Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	books = append(books, newBook)
	json.NewEncoder(w).Encode(newBook)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /books/", getBooks)
	mux.HandleFunc("GET /books/{id}/", getBookById)
	mux.HandleFunc("POST /books/", createBook)

	fmt.Printf("Starting web server on port %s", port)
	log.Fatal(http.ListenAndServe("localhost"+port, mux))
}
