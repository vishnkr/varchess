import {ChatMessage, BoardState, GameInfo, MoveInfo, PiecePosition, MovePatterns, ServerStatus, UserInfo } from "@/types"
import { STANDARD_FEN } from "../utils/constants";
import { convertFENtoBoardState } from "../utils/fen";


export interface RootState{
    board : BoardState,
    userInfo:UserInfo,
    chatMessages : Record<string,ChatMessage[]>,
    gameInfo : GameInfo | null,
    curStartPos : PiecePosition | null,
    curDestPos : PiecePosition | null,
    currentMove: MoveInfo | null,
    serverStatus: ServerStatus,
    movePatterns: MovePatterns,

}

const state: RootState = {
    board: convertFENtoBoardState(STANDARD_FEN),
    userInfo:{isAuthenticated:false},
    chatMessages: {},
    gameInfo: null,
    curStartPos: null, 
    curDestPos: null,
    currentMove: null,
    serverStatus: {isOnline:null,errorMessage:null},
    movePatterns:null,
}

export default state;