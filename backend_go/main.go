package main

import (
	"cryt1c/nyata/models"
	"cryt1c/nyata/todosdb"
	"encoding/json"
	"io"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var idCounter int32 = 0

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
	enableCors(&w)
	todosDB, err := todosdb.OpenDB()
	if err != nil {
		log.Println("Error opening database")
		log.Println(err)
		return
	}
	defer todosDB.Close()

	todos, err := todosDB.GetTodos()
	if err != nil {
		log.Println("Error getting todos")
		log.Println(err)
		return
	}

	json, err := json.Marshal(todos)
	w.Write([]byte(json))
}

func createTodo(w http.ResponseWriter, r *http.Request) {
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

	todosDB, err := todosdb.OpenDB()
	if err != nil {
		log.Println("Error opening database")
		log.Println(err)
		return
	}
	defer todosDB.Close()
	result, err := todosDB.Exec(
		"INSERT INTO todos(Name, Completed) VALUES( ?, ?)",
		todo.Name,
		todo.Completed,
	)
	if err != nil {
		log.Println("Error inserting todo")
		log.Println(err)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Error retrieving last insert id")
		log.Println(err)
		return
	}
	log.Println("Inserted ID:", id)
	todo.Id = int(id)

	// respond with a 201 Created and the created Todo in the body
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}


