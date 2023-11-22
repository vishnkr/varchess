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

export interface PieceEditorState{
	movePatterns: Record<string,MovePattern>
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