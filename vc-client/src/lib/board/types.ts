import type { EditorSubType } from "$lib/components/types";

export enum Color {
	BLACK = 'black',
	WHITE = 'white'
}

export interface IPiece {
	pieceType: string;
	color: Color;
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

export enum BoardType{
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

export type SquareNotation = `${File}${Rank}`;
export type SquareMaps = {
	coordToIdMap: CoordinatetoIDMap;
	squares: Record<SquareIdx, SquareInfo>;
};

export const SQUARE_COLORS = ['light', 'dark'] as const;
export type SquareColor = typeof SQUARE_COLORS[number];
export function getSquareColor(row: number, col: number, isFlipped?: boolean,overrideColorType?:string): SquareColor {
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
	//editable - allowing piece additions/deletions to board squares
	editable: boolean;
	//interactive - allowing moves to be made through clicks/drags
	interactive: boolean;
	isFlipped?: boolean;
	boardType: BoardType;
}
export enum GameType {
	CUSTOM,
	PREDEFINED
}
export interface EditorSettings {
	pieceSelection: IPiece | null;
	editorSubTypeSelected: EditorSubType;
}
/*

16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31
00 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15 



*/
