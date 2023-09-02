package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
	Id        int    `json:"id"`
}

type todosDB struct {
	db *sql.DB
}

var idCounter int32 = 0

var todos = [][]Todo{
	{
		{Id: 0, Completed: false, Name: "Do dishes"},
		{Id: 1, Completed: false, Name: "Do laundry"},
		{Id: 2, Completed: false, Name: "Do homework"},
	},
	{
		{Id: 3, Completed: false, Name: "Do homework"},
		{Id: 4, Completed: false, Name: "Do homework"},
		{Id: 5, Completed: false, Name: "Do homework"},
	},
}

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
	todosDB, err := openDB()
	if err != nil {
		log.Println("Error opening database")
		log.Println(err)
		return
	}
	defer todosDB.db.Close()

	todos, err := todosDB.getTodos()
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

	var todo Todo
	err = json.Unmarshal(body, &todo)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	// assign an id to the todo
	// todo.Id = int(atomic.AddInt32(&idCounter, 1))
	// TODO: Take list and position for todo.
	// todos = append(todos, []Todo{todo})

	// We don't care about the returned values, so we're using Exec. If we
	// wanted to reuse these statements, it would be more efficient to use
	// prepared statements. Learn more:
	// https://go.dev/doc/database/prepared-statements

	todosDB, err := openDB()
	if err != nil {
		log.Println("Error opening database")
		log.Println(err)
		return
	}
	defer todosDB.db.Close()
	result, err := todosDB.db.Exec(
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

func (t *todosDB) getTodos() ([]Todo, error) {
	var todos []Todo
	rows, err := t.db.Query("SELECT * FROM todos")
	if err != nil {
		return todos, fmt.Errorf("unable to get values: %w", err)
	}
	for rows.Next() {
		var todo Todo
		err = rows.Scan(
			&todo.Id,
			&todo.Name,
			&todo.Completed,
		)
		if err != nil {
			return todos, err
		}
		todos = append(todos, todo)
	}
	return todos, err
}

// openDB opens a SQLite database and stores that database in our special spot.
func openDB() (*todosDB, error) {
	db, err := sql.Open("sqlite3", "./todos.db")
	if err != nil {
		log.Println("Error opening database")
		log.Println(err)
		return nil, err
	}
	log.Println("Opened database")
	t := todosDB{db}
	log.Println("Checking if table exists")
	if !t.tableExists("todos") {
		err := t.createTable()
		if err != nil {
			return nil, err
		}
	}
	return &t, nil
}

func (t *todosDB) tableExists(name string) bool {
	if _, err := t.db.Query("SELECT * FROM todos"); err == nil {
		return true
	}
	return false
}

func (t *todosDB) createTable() error {
	_, err := t.db.Exec(`CREATE TABLE "todos" ( "id" INTEGER, "name" TEXT NOT NULL, "project" TEXT, "status" TEXT, "created" DATETIME, PRIMARY KEY("id" AUTOINCREMENT))`)
	return err
}

func initTodosDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(path, 0o770)
		}
		return err
	}
	return nil
}
