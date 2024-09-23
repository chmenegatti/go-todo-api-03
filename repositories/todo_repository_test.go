package repositories

import (
	"go-todo-api-03/database"
	"go-todo-api-03/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTodo_Success(t *testing.T) {
	database.Connect() // Conectar ao banco de dados SQLite
	repo := NewTodoRepository(database.DB)

	todo := models.Todo{Title: "Test Todo"}
	_, err := repo.Create(todo)

	assert.NoError(t, err)
}

func TestFindAllTodos_Success(t *testing.T) {
	database.Connect()
	repo := NewTodoRepository(database.DB)

	todos, err := repo.FindAll()

	assert.NoError(t, err)
	assert.NotNil(t, todos)
}
