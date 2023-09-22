package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type TodosDB struct {
	DB *sql.DB
}

type Todo struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Completed  bool   `json:"completed"`
	PositionId int64  `json:"positionId"`
	ListId     int64  `json:"listId"`
}

type List struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Todos []Todo `json:"todos"`
}

func (m *TodosDB) GetTodos() ([]Todo, error) {
	var todos []Todo
	rows, err := m.DB.Query("SELECT * FROM todos")
	if err != nil {
		return todos, fmt.Errorf("Unable to get values: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		err = rows.Scan(
			&todo.Id,
			&todo.Name,
			&todo.Completed,
			&todo.PositionId,
			&todo.ListId,
		)
		if err != nil {
			return todos, err
		}
		todos = append(todos, todo)
	}
	return todos, err
}

func (m *TodosDB) InsertTodo(todo Todo) (Todo, error) {
	result, err := m.DB.Exec(
		`INSERT INTO todos(name, completed, position_id, list_id) 
		 VALUES(?, ?, (SELECT COALESCE(MAX(position_id), 0) + 10 FROM todos WHERE list_id = ?), ?)
		`,
		todo.Name,
		todo.Completed,
		todo.ListId,
		todo.ListId,
	)
	if err != nil {
		log.Println("Error inserting todo")
		log.Println(err)
		return Todo{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Error retrieving last insert id")
		log.Println(err)
		return Todo{}, err
	}

	rows, err := m.DB.Query("SELECT * FROM todos WHERE id = ?", id)
	if err != nil {
		return Todo{}, fmt.Errorf("Unable to get values: %w", err)
	}
	defer rows.Close()

	rows.Next()
	var newTodo Todo
	err = rows.Scan(
		&newTodo.Id,
		&newTodo.Name,
		&newTodo.Completed,
		&newTodo.PositionId,
		&newTodo.ListId,
	)
	if err != nil {
		return Todo{}, err
	}
	return newTodo, nil
}

func (m *TodosDB) DeleteTodoById(id int64) (int64, error) {
	result, err := m.DB.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		log.Printf("Error deleting todo with ID %d: %v", id, err)
		return 0, fmt.Errorf("failed to delete todo with ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected after deleting todo with ID %d: %v", id, err)
		return 0, fmt.Errorf("Failed to get rows affected after deleting todo with ID %d: %w", id, err)
	}

	if rowsAffected == 0 {
		log.Printf("No todo found with ID %d", id)
		return 0, errors.New(fmt.Sprintf("No todo found with ID %d", id))
	}

	return rowsAffected, nil
}

func (m *TodosDB) UpdateTodo(todo Todo) (*Todo, error) {
	result, err := m.DB.Exec("UPDATE todos SET name = ?, completed = ?, position_id = ?, list_id = ? WHERE id = ?", todo.Name, todo.Completed, todo.PositionId, todo.ListId, todo.Id)
	if err != nil {
		log.Printf("Error updating todo with ID %d: %v", todo.Id, err)
		return nil, fmt.Errorf("Failed to update todo with ID %d: %w", todo.Id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected after deleting todo with ID %d: %v", todo.Id, err)
		return nil, fmt.Errorf("Failed to get rows affected after deleting todo with ID %d: %w", todo.Id, err)
	}

	if rowsAffected == 0 {
		log.Printf("No todo found with ID %d", todo.Id)
		return nil, errors.New(fmt.Sprintf("No todo found with ID %d", todo.Id))
	}

	return &todo, nil
}

func (m *TodosDB) GetLists() ([]List, error) {
	var todos []List
	rows, err := m.DB.Query("SELECT * FROM lists")
	if err != nil {
		return todos, fmt.Errorf("Unable to get values: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var todo List
		err = rows.Scan(
			&todo.Id,
			&todo.Name,
			&todo.Todos,
		)
		if err != nil {
			return todos, err
		}
		todos = append(todos, todo)
	}
	return todos, err
}
func (m *TodosDB) tableExists(name string) bool {
	if _, err := m.DB.Query("SELECT * FROM todos"); err == nil {
		return true
	}
	return false
}

func (m *TodosDB) createTodosTable() error {
	_, err := m.DB.Exec(`CREATE TABLE "todos" ( "id" INTEGER, "name" TEXT NOT NULL, "completed" INTEGER NOT NULL, "position_id" INTEGER NOT NULL, "list_id" INTEGER NOT NULL, PRIMARY KEY("id" AUTOINCREMENT))`)
	if err != nil {
		return err
	}
	return nil
}

func (m *TodosDB) createListsTable() error {
	_, err := m.DB.Exec(`CREATE TABLE "lists" ( "id" INTEGER, "name" TEXT NOT NULL, "todos" INTEGER NOT NULL, PRIMARY KEY("id" AUTOINCREMENT))`)
	if err != nil {
		return err
	}
	return nil
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

func OpenDB(file string) (*TodosDB, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Println("Error opening database")
		log.Println(err)
		return nil, err
	}
	log.Println("Opened database")
	t := TodosDB{db}
	log.Println("Checking if table exists")
	if !t.tableExists("todos") {
		err := t.createTodosTable()
		if err != nil {
			return nil, err
		}
	}
	if !t.tableExists("lists") {
		err := t.createListsTable()
		if err != nil {
			return nil, err
		}
	}
	return &t, nil
}
