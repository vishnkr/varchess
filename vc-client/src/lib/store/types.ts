import type { Color } from '$lib/board/types';

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
	Wormhole = 'Wormhole'
}

export enum MoveType {
	Jump,
	Slide
}
export interface RuleEditorState {
	variantType: VariantType;
}

export interface MovePattern {
	slideDirections: number[][];
	jumpOffsets: number[][];
}

export type PieceSelection = {
	pieceType: string;
	color: Color;
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

type EventType =
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
	| 'game.start_game'
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
export const EventStartGame: EventType = 'game.start_game';
export const EventError: EventType = 'Error';
