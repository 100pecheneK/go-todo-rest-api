package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/100pecheneK/go-todo-rest-api.git/internal/models"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db}
}

func (r *AuthPostgres) CreateUser(user models.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", userTable)

	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)

	if err := row.Scan(&id); err != nil {
		if strings.HasPrefix(err.Error(), "pq: duplicate key value violates unique constraint \"users_username_key\"") {
			return 0, errors.New("User with this username already exist")
		}
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", userTable)
	err := r.db.Get(&user, query, username, password)
	return user, err
}
func (r *AuthPostgres) SetSession(id int, refreshToken string, expiresAt time.Time) error {
	query := fmt.Sprintf("UPDATE %s SET refreshToken=$1, expiresAt=$2 WHERE id=$3", userTable)
	_, err := r.db.Exec(query, refreshToken, expiresAt.String(), id)
	return err
}
func (r *AuthPostgres) GetByRefreshToken(refreshToken string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE refreshToken=$1", userTable)
	err := r.db.Get(&user, query, refreshToken)
	return user, err
}
