package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type TodosModel struct {
	DB *sql.DB
}

type Todo struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
	Id        int    `json:"id"`
}

func (m *TodosModel) GetTodos() ([]Todo, error) {
	var todos []Todo
	rows, err := m.DB.Query("SELECT * FROM todos")
	if err != nil {
		return todos, fmt.Errorf("unable to get values: %w", err)
	}
	defer rows.Close()

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

func (m *TodosModel) InsertTodo(todo Todo) (int64, error) {
	result, err := m.DB.Exec(
		"INSERT INTO todos(Name, Completed) VALUES(?, ?)",
		todo.Name,
		todo.Completed,
	)
	if err != nil {
		log.Println("Error inserting todo")
		log.Println(err)
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Error retrieving last insert id")
		log.Println(err)
		return 0, err
	}
	return id, nil
}

func (m *TodosModel) tableExists(name string) bool {
	if _, err := m.DB.Query("SELECT * FROM todos"); err == nil {
		return true
	}
	return false
}

func (m *TodosModel) createTable() error {
	_, err := m.DB.Exec(`CREATE TABLE "todos" ( "id" INTEGER, "name" TEXT NOT NULL, "project" TEXT, "status" TEXT, "created" DATETIME, PRIMARY KEY("id" AUTOINCREMENT))`)
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

func OpenDB() (*TodosModel, error) {
	db, err := sql.Open("sqlite3", "./todos.db")
	if err != nil {
		log.Println("Error opening database")
		log.Println(err)
		return nil, err
	}
	log.Println("Opened database")
	t := TodosModel{db}
	log.Println("Checking if table exists")
	if !t.tableExists("todos") {
		err := t.createTable()
		if err != nil {
			return nil, err
		}
	}
	return &t, nil
}
