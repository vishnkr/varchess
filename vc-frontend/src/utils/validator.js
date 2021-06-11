
const lib = require("./fen");
console.log(lib)
function validateStartSetup(fen){
    // does each side have exactly 1 king? might change this to allow more customizations later
    //no checks/checkmate in start pos
    //if board row/col =2 then col/row should be higher than 2 (2 kings cant be placed in 2x2 or less)
    const boardState = convertFENtoBoardState(fen)
    console.log(countKings(boardState))
}

function countKings(boardState){
    var wKing=0, bKing=0;
    for(row of boardState.tiles){
        for (cell of row){
            if (cell.isPiecePresent && cell.pieceType==='k'){
                cell.pieceColor === 'black'? bKing+=1 : wKing+=1;
            }
        }
    }
    return wKing==1 && bKing==1;
}
validateStartSetup("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
