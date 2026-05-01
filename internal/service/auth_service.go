package service

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"alumnihub/internal/auth"
	"alumnihub/internal/models"
	"alumnihub/internal/repository"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// TokenExpiryMinutes defines how long a password reset token is valid
const TokenExpiryMinutes = 30

type AuthService struct {
	Repo         repository.DatabaseRepo
	Auth         *auth.Auth
	EmailService *EmailService
}

func NewAuthService(repo repository.DatabaseRepo, authSvc *auth.Auth, emailSvc *EmailService) *AuthService {
	return &AuthService{
		Repo:         repo,
		Auth:         authSvc,
		EmailService: emailSvc,
	}
}

type RegisterPayload struct {
	Username string
	Email    string
	Password string
	AlumniID int
}

func (s *AuthService) RegisterCheck(nisn string) (*models.Alumni, error) {
	alumni, err := s.Repo.GetAlumniByNISN(nisn)
	if err != nil {
		return nil, errors.New("nisn doesn't match any record")
	}

	_, err = s.Repo.GetProfileByAlumniID(alumni.ID)

	if err == nil {
		return nil, errors.New("account already registered")
	}

	if err != sql.ErrNoRows {
		return nil, errors.New("error checking profile")
	}

	return alumni, nil
}

func (s *AuthService) Register(payload RegisterPayload) error {
	if payload.Username == "" || payload.Email == "" || payload.Password == "" || payload.AlumniID == 0 {
		return errors.New("all fields are required")
	}

	_, err := s.Repo.GetProfileByAlumniID(payload.AlumniID)
	if err == nil {
		return errors.New("account already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Username:  payload.Username,
		Email:     payload.Email,
		Password:  string(hashedPassword),
		IsAdmin:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	userID, err := s.Repo.InsertUser(user)
	if err != nil {
		return err
	}

	profile := models.Profile{
		AlumniID: payload.AlumniID,
		UserID:   userID,
	}

	_, err = s.Repo.InsertProfile(profile)
	if err != nil {
		_ = s.Repo.DeleteUser(userID)
		return err
	}

	return nil
}

func (s *AuthService) Authenticate(email, password string) (auth.TokenPairs, error) {
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return auth.TokenPairs{}, errors.New("invalid email")
	}

	valid, err := user.PasswordMatches(password)
	if err != nil {
		return auth.TokenPairs{}, err
	}
	if !valid {
		return auth.TokenPairs{}, errors.New("invalid password")
	}

	tokens, err := s.Auth.GenerateTokenPair(&auth.JwtUser{ID: user.ID, Username: user.Username, Role: user.IsAdmin})
	if err != nil {
		return auth.TokenPairs{}, err
	}

	return tokens, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (auth.TokenPairs, error) {
	claims := &auth.Claims{}
	_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.Auth.Secret), nil
	})
	if err != nil {
		return auth.TokenPairs{}, errors.New("unauthorized")
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return auth.TokenPairs{}, errors.New("unknown user")
	}

	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return auth.TokenPairs{}, errors.New("unknown user")
	}

	return s.Auth.GenerateTokenPair(&auth.JwtUser{ID: user.ID, Username: user.Username, Role: user.IsAdmin})
}

// generateResetToken creates a cryptographically secure random token
func (s *AuthService) generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ForgotPassword handles the forgot password request
// Returns no error for security reasons (generic response)
func (s *AuthService) ForgotPassword(email string) error {
	// Validate email format
	if email == "" {
		return errors.New("email is required")
	}

	// Check if user exists
	_, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			// User doesn't exist - return success anyway for security
			// Do not reveal whether email exists
			return nil
		}
		// Database error - still return success for security
		return nil
	}

	// Generate secure reset token
	token, err := s.generateResetToken()
	if err != nil {
		return err
	}

	// Delete any existing reset tokens for this email (single token policy)
	_ = s.Repo.DeletePasswordResetsByEmail(email)

	// Save token to database with expiration
	pr := models.PasswordReset{
		Email:     email,
		Token:     token,
		ExpiresAt: time.Now().Add(TokenExpiryMinutes * time.Minute),
		CreatedAt: time.Now(),
	}

	if err := s.Repo.InsertPasswordReset(pr); err != nil {
		return err
	}

	// Send password reset email
	if err := s.EmailService.SendPasswordResetEmail(email, token); err != nil {
		fmt.Printf("Failed to send password reset email: %v\n", err)
		return err
	}

	return nil
}

// ResetPassword handles the password reset request
func (s *AuthService) ResetPassword(token, newPassword string) error {
	// Validate input
	if token == "" {
		return errors.New("token is required")
	}
	if newPassword == "" {
		return errors.New("new password is required")
	}
	if len(newPassword) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	// Get token from database
	pr, err := s.Repo.GetPasswordResetByToken(token)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("invalid or expired token")
		}
		return errors.New("error validating token")
	}

	// Check if token has expired
	if time.Now().After(pr.ExpiresAt) {
		// Clean up expired token
		_ = s.Repo.DeletePasswordResetByToken(token)
		return errors.New("token has expired")
	}

	// Get user by email
	user, err := s.Repo.GetUserByEmail(pr.Email)
	if err != nil {
		return errors.New("user not found")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update user password
	if err := s.Repo.UpdateUserPassword(user.ID, string(hashedPassword)); err != nil {
		return err
	}

	// Delete the used token (single-use)
	if err := s.Repo.DeletePasswordResetByToken(token); err != nil {
		// Log error but don't fail - password was already changed
		fmt.Printf("Error deleting used token: %v\n", err)
	}

	// Send confirmation email
	if err := s.EmailService.SendPasswordResetSuccessEmail(pr.Email); err != nil {
		// Log the error but don't fail the password reset
		fmt.Printf("Failed to send password reset confirmation email: %v\n", err)
	}

	return nil
}
