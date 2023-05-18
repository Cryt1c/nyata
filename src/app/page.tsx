"use client";

import { Item } from "@/components/List/Item";
import { useState } from "react";

const handleDrop = (event, setItems, items) => {
  const positionOrigin = +event.dataTransfer.getData("positionItem");
  const positionTarget = +event.target.dataset.positionItem;

  const positionOriginList = +event.dataTransfer.getData("positionList");
  const positionTargetList = +event.target.dataset.positionList;

  console.log("event.target.dataset", event.target.dataset);
  console.log("positionTarget", positionTarget);
  console.log("positionOrigin", positionOrigin);
  console.log("positionTargetList", positionTargetList);
  console.log("positionOriginList", positionOriginList);
  console.log("items", items);

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

const handleDragOver = (event) => {
  event.preventDefault();
};

const Page = () => {
  const [items, setItems] = useState([
    [
      { id: 0, name: "Item 1" },
      { id: 1, name: "Item 2" },
    ],
    [
      { id: 2, name: "Item 3" },
      { id: 3, name: "Item 4" },
      { id: 4, name: "Item 5" },
    ],
  ]);

  return (
    <div className={`flex justify-center`}>
      <ul
        className={``}
        onDrop={(event) => handleDrop(event, setItems, items)}
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
        onDrop={(event) => handleDrop(event, setItems, items)}
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

export default Page;
