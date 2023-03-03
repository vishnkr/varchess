import { ClientInfo, ChatMessage, BoardState, GameInfo, MoveInfo, PiecePosition, MovePatterns } from "@/types"


export interface RootState{
    boards : Record<string,BoardState>,
    chatMessages : Record<string,ChatMessage[]>,
    gameInfo : GameInfo | null,
    curStartPos : PiecePosition | null,
    curDestPos : PiecePosition | null,
    currentMove: MoveInfo | null,
    clientInfo: ClientInfo | null,
    turn: 'w' | 'b',
    errorMessage: string | null,
    movePatterns: MovePatterns

}

const state: RootState = {
    boards: {},
    chatMessages: {},
    gameInfo: null,
    curStartPos: null, 
    curDestPos: null,
    currentMove: null,
    clientInfo: null,
    turn: 'w',
    errorMessage:null,
    movePatterns:null,
}

export default state;