import {ChatMessage, BoardState, GameInfo, MoveInfo, PiecePosition, MovePatterns, ServerStatus } from "@/types"


export interface RootState{
    board : BoardState | null,
    chatMessages : Record<string,ChatMessage[]>,
    gameInfo : GameInfo | null,
    curStartPos : PiecePosition | null,
    curDestPos : PiecePosition | null,
    currentMove: MoveInfo | null,
    serverStatus: ServerStatus,
    movePatterns: MovePatterns,

}

const state: RootState = {
    board: null,
    chatMessages: {},
    gameInfo: null,
    curStartPos: null, 
    curDestPos: null,
    currentMove: null,
    serverStatus: {isOnline:null,errorMessage:null},
    movePatterns:null,
}

export default state;