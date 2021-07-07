/*
* This file contains helper functions to convert a board state to its equivalent serialized version aka. 
* FEN (Forsyth–Edwards Notation) and vice-versa.
*/

function convertBoardStateToFEN(boardState,turn,castlingAvailability,enPassant){
    var cell,row,empty,fen = '';
    for (row of boardState.tiles){
        empty = 0;
        for (cell of row){
            if (cell.isPiecePresent){
                if(empty>0){
                    fen+=empty;
                    empty=0;
                }
                if(cell.pieceColor=='black'){
                    fen+= cell.pieceType!="Knight" || cell.pieceType!="knight"  ? cell.pieceType[0].toLowerCase() : 'n';
                }
                else{
                    fen+=  cell.pieceType!="Knight" || cell.pieceType!="knight" ? cell.pieceType[0].toUpperCase() : 'N';
                }
            }
            else{
                empty+=1;
            }
        }
        if(empty>0){
            fen+=empty;
            empty=0;
        }
        fen+='/';
    }
    fen = fen.substring(0, fen.length - 1);
    fen+= ` ${turn} ${castlingAvailability} ${enPassant} 0 1`;
    return fen
}

function convertFENtoBoardState(fen){
    console.log(fen)
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

export {convertBoardStateToFEN,convertFENtoBoardState};