package service

import (
	"github.com/100pecheneK/go-todo-rest-api.git/internal/models"
	"github.com/100pecheneK/go-todo-rest-api.git/pkg/repository"
)

type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{repo, listRepo}
}

func (s *TodoItemService) Create(userId, listId int, item models.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return 0, err
	}
	return s.repo.Create(listId, item)
}

func (s *TodoItemService) GetAll(userId, listId int) ([]models.TodoItem, error) {
	return s.repo.GetAll(userId, listId)
}

func (s *TodoItemService) GetById(userId, id int) (models.TodoItem, error) {
	return s.repo.GetById(userId, id)
}

func (s *TodoItemService) Delete(userId, id int) error {
	return s.repo.Delete(userId, id)
}

func (s *TodoItemService) Update(userId, id int, input models.UpdateItemInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, id, input)
}
