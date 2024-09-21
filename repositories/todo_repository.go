package repositories

import (
	"go-todo-api-03/models"

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
