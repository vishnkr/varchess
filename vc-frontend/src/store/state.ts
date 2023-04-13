import {IChatMessage, IMoveInfo, IPiecePosition, IMovePatterns, ServerStatus } from "@/types"
import { STANDARD_FEN } from "../utils/constants";
import { convertFENtoBoardState } from "../utils/fen";
import {User, BoardState, GameInfo} from "../classes";


export interface RootState{
    board : BoardState,
    userInfo: User ,
    chatMessages : Record<string,IChatMessage[]>,
    gameInfo : GameInfo | null,
    curStartPos : IPiecePosition | null,
    curDestPos : IPiecePosition | null,
    currentMove: IMoveInfo | null,
    serverStatus: ServerStatus,
    movePatterns: IMovePatterns | null,
}

const state: RootState = {
    board: convertFENtoBoardState(STANDARD_FEN),
    userInfo: new User(),
    chatMessages: {},
    gameInfo: null,
    curStartPos: null, 
    curDestPos: null,
    currentMove: null,
    serverStatus: {isOnline:null,errorMessage:null},
    movePatterns:null,
}

export default state;