"use client";

import { useMemo } from "react";
import { List } from "./List";
import { BoardProvider } from "@/hooks/useSelection";
import { Todo } from "./List/utils/todo";

type BoardProps = {
  items: Todo[];
};
export const Board = ({ items }: BoardProps) => {
  const listAmount: number = useMemo(() => {
    return (
      items.reduce((acc, item) => {
        return item.listId > acc ? item.listId : acc;
      }, 0) + 1
    );
  }, [items]);

  return (
    <div className={`flex justify-center`}>
      <BoardProvider items={items}>
        {Array.from({ length: listAmount }).map((_, index) => {
          const thisListItems: Todo[] = items.filter(
            (item) => item.listId === index
          );
          return <List key={index} items={thisListItems} />;
        })}
      </BoardProvider>
    </div>
  );
};
