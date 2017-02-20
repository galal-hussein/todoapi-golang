package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/todoapi/todo"
)

// Index route
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world, %q", html.EscapeString(r.URL.Path))
}

// TodoIndex Route
func TodoIndex(w http.ResponseWriter, r *http.Request) {
	todos := model.GetAllTodos()
	w.Header().Set("Conetent-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		panic(err)
	}
}

// TodoShow Route
func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoID := vars["todoID"]
	showTodo := model.FindTodoByID(todoID)
	if showTodo != (todo.Todo{}) {
		w.Header().Set("Conetent-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(showTodo); err != nil {
			panic(err)
		}
	} else {
		w.Header().Set("Conetent-Type", "application/json; charset=UTF-8")
		w.WriteHeader(404)
	}
}

// TodoCreate Route
func TodoCreate(w http.ResponseWriter, r *http.Request) {
	var todo todo.Todo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &todo); err != nil {
		w.Header().Set("Conetent-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
	}
	fmt.Println(todo)
	model.CreateTodo(todo)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		panic(err)
	}

}

// TodoUpdate route
func TodoUpdate(w http.ResponseWriter, r *http.Request) {
	var newTodo todo.Todo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	vars := mux.Vars(r)
	todoId := vars["todoId"]

	if err := json.Unmarshal(body, &newTodo); err != nil {
		w.Header().Set("Conetent-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
	}

	model.UpdateTodo(newTodo, todoId)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newTodo); err != nil {
		panic(err)
	}
}

// TodoDelete route
func TodoDelete(w http.ResponseWriter, r *http.Request) {
	var newTodo todo.Todo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	vars := mux.Vars(r)
	todoId := vars["todoId"]

	if err := json.Unmarshal(body, &newTodo); err != nil {
		w.Header().Set("Conetent-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
	}

	model.DeleteTodo(todoId)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
}
