export type Todo = {
  id?: number;
  name: string;
  completed: boolean;
  listId: number;
  positionId?: number;
};

export const getTodos = async (): Promise<Todo[]> => {
  const response = await fetch("http://localhost:8080/todos");

  if (!response.ok || response.status !== 200) {
    throw new Error("Error getting todos");
  }

  let items: Todo[] = response.json();
  // items = [
  //   { id: 0, name: "Item 1", listId: 0, positionId: 0 },
  //   { id: 1, name: "Item 2", listId: 0, positionId: 1 },
  //   { id: 2, name: "Item 3", listId: 1, positionId: 0 },
  //   { id: 3, name: "Item 4", listId: 1, positionId: 1 },
  //   { id: 4, name: "Item 5", listId: 1, positionId: 2 },
  // ];
  return items;
};

export const createTodo = async (newTodo: Todo): Promise<Todo> => {
  const response = await fetch("http://localhost:8080/todo", {
    method: "POST",
    body: JSON.stringify(newTodo),
  });

  if (!response.ok || response.status !== 201) {
    throw new Error("Error creating todo");
  }

  const newTodoWithId = await response.json();
  return newTodoWithId;
};

export const updateTodo = async (updatedTodo: Todo): Promise<Todo> => {
  const response = await fetch("http://localhost:8080/todo", {
    method: "PUT",
    body: JSON.stringify(updatedTodo),
  });

  if (!response.ok || response.status !== 200) {
    throw new Error("Error updating todo");
  }

  const updatedTodoResult = await response.json();
  return updatedTodoResult;
};

export const deleteTodo = async (deletedTodo: Todo): Promise<void> => {
  const response = await fetch(`http://localhost:8080/todo/${deletedTodo.id}`, {
    method: "DELETE",
  });

  console.log("deletedTodo", deletedTodo);
  if (!response.ok || response.status !== 200) {
    throw new Error("Error deleting todo");
  }
};

export const reorderTodos = async (origin: Todo, target: Todo) => {
  const response = await fetch("http://localhost:8080/reorder", {
    method: "PUT",
    body: JSON.stringify({ origin, target }),
  });

  if(!response.ok || response.status !== 200) {
    throw new Error("Error reordering todos");
  }
  const reorderedTodos = await response.json();
  return reorderedTodos;
}
