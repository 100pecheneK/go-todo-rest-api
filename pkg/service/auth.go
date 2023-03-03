package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/100pecheneK/go-todo-rest-api.git/internal/models"
	"github.com/100pecheneK/go-todo-rest-api.git/pkg/repository"
	"github.com/dgrijalva/jwt-go"
)

const (
	salt       = "asdfghqwer"
	signingKey = "lkjhoiuyqw"
	tokenExp   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}
type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo}
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

func (*AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}
	return claims.UserId, nil
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

// manager
func NewJWT(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExp).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		id,
	})
	return token.SignedString([]byte(signingKey))
}
func NewRefreshToken() (string, error) {
	b := make([]byte, 32)
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	_, err := r.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

func (s *AuthService) createSession(id int) (accessToken, refreshToken string, err error) {
	accessToken, err = NewJWT(id)
	if err != nil {
		return
	}
	refreshToken, err = NewRefreshToken()
	if err != nil {
		return
	}
	err = s.repo.SetSession(id, refreshToken, time.Now().Add(tokenExp))
	if err != nil {
		return
	}
	return
}
