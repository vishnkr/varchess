import { BoardState, Square } from "@/types"

export const createDefaultMaxBoardStateSquares = () => {
    let squares :Square[][] =[];
    for(let row=0; row<16;row++){
        squares.push([])
        for(let col=0;col<16;col++){
            let square: Square = {
                disabled:false, 
                squareId: (row*16)+col,
                squareInfo: {
                    isPiecePresent: false,
                    row: row,
                    col: col,
                    squareColor: isLight(row,col) ? 'light' : 'dark'
                }
            };
            if( (col===0||col==7) && (row===0||row===7)) {
                square.squareInfo.pieceType='r';
                square.squareInfo.isPiecePresent=true
            }
            else if((col===1||col==6) && (row===0||row===7)){
                square.squareInfo.pieceType='n';
                square.squareInfo.isPiecePresent=true
            }
            else if((col===2||col==5) && (row===0||row===7)){
                square.squareInfo.pieceType='b';
                square.squareInfo.isPiecePresent=true
            }
            else if((col===3) && (row===0||row===7)){
                square.squareInfo.pieceType='q';
                square.squareInfo.isPiecePresent=true
            }
            else if((col===4) && (row===0||row===7)){
                square.squareInfo.pieceType='k';
                square.squareInfo.isPiecePresent=true
            }
            else if((row===1||row===6) && col<8){
                square.squareInfo.pieceType='p';
                square.squareInfo.isPiecePresent=true
            }
            if(row==0||row==1){square.squareInfo.pieceColor='black'}
            else if(row==6||row==7){square.squareInfo.pieceColor='white'}
            squares[row].push(square)
        }
    }
    return squares
}

export const isEven = (val:number) => {return val%2==0}

export const isLight =(row:number,col:number) => {
  return isEven(row)&&isEven(col)|| (!isEven(row)&&!isEven(col))
}


export const setupEmptyMaxSizeBoard = ()=>{
    let squares :Square[][] =[];
    for(let row=0; row<16;row++){
        squares.push([])
        for(let col=0;col<16;col++){
            let square: Square = {
                squareId: (row*16)+col,
                disabled:false, 
                squareInfo: {
                    isPiecePresent: false,
                    row: row,
                    col: col,
                    squareColor: isLight(row,col) ? 'light' : 'dark'
                }
            };
            squares[row].push(square)
        }
  }
  return squares;
}

export const withId = (boardState:BoardState):BoardState =>{
    for(let row=0; row<boardState.squares.length;row++){
        for(let col=0;col<boardState.squares[row].length;col++){
            boardState.squares[row][col].squareId = (row*boardState.dimensions.cols)+col
        }
    }
    return boardState
}