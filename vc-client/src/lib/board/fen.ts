import { Color, type Coordinate, type Dimensions, type IPiece, type PiecePresentInfo, type Position, type SquareIdx, type SquareNotation } from "./types";

export const convertFenToPosition = (fen:string):{
  position:Position,
  dimensions:Dimensions,
  maxBoardState:PiecePresentInfo[][]
  }|undefined =>{
  let position:Position = {
    piecePositions:{},
    disabled:{}
  };
  let dimensions = {ranks:8,files:8};
  let maxBoardState = createEmptyMaxBoardState();
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
  let row=0;
  for (let i=0;i<ranks.length;i++){
      secDigit = 0;
      colCount = 0;
      let j = 0; 
      let col=0;
      while(j < ranks[i].length){
          char = ranks[i].charAt(j);
         if (char === "."){
              maxBoardState[i][col]={isPiecePresent:false,disabled:true}
              colCount+=1;
              idx+=1;
              col+=1;
          }
          else if (/\d/.test(char)){
            if(j+1<ranks[i].length && (/\d/.test(ranks[i].charAt(j+1)))){
                secDigit=parseInt(char)
            } else{
                if(secDigit!=0){
                    colEnd = secDigit*10+parseInt(char)
                } else {colEnd=parseInt(char)}
                for(let empty=0;empty<colEnd;empty++){
                  maxBoardState[row][col] = {isPiecePresent:false}
                  col+=1;
                }
                colCount+=colEnd
                idx+=colEnd;
            }
          }
          else{
            let piece = {
              color: char.toLowerCase() === char ? Color.BLACK : Color.WHITE,
              pieceType: char
            };
            maxBoardState[row][col]={isPiecePresent:true,piece}
            colCount+=1;
            col+=1;

          }
          j+=1;
        }
        row+=1
  }
  dimensions.files = colCount;  
  console.log(maxBoardState,position)
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