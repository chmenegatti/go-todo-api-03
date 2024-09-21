package main

import (
	"go-todo-api-03/routes"
)

func main() {

	// Instancia as rotas
	router := routes.SetupRoutes()

	// Sobe o servidor na porta 8080
	router.Run(":8080")
}
