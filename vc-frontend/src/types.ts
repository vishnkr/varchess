import { Players } from "./classes"

export type PieceColor = "black" | "white"
export type GameRole = "p1" | "p2" | "member"
export type SquareColor = 'dark' | 'light' | 'disabled' | 'jump' | 'slide' | 'to' | 'from'

export interface ISquareInfo{
    isPiecePresent: Boolean,
    pieceColor?:PieceColor,
    pieceType?:string,
    row:number,
    col:number
    squareColor: SquareColor
    tempSquareColor?: SquareColor | null
}

export interface IDimensions{
    rows: number,
    cols: number
}

export interface IChatMessage{
    username: string,
    message: string
}


export interface IWsMessage{
    type:string,
    data:string,
}

export interface IPiecePosition{
    piece: string,
    row: number,
    col: number
}

export interface IMoveInfo{
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

export interface IMoveInfoPayload extends IMoveInfo{
    roomId: string;
}

export type MPTuple = [x: number, y: number];

export interface IMovePattern{
    piece: string,
    jumpPatterns: MPTuple[],
    slidePatterns: MPTuple[]
}

export type IMovePatterns = Record<string,IMovePattern>


export interface ServerStatus{
    isOnline: boolean | null,
    errorMessage: string | null
}

export interface PossibleSquaresResponse {
    moves: number[];
  }
  
  export interface RoomState{
    fen:string,
    movePatterns:IMovePattern[],
    roomId:string,
    p1:string | undefined,
    p2: string | undefined,
    members : string[],
    turn: string,
  }
  
  export interface ICreateRoomResponse{
    roomId:string
  }
  
  export interface ISquareClick{
    clickType:string,
    row: number,
    col: number,
    piece?:string
  }

  type PiecesInPlay = Record<string, {
    isAddedToBoard: boolean,
    isMPDefined: boolean,
    movePattern?:IMovePattern
  }>

  export type EditorModeType = 'MP' | 'Game'

  export interface IEditorState{
    curPiece: string,
    curPieceColor: PieceColor,
    isDisableTileOn:boolean,
    piecesInPlay: PiecesInPlay,
    editorType:EditorModeType,
    curCustomPiece: string | null
  }

export type MoveType = 'jump' | 'slide'
 
export interface IMPEditorState extends Omit<IEditorState, 'curCustomPiece'> {
    moveType: MoveType,
    curCustomPiece: string
}

export function isMPEditor(editorState: IMPEditorState | IEditorState): editorState is IMPEditorState {
    return (editorState as IMPEditorState).moveType !== undefined;
}

export interface IChatMessage{
    id: number,
    username: string,
    message: string
}

export interface UserInfo {
    username?: string,
    isAuthenticated: boolean,
    curGameRole: GameRole
}


export interface GameInfo {
  players: Players,
  members: Record<string,UserInfo>,
  result?: string
}

export { Players }
