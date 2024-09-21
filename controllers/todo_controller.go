package controllers

import (
	"go-todo-api-03/models"
	"go-todo-api-03/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TodoController struct {
	service services.TodoService
}

func NewTodoController(service services.TodoService) *TodoController {
	return &TodoController{service: service}
}

func (ctrl *TodoController) GetAllTodos(c *gin.Context) {
	todos, err := ctrl.service.GetAllTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar todos os todos"})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (ctrl *TodoController) GetTodoByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	todo, err := ctrl.service.GetTodoByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo não encontrado"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (ctrl *TodoController) CreateTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdTodo, err := ctrl.service.CreateTodo(todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar o todo"})
		return
	}
	c.JSON(http.StatusCreated, createdTodo)
}

func (ctrl *TodoController) UpdateTodo(c *gin.Context) {
	var todo models.Todo
	id, _ := strconv.Atoi(c.Param("id"))
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo.ID = uint(id)
	updatedTodo, err := ctrl.service.UpdateTodo(todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar o todo"})
		return
	}
	c.JSON(http.StatusOK, updatedTodo)
}

func (ctrl *TodoController) DeleteTodo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctrl.service.DeleteTodo(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Erro ao deletar o todo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo deletado com sucesso"})
}
