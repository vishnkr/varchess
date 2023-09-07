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
	Editor,
	MovePatternEditor,
	MovePatternView,
	View,
	GameBoard
}

export type SquareNotation = `${File}${Rank}`;
export type SquareMaps = {
	coordToIdMap: CoordinatetoIDMap;
	squares: Record<SquareIdx, SquareInfo>;
};

export const SQUARE_COLORS = ['light', 'dark'] as const;
export type SquareColor = typeof SQUARE_COLORS[number];
export function getSquareColor(row: number, col: number, isFlipped?: boolean): SquareColor {
	const idx = row + col;
	return idx % 2 === 1 ? 'dark' : 'light';
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
export type DisabledSquares = Record<SquareIdx, boolean>;
export interface Position {
	piecePositions: PiecePositions;
	disabled: DisabledSquares;
}

export interface PiecePresentInfo {
	isPiecePresent: boolean;
	piece?: IPiece | null;
	disabled?: boolean;
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
	disableSelected: boolean;
	pieceSelection: IPiece | null;
	editorSubTypeSelected: EditorSubType;
}
/*

16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31
00 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15 



*/
