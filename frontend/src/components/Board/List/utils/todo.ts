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
  console.log(updatedTodoResult);
  return updatedTodoResult;
};
