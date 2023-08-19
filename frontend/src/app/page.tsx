import { Board } from "@/components/Board";

const getTodos = async () => {
  const response = await fetch("http://localhost:8080/todos");
  return await response.json();
};

const Page = async () => {
  let items = await getTodos();
  // items = [
  //   [
  //     { id: 0, name: "Item 1" },
  //     { id: 1, name: "Item 2" },
  //   ],
  //   [
  //     { id: 2, name: "Item 3" },
  //     { id: 3, name: "Item 4" },
  //     { id: 4, name: "Item 5" },
  //   ],
  // ];
  return <Board items={items} />;
};

export default Page;
