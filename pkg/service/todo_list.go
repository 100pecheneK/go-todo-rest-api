package service

import (
	"github.com/100pecheneK/go-todo-rest-api.git/internal/models"
	"github.com/100pecheneK/go-todo-rest-api.git/pkg/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo}
}

func (s *TodoListService) Create(userId int, list models.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}
