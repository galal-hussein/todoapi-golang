package todo

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Todo Struct
type Todo struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Name      string        `json:"name"`
	Completed bool          `json:"completed"`
	Due       time.Time     `json:"due"`
}

// Todos list
type Todos []Todo
