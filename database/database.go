package database

import (
	"go-todo-api-03/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Conecta ao banco de dados e faz as migrações
func Connect() {
	// Conecta ao banco de dados SQLite
	database, err := gorm.Open(sqlite.Open("todos.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Falha ao conectar ao banco de dados:", err)
	}

	// Auto-migração para criar a tabela no banco
	database.AutoMigrate(&models.Todo{})

	DB = database
}
