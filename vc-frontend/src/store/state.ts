import { ClientInfo, ChatMessage, BoardState, GameInfo, MoveInfo, PiecePosition, MovePatterns, ServerStatus } from "@/types"


export interface RootState{
    boards : Record<string,BoardState>,
    chatMessages : Record<string,ChatMessage[]>,
    gameInfo : GameInfo | null,
    curStartPos : PiecePosition | null,
    curDestPos : PiecePosition | null,
    currentMove: MoveInfo | null,
    clientInfo: ClientInfo | null,
    turn: 'w' | 'b',
    serverStatus: ServerStatus,
    movePatterns: MovePatterns,

}
export interface WebSocketState {
    ws: WebSocket | null;
    userId: string | null;
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
    serverStatus: {isOnline:null,errorMessage:null},
    movePatterns:null,
}

export default state;