
import type { Coordinate, CoordinatetoIDMap, Dimensions, DisabledSquares, IPiece, Position, PiecePresentInfo, SquareIdx, SquareInfo, SquareMaps, SquareNotation, PiecePositions } from "./types";


export const generateSquareMaps = (dimensions:Dimensions,isFlipped:boolean)=> {
    const columns = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p'].slice(0, dimensions.files);
    const rows = Array.from({ length: dimensions.ranks }, (_, i) => i);
  
    //const sqToIdMap: { [key: string]: number } = {};
    //const idToSqMap: { [key: number]: SquareNotation } = {};
    const coordToIdMap:CoordinatetoIDMap = {};
    let squareIndex = 0;

    let squares:Record<SquareIdx,SquareInfo> = {};
      for (const row of rows) {
        columns.forEach((column,colIdx)=>{
            const squareNotation = `${column}${row+1}` as SquareNotation;
            //sqToIdMap[squareNotation] = squareIndex;
            //idToSqMap[squareIndex] = squareNotation;
            coordToIdMap[`${row}:${column}`] = squareIndex;
            squares[squareIndex]=  {
                gridX: isFlipped ? row+1 : dimensions.ranks-row,
                gridY: colIdx+1,
                squareIndex,
                squareNotation,
                row,
                column:colIdx
            } 
            squareIndex++;
        })
      }

      return {/*sqToIdMap,idToSqMap,*/coordToIdMap,squares};
}

export const updatePiecePositions = (maxBoardState:PiecePresentInfo[][],dimensions:Dimensions):Position=>{
    let piecePositions:PiecePositions = {};
    let disabled:DisabledSquares = {}
    for(let row=0;row<dimensions.ranks;row++){
        for(let col=0;col<dimensions.files;col++){
            if(maxBoardState[row][col].isPiecePresent){
                piecePositions[((dimensions.ranks-row-1)*(dimensions.files))+col] = maxBoardState[row][col].piece as IPiece;
            } else if(maxBoardState[row][col].disabled){
                disabled[((dimensions.ranks-row-1)*(dimensions.files))+col] = true;
            }
        }
    }
    return {piecePositions,disabled};
}


