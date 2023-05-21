export const getHighestIdonBoard = (board) => {
  const highestIdsOnEachList = board.map((list) => {
    console.log(list);
    return list.reduce((acc, item) => {
      if (item.id > acc) {
        console.log(item.id);
        return item.id;
      }
      return acc;
    }, 0);
  });
  console.log(highestIdsOnEachList);
  return Math.max(...highestIdsOnEachList);
};
