import type { Color } from "$lib/board/types";

export interface BoardEditorState {
	ranks: number;
	files: number;
	theme: string;
}

export enum VariantType{
	Standard = "Standard",
	DuckChess = "Duck",
	ArcherChess = "Archer",
	Wormhole = "Wormhole"
}


export enum Objective{
	Checkmate,
	Antichess,
	NCheck,
}

export enum MoveType{
	Jump,
	Slide
}
export interface RuleEditorState{
	variantType: VariantType,
	objective: Objective
}

interface MovePattern{
	slideOffsets: number[][],
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