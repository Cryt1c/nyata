"use client";

import {
  DragEvent,
  DragEventHandler,
  useEffect,
  useMemo,
  useState,
} from "react";
import { Item } from "@/components/Board/Item";
import { List } from "./List";
import { useSelection } from "@/hooks/useSelection";

export const Board = ({ items }) => {
  const [board, setBoard] = useState(items);
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
