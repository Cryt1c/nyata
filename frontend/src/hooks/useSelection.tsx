import { Todo } from "@/components/Board/List/utils/todo";
import {
  Dispatch,
  SetStateAction,
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";

export interface Position {
  x: number;
  y: number;
}

const BoardContext = createContext<{
  board: Todo[];
  setBoard: Dispatch<SetStateAction<Todo[]>>;
}>({ board: [], setBoard: () => {} });

const SelectionContext = createContext<{
  selection: Position;
  setSelection: Dispatch<SetStateAction<Position>>;
}>({ selection: { x: 0, y: 0 }, setSelection: () => {} });

const ChangingContext = createContext<{
  changing: boolean;
  setChanging: Dispatch<SetStateAction<boolean>>;
}>({ changing: false, setChanging: () => {} });

interface SelectionProviderProps {
  children: React.ReactNode;
  items: Todo[];
}

export const BoardProvider = ({ children, items }: SelectionProviderProps) => {
  const [selection, setSelection] = useState<Position>({ x: 0, y: 0 });
  const [changing, setChanging] = useState(false);
  const [board, setBoard] = useState<Todo[]>(items);
  return (
    <BoardContext.Provider value={{ board, setBoard }}>
      <SelectionContext.Provider value={{ selection, setSelection }}>
        <ChangingContext.Provider value={{ changing, setChanging }}>
          {children}
        </ChangingContext.Provider>
      </SelectionContext.Provider>
    </BoardContext.Provider>
  );
};

export const useSelection = () => {
  const { selection, setSelection } = useContext(SelectionContext);
  const { changing, setChanging } = useContext(ChangingContext);
  const { board, setBoard } = useContext(BoardContext);

  const boardLimits = useMemo(() => {
    return board.reduce((acc: number[], item) => {
      if (acc[item.listId]) {
        acc[item.listId] = acc[item.listId] + 1;
      } else {
        acc[item.listId] = 1;
      }
      return acc;
    }, []);
  }, [board]);

  const handleKeyDown = (event: KeyboardEvent) => {
    switch (event.key) {
      // Up.
      case "k":
        if (changing) return;
        setSelection((oldSelection: Position) => ({
          x: oldSelection.x,
          y: Math.max(oldSelection.y - 1, 0),
        }));
        break;
      // Down.
      case "j":
        if (changing) return;
        setSelection((oldSelection: Position) => ({
          x: oldSelection.x,
          y: Math.min(oldSelection.y + 1, boardLimits[oldSelection.x] - 1),
        }));
        break;
      case "l":
        if (changing) return;
        // Right.
        setSelection((oldSelection: Position) => ({
          x: Math.min(oldSelection.x + 1, boardLimits.length - 1),
          y: oldSelection.y,
        }));
        break;
      // Left.
      case "h":
        if (changing) return;
        setSelection((oldSelection: Position) => ({
          x: Math.max(oldSelection.x - 1, 0),
          y: oldSelection.y,
        }));
        break;
      case "c":
        if (changing) return;
        setChanging(true);
        break;
      case "Escape":
        setChanging(false);
        document.activeElement.blur();
        break;
      default:
        console.log("unhandled key", event.key);
        break;
    }
  };

  useEffect(() => {
    window.addEventListener("keyup", handleKeyDown);
    return () => {
      removeEventListener("keyup", handleKeyDown);
    };
  }, []);

  return { selection, changing, board, setBoard };
};
