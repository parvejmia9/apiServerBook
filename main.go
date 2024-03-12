package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwa"
	"log"
	"net/http"
	"strconv"
	"time"
)

var tokenAuth *jwtauth.JWTAuth
var Secret = []byte("this_is_a_secret_key")

type Book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Genre    string `json:"genre"`
	AuthorID string `json:"authorID"`
}
type Author struct {
	ID        string `json:"id"`
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
}

type Credential struct {
	UserName string `json:"uname"`
	Password string `json:"password"`
}

type BookDB map[string]Book
type AuthorDB map[string]Author
type CredDB map[string]Credential

var bookList BookDB
var authorList AuthorDB
var credList CredDB

func Init() {
	// initializing the containers
	bookList = make(BookDB)
	authorList = make(AuthorDB)
	credList = make(CredDB)

	// Sample data for authors
	authors := []Author{
		{ID: "1", FirstName: "Stephen", LastName: "King"},
		{ID: "2", FirstName: "J.K.", LastName: "Rowling"},
		{ID: "3", FirstName: "Agatha", LastName: "Christie"},
		{ID: "4", FirstName: "George", LastName: "Orwell"},
		{ID: "5", FirstName: "Ernest", LastName: "Hemingway"},
		{ID: "6", FirstName: "William", LastName: "Shakespeare"},
		{ID: "7", FirstName: "Mark", LastName: "Twain"},
		{ID: "8", FirstName: "Harper", LastName: "Lee"},
		{ID: "9", FirstName: "J.R.R.", LastName: "Tolkien"},
		{ID: "10", FirstName: "Jane", LastName: "Austen"},
	}
	for _, val := range authors {
		authorList[val.ID] = val
	}

	// Sample data for books
	books := []Book{
		{ID: "1", Title: "The Shining", Genre: "Horror", AuthorID: "1"},
		{ID: "2", Title: "Harry Potter and the Sorcerer's Stone", Genre: "Fantasy", AuthorID: "2"},
		{ID: "3", Title: "Murder on the Orient Express", Genre: "Mystery", AuthorID: "3"},
		{ID: "4", Title: "1984", Genre: "Dystopian", AuthorID: "4"},
		{ID: "5", Title: "The Old Man and the Sea", Genre: "Fiction", AuthorID: "5"},
		//{ID: "6", Title: "Hamlet", Genre: "Tragedy", AuthorID: "6"},
		//{ID: "7", Title: "The Adventures of Tom Sawyer", Genre: "Adventure", AuthorID: "7"},
		//{ID: "8", Title: "To Kill a Mockingbird", Genre: "Fiction", AuthorID: "8"},
		//{ID: "9", Title: "The Hobbit", Genre: "Fantasy", AuthorID: "9"},
		//{ID: "10", Title: "Pride and Prejudice", Genre: "Romance", AuthorID: "10"},
		//{ID: "11", Title: "It", Genre: "Horror", AuthorID: "1"},
		//{ID: "12", Title: "The Stand", Genre: "Post-apocalyptic", AuthorID: "1"},
		//{ID: "13", Title: "The Casual Vacancy", Genre: "Fiction", AuthorID: "2"},
		//{ID: "14", Title: "Harry Potter and the Chamber of Secrets", Genre: "Fantasy", AuthorID: "2"},
		//{ID: "15", Title: "Death on the Nile", Genre: "Mystery", AuthorID: "3"},
		//{ID: "16", Title: "Animal Farm", Genre: "Dystopian", AuthorID: "4"},
		//{ID: "17", Title: "A Farewell to Arms", Genre: "War", AuthorID: "5"},
	}
	for _, val := range books {
		bookList[val.ID] = val
	}
	Creds := []Credential{
		{UserName: "parvej", Password: "1234"},
		{UserName: "sabbir", Password: "1234"},
		{UserName: "zayed", Password: "1234"},
	}
	for _, val := range Creds {
		credList[val.UserName] = val
	}
	InitToken()
}
func InitToken() {
	tokenAuth = jwtauth.New(string(jwa.HS256), Secret, nil)
}

