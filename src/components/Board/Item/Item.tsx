"use client";

import { useState } from "react";

const handleDragStart = (e, positionItem, positionList) => {
  event.dataTransfer.setData("positionItem", positionItem);
  event.dataTransfer.setData("positionList", positionList);
};

// https://codesandbox.io/s/framer-motion-drag-to-reorder-pkm1k?file=/src/Example.tsx:1479-1525
export const Item = ({ children, positionItem, positionList }) => {
  const [text, setText] = useState(children);
  return (
    <li
      draggable
      onDragStart={(e) => handleDragStart(e, positionItem, positionList)}
      className={`p-5 border-2`}
      data-position-item={positionItem}
      data-position-list={positionList}
    >
      <input
        type="text"
        value={text}
        onChange={(e) => setText(e.target.value)}
        // @todo: Find another solution
        data-position-item={positionItem}
        data-position-list={positionList}
      />
    </li>
  );
};
