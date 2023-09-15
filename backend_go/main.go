package main

import (
	"cryt1c/nyata/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Env struct {
	todos *models.TodosModel
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

	http.HandleFunc("/todos", env.readTodos)
	http.HandleFunc("/todo", env.createTodo)
	http.ListenAndServe(":8080", nil)
}

func (env *Env) readTodos(w http.ResponseWriter, r *http.Request) {
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

func (env *Env) createTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	id, err := env.todos.InsertTodo(todo)
	if err != nil {
		http.Error(w, "Error inserting todo", http.StatusBadRequest)
		return
	}

	log.Println("Inserted ID:", id)
	todo.Id = int64(id)

	// respond with a 201 Created and the created Todo in the body
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
