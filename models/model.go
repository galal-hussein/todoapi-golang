package models

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/todoapi/todo"
)

// Model interface
type Model interface {
	Init() error
	CreateTodo(todo todo.Todo)
	UpdateTodo(todo todo.Todo, todoID string)
	DeleteTodo(todoID string)
	GetAllTodos() todo.Todos
	FindTodoByID(todoID string) todo.Todo
}

var (
	models = make(map[string]Model)
)

// GetDBModel function .. get
func GetDBModel(name string) (Model, error) {
	if model, ok := models[name]; ok {
		if err := model.Init(); err != nil {
			return nil, err
		}
		return model, nil
	}
	return nil, fmt.Errorf("No such db model: %s", name)
}

// RegisterModel function .. set
func RegisterModel(name string, model Model) {
	if _, exists := models[name]; exists {
		logrus.Fatalf("Provider '%s' tried to register twice", name)
	}
	models[name] = model
}
