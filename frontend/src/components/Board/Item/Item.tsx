"use client";

import { useRef } from "react";

const handleDragStart = (event, positionId, listId) => {
  event.dataTransfer.setData("positionId", positionId);
  event.dataTransfer.setData("listId", listId);
};

export const Item = ({ listItem, onChange, selected, changing }) => {
  const inputRef = useRef<HTMLInputElement | null>(null);

  if (changing && selected) {
    if (inputRef.current) {
      inputRef.current.focus();
    }
  }

  return (
    <li
      draggable
      onDragStart={(e) =>
        handleDragStart(e, listItem.positionId, listItem.listId)
      }
      className={`p-5 border-2`}
      data-position-id={listItem.positionId}
      data-list-id={listItem.listId}
      style={{ backgroundColor: selected ? "red" : "white" }}
    >
      <input
        type="text"
        value={listItem.name}
        onChange={onChange}
        // @todo: Find another solution
        data-position-id={listItem.positionId}
        data-list-id={listItem.listId}
        ref={inputRef}
      />
    </li>
  );
};
