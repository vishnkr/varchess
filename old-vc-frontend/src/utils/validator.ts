/*
* This file contains functions that help with the validation of the custom board state created by the user
* in editor mode. This is done in the client side to make sure that the board state sent to th server is
* always a valid start position.
*/

import { BoardState } from "@/types";
import { convertFENtoBoardState } from "./fen";

export function validateStartSetup(fen:string){
    // does each side have exactly 1 king?
    //no checks/checkmate in start pos
    const boardState = convertFENtoBoardState(fen)
    var result = countKings(boardState)
    if(result.isValidKings){
        //return isPlayerInCheck('b',boardState,result.bpos) && isPlayerInCheck('w',boardState,result.wpos)
        return true;
    }
    return false;
}

/*
function isPlayerInCheck(color,board,kingPos){
    var attackedSquares = {}
    for(var row=0;row<board.tiles.length;row++){
        for (var col=0;col<row.length;col++){
            if (board.tiles[row][col].isPiecePresent){
                if(color==='b' && board.tiles[row][col].pieceType === board.tiles[row][col].pieceType.toLowerCase()){
                    
                    //attackedSquares[row+1] = attackedSquares[row+1]? attackedSquares[row+1].push(col+1) : [col+1]
            }
        }
    }
}
}
*/

interface KingCount {
    bpos: number[],
    wpos: number[],
    isValidKings?: boolean
}

export function countKings(boardState:BoardState){
    let wKing=0, bKing=0;
    let returnObj: KingCount = {bpos:[], wpos:[]};
    let rowpos = 1
    let colpos = 1
    for(let row of boardState.tiles){
        for (let cell of row){
            if (cell.isPiecePresent && cell.pieceType==='k'){
                if(cell.pieceColor === 'black'){ 
                    returnObj.bpos = [rowpos,colpos]
                    bKing+=1; 
                } 
                else {
                    wKing+=1;
                    returnObj.wpos = [rowpos,colpos]
                } 
            }
            colpos+=1
        }
        rowpos+=1
        colpos=1
    }
    returnObj.isValidKings = wKing==1 && bKing==1
    return returnObj;
}
