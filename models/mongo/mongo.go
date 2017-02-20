package mongo

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/todoapi/models"
	"github.com/todoapi/todo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// ModelName const
	ModelName = "mongo"
)

// Mongo struct
type Mongo struct {
	mongodb string
	session *mgo.Session
}

func init() {
	models.RegisterModel(ModelName, new(Mongo))
}

// Init function
func (m *Mongo) Init() error {
	// Parse env variables
	mongoHost := os.Getenv("MONGO_HOST")
	if len(mongoHost) == 0 {
		return fmt.Errorf("MONGO_HOST is not set")
	}
	mongoPort := os.Getenv("MONGO_PORT")
	if len(mongoPort) == 0 {
		return fmt.Errorf("MONGO_PORT is not set")
	}
	m.mongodb = os.Getenv("MONGO_DB")
	if len(m.mongodb) == 0 {
		return fmt.Errorf("MONGO_DB is not set")
	}
	logrus.Debugf("Initializing MongoDB model with host: %s, port: %d, database: %s",
		mongoHost, mongoPort, m.mongodb)

	var err error
	m.session, err = mgo.Dial(mongoHost + ":" + string(mongoPort))
	if err != nil {
		logrus.Fatal("Error: ", err)
	}
	m.session.SetMode(mgo.Monotonic, true)
	return nil
}

// CreateTodo function
func (m *Mongo) CreateTodo(todo todo.Todo) {
	newsession := m.session.Copy()
	defer newsession.Close()
	collection := newsession.DB(m.mongodb).C("tasks")
	todo.ID = bson.NewObjectId()
	err := collection.Insert(&todo)
	if err != nil {
		logrus.Fatal("Error: ", err)
	}
}

// UpdateTodo function
func (m *Mongo) UpdateTodo(todo todo.Todo, todoID string) {
	newsession := m.session.Copy()
	defer newsession.Close()
	collection := newsession.DB(m.mongodb).C("tasks")
	err := collection.Update(bson.M{"_id": bson.ObjectIdHex(todoID)}, todo)
	if err != nil {
		logrus.Fatal("Error: ", err)
	}
}

// GetAllTodos function
func (m *Mongo) GetAllTodos() todo.Todos {
	var todos todo.Todos
	newsession := m.session.Copy()
	defer newsession.Close()
	collection := newsession.DB(m.mongodb).C("tasks")
	err := collection.Find(nil).All(&todos)
	if err != nil {
		logrus.Fatal("Error: ", err)
	}
	return todos
}

// DeleteTodo function
func (m *Mongo) DeleteTodo(todoID string) {
	newsession := m.session.Copy()
	defer newsession.Close()
	collection := newsession.DB(m.mongodb).C("tasks")
	err := collection.Remove(bson.M{"_id": bson.ObjectIdHex(todoID)})
	if err != nil {
		logrus.Fatal("Error: ", err)
	}
}

// FindTodoByID function
func (m *Mongo) FindTodoByID(todoID string) todo.Todo {
	var getTodo todo.Todo
	newsession := m.session.Copy()
	defer newsession.Close()
	collection := newsession.DB(m.mongodb).C("tasks")
	err := collection.Find(bson.M{"_id": bson.ObjectIdHex(todoID)}).One(&getTodo)
	if err.Error() == "not found" {
		return todo.Todo{}
	}
	if err != nil {
		logrus.Fatal(err)
	}
	return getTodo
}
