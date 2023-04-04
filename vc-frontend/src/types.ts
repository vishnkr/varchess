export type PieceColor = "black" | "white"

export type SquareColor = 'dark' | 'light' | 'disabled' | 'jump' | 'slide' | 'to' | 'from'
export interface SquareInfo{
    isPiecePresent: Boolean,
    pieceColor?:PieceColor,
    pieceType?:String,
    row:number,
    col:number
    squareColor: SquareColor
    tempSquareColor?: SquareColor | null
}

export interface Square{
    squareId?: number,
    x?: number,
    y?: number,
    disabled: boolean,
    squareInfo: SquareInfo
}
export interface Dimensions{
    rows: number,
    cols: number
}
export interface BoardState{
    squares : Square[][],
    dimensions: Dimensions,
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

export type MPTuple = [x: number, y: number];

export interface MovePattern{
    piece: string,
    jumpPatterns: MPTuple[],
    slidePatterns: MPTuple[]
}

export type MovePatterns = MovePattern[] | null;

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
  
  type PiecesInPlay = Record<string, {
    isAddedToBoard: boolean,
    isMPDefined: boolean,
    movePattern?:MovePattern
  }>

  export type EditorModeType = 'MP' | 'Game'

  export interface EditorState{
    curPiece: string,
    curPieceColor: PieceColor,
    isDisableTileOn:boolean,
    piecesInPlay: PiecesInPlay,
    editorType:EditorModeType,
    curCustomPiece: string | null
  }

export type MoveType = 'jump' | 'slide'
 
export interface MPEditorState extends Omit<EditorState, 'curCustomPiece'> {
    moveType: MoveType,
    curCustomPiece: string
}

export function isMPEditor(editorState: MPEditorState | EditorState): editorState is MPEditorState {
    return (editorState as MPEditorState).moveType !== undefined;
  }