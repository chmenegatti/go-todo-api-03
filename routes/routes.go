package routes

import (
	"go-todo-api-03/controllers"
	"go-todo-api-03/database"
	"go-todo-api-03/repositories"
	"go-todo-api-03/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	// Conectando ao banco de dados
	database.Connect()

	// Inicializando o repositório passadno a conexão com o banco de dados como dependência
	todoRepo := repositories.NewTodoRepository(database.DB)

	// Inicializando o serviço passando o repositório como dependência
	todoService := services.NewTodoService(todoRepo)

	// Inicializando o controller passando o serviço como dependência
	todoController := controllers.NewTodoController(todoService)

	// Criando o router do Gin
	router := gin.Default()

	// Definindo as rotas
	router.GET("/todos", todoController.GetAllTodos)
	router.POST("/todos", todoController.CreateTodo)
	router.GET("/todos/:id", todoController.GetTodoByID)
	router.PUT("/todos/:id", todoController.UpdateTodo)
	router.DELETE("/todos/:id", todoController.DeleteTodo)

	return router
}
