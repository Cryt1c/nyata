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
	_, err := todosDb.InsertTodo(
		models.Todo{
			Name:      "Test Todo",
			Completed: false,
			ListId:    1,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = todosDb.InsertTodo(
		models.Todo{
			Name:      "Test Todo",
			Completed: false,
			ListId:    1,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	cmpToDBState([]models.Todo{
		{
			Id:         1,
			Name:       "Test Todo",
			Completed:  false,
			PositionId: 10,
			ListId:     1,
		},
		{
			Id:         2,
			Name:       "Test Todo",
			Completed:  false,
			PositionId: 20,
			ListId:     1,
		},
	}, t)
}

func TestUpdate(t *testing.T) {
	_, err := todosDb.UpdateTodo(
		models.Todo{
			Id:         1,
			Name:       "Update todo",
			Completed:  true,
			ListId:     2,
			PositionId: 10,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	cmpToDBState([]models.Todo{
		{
			Id:         1,
			Name:       "Update todo",
			Completed:  true,
			PositionId: 10,
			ListId:     2,
		},
		{
			Id:         2,
			Name:       "Test Todo",
			Completed:  false,
			PositionId: 20,
			ListId:     1,
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

	cmpToDBState([]models.Todo{
		{
			Id:         2,
			Name:       "Test Todo",
			Completed:  false,
			PositionId: 20,
			ListId:     1,
		},
	}, t)
}

func TestGetByID(t *testing.T) {
	todo, err := todosDb.GetTodoById(2)
	if err != nil {
		t.Fatal(err)
	}

	if todo.Id != 2 {
		t.Fatal("Expected todo id to be 2")
	}
	if todo.Name != "Test Todo" {
		t.Fatal("Expected todo name to be 'Test Todo'")
	}
	if todo.Completed != false {
		t.Fatal("Expected todo completed to be false")
	}
	if todo.PositionId != 20 {
		t.Fatal("Expected todo position to be 20")
	}
	if todo.ListId != 1 {
		t.Fatal("Expected todo list id to be 1")
	}
}

func TestReorder(t *testing.T) {
	_, err := todosDb.InsertTodo(
		models.Todo{
			Name:      "Test Todo",
			Completed: false,
			ListId:    1,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	origin, err := todosDb.InsertTodo(
		models.Todo{
			Name:      "Reorder Todo",
			Completed: false,
			ListId:    1,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	todos, _ := todosDb.GetTodos(models.GetOptions{Sorted: true})
	t.Log(todos)

	_, err = todosDb.ResetListOrder(1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(todos)

	cmpToDBState([]models.Todo{
		{
			Id:         2,
			Name:       "Test Todo",
			Completed:  false,
			PositionId: 0,
			ListId:     1,
		},
		{
			Id:         3,
			Name:       "Test Todo",
			Completed:  false,
			PositionId: 10,
			ListId:     1,
		},
		{
			Id:         4,
			Name:       "Reorder Todo",
			Completed:  false,
			PositionId: 20,
			ListId:     1,
		},
	}, t)

	target := models.Todo{
		PositionId: 0,
		ListId:     1,
	}

	todos, err = todosDb.ReorderTodos(models.Reorder{origin, target})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(todos)

	cmpToDBState([]models.Todo{
		{
			Id:         2,
			Name:       "Test Todo",
			Completed:  false,
			PositionId: 0,
			ListId:     1,
		},
		{
			Id:         3,
			Name:       "Test Todo",
			Completed:  false,
			PositionId: 10,
			ListId:     1,
		},
		{
			Id:         4,
			Name:       "Reorder Todo",
			Completed:  false,
			PositionId: 5,
			ListId:     1,
		},
	}, t)

	origin = models.Todo{
		Id:         3,
		Name:       "Test Todo",
		PositionId: 10,
		ListId:     1,
	}
	target = models.Todo{
		PositionId: 0,
		ListId:     1,
	}

	todos, err = todosDb.ReorderTodos(models.Reorder{origin, target})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(todos)

	origin = models.Todo{
		Id:         4,
		Name:       "Reorder Todo",
		PositionId: 5,
		ListId:     1,
	}
	target = models.Todo{
		PositionId: 0,
		ListId:     1,
	}

	todos, err = todosDb.ReorderTodos(models.Reorder{origin, target})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(todos)

	origin = models.Todo{
		Id:         3,
		Name:       "Test Todo",
		PositionId: 2,
		ListId:     1,
	}
	target = models.Todo{
		PositionId: 0,
		ListId:     1,
	}

	todos, _ = todosDb.GetTodos(models.GetOptions{Sorted: true})
	t.Log(todos)

	todos, err = todosDb.ReorderTodos(models.Reorder{origin, target})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(todos)
}

func cmpToDBState(expected []models.Todo, t *testing.T) {
	actual, err := todosDb.GetTodos(models.GetOptions{Sorted: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(actual) != len(expected) {
		t.Logf("Actual: %v", actual)
		t.Logf("Expected: %v", expected)
		t.Fatal("Different length of todos in db")
	}

	for i, todo1 := range actual {
		if todo1.Id != expected[i].Id {
			t.Logf("Actual: %v", actual)
			t.Logf("Expected: %v", expected)
			t.Fatalf("Different id of todos in db: %d != %d", todo1.Id, expected[i].Id)
		}
		if todo1.Name != expected[i].Name {
			t.Logf("Actual: %v", actual)
			t.Logf("Expected: %v", expected)
			t.Fatalf("Different name of todos in db: %s != %s", todo1.Name, expected[i].Name)
		}
		if todo1.Completed != expected[i].Completed {
			t.Logf("Actual: %v", actual)
			t.Logf("Expected: %v", expected)
			t.Fatalf("Different completed of todos in db: %t != %t", todo1.Completed, expected[i].Completed)
		}
		if todo1.ListId != expected[i].ListId {
			t.Logf("Actual: %v", actual)
			t.Logf("Expected: %v", expected)
			t.Fatalf("Different list id of todos in db: %d != %d", todo1.ListId, expected[i].ListId)
		}
		if todo1.PositionId != expected[i].PositionId {
			t.Logf("Actual: %v", actual)
			t.Logf("Expected: %v", expected)
			t.Fatalf("Different position id of todos in db: %d != %d", todo1.PositionId, expected[i].PositionId)
		}
	}
}
