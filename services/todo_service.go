package services

import (
	"go-todo-api-03/models"
	"go-todo-api-03/repositories"
)

type TodoService interface {
	GetAllTodos() ([]models.Todo, error)
	GetTodoByID(id uint) (models.Todo, error)
	CreateTodo(todo models.Todo) (models.Todo, error)
	UpdateTodo(todo models.Todo) (models.Todo, error)
	DeleteTodo(id uint) error
}

type todoService struct {
	repo repositories.TodoRepository
}

func NewTodoService(repo repositories.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) GetAllTodos() ([]models.Todo, error) {
	return s.repo.FindAll()
}

func (s *todoService) GetTodoByID(id uint) (models.Todo, error) {
	return s.repo.FindByID(id)
}

func (s *todoService) CreateTodo(todo models.Todo) (models.Todo, error) {
	return s.repo.Create(todo)
}

func (s *todoService) UpdateTodo(todo models.Todo) (models.Todo, error) {
	return s.repo.Update(todo)
}

func (s *todoService) DeleteTodo(id uint) error {
	todo, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(todo)
}
