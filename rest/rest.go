package rest

import (
	"apiServerBook/data"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"strings"
)

func GetAllBooks(w http.ResponseWriter, _ *http.Request) {
	var books []data.Book
	for _, val := range data.BookList {
		books = append(books, val)
	}
	err := json.NewEncoder(w).Encode(books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func GetAllAuthors(w http.ResponseWriter, _ *http.Request) {
	var authors []data.Author
	for _, val := range data.AuthorList {
		authors = append(authors, val)
	}
	err := json.NewEncoder(w).Encode(authors)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
func GetOneBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	val, ok := data.BookList[id]
	if !ok {
		http.Error(w, "Book Not Found", http.StatusNotFound)
		return
	}
	err := json.NewEncoder(w).Encode(val)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
func GetOneAuthor(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	val, ok := data.AuthorList[id]
	if !ok {
		http.Error(w, "Author Not Found", http.StatusNotFound)
		return
	}
	err := json.NewEncoder(w).Encode(val)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
func isValidFormat(w http.ResponseWriter, bk data.Book) bool {
	if len(bk.ID) == 0 {
		http.Error(w, "Book ID cannot be empty", http.StatusBadRequest)
		return false
	}

	// checking for integer book ID
	_, err := strconv.ParseInt(bk.ID, 10, 64)
	if err != nil {
		http.Error(w, "Book ID must be an integer", http.StatusBadRequest)
		return false
	}

	// checking for the existence of the Book ID

	if len(bk.Title) == 0 {
		http.Error(w, "Title cannot be empty", http.StatusBadRequest)
		return false
	}

	if len(bk.AuthorID) == 0 {
		http.Error(w, "Author ID cannot be empty", http.StatusBadRequest)
		return false
	}

	// checking if the author is registered or not
	_, ok := data.AuthorList[bk.AuthorID]
	if !ok {
		http.Error(w, "Author ID is not registered. First register!", http.StatusBadRequest)
		return false
	}
	if len(bk.Genre) == 0 {
		http.Error(w, "Genre cannot be empty", http.StatusBadRequest)
		return false
	}
	return true
}
func NewBook(w http.ResponseWriter, r *http.Request) {
	var bk data.Book
	err := json.NewDecoder(r.Body).Decode(&bk)
	if err != nil {
		http.Error(w, "Bad Format", http.StatusBadRequest)
		return
	}
	if !isValidFormat(w, bk) {
		return
	}
	_, ok := data.BookList[bk.ID]
	if ok {
		http.Error(w, "Book ID Already Exists try another one", http.StatusBadRequest)
		return
	}

	data.BookList[bk.ID] = bk
	_, err = w.Write([]byte("data added successfully"))
	if err != nil {
		http.Error(w, "Can not Write data", http.StatusInternalServerError)
		return
	}

}
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	oldID := chi.URLParam(r, "id")
	if len(oldID) == 0 {
		http.Error(w, "Give a ID", http.StatusBadRequest)
		return
	}
	var bk data.Book
	err := json.NewDecoder(r.Body).Decode(&bk)
	if err != nil {
		http.Error(w, "Bad Format", http.StatusBadRequest)
		return
	}
	if !isValidFormat(w, bk) {
		return
	}
	_, ok := data.BookList[oldID]
	if !ok {
		http.Error(w, "Book with the given ID not exists.", http.StatusBadRequest)
		return
	}
	if oldID != bk.ID {
		http.Error(w, "You cannot change ID of a book.", http.StatusBadRequest)
		return
	}
	data.BookList[oldID] = bk
	_, err = w.Write([]byte("data updated successfully"))
	if err != nil {
		http.Error(w, "Can not Write data", http.StatusInternalServerError)
		return
	}

}
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if len(id) == 0 {
		http.Error(w, "Give a ID", http.StatusBadRequest)
		return
	}
	_, ok := data.BookList[id]
	if !ok {
		http.Error(w, "Book Not Found", http.StatusBadRequest)
		return
	}
	delete(data.BookList, id)
	_, err := w.Write([]byte("data Deleted successfully"))
	if err != nil {
		http.Error(w, "Can not Write data", http.StatusInternalServerError)
		return
	}
}
func Search(w http.ResponseWriter, r *http.Request) {
	type sItems struct {
		Authors []data.Author `json:"authors"`
		Books   []data.Book   `json:"books"`
	}
	var res sItems
	sToken := chi.URLParam(r, "sToken")
	for _, val := range data.BookList {
		if strings.Contains(val.Title, sToken) {
			res.Books = append(res.Books, val)
		}
	}
	for _, val := range data.AuthorList {
		if strings.Contains(val.FirstName, sToken) || strings.Contains(val.LastName, sToken) {
			res.Authors = append(res.Authors, val)
		}
	}
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
