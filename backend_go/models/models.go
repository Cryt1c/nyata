package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type TodosModel struct {
	DB *sql.DB
}

type Todo struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Completed  bool   `json:"completed"`
	ListId     int64  `json:"listId"`
	PositionId int64  `json:"positionId"`
}

func (m *TodosModel) GetTodos() ([]Todo, error) {
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

func (m *TodosModel) InsertTodo(todo Todo) (int64, error) {
	result, err := m.DB.Exec(
		"INSERT INTO todos(name, completed, position_id, list_id) VALUES(?, ?, ?, ?)",
		todo.Name,
		todo.Completed,
		todo.PositionId,
		todo.ListId,
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

func (m *TodosModel) DeleteTodoById(id int64) (int64, error) {
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

func (m *TodosModel) UpdateTodo(todo Todo) (int64, error) {
	result, err := m.DB.Exec("UPDATE todos SET name = ?, completed = ?, position_id = ?, list_id = ? WHERE id = ?", todo.Name, todo.Completed, todo.PositionId, todo.ListId, todo.Id)
	if err != nil {
		log.Printf("Error updating todo with ID %d: %v", todo.Id, err)
		return 0, fmt.Errorf("Failed to update todo with ID %d: %w", todo.Id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected after deleting todo with ID %d: %v", todo.Id, err)
		return 0, fmt.Errorf("Failed to get rows affected after deleting todo with ID %d: %w", todo.Id, err)
	}

	if rowsAffected == 0 {
		log.Printf("No todo found with ID %d", todo.Id)
		return 0, errors.New(fmt.Sprintf("No todo found with ID %d", todo.Id))
	}

	return rowsAffected, nil
}

func (m *TodosModel) tableExists(name string) bool {
	if _, err := m.DB.Query("SELECT * FROM todos"); err == nil {
		return true
	}
	return false
}

func (m *TodosModel) createTable() error {
	_, err := m.DB.Exec(`CREATE TABLE "todos" ( "id" INTEGER, "name" TEXT NOT NULL, "completed" INTEGER NOT NULL, "position_id" INTEGER NOT NULL, "list_id" INTEGER NOT NULL, PRIMARY KEY("id" AUTOINCREMENT))`)
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

func OpenDB(file string) (*TodosModel, error) {
	db, err := sql.Open("sqlite3", file)
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
