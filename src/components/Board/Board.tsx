"use client";

import { DragEventHandler, useState } from "react";
import { Item } from "@/components/Board/Item";

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
  // console.log("items", items);

  setItems((items) => {
    const newItems = [...items];
    // Remove item from origin and insert it in the target.
    newItems[positionTargetList].splice(
      positionTarget,
      0,
      newItems[positionOriginList].splice(positionOrigin, 1)[0]
    );
    return newItems;
  });
};

const handleDragOver: DragEventHandler = (event: DragEvent) => {
  event.preventDefault();
};

export const Board = ({ items }) => {
  const [board, setBoard] = useState(items);
  return (
    <div className={`flex justify-center`}>
      <ul
        className={``}
        onDrop={(event) => handleDrop(event, setBoard)}
        onDragOver={handleDragOver}
      >
        {items[0].map((item, index) => (
          <Item key={item.id} positionItem={index} positionList={0}>
            {item.name}
          </Item>
        ))}
      </ul>
      <ul
        className={``}
        onDrop={(event) => handleDrop(event, setBoard)}
        onDragOver={handleDragOver}
      >
        {items[1].map((item, index) => (
          <Item key={item.id} positionItem={index} positionList={1}>
            {item.name}
          </Item>
        ))}
      </ul>
    </div>
  );
};
