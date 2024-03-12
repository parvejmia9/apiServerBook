package main

import (
	"apiServerBook/auth"
	"apiServerBook/data"
	"apiServerBook/rest"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"log"
	"net/http"
)

func main() {
	data.Init()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Post("/login", auth.LogIn)
	r.Post("/logout", auth.LogOut)

	r.Group(func(r chi.Router) {
		r.Route("/books", func(r chi.Router) {
			r.Get("/", rest.GetAllBooks)
			r.Get("/{id}", rest.GetOneBook)
			r.Group(func(r chi.Router) {
				// need to add authentication
				r.Use(jwtauth.Verifier(data.TokenAuth))
				r.Use(jwtauth.Authenticator(data.TokenAuth))

				r.Post("/", rest.NewBook)
				r.Put("/{id}", rest.UpdateBook)
				r.Delete("/{id}", rest.DeleteBook)
			})
		})
		r.Route("/authors", func(r chi.Router) {
			r.Get("/", rest.GetAllAuthors)
			r.Get("/{id}", rest.GetOneAuthor)
		})
		r.Get("/search/{sToken}", rest.Search)
	})

	fmt.Println("Listening and Serving to 9090")
	err := http.ListenAndServe("localhost:9090", r)
	if err != nil {
		log.Fatalln(err)
		return
	}
}
