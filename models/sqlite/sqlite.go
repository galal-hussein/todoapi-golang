package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/Sirupsen/logrus"
	_ "github.com/mxk/go-sqlite/sqlite3"
	"github.com/todoapi/models"
	"github.com/todoapi/todo"
)

const (
	// ModelName const
	ModelName = "sqlite"
	timeForm  = "0001-01-01 00:00:00 +0000 UTC"
)

// SQLite struct
type SQLite struct {
	db *sql.DB
}

func init() {
	models.RegisterModel(ModelName, new(SQLite))
}

// Init function
func (s *SQLite) Init() error {
	// Parse env variables
	sqlDB := os.Getenv("SQL_DB")
	if len(sqlDB) == 0 {
		return fmt.Errorf("SQL_DB is not set")
	}

	var err error
	s.db, err = sql.Open("sqlite3", sqlDB)
	if err != nil {
		logrus.Fatal("Error: ", err)
	}
	return nil
}

// CreateTodo function
func (s *SQLite) CreateTodo(todo todo.Todo) {
	todo.ID = bson.NewObjectId()
	timeString := todo.Due.String()
	createItemSQL := "INSERT INTO tasks (ID, Name, Completed, Due) values('" +
		todo.ID.Hex() + "','" + todo.Name + "','" + strconv.FormatBool(todo.Completed) +
		"','" + timeString + "')"
	fmt.Println(createItemSQL)
	_, err := s.db.Exec(createItemSQL)
	if err != nil {
		logrus.Fatal("Error: ", err)
	}
}

// UpdateTodo function
func (s *SQLite) UpdateTodo(t todo.Todo, todoID string) {
	var todoNew todo.Todo
	todoNew = s.FindTodoByID(todoID)
	if todoNew.Completed != t.Completed {
		todoNew.Completed = t.Completed
	}
	if todoNew.Name != t.Name {
		todoNew.Name = t.Name
	}
	if todoNew.Due != t.Due {
		todoNew.Due = t.Due
	}
	updateItemSQL := "UPDATE tasks set Name = '" + todoNew.Name +
		"', Completed = '" + strconv.FormatBool(todoNew.Completed) +
		"', Due = '" + todoNew.Due.String() + "' where ID='" + todoID + "'"
	_, err := s.db.Exec(updateItemSQL)
	if err != nil {
		logrus.Fatal("Error: ", err)
	}

}

// GetAllTodos function
func (s *SQLite) GetAllTodos() todo.Todos {
	var todos todo.Todos
	var todo todo.Todo
	rows, err := s.db.Query("SELECT * FROM tasks")
	defer rows.Close()
	if err != nil {
		logrus.Fatal("Error: ", err)
	}
	var todoid string
	var name string
	var completed bool
	var due string
	for rows.Next() {
		err = rows.Scan(&todoid, &name, &completed, &due)
		if err != nil {
			logrus.Fatal("Error: ", err)
		}
		todo.ID = bson.ObjectIdHex(todoid)
		todo.Completed = completed
		todo.Name = name
		todo.Due, _ = time.Parse(timeForm, due)
		todos = append(todos, todo)
	}
	return todos
}

// DeleteTodo function
func (s *SQLite) DeleteTodo(todoID string) {
	_, err := s.db.Exec("DELETE FROM tasks where ID=" + todoID)
	if err != nil {
		logrus.Fatal("Error: ", err)
	}
}

// FindTodoByID function
func (s *SQLite) FindTodoByID(todoID string) todo.Todo {
	var getTodo todo.Todo
	rows, err := s.db.Query("SELECT * FROM tasks where ID='" + todoID + "'")
	defer rows.Close()
	if err != nil {
		logrus.Fatal("Error: ", err)
	}
	var todoid string
	var name string
	var completed bool
	var due string
	found := false
	for rows.Next() {
		found = true
		err = rows.Scan(&todoid, &name, &completed, &due)
		if err != nil {
			logrus.Fatal("Error: ", err)
		}
	}
	if !found {
		return todo.Todo{}
	}
	getTodo.ID = bson.ObjectIdHex(todoid)
	getTodo.Completed = completed
	getTodo.Name = name
	getTodo.Due, _ = time.Parse(timeForm, due)
	return getTodo
}
