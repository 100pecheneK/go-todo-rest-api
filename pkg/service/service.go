package service

import (
	"github.com/100pecheneK/go-todo-rest-api.git/internal/models"
	"github.com/100pecheneK/go-todo-rest-api.git/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, string, error)
	ParseToken(token string) (int, error)
	RefreshTokens(refreshToken string) (string, string, error)
}
type TodoList interface {
	Create(userId int, list models.TodoList) (int, error)
	GetAll(userId int) ([]models.TodoList, error)
	GetById(userId, id int) (models.TodoList, error)
	Delete(userId, id int) error
	Update(userId, id int, input models.UpdateListInput) error
}
type TodoItem interface {
	Create(userId, listId int, input models.TodoItem) (int, error)
	GetAll(userId, listId int) ([]models.TodoItem, error)
	GetById(userId, id int) (models.TodoItem, error)
	Delete(userId, id int) error
	Update(userId, id int, input models.UpdateItemInput) error
}
type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
