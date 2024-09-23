package services

import (
	"errors"
	"go-todo-api-03/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock do repositório
type MockTodoRepository struct {
	mock.Mock
}

func (m *MockTodoRepository) FindAll() ([]models.Todo, error) {
	args := m.Called()
	return args.Get(0).([]models.Todo), args.Error(1)
}

func (m *MockTodoRepository) FindByID(id uint) (models.Todo, error) {
	args := m.Called(id)
	return args.Get(0).(models.Todo), args.Error(1)
}

func (m *MockTodoRepository) Create(todo models.Todo) (models.Todo, error) {
	args := m.Called(todo)
	return args.Get(0).(models.Todo), args.Error(1)
}

func (m *MockTodoRepository) Update(todo models.Todo) (models.Todo, error) {
	args := m.Called(todo)
	return args.Get(0).(models.Todo), args.Error(1)
}

func (m *MockTodoRepository) Delete(todo models.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func TestGetAllTodos_Success(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	service := NewTodoService(mockRepo)

	// Simulando retorno do mock
	mockRepo.On("FindAll").Return([]models.Todo{
		{Title: "Task 1"},
		{Title: "Task 2"},
	}, nil)

	todos, err := service.GetAllTodos()

	assert.NoError(t, err)
	assert.Len(t, todos, 2)
	assert.Equal(t, "Task 1", todos[0].Title)
	mockRepo.AssertExpectations(t)
}

func TestGetTodoByID_NotFound(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	service := NewTodoService(mockRepo)

	// Simulando retorno do mock
	mockRepo.On("FindByID", uint(1)).Return(models.Todo{}, errors.New("not found"))

	todo, err := service.GetTodoByID(1)

	assert.Error(t, err)
	assert.Equal(t, "not found", err.Error())
	assert.Empty(t, todo)
	mockRepo.AssertExpectations(t)
}
