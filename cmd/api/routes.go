package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Get("/", app.Home)
	mux.Post("/authenticate", app.authenticate)
	mux.Get("/refresh", app.refreshToken)
	mux.Get("/logout", app.logout)

	mux.Route("/", func(mux chi.Router) {
		mux.Use(app.authRequired)
		mux.Get("/alumni", app.AllAlumni)
		mux.Get("/alumni/{id}", app.Alumni)
		mux.Get("/articles", app.allArticles)
		mux.Get("/articles/{id}", app.article)
		mux.Post("/articles/create", app.insertArticle)
		mux.Get("/forms", app.allForms)
		mux.Get("/forms/{id}", app.form)
		mux.Put("/forms/create", app.insertForm)
		mux.Get("/questions/{id}", app.question)
		mux.Post("/questions/create", app.insertQuestion)
	})

	return mux
}
