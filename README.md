### Construindo uma REST API com SOLID em Go - Parte 3: Camada de Serviço, Repositórios e Injeção de Dependência

Nesta etapa da nossa série, vamos refatorar a estrutura da API, movendo as regras de negócio para a camada de serviço, implementando repositórios para lidar com a persistência e aplicando injeção de dependência para manter nosso código modular e testável.

#### Por que Camada de Serviço e Repositórios?

Ao separar a lógica de negócio em uma camada de serviço, tornamos o código mais organizado e fácil de testar. O repositório encapsula as operações de acesso ao banco de dados, permitindo uma possível troca futura de tecnologia de persistência sem impacto na camada de negócio.

#### Estrutura do Projeto

Aqui está uma visão da nova estrutura do projeto:

```
/controllers
  todo_controller.go
/services
  todo_service.go
/repositories
  todo_repository.go
/database
  database.go
/models
  todo.go
```

#### Passo 1: Criando a Camada de Repositório

O repositório será responsável por todas as operações relacionadas à persistência de dados. Vamos criar um arquivo chamado `todo_repository.go` dentro da pasta `repositories`.

**repositories/todo_repository.go**:
```go
package repositories

import (
	"go-todo-api/models"
	"gorm.io/gorm"
)

type TodoRepository interface {
	FindAll() ([]models.Todo, error)
	FindByID(id uint) (models.Todo, error)
	Create(todo models.Todo) (models.Todo, error)
	Update(todo models.Todo) (models.Todo, error)
	Delete(todo models.Todo) error
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) FindAll() ([]models.Todo, error) {
	var todos []models.Todo
	err := r.db.Find(&todos).Error
	return todos, err
}

func (r *todoRepository) FindByID(id uint) (models.Todo, error) {
	var todo models.Todo
	err := r.db.First(&todo, id).Error
	return todo, err
}

func (r *todoRepository) Create(todo models.Todo) (models.Todo, error) {
	err := r.db.Create(&todo).Error
	return todo, err
}

func (r *todoRepository) Update(todo models.Todo) (models.Todo, error) {
	err := r.db.Save(&todo).Error
	return todo, err
}

func (r *todoRepository) Delete(todo models.Todo) error {
	return r.db.Delete(&todo).Error
}
```

#### Passo 2: Criando a Camada de Serviço

Agora vamos criar a camada de serviço, que conterá as regras de negócio. O serviço vai interagir com o repositório para persistir os dados.

**services/todo_service.go**:
```go
package services

import (
	"go-todo-api/models"
	"go-todo-api/repositories"
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
```

#### Passo 3: Refatorando o Controlador para Usar o Serviço

Agora vamos refatorar o controlador `Todo` para delegar a lógica de negócio ao serviço.

**controllers/todo_controller.go**:
```go
package controllers

import (
	"go-todo-api/models"
	"go-todo-api/services"
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
```
#### Passo 4: Atualizando o `routes.go` para Usar o Controlador e a injeção de Dependências

Com a injenção de dependências, precisamos atualizar o arquivo `routes.go` para usar o controlador e a camada de serviço.

**routes/routes.go**:
```go
package routes

import (
	"go-todo-api/controllers"
	"go-todo-api/services"
	"go-todo-api/repositories"
	"go-todo-api/database"

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

```

#### Passo 5: Atualizando o `main.go`

Por fim, vamos atualizar o arquivo `main.go` para usar a nova estrutura da aplicação.

**main.go**:
```go
package main

import (
	"go-todo-api/routes"
)

func main() {
	// Configurando e iniciando o servidor
	router := routes.SetupRoutes()
	router.Run(":8080")
}

```

#### Conclusão

Neste artigo, movemos a lógica de negócio para uma camada de serviço e implementamos o padrão de repositórios para a persistência de dados, aplicando injeção de dependência para melhorar a modularidade e testabilidade da nossa aplicação. Essa separação é essencial para seguir os princípios do SOLID, mantendo o código coeso e desacoplado.

No próximo artigo, abordaremos a implementação de testes unitários e de integração, seguindo a abordagem TDD.

Fique atento para a próxima parte da série!

--- 

Nos vemos no próximo post!

### [Testes Unitários e de Integração](Testes.md)



