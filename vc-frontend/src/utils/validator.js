/*
* This file contains functions that help with the validation of the custom board state created by the user
* in editor mode. This is done in the client side to make sure that the board state sent to th server is
* always a valid start position.
*/

function convertFENtoBoardState(fen){
    var splitFen = fen.split(' ');
    var boardState = {tiles:[], castlingAvailability: splitFen[2], turn: splitFen[1], enPassant: splitFen[3]}
    var rows = splitFen[0].split('/');
    var char;
    var secDigit,colEnd = 0;
    for (var i=0; i < rows.length;i++){
        boardState.tiles.push([]);
        secDigit = 0
        for(var j = 0; j < rows[i].length; j++){
            char = rows[i].charAt(j);
            if (/\d/.test(char)){
                if(j+1<rows[i].length && (/\d/.test(rows[i].charAt(j+1)))){
                    secDigit=char
                } else{
                    if(secDigit!=0){
                        colEnd = parseInt(secDigit)*10+parseInt(char)
                    } else {colEnd=parseInt(char)}
                    for(var empty=0; empty<colEnd;empty++){
                        boardState.tiles[i].push({isPiecePresent:false});
                    }
                }
            }
            else{
                boardState.tiles[i].push(
                    {
                        isPiecePresent: true,
                        pieceColor: (char == char.toLowerCase() && char != char.toUpperCase())? 'black': 'white',
                        pieceType: char.toLowerCase()
                    }
                    )
            }
        }
    }
    return boardState;
}

export function validateStartSetup(fen){
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
export function countKings(boardState){
    var wKing=0, bKing=0;
    var returnObj = {};
    var rowpos = 1
    var colpos = 1
    for(var row of boardState.tiles){
        for (var cell of row){
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
