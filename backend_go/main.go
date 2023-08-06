package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
)

type Todo struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
	Id        int    `json:"id"`
}

var idCounter int32 = 0
var todo = Todo{Id: 0, Completed: false, Name: "Do dishes"}
var todos = []Todo{todo}

func main() {
	http.HandleFunc("/todos", readTodos)
	http.HandleFunc("/todo", createTodo)
	http.ListenAndServe(":8080", nil)
}

func readTodos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	json, _ := json.Marshal(todos)
	w.Write([]byte(json))
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
		return
	}

	var todo Todo
	err = json.Unmarshal(body, &todo)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	// assign an id to the todo
	todo.Id = int(atomic.AddInt32(&idCounter, 1))
	todos = append(todos, todo)	

	fmt.Printf("Received a new todo: %+v\n", todo)

	// respond with a 201 Created and the created Todo in the body
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}