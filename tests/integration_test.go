package tests

import (
	"go-todo-api-03/database"
	"go-todo-api-03/routes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTodoIntegration(t *testing.T) {
	database.Connect() // Conectando ao banco de dados

	router := routes.SetupRoutes()
	w := httptest.NewRecorder()

	body := strings.NewReader(`{"title":"New Todo"}`)
	req, _ := http.NewRequest("POST", "/todos", body)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
