"use client";

const handleDragStart = (event, positionId, listId) => {
  event.dataTransfer.setData("positionItem", positionId);
  event.dataTransfer.setData("positionList", listId);
};

export const Item = ({ listItem, onChange, selected }) => {
  return (
    <li
      draggable
      onDragStart={(e) =>
        handleDragStart(e, listItem.positionId, listItem.listId)
      }
      className={`p-5 border-2`}
      data-position-item={listItem.positionId}
      data-position-list={listItem.listId}
      style={{ backgroundColor: selected ? "red" : "white" }}
    >
      <input
        type="text"
        value={listItem.name}
        onChange={onChange}
        // @todo: Find another solution
        data-position-item={listItem.positionId}
        data-position-list={listItem.listId}
      />
    </li>
  );
};
