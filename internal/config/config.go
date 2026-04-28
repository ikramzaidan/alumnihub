package config

type AppConfig struct {
	Port         int
	DSN          string
	Domain       string
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
}
