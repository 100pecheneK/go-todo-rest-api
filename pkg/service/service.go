package service

import "github.com/100pecheneK/go-todo-rest-api.git/pkg/repository"

type Authorization interface {
}
type TodoList interface {
}
type TodoItem interface {
}
type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
