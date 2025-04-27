import { SvelteComponent } from "svelte";

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
	// GameBoard: Playable board
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
  

export const doesSupportDragDrop = (bType: BoardType): boolean => isGameBoard(bType) || isEditor(bType);

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


export interface BoardEditorState {
	ranks: number;
	files: number;
	theme: string;
	isWallSelectorOn: boolean;
}

export enum VariantType {
	Checkmate = 'Checkmate',
	Antichess = 'AntiChess',
	NCheck = 'NCheck',
	DuckChess = 'Duck',
	ArcherChess = 'ArcherChess',
	Wormhole = 'Wormhole',
	GoalChess = 'GoalChess',
}

export enum MoveType {
	Jump,
	Slide
}
export type RuleEditorState = {
	variantType: VariantType;
	isViewVariantRulesOn: boolean;
	ruleComponent?: typeof SvelteComponent<any> | null
}

export interface MovePattern {
	slideDirections: number[][];
	jumpOffsets: number[][];
}

export type PieceSelection = {
	piece: IPiece;
	group: string;
};
export interface PieceEditorState {
	movePatterns: Record<string, MovePattern>;
	pieceSelection: PieceSelection;
}

export interface EditorState {
	boardEditor: BoardEditorState;
	rulesEditor: RuleEditorState;
	pieceEditor: PieceEditorState;
	gameSettings: {
		showPossibleMoves: boolean;
		disableChat: boolean;
	};
}

export type EventType =
	| 'chat.message'
	| 'game.connect_user'
	| 'game.create_game'
	| 'game.join_game'
	| 'game.disconnect_user'
	| 'game.set_players'
	| 'game.result'
	| 'game.make_move'
	| 'game.offer_draw'
	| 'game.draw_result'
	| 'game.resign'
	| 'start'
	| 'Error';

export const EventChatMessage: EventType = 'chat.message';
export const EventUserConnect: EventType = 'game.connect_user';
export const EventCreateGame: EventType = 'game.create_game';
export const EventJoinGame: EventType = 'game.join_game';
export const EventUserDisconnect: EventType = 'game.disconnect_user';
export const EventSetPlayers: EventType = 'game.set_players';
export const EventGameResult: EventType = 'game.result';
export const EventGameMakeMove: EventType = 'game.make_move';
export const EventGameDrawOffer: EventType = 'game.offer_draw';
export const EventGameDrawResult: EventType = 'game.draw_result';
export const EventGameResign: EventType = 'game.resign';
export const EventStartGame: EventType = 'start';
export const EventError: EventType = 'Error';

export interface WebSocketMessage{
	event: EventType
	params: WSParams
}

export interface WSParams{
	gameId: string,
}

export interface ChatParams extends WSParams{
	message: string
}

export interface MoveParams extends WSParams{
	move: Move
}

export interface Game {
    id: string;
    name: string;
    type: string;
    createdAt: string;
}
export type Offset = { x: number; y: number };

export type PieceProps = {
	slideOffsets: Offset[];
	jumpProps: Offset[];
};

export type Template = {
	name?: string,
	variantType: VariantType;
	dimensions: {
		ranks: number;
		files: number;
	};
	fen: string;
	pieceProps: Record<string, MovePattern>;
	customData?: Record<string,any>;
}