
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
    p1: string,
    p2?: string,
    members?: string[],
    turn: string,
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
}

export interface MoveInfoPayload extends MoveInfo{
    roomId: string;
    color: string;
}

export interface ClientInfo{
    username: string,
    isPlayer: boolean,
    color?: string,
    ws: WebSocket
}

export interface MovePattern{
    piece: string,
}
export type MovePatterns = MovePattern[] | null;

export interface EditorRouteParams extends Dictionary<string>{
    username: string;
    roomId: string;
}
