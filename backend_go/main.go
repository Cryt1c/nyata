package main

import (
	"cryt1c/nyata/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

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
	r.HandleFunc("/todo/{todoId}", env.todoIdHandler).Methods(http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/reorder", env.reorderHandler).Methods(http.MethodPut, http.MethodOptions)
	r.HandleFunc("/reset", env.resetHandler).Methods(http.MethodPut, http.MethodOptions)

	r.Use(mux.CORSMethodMiddleware(r))
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func (env *Env) reorderHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)
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

	todos, err := env.todos.ReorderTodos(reorder)
	if err != nil {
		log.Println("Error reordering todos %v", err)
		http.Error(w, "Error reordering todos", http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(todos)
	w.Write([]byte(json))
}

func (env *Env) resetHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	affectedRows, err := env.todos.ResetListOrder(2)
	if err != nil {
		http.Error(w, "Error reordering todos", http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(affectedRows)
	w.Write([]byte(json))
}

func (env *Env) todosHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)

	todos, err := env.todos.GetTodos(models.GetOptions{Sorted: true})
	if err != nil {
		log.Println("Error getting todos")
		log.Println(err)
		return
	}
	json, err := json.Marshal(todos)
	w.Write([]byte(json))
}

func (env *Env) todoHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)
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
		returnedTodo, err = env.createTodoHandler(w, todo)
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

func (env *Env) todoIdHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var affectedRows int64
	switch r.Method {
	case http.MethodDelete:
		vars := mux.Vars(r)
		idString := vars["todoId"]
		todoId, _ := strconv.ParseInt(idString, 10, 64)
		var err error
		affectedRows, err = env.deleteTodoHandler(w, todoId)
		if err != nil {
			return
		}
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// respond with a 201 Created and the created Todo in the body
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(affectedRows)
}

func (env *Env) createTodoHandler(w http.ResponseWriter, todo models.Todo) (models.Todo, error) {
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

func (env *Env) deleteTodoHandler(w http.ResponseWriter, todoId int64) (int64, error) {
	rowsAffected, err := env.todos.DeleteTodoById(todoId)
	if err != nil {
		http.Error(w, "Error deleting todo", http.StatusBadRequest)
		return rowsAffected, err
	}

	log.Println("Deleted todo:", rowsAffected)
	w.WriteHeader(http.StatusOK)
	return rowsAffected, nil
}

func setHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Content-Type", "application/json")
}
