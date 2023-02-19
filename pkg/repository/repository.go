package repository

import (
	"github.com/100pecheneK/go-todo-rest-api.git/internal/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}
type TodoList interface {
	Create(userId int, list models.TodoList) (int, error)
	GetAll(userId int) ([]models.TodoList, error)
	GetById(userId, id int) (models.TodoList, error)
	Delete(userID, id int) error
	Update(userId, id int, input models.UpdateListInput) error
}
type TodoItem interface {
	Create(id int, input models.TodoItem) (int, error)
	GetAll(userId, listId int) ([]models.TodoItem, error)
	GetById(userId, id int) (models.TodoItem, error)
	Delete(userId, id int) error
}
type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListRepository(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
