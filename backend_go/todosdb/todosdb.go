package todosdb

import (
	"cryt1c/nyata/models"
	"database/sql"
	"fmt"
	"log"
	"os"
)

type TodosDB struct {
	db *sql.DB
}

func (db *TodosDB) Close() error {
	return db.db.Close()
}

func (db *TodosDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.db.Exec(query, args...)
}

func (db *TodosDB) GetTodos() ([]models.Todo, error) {
	var todos []models.Todo
	rows, err := db.db.Query("SELECT * FROM todos")
	if err != nil {
		return todos, fmt.Errorf("unable to get values: %w", err)
	}
	for rows.Next() {
		var todo models.Todo
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

func (db *TodosDB) insertTodo(todo models.Todo) {

}

func (t *TodosDB) tableExists(name string) bool {
	if _, err := t.db.Query("SELECT * FROM todos"); err == nil {
		return true
	}
	return false
}

func (t *TodosDB) createTable() error {
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

// openDB opens a SQLite database and stores that database in our special spot.
func OpenDB() (*TodosDB, error) {
	db, err := sql.Open("sqlite3", "./todos.db")
	if err != nil {
		log.Println("Error opening database")
		log.Println(err)
		return nil, err
	}
	log.Println("Opened database")
	t := TodosDB{db}
	log.Println("Checking if table exists")
	if !t.tableExists("todos") {
		err := t.createTable()
		if err != nil {
			return nil, err
		}
	}
	return &t, nil
}
