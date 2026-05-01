package config

type AppConfig struct {
	Port         int
	DSN          string
	Domain       string
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string

	// SMTP Configuration
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
}
