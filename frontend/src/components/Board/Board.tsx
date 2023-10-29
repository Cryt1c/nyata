"use client";

import { useMemo } from "react";
import { List } from "./List";
import { Todo } from "./List/utils/todo";
import { BoardProvider, useSelection } from "@/hooks/useSelection";

export const Board = () => {
  const { board } = useSelection();

  const listAmount: number = useMemo(() => {
    return board.reduce((acc, item) => {
      return item.listId > acc ? item.listId : acc;
    }, 0) + 1;
  }, [board]);

  return (
    <div className={`flex justify-center`}>
        {Array.from({ length: listAmount }).map((_, index) => {
          const thisListItems: Todo[] = board.filter(
            (item) => item.listId === index
          );
          return <List key={index} items={thisListItems} />;
        })}
    </div>
  );
};
