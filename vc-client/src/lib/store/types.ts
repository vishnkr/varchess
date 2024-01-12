import type { Color } from "$lib/board/types";

export interface BoardEditorState {
	ranks: number;
	files: number;
	theme: string;
	isWallSelectorOn: boolean;
}

export enum VariantType{
	Checkmate = "Checkmate",
	Antichess = "AntiChess",
	NCheck = "NCheck",
	DuckChess = "Duck",
	ArcherChess = "ArcherChess",
	Wormhole = "Wormhole"
}

export enum MoveType{
	Jump,
	Slide
}
export interface RuleEditorState{
	variantType: VariantType,
}

interface MovePattern{
	slideDirections: number[][],
	jumpOffsets: number[][],
}

export type PieceSelection = {
	pieceType: string,
	color: Color,
	group: string,
}
export interface PieceEditorState{
	movePatterns: Record<string,MovePattern>
	pieceSelection: PieceSelection
}

export interface EditorState{
    boardEditor: BoardEditorState,
    rulesEditor: RuleEditorState,
	pieceEditor: PieceEditorState,
    gameSettings: {
        showPossibleMoves: boolean,
        disableChat: boolean
    }
}