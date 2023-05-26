import type { Coordinate, CoordinatetoIDMap, Dimensions, IPiece, PiecePositions, PiecePresentInfo, SquareIdx, SquareInfo, SquareMaps, SquareNotation } from "./types";


export const generateSquareMaps = (dimensions:Dimensions,isFlipped:boolean)=> {
    const columns = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p'].slice(0, dimensions.files);
    const rows = Array.from({ length: dimensions.ranks }, (_, i) => i);
  
    const sqToIdMap: { [key: string]: number } = {};
    const idToSqMap: { [key: number]: SquareNotation } = {};
    const coordToIdMap:CoordinatetoIDMap = {};
    let squareIndex = 0;

    let squares:Record<SquareIdx,SquareInfo> = {};
      for (const row of rows) {
        columns.forEach((column,colIdx)=>{
            const square = `${column}${row+1}` as SquareNotation;
            sqToIdMap[square] = squareIndex;
            idToSqMap[squareIndex] = square;
            coordToIdMap[`${row}:${column}`] = squareIndex;
            squares[squareIndex]= isFlipped ? {
                gridX:row+1,
                gridY:colIdx+1,
                squareIndex,row,column:colIdx} 
                : {
                gridX:dimensions.ranks-row,
                gridY:colIdx+1,
                squareIndex,row,column:colIdx}  
            squareIndex++;
        })
      }

      return {sqToIdMap,idToSqMap,coordToIdMap,squares};
}

export const updatePiecePositions = (maxBoardState:PiecePresentInfo[][],dimensions:Dimensions):PiecePositions=>{
    let piecePositions:PiecePositions = {};
    for(let row=0;row<dimensions.ranks;row++){
        for(let col=0;col<dimensions.files;col++){
            if(maxBoardState[row][col].isPiecePresent){
                piecePositions[((dimensions.ranks-row-1)*(dimensions.files))+col] = maxBoardState[row][col].piece as IPiece;
            }
        }
    }
    return piecePositions;
}


