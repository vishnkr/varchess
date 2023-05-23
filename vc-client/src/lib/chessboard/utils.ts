
export type File = 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p';
export type Rank = `${1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16}`;

export type SquareNotation = `${File}${Rank}`;

export const generateSquaresMap = (width: number, height: number)=> {
    const columns = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p'].slice(0, width);
    const rows = Array.from({ length: height }, (_, i) => i + 1);
  
    const sqToIdMap: { [key: string]: number } = {};
    const idToSqMap: { [key: number]: SquareNotation } = {};
    let squareIndex = 0;
    for (const column of columns) {
      for (const row of rows) {
        const square = `${column}${row}` as SquareNotation;
        sqToIdMap[square] = squareIndex;
        idToSqMap[squareIndex] = square;
        squareIndex++;
      }
    }
    return {sqToIdMap,idToSqMap};
  }

