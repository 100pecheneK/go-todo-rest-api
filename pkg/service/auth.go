package service

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/100pecheneK/go-todo-rest-api.git/internal/models"
	"github.com/100pecheneK/go-todo-rest-api.git/internal/token_manager"
	"github.com/100pecheneK/go-todo-rest-api.git/pkg/repository"
)

const (
	salt            = "asdfghqwer"
	signingKey      = "lkjhoiuyqw"
	tokenExp        = 12 * time.Hour
	refreshTokenExp = 30 * 24 * time.Hour
)

type TokenManager interface {
	NewJWT(id int, exp time.Duration) (string, error)
	NewRefreshToken() (string, error)
	Parse(accessToken string) (int, error)
}
type AuthService struct {
	repo         repository.Authorization
	tokenManager TokenManager
}

func NewAuthService(repo repository.Authorization) *AuthService {
	manager, _ := token_manager.NewManager(signingKey)
	return &AuthService{repo, manager}
}

func (s *AuthService) CreateUser(user models.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, string, error) {
	user, err := s.repo.GetUser(username, s.generatePasswordHash(password))
	if err != nil {
		return "", "", err
	}
	return s.createSession(user.Id)
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	return s.tokenManager.Parse(accessToken)
}

func (s *AuthService) RefreshTokens(refreshToken string) (string, string, error) {
	user, err := s.repo.GetByRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}
	return s.createSession(user.Id)
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) createSession(id int) (accessToken, refreshToken string, err error) {
	accessToken, err = s.tokenManager.NewJWT(id, tokenExp)
	if err != nil {
		return
	}
	refreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return
	}
	err = s.repo.SetSession(id, refreshToken, time.Now().Add(refreshTokenExp))
	if err != nil {
		return
	}
	return
}
