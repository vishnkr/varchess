export type PieceColor = "black" | "white"
export type GameRole = "p1" | "p2" | "member"
export type SquareColor = 'dark' | 'light' | 'disabled' | 'jump' | 'slide' | 'to' | 'from'

export interface SquareInfo{
    isPiecePresent: Boolean,
    pieceColor?:PieceColor,
    pieceType?:string,
    row:number,
    col:number
    squareColor: SquareColor
    tempSquareColor?: SquareColor | null
}

export class Square{
    squareId?: number;
    x?: number;
    y?: number;
    disabled: boolean;
    squareInfo: SquareInfo

    constructor({ disabled, squareInfo,squareId,x,y }:{squareInfo: SquareInfo,disabled?:boolean, squareId?: number, x?: number, y?: number}) {
        this.squareId = squareId;
        this.disabled = disabled ?? false;
        this.x = x;
        this.y = y;
        this.squareInfo = squareInfo;
    }

    updateSquareInfo(newInfo: Partial<SquareInfo>) {
        Object.assign(this.squareInfo, newInfo);
      }
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

export interface UserInfo{
    username?:String,
    isAuthenticated?:boolean,
    curGameRole?: GameRole
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
  
  export interface SquareClick{
    clickType:string,
    row: number,
    col: number,
    piece?:string
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

export interface ChatMessage{
    id: number,
    username: string,
    message: string
}