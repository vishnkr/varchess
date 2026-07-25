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

/** Matches vc-server/chess classicMoveType iota values. */
export enum ClassicMoveType {
	Null = 0,
	Castle = 1,
	Capture = 2,
	Quiet = 3,
	DoublePawnPush = 4,
	EnPassant = 5,
	Promotion = 6
}

export interface Move {
	from: number;
	to: number;
	piece: string;
	capture?: boolean;
	promotion?: string;
	classicMoveType?: ClassicMoveType;
	variantMoveType?: string;
	additionalData?: Record<string, unknown>;
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
	Checkmate = 'checkmate',
	Antichess = 'antichess',
	NCheck = 'ncheck',
	ArcherChess = 'archerchess',
	Wormhole = 'wormhole',
	PoisonedPawn = 'poisonedpawn',
	HiddenQueen = 'hiddenqueen',
	FootballChess = 'footballchess',
	SpyChess = 'spychess'
}

export enum MoveType {
	Jump,
	Slide
}
export type RuleEditorState = {
	variantType: VariantType;
	isViewVariantRulesOn: boolean;
	ruleComponent?: typeof SvelteComponent<any> | null;
	/** Variant-specific config (e.g. NCheck targetChecks). */
	customData: Record<string, unknown>;
};

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

/** Wire event types — must match vc-ws/vc-core EventType short strings. */
export type EventType =
	| 'start'
	| 'move'
	| 'join'
	| 'resign'
	| 'draw-offer'
	| 'draw-accept'
	| 'draw-reject'
	| 'over'
	| 'presence'
	| 'chat'
	| 'hints';

export const EventStartGame: EventType = 'start';
export const EventGameMakeMove: EventType = 'move';
export const EventJoinGame: EventType = 'join';
export const EventGameResign: EventType = 'resign';
export const EventGameDrawOffer: EventType = 'draw-offer';
export const EventGameDrawAccept: EventType = 'draw-accept';
export const EventGameDrawReject: EventType = 'draw-reject';
export const EventGameOver: EventType = 'over';
export const EventPresence: EventType = 'presence';
export const EventChat: EventType = 'chat';
export const EventHints: EventType = 'hints';

export interface HintsPayload {
	from: number;
}

export interface HintsResultPayload {
	from: number;
	to: number[];
}

export interface PresencePayload {
	userId: string;
	online: boolean;
}

export interface ChatPayload {
	gid?: string;
	sender: string;
	msg: string;
}

export interface WSMessage<T = unknown> {
	t: EventType;
	p: T;
}

export interface ActiveGamePlayer {
	userId: string;
	name: string;
}

export interface ActiveGamePayload {
	id: string;
	players: Record<string, ActiveGamePlayer>;
	gameConfig: {
		variantType?: string;
		fen?: string;
		dimensions?: { ranks: number; files: number };
		pieceProps?: Record<string, PieceProps>;
		customData?: Record<string, unknown>;
	};
	moves?: string[];
	state: string;
	turn: string;
	variantState?: Record<string, unknown> | null;
}

export interface MoveWirePayload {
	m: Move;
	/** True when the side to move after this ply is in check. */
	check?: boolean;
	/** Live variant state after the move (e.g. wormhole cooldowns). */
	variantState?: Record<string, unknown> | null;
}

export interface GameOverPayload {
	winner?: string;
	reason?: string;
}

export interface GameConfigSnapshot {
	variantType?: string;
	name?: string;
	fen?: string;
	dimensions?: { ranks: number; files: number };
	pieceProps?: Record<string, PieceProps>;
	customData?: Record<string, unknown>;
}

export interface Game {
	id?: string;
	_id?: string;
	shortId?: string;
	players?: Record<string, string>;
	playerNames?: Record<string, string>;
	templateId?: string;
	config?: GameConfigSnapshot;
	moves?: string[];
	result?: { winner: string; reason: string };
	createdAt?: string;
	created_at?: string;
	endedAt?: string;
}
export type Offset = { x: number; y: number };

/** Wire piece props — matches vc-server chess.PieceProps JSON. */
export type PieceProps = {
	slideOffsets: Offset[];
	jumpOffsets: Offset[];
};

export type Template = {
	name?: string;
	variantType: VariantType;
	dimensions: {
		ranks: number;
		files: number;
	};
	fen: string;
	pieceProps: Record<string, PieceProps>;
	customData?: Record<string, unknown>;
	/** Nested form returned by some API paths; prefer flat fields above. */
	position?: {
		dimensions?: { ranks: number; files: number };
		fen?: string;
		pieceProps?: Record<string, PieceProps>;
	};
};