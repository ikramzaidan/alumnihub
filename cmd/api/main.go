package main

import (
	"alumnihub/internal/auth"
	"alumnihub/internal/config"
	"alumnihub/internal/handler"
	"alumnihub/internal/repository"
	"alumnihub/internal/service"
	"flag"
	"fmt"
	"log"
	"net/http"
)

const port = 8080

type application struct {
	Config   config.AppConfig
	DB       repository.DatabaseRepo
	Auth     *auth.Auth
	Handlers *handler.Handler
}

func main() {
	var app application

	flag.StringVar(&app.Config.DSN, "dsn", "postgres://alumnihub:alumnihub@postgres:5432/alumnihub?sslmode=disable", "Postgres connection")
	flag.StringVar(&app.Config.JWTSecret, "jwt-secret", "verysecret", "signing secret")
	flag.StringVar(&app.Config.JWTIssuer, "jwt-issuer", "example.com", "signing issuer")
	flag.StringVar(&app.Config.JWTAudience, "jwt-audience", "example.com", "signing audience")
	flag.StringVar(&app.Config.CookieDomain, "cookie-domain", "localhost", "cookie domain for local development")
	flag.StringVar(&app.Config.Domain, "domain", "example.com", "Domain")
	flag.Parse()

	conn, err := openDB(app.Config.DSN)
	if err != nil {
		log.Fatal(err)
	}

	app.DB = &repository.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	app.Auth = auth.New(app.Config)

	authService := service.NewAuthService(app.DB, app.Auth)
	profileService := service.NewProfileService(app.DB)
	alumniService := service.NewAlumniService(app.DB)
	articleService := service.NewArticleService(app.DB)
	formService := service.NewFormService(app.DB)
	forumService := service.NewForumService(app.DB)
	jobService := service.NewJobService(app.DB)
	dashboardService := service.NewDashboardService(app.DB)

	app.Handlers = handler.NewHandler(
		authService,
		profileService,
		alumniService,
		articleService,
		formService,
		forumService,
		jobService,
		dashboardService,
		app.Auth,
	)

	log.Println("Starting application on", port)

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("/home/ikramzaidann/alumnihub/public"))))

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), routes(app.Handlers, app.Auth))
	if err != nil {
		log.Fatal(err)
	}
}
