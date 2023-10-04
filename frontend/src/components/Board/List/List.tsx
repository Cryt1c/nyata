import { DragEvent, DragEventHandler } from "react";
import { Item } from "../Item";
import { createTodo, updateTodo, Todo, reorderTodos } from "./utils/todo";

const handleDrop = async (
  event: DragEvent<HTMLUListElement>,
  board: Todo[],
  setBoard: (board: Todo[]) => void
) => {
  const originPositionId = +event.dataTransfer.getData("positionId");
  const targetPositionId = +event.target.dataset.positionId;

  const originListId = +event.dataTransfer.getData("listId");
  const targetListId = +event.target.dataset.listId;

  console.log("event.target.dataset", event.target.dataset);
  console.log("originPositionId", originPositionId);
  console.log("targetPositionId", targetPositionId);
  console.log("originListId", originListId);
  console.log("targetListId", targetListId);

  if (
    isNaN(targetPositionId) ||
    isNaN(targetListId) ||
    isNaN(originPositionId) ||
    isNaN(originListId)
  ) {
    return;
  }

  const origin = board.filter(
    (item) =>
      item.positionId === originPositionId && item.listId === originListId
  )[0];
  const target = board.filter(
    (item) =>
      item.positionId === targetPositionId && item.listId === targetListId
  )[0];

  const newBoard = await reorderTodos(origin, target);
  setBoard((board: Todo[]) => {
    return newBoard;
    //   const newBoard = board.map((item: Todo) => {
    //     if (
    //       item.positionId === originPositionId &&
    //       item.listId === originListId
    //     ) {
    //       item.positionId = targetPositionId;
    //       item.listId = targetListId;
    //       return item;
    //     }
    //     if (item.listId === targetListId && item.positionId >= targetPositionId) {
    //       item.positionId = item.positionId + 1;
    //       return item;
    //     }
    //     if (item.listId === originListId && item.positionId >= originPositionId) {
    //       item.positionId = item.positionId - 1;
    //       return item;
    //     }
    //     return item;
    //   });
    //   console.log("newBoard", newBoard);
    //   return newBoard;
  });
};

const handleDragOver: DragEventHandler = (event: DragEvent) => {
  event.preventDefault();
};

const handleOnClick = async (listId: number, setBoard) => {
  const newTodo = await createTodo({
    name: "",
    completed: false,
    listId: listId,
  });
  setBoard((board: Todo[]) => {
    const newBoard = [...board, newTodo];
    return newBoard;
  });
};

const handleItemTextChange = async (event, item, setBoard) => {
  item.name = event.target.value;
  const updatedItem = await updateTodo(item);
  setBoard((board: Todo[]) => {
    const newBoard = board.map((item: Todo) => {
      if (item.id === updatedItem.id) {
        return updatedItem;
      }
      return item;
    });
    return newBoard;
  });
};

type ListProps = {
  items: Todo[];
  board: Todo[];
  setBoard: (board: Todo[]) => void;
  selection: { x: number; y: number };
};

export const List = ({ items, board, setBoard, selection }: ListProps) => {
  const listId = items[0]?.listId || 0;
  const sortedListItems = items.sort((a, b) => a.positionId - b.positionId);
  return (
    <ul
      className={`m-10`}
      onDrop={(event) => handleDrop(event, board, setBoard)}
      onDragOver={handleDragOver}
    >
      {sortedListItems.map((item: Todo) => (
        <Item
          listItem={item}
          key={item.id}
          onChange={(event: Event) =>
            handleItemTextChange(event, item, setBoard)
          }
          selected={
            selection.x == item.listId && selection.y == item.positionId
          }
        ></Item>
      ))}
      <button onClick={() => handleOnClick(listId, setBoard)}>Add item</button>
    </ul>
  );
};
