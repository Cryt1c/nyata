"use client";

import { DragEvent, useMemo, useState } from "react";
import { List } from "./List";
import { useSelection } from "@/hooks/useSelection";
import { Todo } from "@/app/page";

type BoardProps = {
  items: Todo[];
};
export const Board = ({ items }: BoardProps) => {
  console.log("items", items);
  const [board, setBoard] = useState<Todo[]>(items);
  const boardLimits = useMemo(() => {
    return board.map((list) => list.length);
  }, [board]);
  const listAmount: number = useMemo(() => {
    return (
      board.reduce((acc, item) => {
        return item.listId > acc ? item.listId : acc;
      }, 0) + 1
    );
  }, [board]);
  console.log("boardLimits", boardLimits);
  console.log("listAmount", listAmount);

  const handleAccept = (event: DragEvent) => {
    console.log("accept", event);
  };
  const [selection] = useSelection(handleAccept, boardLimits);

  return (
    <div className={`flex justify-center`}>
      {Array.from({ length: listAmount }).map((_, index) => {
        const listItems: Todo[] = board.filter((item) => item.listId === index);
        return (
          <List
            positionList={index}
            listItems={listItems}
            board={board}
            setBoard={setBoard}
            selection={selection}
          />
        );
      })}
    </div>
  );
};