func getAllBooks(w http.ResponseWriter, _ *http.Request) {
	var books []Book
	for _, val := range bookList {
		books = append(books, val)
	}
	err := json.NewEncoder(w).Encode(books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func getAllAuthors(w http.ResponseWriter, _ *http.Request) {
	var authors []Author
	for _, val := range authorList {
		authors = append(authors, val)
	}
	err := json.NewEncoder(w).Encode(authors)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
func getOneBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	val, ok := bookList[id]
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
func getOneAuthor(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	val, ok := authorList[id]
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
func isValidFormat(w http.ResponseWriter, bk Book) bool {
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
	_, ok := authorList[bk.AuthorID]
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
func newBook(w http.ResponseWriter, r *http.Request) {
	var bk Book
	err := json.NewDecoder(r.Body).Decode(&bk)
	if err != nil {
		http.Error(w, "Bad Format", http.StatusBadRequest)
		return
	}
	if !isValidFormat(w, bk) {
		return
	}
	_, ok := bookList[bk.ID]
	if ok {
		http.Error(w, "Book ID Already Exists try another one", http.StatusBadRequest)
		return
	}

	bookList[bk.ID] = bk
	_, err = w.Write([]byte("Data added successfully"))
	if err != nil {
		http.Error(w, "Can not Write Data", http.StatusInternalServerError)
		return
	}

}
func updateBook(w http.ResponseWriter, r *http.Request) {
	oldID := chi.URLParam(r, "id")
	if len(oldID) == 0 {
		http.Error(w, "Give a ID", http.StatusBadRequest)
		return
	}
	var bk Book
	err := json.NewDecoder(r.Body).Decode(&bk)
	if err != nil {
		http.Error(w, "Bad Format", http.StatusBadRequest)
		return
	}
	if !isValidFormat(w, bk) {
		return
	}
	_, ok := bookList[oldID]
	if !ok {
		http.Error(w, "Book with the given ID not exists.", http.StatusBadRequest)
		return
	}
	if oldID != bk.ID {
		http.Error(w, "You cannot change ID of a book.", http.StatusBadRequest)
		return
	}
	bookList[oldID] = bk
	_, err = w.Write([]byte("Data updated successfully"))
	if err != nil {
		http.Error(w, "Can not Write Data", http.StatusInternalServerError)
		return
	}

}
func deleteBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if len(id) == 0 {
		http.Error(w, "Give a ID", http.StatusBadRequest)
		return
	}
	_, ok := bookList[id]
	if !ok {
		http.Error(w, "Book Not Found", http.StatusBadRequest)
		return
	}
	delete(bookList, id)
	_, err := w.Write([]byte("Data Deleted successfully"))
	if err != nil {
		http.Error(w, "Can not Write Data", http.StatusInternalServerError)
		return
	}
}
func logIn(w http.ResponseWriter, r *http.Request) {
	var tmp Credential
	fmt.Println(r.Body)
	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		http.Error(w, "Cannot Decode", http.StatusBadRequest)
		return
	}
	cred, okay := credList[tmp.UserName]
	if !okay {
		http.Error(w, "User Name doesn't Exist", http.StatusBadRequest)
		return
	}
	if cred.Password != tmp.Password {
		http.Error(w, "Password didn't matched", http.StatusBadRequest)
		return
	}
	et := time.Now().Add(15 * time.Minute)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"aud": "Parvej Mia",
		"exp": et.Unix(),
	})
	fmt.Println(tokenString)
	if err != nil {
		http.Error(w, "Can not generate jwt", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Value:   tokenString,
		Expires: et,
	})
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Successfully Logged In"))
	if err != nil {
		http.Error(w, "Can not Write Data", http.StatusInternalServerError)
		return
	}

}
func logOut(w http.ResponseWriter, _ *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Expires: time.Now(),
	})
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Successfully Logged Out"))
	if err != nil {
		http.Error(w, "Can not Write Data", http.StatusInternalServerError)
		return
	}
}
func main() {
	Init()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Post("/login", logIn)
	r.Post("/logout", logOut)

	r.Group(func(r chi.Router) {
		r.Route("/books", func(r chi.Router) {
			r.Get("/", getAllBooks)
			r.Get("/{id}", getOneBook)
			r.Group(func(r chi.Router) {
				// need to add authentication
				r.Use(jwtauth.Verifier(tokenAuth))
				r.Use(jwtauth.Authenticator(tokenAuth))

				r.Post("/", newBook)
				r.Put("/{id}", updateBook)
				r.Delete("/{id}", deleteBook)
			})
		})
		r.Route("/authors", func(r chi.Router) {
			r.Get("/", getAllAuthors)
			r.Get("/{id}", getOneAuthor)
		})
	})

	fmt.Println("Listening and Serving to 9090")
	err := http.ListenAndServe("localhost:9090", r)
	if err != nil {
		log.Fatalln(err)
		return
	}
}
