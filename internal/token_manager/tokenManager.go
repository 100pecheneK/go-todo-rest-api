package token_manager

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type Manager struct {
	signingKey string
}

func NewManager(signinKey string) (*Manager, error) {
	if signinKey == "" {
		return nil, errors.New("empty signin key")
	}
	return &Manager{signinKey}, nil
}
func (m *Manager) NewJWT(id int, exp time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(exp).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		id,
	})
	return token.SignedString([]byte(m.signingKey))
}
func (Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	_, err := r.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
func (m *Manager) Parse(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(m.signingKey), nil
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
