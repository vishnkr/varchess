import { BoardState, Square } from "@/types"

export const createDefaultMaxBoardStateSquares = () => {
    let squares :Square[][] =[];
    for(let row=0; row<16;row++){
        squares.push([])
        for(let col=0;col<16;col++){
            let square: Square = {
                disabled:false, 
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
            else{square.squareInfo.isPiecePresent=false}
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

/*
setupDefaultBoardMaxSize(){
    for(let row =0;row<this.maxBoardState.rows;row++){
      this.maxBoardState.tiles.push([])
      for(let col=0;col<this.maxBoardState.cols;col++){
        var tile = {}
        tile.tileType = this.isLight(row,col)? 'l' : 'd';
        if(this.isEmpty){
          tile.isPiecePresent=false
        }
        else if( (col===0||col==7) && (row===0||row===7)) {
            tile.pieceType='r';
            tile.isPiecePresent=true
        }
        else if((col===1||col==6) && (row===0||row===7)){
            tile.pieceType='n';
            tile.isPiecePresent=true
        }
        else if((col===2||col==5) && (row===0||row===7)){
            tile.pieceType='b';
            tile.isPiecePresent=true
        }
        else if((col===3) && (row===0||row===7)){
            tile.pieceType='q';
            tile.isPiecePresent=true
        }
        else if((col===4) && (row===0||row===7)){
            tile.pieceType='k';
            tile.isPiecePresent=true
        }
        else if((row===1||row===6) && col<8){
            tile.pieceType='p';
            tile.isPiecePresent=true
        }
        else{tile.isPiecePresent=false}
        if(row==0||row==1){tile.pieceColor='black'}
        else if(row==6||row==7){tile.pieceColor='white'}
        this.maxBoardState.tiles[row].push(tile)
      }
    }
    this.boardState.rows = this.rows
    this.boardState.cols = this.cols
    this.formatBoardState(this.maxBoardState)
  }
},
*/