package models

import "gorm.io/gorm"

// Todo representa a estrutura da tabela no banco de dados
type Todo struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}
