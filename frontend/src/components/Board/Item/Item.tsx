"use client";

import { useState } from "react";

const handleDragStart = (e, positionItem, positionList) => {
  event.dataTransfer.setData("positionItem", positionItem);
  event.dataTransfer.setData("positionList", positionList);
};

export const Item = ({ text, positionItem, positionList, onChange, selected }) => {
  return (
    <li
      draggable
      onDragStart={(e) => handleDragStart(e, positionItem, positionList)}
      className={`p-5 border-2`}
      data-position-item={positionItem}
      data-position-list={positionList}
      style={{ backgroundColor: selected ? "red" : "white" }}
    >
      <input
        type="text"
        value={text}
        onChange={onChange}
        // @todo: Find another solution
        data-position-item={positionItem}
        data-position-list={positionList}
      />
    </li>
  );
};
