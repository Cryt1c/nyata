import { Board } from "@/components/Board";

const items = [
  [
    { id: 0, name: "Item 1" },
    { id: 1, name: "Item 2" },
  ],
  [
    { id: 2, name: "Item 3" },
    { id: 3, name: "Item 4" },
    { id: 4, name: "Item 5" },
  ],
];

const Page = () => {
  return <Board items={items} />;
};

export default Page;
