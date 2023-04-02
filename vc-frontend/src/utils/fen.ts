/*
* This file contains helper functions to convert a board state to its equivalent serialized version aka. 
* FEN (Forsyth–Edwards Notation) and vice-versa.
*/

import { BoardState } from "@/types";
import { isLight } from ".";
function convertBoardStateToFEN(boardState:BoardState,turn?:string,castlingAvailability?:string,enPassant?:string){
    var cell,row,empty,fen = '';
    for (row of boardState.squares){
        empty = 0;
        for (cell of row){
            if (cell.squareInfo.isPiecePresent || cell.disabled){
                if(empty>0){
                    fen+=empty;
                    empty=0;
                }
                if (cell.disabled){
                    fen+="."
                    continue
                }
                const pieceType = cell.squareInfo.pieceType ?? '';
                if(cell.squareInfo.pieceColor=='black'){
                    fen+= pieceType.toLowerCase() != "knight"  ? pieceType[0].toLowerCase() : 'n';
                }
                else{
                    fen+= pieceType.toLowerCase() != "knight" ? pieceType[0].toUpperCase() : 'N';
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
    fen+= ` ${turn ?? boardState.turn} ${castlingAvailability ?? boardState.castlingAvailability} ${enPassant ?? boardState.enPassant} 0 1`;
    return fen
}

function convertFENtoBoardState(fen:string){
    console.log("converting fen",fen)
    let splitFen = fen.split(' ');
    let boardState:BoardState = {squares:[], castlingAvailability: splitFen[2], turn: splitFen[1], enPassant: splitFen[3], dimensions:{rows:8,cols:8}}
    let rows = splitFen[0].split('/');
    let char;
    let secDigit = 0;
    let colEnd = 0;
    for (var i=0; i < rows.length;i++){
        boardState.squares.push([]);
        secDigit = 0;
        for(var j = 0; j < rows[i].length; j++){
            char = rows[i].charAt(j);
            if (char === "."){
                boardState.squares[i].push({squareInfo:{isPiecePresent:false,row:i,col:j, squareColor: isLight(i,j) ? 'light' : 'dark'},disabled:true});
            }
            else if (/\d/.test(char)){
                if(j+1<rows[i].length && (/\d/.test(rows[i].charAt(j+1)))){
                    secDigit=parseInt(char);
                } else{
                    if(secDigit!=0){
                        colEnd = secDigit*10+parseInt(char)
                    } else {colEnd=parseInt(char)}
                    for(var empty=0; empty<colEnd;empty++){
                        boardState.squares[i].push({squareInfo:{isPiecePresent:false, row:i, col: j, squareColor: isLight(i,j) ? 'light' : 'dark'}, disabled: false});
                    }
                }
            }
            else{
                boardState.squares[i].push(
                    {
                        squareInfo:{
                            isPiecePresent: true,
                            pieceColor: (char == char.toLowerCase() && char != char.toUpperCase())? 'black': 'white',
                            pieceType: char.toLowerCase(),
                            row:i,
                            col:j,
                            squareColor:isLight(i,j) ? 'light' : 'dark'
                        },
                        disabled: false
                    }
                    )
            }
        }
    }
    boardState.dimensions.rows = rows.length;
    boardState.dimensions.cols = boardState.squares[0].length;
    return boardState;
}

export {convertBoardStateToFEN,convertFENtoBoardState};