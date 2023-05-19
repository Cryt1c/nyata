"use client";

import { useState } from "react";

const handleDragStart = (e, positionItem, positionList) => {
  event.dataTransfer.setData("positionItem", positionItem);
  event.dataTransfer.setData("positionList", positionList);
};

// https://codesandbox.io/s/framer-motion-drag-to-reorder-pkm1k?file=/src/Example.tsx:1479-1525
export const Item = ({ children, positionItem, positionList }) => {
  return (
    <li
      draggable
      onDragStart={(e) => handleDragStart(e, positionItem, positionList)}
      className={`p-5 border-black border-solid border-1`}
      data-position-item={positionItem}
      data-position-list={positionList}
    >
      {children}
    </li>
  );
};
