import { DragEvent, DragEventHandler } from "react";
import { Item } from "../Item";

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
    // const result = board.reduce((outerAcc, outerItem) => {
    //   list.reduce((innerAcc, innerItem) => if (innerItem.id > innerAcc) innerItem.id, 0);
    //   if (outerItem.id > outerAcc) outerItem.id, 0);
    //
    // console.log("result", result);

    newBoard[positionList].push({ id: board.length, name: `${board.length}` });
    return newBoard;
  });
};

export const List = ({ positionList, board, setBoard }) => {
  return (
    <ul
      className={`m-10`}
      onDrop={(event) => handleDrop(event, setBoard)}
      onDragOver={handleDragOver}
    >
      {board[positionList].map((item, index) => (
        <Item key={item.id} positionItem={index} positionList={positionList}>
          {item.name}
        </Item>
      ))}
      <button onClick={(event) => handleOnClick(event, positionList, setBoard)}>
        Add item
      </button>
    </ul>
  );
};
