package main

import (
	"cryt1c/nyata/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Env struct {
	todos *models.TodosDB
}

func main() {
	todosDB, err := models.OpenDB("./todos.db")
	if err != nil {
		log.Println("Error opening database")
		log.Println(err)
		return
	}

	env := &Env{todos: todosDB}
	defer env.todos.DB.Close()

	http.HandleFunc("/todos", env.todoHandler)
	http.HandleFunc("/todo", env.todosHandler)
	http.ListenAndServe(":8080", nil)
}

func (env *Env) todoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	enableCors(&w)

	todos, err := env.todos.GetTodos()
	if err != nil {
		log.Println("Error getting todos")
		log.Println(err)
		return
	}

	json, err := json.Marshal(todos)
	w.Write([]byte(json))
}

func (env *Env) todosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	enableCors(&w)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
		return
	}

	var todo models.Todo
	err = json.Unmarshal(body, &todo)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	var returnedTodo models.Todo
	switch r.Method {
	case http.MethodPost:
		returnedTodo, err = env.insertTodoHandler(w, todo)
		if err != nil {
			return
		}
	case http.MethodPut:
		returnedTodo, err = env.updateTodoHandler(w, todo)
		if err != nil {
			return
		}
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// respond with a 201 Created and the created Todo in the body
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returnedTodo)
}

func (env *Env) insertTodoHandler(w http.ResponseWriter, todo models.Todo) (models.Todo, error) {
	id, err := env.todos.InsertTodo(todo)
	if err != nil {
		http.Error(w, "Error inserting todo", http.StatusBadRequest)
		return todo, err
	}

	log.Println("Inserted ID:", id)
	w.WriteHeader(http.StatusCreated)
	todo.Id = int64(id)
	return todo, nil
}

func (env *Env) updateTodoHandler(w http.ResponseWriter, todo models.Todo) (models.Todo, error) {
	_, err := env.todos.UpdateTodo(todo)
	if err != nil {
		http.Error(w, "Error updating todo", http.StatusBadRequest)
		return todo, err
	}

	log.Println("Updated todo:", todo)
	w.WriteHeader(http.StatusOK)
	return todo, nil
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
