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
    console.log('result fen-',fen)
}

function convertFENtoBoardState(fen){
    var splitFen = fen.split(' ');
    var boardState = {tiles:[], castlingAvailability: splitFen[2], turn: splitFen[1], enPassant: splitFen[3]}
    var rows = splitFen[0].split('/');
    var char;
    for (var i=0; i < rows.length;i++){
        boardState.tiles.push([]);
        for(var j = 0; j < rows[i].length; j++){
            char = rows[i].charAt(j);
            if (/\d/.test(char)){
                for(var empty=0; empty<parseInt(char);empty++){
                    boardState.tiles[i].push({isPiecePresent:false});
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
//convertFENtoBoardState("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
//convertFENtoBoardState("r2q2nr/p4ppp/8/8/4p2P/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
//convertFENtoBoardState("3/r1P/Qk1/1N1 b KQkq - 0 1")

export {convertBoardStateToFEN,convertFENtoBoardState};