"use client"

import { Board } from "@/components/Board";
import { getTodos } from "@/components/Board/List/utils/todo";
import { BoardProvider } from "@/hooks/useSelection";

const Page = async () => {
  let items = await getTodos();
  return (
    <BoardProvider items={items}>
      <Board />
    </BoardProvider>
  );
};

export default Page;
