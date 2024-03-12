package data

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwa"
)

var TokenAuth *jwtauth.JWTAuth
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

var BookList BookDB
var AuthorList AuthorDB
var CredList CredDB

func Init() {
	// initializing the containers
	BookList = make(BookDB)
	AuthorList = make(AuthorDB)
	CredList = make(CredDB)

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
		AuthorList[val.ID] = val
	}

	// Sample data for books
	books := []Book{
		{ID: "1", Title: "The Shining", Genre: "Horror", AuthorID: "1"},
		{ID: "2", Title: "Harry Potter and the Sorcerer's Stone", Genre: "Fantasy", AuthorID: "2"},
		{ID: "3", Title: "Murder on the Orient Express", Genre: "Mystery", AuthorID: "3"},
		{ID: "4", Title: "1984", Genre: "Dystopian", AuthorID: "4"},
		{ID: "5", Title: "The Old Man and the Sea", Genre: "Fiction", AuthorID: "5"},
		{ID: "6", Title: "Hamlet", Genre: "Tragedy", AuthorID: "6"},
		{ID: "7", Title: "The Adventures of Tom Sawyer", Genre: "Adventure", AuthorID: "7"},
		{ID: "8", Title: "To Kill a Mockingbird", Genre: "Fiction", AuthorID: "8"},
		{ID: "9", Title: "The Hobbit", Genre: "Fantasy", AuthorID: "9"},
		{ID: "10", Title: "Pride and Prejudice", Genre: "Romance", AuthorID: "10"},
		{ID: "11", Title: "It", Genre: "Horror", AuthorID: "1"},
		{ID: "12", Title: "The Stand", Genre: "Post-apocalyptic", AuthorID: "1"},
		{ID: "13", Title: "The Casual Vacancy", Genre: "Fiction", AuthorID: "2"},
		{ID: "14", Title: "Harry Potter and the Chamber of Secrets", Genre: "Fantasy", AuthorID: "2"},
		{ID: "15", Title: "Death on the Nile", Genre: "Mystery", AuthorID: "3"},
		{ID: "16", Title: "Animal Farm", Genre: "Dystopian", AuthorID: "4"},
		{ID: "17", Title: "A Farewell to Arms", Genre: "War", AuthorID: "5"},
	}
	for _, val := range books {
		BookList[val.ID] = val
	}
	Creds := []Credential{
		{UserName: "parvej", Password: "1234"},
		{UserName: "sabbir", Password: "1234"},
		{UserName: "zayed", Password: "1234"},
	}
	for _, val := range Creds {
		CredList[val.UserName] = val
	}
	InitToken()
}
func InitToken() {
	TokenAuth = jwtauth.New(string(jwa.HS256), Secret, nil)
}
