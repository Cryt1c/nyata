import { useEffect, useState } from "react";

export interface Position {
  x: number;
  y: number;
}

export const useSelection = (handleAccept, boardLimits) => {
  const [selection, setSelection] = useState<Position>({ x: 0, y: 0 });
  const handleKeyDown = (event) => {
    switch (event.key) {
      // Up.
      case "k":
        setSelection((oldSelection) => ({
          x: oldSelection.x,
          y: Math.max(oldSelection.y - 1, 0),
        }));
        break;
      // Down.
      case "j":
        setSelection((oldSelection) => ({
          x: oldSelection.x,
          y: Math.min(oldSelection.y + 1, boardLimits[oldSelection.x] - 1),
        }));
        break;
      case "l":
        // Right.
        setSelection((oldSelection) => ({
          x: Math.min(oldSelection.x + 1, boardLimits.length - 1),
          y: oldSelection.y,
        }));
        break;
      // Left.
      case "h":
        setSelection((oldSelection) => ({
          x: Math.max(oldSelection.x - 1, 0),
          y: oldSelection.y,
        }));
        break;
      case "c":
        handleAccept(selection);
        break;
    }
  };

  useEffect(() => {
    window.addEventListener("keyup", handleKeyDown);
    return () => removeEventListener("keyup", handleKeyDown);
  }, []);

  return [selection];
};
