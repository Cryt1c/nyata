import { Board } from "@/components/Board";
import { getTodos } from "@/components/Board/List/utils/list";

const Page = async () => {
  let items = await getTodos();
  return <Board items={items} />;
};

export default Page;
