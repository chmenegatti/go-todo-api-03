Para finalizar nossa série de artigos sobre Go, vamos implementar testes unitários e de integração seguindo os princípios de TDD (Test-Driven Development). Com isso, garantiremos que todos os métodos e chamadas da nossa REST API para o CRUD de Todo's estejam devidamente testados, tanto de forma isolada (unitários) quanto em conjunto (integração).

### O que é TDD?

TDD é uma abordagem de desenvolvimento de software onde os testes são escritos antes da implementação. O ciclo típico de TDD consiste em três etapas principais:

1. **Escrever um teste** que falha, pois a funcionalidade ainda não foi implementada.
2. **Implementar o código** que faz o teste passar.
3. **Refatorar** o código para melhorar sua qualidade e design, sem alterar a funcionalidade.

### Ferramentas para Testes

- **`testing`**: Pacote nativo do Go para escrever testes unitários.
- **`testify`**: Uma biblioteca popular que fornece assertivas mais robustas e mock de dependências.
- **`httptest`**: Pacote nativo para testar rotas HTTP.

### Estrutura do Projeto

Vamos garantir que cada camada (controladores, serviços, repositórios) tenha seus próprios testes, com o seguinte foco:

- **Testes unitários**: Testar as funções de forma isolada, usando mocks para simular dependências.
- **Testes de integração**: Testar a integração entre diferentes camadas e a API como um todo.

A estrutura do projeto será similar à seguinte:

```
/controllers
    todo_controller.go
    todo_controller_test.go
/database
    database.go
/models
    todo.go
/repositories
    todo_repository.go
    todo_repository_test.go
/routes
    routes.go
/services
    todo_service.go
    todo_service_test.go
/tests
    integration_test.go
/main.go
```

### Testando a Camada de Serviços

Primeiro, vamos criar um exemplo de teste unitário para o serviço de `Todo`.

**services/todo_service_test.go**:

```go
package services

import (
	"errors"
	"go-todo-api/models"
	"go-todo-api/repositories"
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

func (m *MockTodoRepository) Create(todo models.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockTodoRepository) Update(todo models.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockTodoRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetAllTodos_Success(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	service := NewTodoService(mockRepo)

	// Simulando retorno do mock
	mockRepo.On("FindAll").Return([]models.Todo{
		{ID: 1, Title: "Task 1"},
		{ID: 2, Title: "Task 2"},
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
```

Esse exemplo de teste unitário testa a camada de serviço com um mock do repositório usando o `testify`. As funções `GetAllTodos` e `GetTodoByID` são testadas para garantir que o serviço retorne corretamente os dados.

### Testando a Camada de Repositório

Agora, vamos criar testes para o repositório, onde faremos interações reais com o banco de dados.

**repositories/todo_repository_test.go**:

```go
package repositories

import (
	"go-todo-api/database"
	"go-todo-api/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTodo_Success(t *testing.T) {
	database.Connect() // Conectar ao banco de dados SQLite
	repo := NewTodoRepository(database.DB)

	todo := models.Todo{Title: "Test Todo"}
	err := repo.Create(todo)

	assert.NoError(t, err)
}

func TestFindAllTodos_Success(t *testing.T) {
	database.Connect()
	repo := NewTodoRepository(database.DB)

	todos, err := repo.FindAll()

	assert.NoError(t, err)
	assert.NotNil(t, todos)
}
```

Aqui, testamos a camada de repositório interagindo diretamente com o banco de dados SQLite. O objetivo é garantir que os dados possam ser inseridos e recuperados corretamente.

### Testando a Camada de Controladores

Agora, vamos testar as rotas HTTP com o pacote `httptest`.

**controllers/todo_controller_test.go**:

```go
package controllers

import (
	"go-todo-api/models"
	"go-todo-api/services"
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

func TestGetAllTodosHandler_Success(t *testing.T) {
	mockService := new(MockTodoService)
	controller := NewTodoController(mockService)

	// Simulando retorno do mock
	mockService.On("GetAllTodos").Return([]models.Todo{
		{ID: 1, Title: "Task 1"},
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
```

