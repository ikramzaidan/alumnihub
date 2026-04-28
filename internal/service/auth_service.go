package service

import (
	"database/sql"
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

type AuthService struct {
	Repo repository.DatabaseRepo
	Auth *auth.Auth
}

func NewAuthService(repo repository.DatabaseRepo, authSvc *auth.Auth) *AuthService {
	return &AuthService{
		Repo: repo,
		Auth: authSvc,
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
