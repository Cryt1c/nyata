export const getHighestIdonBoard = (board) => {
  const highestIdsOnEachList = board.map((list) => {
    return list.reduce((acc, item) => {
      if (item.id > acc) {
        console.log(item.id);
        return item.id;
      }
      return acc;
    }, 0);
  });
  return Math.max(...highestIdsOnEachList);
};