Aqui usamos o `httptest` para simular chamadas HTTP e verificar se a rota está retornando o que esperamos.

### Teste de Integração

Agora, vamos criar um teste de integração para verificar o fluxo completo da aplicação, desde a API até a persistência de dados.

**tests/integration_test.go**:

```go
package tests

import (
	"go-todo-api/database"
	"go-todo-api/models"
	"go-todo-api/routes"
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
```

Nesse teste de integração, testamos a criação de um Todo, simulando uma requisição completa à API.

### Executando os Testes

Para executar os testes, basta rodar o comando:

```bash
$ go test -v go-todo-api/...
```	
A saída do comando mostrará os resultados dos testes, indicando se todos passaram ou se houve algum erro.

**Exemplo de saída:**
```bash
?   	go-todo-api-03	[no test files]
?   	go-todo-api-03/database	[no test files]
?   	go-todo-api-03/models	[no test files]
=== RUN   TestGetAllTodosHandler_Success
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /todos                    --> go-todo-api-03/controllers.(*TodoController).GetAllTodos-fm (3 handlers)
[GIN] 2024/09/23 - 08:11:17 | 200 |     147.947µs |                 | GET      "/todos"
--- PASS: TestGetAllTodosHandler_Success (0.00s)
PASS
ok  	go-todo-api-03/controllers	0.007s
?   	go-todo-api-03/routes	[no test files]
=== RUN   TestCreateTodo_Success
--- PASS: TestCreateTodo_Success (0.01s)
=== RUN   TestFindAllTodos_Success
--- PASS: TestFindAllTodos_Success (0.00s)
PASS
ok  	go-todo-api-03/repositories	0.012s
=== RUN   TestGetAllTodos_Success
--- PASS: TestGetAllTodos_Success (0.00s)
=== RUN   TestGetTodoByID_NotFound
--- PASS: TestGetTodoByID_NotFound (0.00s)
PASS
ok  	go-todo-api-03/services	(cached)
=== RUN   TestCreateTodoIntegration
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /todos                    --> go-todo-api-03/controllers.(*TodoController).GetAllTodos-fm (3 handlers)
[GIN-debug] POST   /todos                    --> go-todo-api-03/controllers.(*TodoController).CreateTodo-fm (3 handlers)
[GIN-debug] GET    /todos/:id                --> go-todo-api-03/controllers.(*TodoController).GetTodoByID-fm (3 handlers)
[GIN-debug] PUT    /todos/:id                --> go-todo-api-03/controllers.(*TodoController).UpdateTodo-fm (3 handlers)
[GIN-debug] DELETE /todos/:id                --> go-todo-api-03/controllers.(*TodoController).DeleteTodo-fm (3 handlers)
[GIN] 2024/09/23 - 08:11:17 | 201 |     1.72825ms |                 | POST     "/todos"
--- PASS: TestCreateTodoIntegration (0.00s)
PASS
ok  	go-todo-api-03/tests	0.008s
```

Você também pode rodar testes de uma camada específica, como os testes de serviço:

```bash
$ go test -v go-todo-api/services
```

Isso executará apenas os testes da camada de serviço.

**Exemplo de saída:**
```bash
=== RUN   TestGetAllTodos_Success
--- PASS: TestGetAllTodos_Success (0.00s)
PASS
ok      go-todo-api-03/services 0.004s
```

### Conclusão

Com esses exemplos, você tem uma base sólida de como implementar testes unitários e de integração para uma aplicação REST em Go, utilizando os conceitos de TDD. Cada camada da aplicação foi testada de maneira isolada e integrada, garantindo que as funcionalidades estejam bem validadas. Isso ajuda a aumentar a confiabilidade do código e a qualidade do software.

Agora sua API está pronta para ser testada e mantida de maneira mais eficiente!