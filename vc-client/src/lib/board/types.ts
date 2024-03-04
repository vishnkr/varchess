export enum Color {
	BLACK = 'black',
	WHITE = 'white'
}

export interface IPiece {
	pieceType: string;
	color: Color;
	notation: string;
}

export type SquareIdx = number;
export type File =
	| 'a'
	| 'b'
	| 'c'
	| 'd'
	| 'e'
	| 'f'
	| 'g'
	| 'h'
	| 'i'
	| 'j'
	| 'k'
	| 'l'
	| 'm'
	| 'n'
	| 'o'
	| 'p';
export type Rank = `${1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16}`;

export enum ClassicMoveType {
	Quiet = "Quiet",
	Capture = "Capture",
	DoublePawnPush = "DoublePawnPush",
	EnPassant = "EnPassant",
	Castle = "Castle",
	Promotion = "Promotion",
  }


export enum VariantMoveType {
	Duck = "Duck"
  }
  
export interface Move {
	src: number,
	dest: number,
	classic_move_type: ClassicMoveType,
	variant_move_type?: VariantMoveType,
	piece: IPiece
}

export enum BoardType {
	// Editor: clicks on squares adds/removes pieces/walls
	Editor,
	// MovePatternEditor: clicks on squares adds jump patterns and used to display selected slide pattern
	MovePatternEditor,
	// MovePatternView: View only board with piece and colored move pattern squares
	MovePatternView,
	// View: View only board to view current game state
	View,
	// GameBoard: Playable board, can trigger websocket messages for move validation
	GameBoard
}
export function isEditor(boardType: BoardType): boolean {
	return boardType === BoardType.Editor;
  }
  
  export function isMovePatternEditor(boardType: BoardType): boolean {
	return boardType === BoardType.MovePatternEditor;
  }
  
  export function isMovePatternView(boardType: BoardType): boolean {
	return boardType === BoardType.MovePatternView;
  }
  
  export function isView(boardType: BoardType): boolean {
	return boardType === BoardType.View;
  }
  
  export function isGameBoard(boardType: BoardType): boolean {
	return boardType === BoardType.GameBoard;
  }
  

export function doesSupportDragDrop(bType: BoardType): boolean {
	return bType === BoardType.GameBoard || bType === BoardType.Editor;
}

export type SquareNotation = `${File}${Rank}`;
export type SquareMaps = {
	coordToIdMap: CoordinatetoIDMap;
	squares: Record<SquareIdx, SquareInfo>;
};

export const SQUARE_COLORS = ['light', 'dark'] as const;
export type SquareColor = (typeof SQUARE_COLORS)[number];
export function getSquareColor(
	row: number,
	col: number,
	isFlipped?: boolean,
	overrideColorType?: string
): SquareColor {
	const idx = row + col;
	return overrideColorType ?? idx % 2 === 1 ? 'dark' : 'light';
}

export interface Dimensions {
	ranks: number;
	files: number;
}

export interface SquareInfo {
	squareIndex: SquareIdx;
	squareNotation: SquareNotation;
	gridX: number;
	gridY: number;
	row: number;
	column: number;
	isMarkedTarget?: boolean;
}
export type Coordinate = [number, number];

export type PiecePositions = Record<SquareIdx, IPiece>;
export type Walls = Record<SquareIdx, boolean>;
export interface Position {
	piecePositions: PiecePositions;
	walls: Walls;
}

export interface PiecePresentInfo {
	isPiecePresent: boolean;
	piece?: IPiece | null;
	wall?: boolean;
}
export type CoordinatetoIDMap = Record<string, SquareIdx>;
export interface BoardConfig {
	dimensions: Dimensions;
	fen: string;
	isFlipped?: boolean;
	boardType: BoardType;
}

/*

16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31
00 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15 



*/
