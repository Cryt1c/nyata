package main

import (
	"cryt1c/nyata/models"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

	r := mux.NewRouter()
	r.HandleFunc("/todos", env.todosHandler).Methods(http.MethodGet)
	r.HandleFunc("/todo", env.todoHandler).Methods(http.MethodPost, http.MethodPut, http.MethodOptions)
	r.HandleFunc("/reorder", env.reorderHandler).Methods(http.MethodPut, http.MethodOptions)

	r.Use(mux.CORSMethodMiddleware(r))
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func (env *Env) reorderHandler(w http.ResponseWriter, r *http.Request) {
	setCorsHeaders(&w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
		return
	}

	var reorder models.Reorder
	err = json.Unmarshal(body, &reorder)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	err = env.todos.ReorderTodos(reorder)
	if err != nil {
		http.Error(w, "Error reordering todos", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (env *Env) todosHandler(w http.ResponseWriter, r *http.Request) {
	setCorsHeaders(&w)

	todos, err := env.todos.GetTodos()
	if err != nil {
		log.Println("Error getting todos")
		log.Println(err)
		return
	}

	json, err := json.Marshal(todos)
	w.Write([]byte(json))
}

func (env *Env) todoHandler(w http.ResponseWriter, r *http.Request) {
	setCorsHeaders(&w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

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
	insertedTodo, err := env.todos.InsertTodo(todo)
	if err != nil {
		http.Error(w, "Error inserting todo", http.StatusBadRequest)
		return todo, err
	}

	log.Println("Inserted ID:", insertedTodo)
	w.WriteHeader(http.StatusCreated)
	return insertedTodo, nil
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

func setCorsHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
