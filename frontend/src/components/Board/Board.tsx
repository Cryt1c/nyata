"use client";

import { DragEvent, useMemo, useState } from "react";
import { List } from "./List";
import { useSelection } from "@/hooks/useSelection";
import { Todo } from "@/app/page";

type BoardProps = {
  items: Todo[];
};
export const Board = ({ items }: BoardProps) => {
  const [board, setBoard] = useState<Todo[]>(items);
  const boardLimits = useMemo(() => {
    return board.map((list) => list.length);
  }, [board]);

  const handleAccept = (event: DragEvent) => {
    console.log("accept", event);
  };
  const [selection] = useSelection(handleAccept, boardLimits);

  return (
    <div className={`flex justify-center`}>
      <List
        positionList={0}
        board={board}
        setBoard={setBoard}
        selection={selection}
      />
      <List
        positionList={1}
        board={board}
        setBoard={setBoard}
        selection={selection}
      />
    </div>
  );
};
