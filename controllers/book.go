package controllers

import (
	"database/sql"
	"encoding/json"
	"go-simple-books-list-api-server/models"
	bookRepository "go-simple-books-list-api-server/repository/book"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Controller struct{}

var books []models.Book

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (c Controller) GetBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book models.Book
		books = []models.Book{}
		bookRepo := bookRepository.BookRepository{}

		books = bookRepo.GetBooks(db, book, books)

		addHeaders(w)
		json.NewEncoder(w).Encode(books)
	}
}

func (c Controller) GetBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book models.Book
		bookRepo := bookRepository.BookRepository{}
		params := mux.Vars(r)

		id, err := strconv.Atoi(params["id"])
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		book, err = bookRepo.GetBook(db, book, id)

		if err != nil {
			log.Println(err)
			writeError(w, err)
			return
		}

		addHeaders(w)
		json.NewEncoder(w).Encode(book)
	}
}

func (c Controller) AddBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book models.Book
		json.NewDecoder(r.Body).Decode(&book)

		bookRepo := bookRepository.BookRepository{}

		var bookID int
		bookID = bookRepo.AddBook(db, book)

		addHeaders(w)
		json.NewEncoder(w).Encode(bookID)
	}
}

func (c Controller) UpdateBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book models.Book
		json.NewDecoder(r.Body).Decode(&book)

		bookRepo := bookRepository.BookRepository{}
		rowsUpdated := bookRepo.UpdateBook(db, book)

		addHeaders(w)
		json.NewEncoder(w).Encode(rowsUpdated)
	}
}

func (c Controller) RemoveBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		id, err := strconv.Atoi(params["id"])
		logFatal(err)

		bookRepo := bookRepository.BookRepository{}
		rowsDeleted := bookRepo.RemoveBook(db, id)

		json.NewEncoder(w).Encode(rowsDeleted)
	}
}

func addHeaders(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
}

func writeError(w http.ResponseWriter, err error) {
	addHeaders(w)
	w.WriteHeader(http.StatusInternalServerError)

	var model models.ErrorModel
	model.ErrorMessage = err.Error()
	json.NewEncoder(w).Encode(model)
}
