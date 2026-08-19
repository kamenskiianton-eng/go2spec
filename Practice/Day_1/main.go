/*
	Добавить реализацию:

	- PUT /books/{id}/
	- DELETE /books/{id}/
	- DELETE /books/

	- go 1.25.8
	- Проверок излишних не добавлял (пустое имя, автор и т.п.)
	- Судя по всему корректнее без слэшей в конце писать URL'ы

	Гугл подсказывает (проверено):
	В роутере Go встроенный http.NewServeMux трактует пути слэшем на конце как префиксы (директории), а без слэша — как точные совпадения.
	Маршрут "GET /books/" будет ловить не только /books, но и /books/anything/else/inside.

	Но не суть, вероятно решается перенапрвлением, например, на уровне httpd / nginx
*/

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
	w.Header().Set("Content-Type", "application/json")

	var newBook Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверим, что книги с таким ид не существует
	idx := slices.IndexFunc(books, func(b Book) bool {
		return b.ID == newBook.ID
	})
	if idx != -1 {
		http.Error(w, fmt.Sprintf("Book with id=%s already exists", newBook.ID), http.StatusBadRequest)
		return
	}

	books = append(books, newBook)
	json.NewEncoder(w).Encode(newBook)
}

func deleteBookById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	idx := slices.IndexFunc(books, func(b Book) bool {
		return b.ID == id
	})
	if idx == -1 {
		http.Error(w, fmt.Sprintf("Book with id=%s is not found", id), http.StatusNotFound)
		return
	}

	books = slices.DeleteFunc(books, func(b Book) bool {
		return b.ID == id
	})

	w.WriteHeader(http.StatusOK)
}

func deleteBooks(w http.ResponseWriter, r *http.Request) {
	books = []Book{}

	w.WriteHeader(http.StatusOK)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")

	idx := slices.IndexFunc(books, func(b Book) bool {
		return b.ID == id
	})
	if idx == -1 {
		http.Error(w, fmt.Sprintf("Book with id=%s is not found", id), http.StatusNotFound)
		return
	}

	var updateBook Book
	err := json.NewDecoder(r.Body).Decode(&updateBook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Предполагаем, что id в url корректный, id в body перезапишем, если есть
	updateBook.ID = id
	books[idx] = updateBook

	json.NewEncoder(w).Encode(updateBook)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /books", getBooks)
	mux.HandleFunc("GET /books/{id}", getBookById)
	mux.HandleFunc("POST /books", createBook)
	mux.HandleFunc("DELETE /books/{id}", deleteBookById)
	mux.HandleFunc("DELETE /books", deleteBooks)
	mux.HandleFunc("PUT /books/{id}", updateBook)

	fmt.Printf("Starting web server on port %s", port)
	log.Fatal(http.ListenAndServe("localhost"+port, mux))
}
