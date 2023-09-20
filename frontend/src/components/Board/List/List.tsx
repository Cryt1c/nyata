import { DragEvent, DragEventHandler } from "react";
import { Item } from "../Item";
import { createTodo, updateTodo, Todo } from "./utils/todo";

const handleDrop = (event: DragEvent<HTMLUListElement>, setItems) => {
  const positionOrigin = +event.dataTransfer.getData("positionItem");
  const positionTarget = +event.target.dataset.positionItem;

  const positionOriginList = +event.dataTransfer.getData("positionList");
  const positionTargetList = +event.target.dataset.positionList;

  // console.log("event.target.dataset", event.target.dataset);
  // console.log("positionTarget", positionTarget);
  // console.log("positionOrigin", positionOrigin);
  // console.log("positionTargetList", positionTargetList);
  // console.log("positionOriginList", positionOriginList);

  setItems((items) => {
    const reorderedItems = [...items];
    const removedItem = reorderedItems[positionOriginList].splice(
      positionOrigin,
      1
    )[0];
    // Remove item from origin and insert it in the target.
    reorderedItems[positionTargetList].splice(positionTarget, 0, removedItem);
    return reorderedItems;
  });
};

const handleDragOver: DragEventHandler = (event: DragEvent) => {
  event.preventDefault();
};

const handleOnClick = (event, listId, setBoard) => {
  setBoard((board) => {
    const newBoard = [...board];
    const newTodo = createTodo({ name: "", listId: listId, positionId: 0 });
    newBoard[listId].push(newTodo);
    return newBoard;
  });
};

const handleItemTextChange = async (event, item, setBoard) => {
  const newText = event.target.value;
  setBoard((board) => {
    const newBoard = [...board];
    item.name = newText;
    return newBoard;
  });
  updateTodo(item);
};

type ListProps = {
  positionList: number;
  listItems: Todo[];
  board: Todo[];
  setBoard: (board: Todo[]) => void;
  selection: { x: number; y: number };
};

export const List = ({ listItems, setBoard, selection }) => {
  const listId = listItems[0].listId;
  const sortedListItems = listItems.sort((a, b) => a.positionId - b.positionId);
  return (
    <ul
      className={`m-10`}
      onDrop={(event) => handleDrop(event, setBoard)}
      onDragOver={handleDragOver}
    >
      {sortedListItems.map((item: Todo) => (
        <Item
          listItem={item}
          key={item.id}
          onChange={(event: Event) => handleItemTextChange(event, item, setBoard)}
          selected={
            selection.x == item.listId && selection.y == item.positionId
          }
        ></Item>
      ))}
      <button onClick={(event) => handleOnClick(event, listId, setBoard)}>
        Add item
      </button>
    </ul>
  );
};
