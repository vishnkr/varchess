interface BoardEditorState {
	ranks: number;
	files: number;
	theme: string;
}

enum VariantType{
	Standard,
	DuckChess,
	ArcherChess,
	Wormhole
}


interface RuleEditorState{
	variantType: VariantType,
	
}

interface MovePattern{
	slideOffsets?: number[][],
	jumpOffsets?: number[][],
}

interface PieceEditorState{
	movePatterns: MovePattern[]
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