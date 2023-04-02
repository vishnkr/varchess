
import { Dictionary } from "vue-router/types/router";

export type PieceColor = "black" | "white"

export interface Tile{
    tileId?: number,
    x?: number,
    y?: number,
    row?: number,
    col?: number,
    disabled: boolean,
    tileType?: string,
    pieceColor?: PieceColor,
    isPiecePresent: boolean,
    pieceType? : string,
}

export interface BoardState{
    tiles : Tile[][],
    castlingAvailability: string,
    enPassant: string,
    turn: string
}

export interface ChatMessage{
    username: string,
    message: string
}

export interface GameInfo{
    result?: string,
    players: Players
    members: string[],
}

export interface Players{
    p1: string,
    p2?: string
}

export interface WsMessage{
    type:string,
    data:string,
}

export interface PiecePosition{
    piece: string,
    row: number,
    col: number
}

export interface MoveInfo{
    piece: string,
    srcRow: number,
    srcCol: number,
    destRow: number,
    destCol: number,
    type?: string,
    promote?: string,
    castle?: boolean,
    isValid?: boolean
    check?: boolean,
    result?: string,
    roomId?:string
    color?:string,
}

export interface MoveInfoPayload extends MoveInfo{
    roomId: string;
    color: string;
}

export interface MovePattern{
    piece: string,
}
export type MovePatterns = MovePattern[] | null;

export interface EditorRouteParams extends Dictionary<string>{
    username: string;
    roomId: string;
}

export interface ServerStatus{
    isOnline: boolean | null,
    errorMessage: string | null
}

export interface PossibleSquaresResponse {
    moves: number[];
  }
  
  export interface RoomState{
    fen:string,
    movePatterns:MovePattern[],
    roomId:string,
    p1:string | undefined,
    p2: string | undefined,
    members : string[],
    turn: string,
  }
  
  export interface CreateRoomResponse{
    roomId:string
  }