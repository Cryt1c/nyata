import { Board } from "@/components/Board";

export type Todo = {
  id: number;
  name: string;
  listId: number;
  positionId: number;
};

const getTodos = async () => {
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

const Page = async () => {
  let items = await getTodos();
  return <Board items={items} />;
};

export default Page;
