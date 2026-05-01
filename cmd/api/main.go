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
	"os"
	"strconv"

	"github.com/joho/godotenv"
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

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Build DSN from environment variables
	dbDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	flag.StringVar(&app.Config.DSN, "dsn", dbDSN, "Postgres connection")
	flag.StringVar(&app.Config.JWTSecret, "jwt-secret", os.Getenv("JWT_SECRET"), "signing secret")
	flag.StringVar(&app.Config.JWTIssuer, "jwt-issuer", os.Getenv("JWT_ISSUER"), "signing issuer")
	flag.StringVar(&app.Config.JWTAudience, "jwt-audience", os.Getenv("JWT_AUDIENCE"), "signing audience")
	flag.StringVar(&app.Config.CookieDomain, "cookie-domain", os.Getenv("COOKIE_DOMAIN"), "cookie domain for local development")
	flag.StringVar(&app.Config.Domain, "domain", os.Getenv("DOMAIN"), "Domain")

	// SMTP flags
	flag.StringVar(&app.Config.SMTPHost, "smtp-host", os.Getenv("SMTP_HOST"), "SMTP host")
	flag.IntVar(&app.Config.SMTPPort, "smtp-port", smtpPort, "SMTP port")
	flag.StringVar(&app.Config.SMTPUsername, "smtp-username", os.Getenv("SMTP_USERNAME"), "SMTP username")
	flag.StringVar(&app.Config.SMTPPassword, "smtp-password", os.Getenv("SMTP_PASSWORD"), "SMTP password")
	flag.StringVar(&app.Config.SMTPFrom, "smtp-from", os.Getenv("SMTP_FROM"), "SMTP from email")
	flag.Parse()

	conn, err := openDB(app.Config.DSN)
	if err != nil {
		log.Fatal(err)
	}

	app.DB = &repository.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	app.Auth = auth.New(app.Config)

	emailService := service.NewEmailService(app.Config)
	authService := service.NewAuthService(app.DB, app.Auth, emailService)
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
		emailService,
		app.Auth,
	)

	log.Println("Starting application on", port)

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("/home/ikramzaidann/alumnihub/public"))))

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), routes(app.Handlers, app.Auth))
	if err != nil {
		log.Fatal(err)
	}
}
