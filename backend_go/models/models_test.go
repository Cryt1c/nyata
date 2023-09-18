package models_test

import (
	"cryt1c/nyata/models"
	"log"
	"os"
	"testing"
)

var todosDb *models.TodosDB

func TestMain(m *testing.M) {
	// os.Exit skips defer calls
	// so we need to call another function
	code, err := run(m)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	todosDb, err = models.OpenDB("./test.db")
	if err != nil {
		log.Fatal(err)
	}

	// Truncates all test data after the tests are run
	defer func() {
		todosDb.DB.Exec("DELETE FROM todos")
		todosDb.DB.Exec("DELETE FROM sqlite_sequence WHERE name='todos'")
		todosDb.DB.Close()
	}()
	// Run the tests
	return m.Run(), nil
}

func TestInsert(t *testing.T) {
	id, err := todosDb.InsertTodo(
		models.Todo{
			Name:       "Test Todo",
			Completed:  false,
			PositionId: 1,
			ListId:     1,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 {
		t.Fatalf("Expected id to be 1, but was %v", id)
	}

	cmpToDBState([]models.Todo{
		{
			Id:         1,
			Name:       "Test Todo",
			Completed:  false,
			PositionId: 1,
			ListId:     1,
		},
	}, t)
}

func TestUpdate(t *testing.T) {
	affectedRows, err := todosDb.UpdateTodo(
		models.Todo{
			Id:         1,
			Name:       "Update todo",
			Completed:  true,
			ListId:     2,
			PositionId: 2,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if affectedRows != 1 {
		t.Fatal("Expected affected rows to be 1")
	}

	cmpToDBState([]models.Todo{
		{
			Id:         1,
			Name:       "Update todo",
			Completed:  true,
			PositionId: 2,
			ListId:     2,
		},
	}, t)
}

func TestDeletion(t *testing.T) {
	affectedRows, err := todosDb.DeleteTodoById(1)
	if err != nil {
		t.Fatal(err)
	}
	if affectedRows != 1 {
		t.Fatal("Expected affected rows to be 1")
	}

	cmpToDBState([]models.Todo{}, t)
}

func cmpToDBState(expected []models.Todo, t *testing.T) {
	actual, err := todosDb.GetTodos()
	if err != nil {
		t.Fatal(err)
	}
	if len(actual) != len(expected) {
		t.Fatal("Different length of todos in db")
	}

	for i, todo1 := range actual {
		if todo1.Id != expected[i].Id {
			t.Fatal("Different id of todos in db")
		}
		if todo1.Name != expected[i].Name {
			t.Fatal("Different name of todos in db")
		}
		if todo1.Completed != expected[i].Completed {
			t.Fatal("Different completed of todos in db")
		}
		if todo1.ListId != expected[i].ListId {
			t.Fatal("Different list id of todos in db")
		}
		if todo1.PositionId != expected[i].PositionId {
			t.Fatal("Different position id of todos in db")
		}
	}
}
