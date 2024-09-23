package controllers

import (
	"go-todo-api-03/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTodoService struct {
	mock.Mock
}

func (m *MockTodoService) GetAllTodos() ([]models.Todo, error) {
	args := m.Called()
	return args.Get(0).([]models.Todo), args.Error(1)
}

func (m *MockTodoService) CreateTodo(todo models.Todo) (models.Todo, error) {
	args := m.Called(todo)
	return args.Get(0).(models.Todo), args.Error(1)
}

func (m *MockTodoService) DeleteTodo(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTodoService) GetTodoByID(id uint) (models.Todo, error) {
	args := m.Called(id)
	return args.Get(0).(models.Todo), args.Error(1)
}

func (m *MockTodoService) UpdateTodo(todo models.Todo) (models.Todo, error) {
	args := m.Called(todo)
	return args.Get(0).(models.Todo), args.Error(1)
}

func TestGetAllTodosHandler_Success(t *testing.T) {
	mockService := new(MockTodoService)
	controller := NewTodoController(mockService)

	// Simulando retorno do mock
	mockService.On("GetAllTodos").Return([]models.Todo{
		{Title: "Task 1"},
	}, nil)

	// Simulando a chamada HTTP
	router := gin.Default()
	router.GET("/todos", controller.GetAllTodos)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/todos", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
