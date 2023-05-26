import { Color, type Coordinate, type Dimensions, type IPiece, type PiecePresentInfo, type Position, type SquareIdx, type SquareNotation } from "./types";

export const convertFenToPosition = (fen:string):{
  position:Position,
  dimensions:Dimensions,
  maxBoardState:PiecePresentInfo[][]
  }|undefined =>{
  let position:Position = {
    piecePositions:{},
  };
  let dimensions = {ranks:8,files:8};
  let maxBoardState = createEmptyMaxBoardState();
  let tempPiecePositions:Map<Coordinate,IPiece> = new Map();
  const fenSplit = fen.split(" ");
  const ranks = fenSplit[0].split("/");
  const rankCount = ranks.length;
  dimensions.ranks = rankCount;
  if(rankCount>16){
    return undefined;
  }
  let idx:SquareIdx=0;
  let secDigit,colEnd = 0;
  let char;
  let colCount=0;
  for (let i=0;i<ranks.length;i++){
      const rank = rankCount-i;
      secDigit = 0;
      colCount = 0;
      for(let j = 0; j < ranks[i].length; j++){
          char = ranks[i].charAt(j);
         if (char === "."){
              colCount+=1;
              idx+=1;
          }
          else if (/\d/.test(char)){
            if(j+1<ranks[i].length && (/\d/.test(ranks[i].charAt(j+1)))){
                secDigit=parseInt(char)
            } else{
                if(secDigit!=0){
                    colEnd = secDigit*10+parseInt(char)
                } else {colEnd=parseInt(char)}
                colCount+=colEnd
                idx+=colEnd;
            }
          }
          else{
            let piece = {
              color: char.toLowerCase() === char ? Color.BLACK : Color.WHITE,
              pieceType: char
            };
            maxBoardState[i][j]={isPiecePresent:true,piece}
            colCount+=1;

          }
        }
  }
  dimensions.files = colCount;
  let sqIndex = 0;
  for(let row=0;row<dimensions.ranks;row++){
    for(let col=0;col<dimensions.files;col++){
      if(maxBoardState[row][col].isPiecePresent && maxBoardState[row][col].piece){
        position.piecePositions[sqIndex] = maxBoardState[dimensions.ranks-1-row][col].piece as IPiece;
      }
      sqIndex+=1;
    }
  }
  
  return {position,dimensions,maxBoardState};
}


export const createEmptyMaxBoardState = ():PiecePresentInfo[][]=>{
  let maxBoard:PiecePresentInfo[][] = [];
  for(let row=0;row<16;row++){
      let boardRow=[];
      for(let col=0;col<16;col++){
              boardRow.push({isPiecePresent:false})
      }
      maxBoard.push(boardRow)
  }
  return maxBoard;
}