import { boardEditor, pieceEditor, ruleEditor } from "$lib/store/editor";
import { editorMaxBoard } from "$lib/board/board";
import type { BoardEditorState, PieceEditorState, RuleEditorState } from "$lib/store/types";

export const generateGameConfigJSON = () =>{
    let boardEditorState:any;
    let pieceEditorState:any;
    let ruleEditorState:any;
    boardEditor.subscribe((value) => {
        boardEditorState = value
    })
    ruleEditor.subscribe((value) => {
        ruleEditorState = value
    })
    pieceEditor.subscribe((value) => {
        pieceEditorState = value
    })
    
    //const piece_props = getPieceProps();
    const fen = getFEN();
    const gameConfig = {
        variant_type: ruleEditorState.variantType,
        position: {
            dimensions: {
                ranks : boardEditorState.ranks,
                files : boardEditorState.files
            },
            fen,
            //piece_props,
        },
        objective: ruleEditorState.objective
    }
    return JSON.stringify(gameConfig);
}

/*export const getPieceProps =()=>{
};*/

export const getFEN = ()=>{
    const position = '';
    let boardEditorVal:any;
    let editorMaxBoardState:any;
    boardEditor.subscribe((value) => boardEditorVal = value)
    editorMaxBoard.subscribe((value)=> editorMaxBoardState = editorMaxBoard) 
    
    const turn = 'w';
    const castleRights = "KQkq";
    const ep="-";
    return `${position} ${turn} ${castleRights} ${ep} 0 0`
}