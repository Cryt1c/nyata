"use client";

import { DragEvent, DragEventHandler, useState } from "react";
import { Item } from "@/components/Board/Item";
import { List } from "./List";

export const Board = ({ items }) => {
  const [board, setBoard] = useState(items);
  return (
    <div className={`flex justify-center`}>
      <List positionList={0} board={board} setBoard={setBoard} />
      <List positionList={1} board={board} setBoard={setBoard} />
    </div>
  );
};
