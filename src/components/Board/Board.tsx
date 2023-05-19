"use client";

import { DragEventHandler, useState } from "react";
import { Item } from "@/components/Board/Item";

const handleDrop = (event: DragEvent<HTMLUListElement>, setItems) => {
  const positionOrigin = +event.dataTransfer.getData("positionItem");
  const positionTarget = +event.target.dataset.positionItem;

  const positionOriginList = +event.dataTransfer.getData("positionList");
  const positionTargetList = +event.target.dataset.positionList;

  console.log("event.target.dataset", event.target.dataset);
  console.log("positionTarget", positionTarget);
  console.log("positionOrigin", positionOrigin);
  console.log("positionTargetList", positionTargetList);
  console.log("positionOriginList", positionOriginList);

  setItems((items) => {
    const newItems = [...items];
    const removedItem = newItems[positionOriginList].splice(
      positionOrigin,
      1
    )[0];
    // Remove item from origin and insert it in the target.
    newItems[positionTargetList].splice(positionTarget, 0, removedItem);
    return newItems;
  });
};

const handleDragOver: DragEventHandler = (event: DragEvent) => {
  event.preventDefault();
};

const handleOnClick = (event, positionList, setBoard) => {
  setBoard((board) => {
    const newBoard = [...board];
    // @todo Create unique ids
    newBoard[positionList].push({ id: 5, name: "Item 6" });
    return newBoard;
  });
};

export const Board = ({ items }) => {
  const [board, setBoard] = useState(items);
  return (
    <div className={`flex justify-center`}>
      <ul
        className={`m-10`}
        onDrop={(event) => handleDrop(event, setBoard)}
        onDragOver={handleDragOver}
      >
        {board[0].map((item, index) => (
          <Item key={item.id} positionItem={index} positionList={0}>
            {item.name}
          </Item>
        ))}
        <button onClick={(event) => handleOnClick(event, 0, setBoard)}>
          Add item
        </button>
      </ul>
      <ul
        className={`m-10`}
        onDrop={(event) => handleDrop(event, setBoard)}
        onDragOver={handleDragOver}
      >
        {board[1].map((item, index) => (
          <Item key={item.id} positionItem={index} positionList={1}>
            {item.name}
          </Item>
        ))}
        <button onClick={(event) => handleOnClick(event, 1, setBoard)}>
          Add item
        </button>
      </ul>
    </div>
  );
};
