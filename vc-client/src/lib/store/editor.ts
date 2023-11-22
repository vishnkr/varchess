import type { BoardEditorState, PieceEditorState } from './types';
import { writable, type Writable } from 'svelte/store';
import { EditorSubType } from '$lib/components/types';
import { type EditorSettings, Color } from '$lib/board/types';

function newPieceEditorStore(){
    const { subscribe, set, update } = writable<PieceEditorState>(
      { movePatterns: {}, pieceSelection: { pieceType: 'p', color: Color.WHITE, group: 'standard' }}
      )
    const addJumpPattern= (piece:string,offset:number[])=>{
      update((editorState)=>{
        const newMovePatterns = { ...editorState.movePatterns };
        if (!newMovePatterns[piece]) {
            newMovePatterns[piece] = {slideOffsets:[],jumpOffsets:[]}
            newMovePatterns[piece].jumpOffsets = [offset];
        } else if(newMovePatterns[piece].jumpOffsets && !newMovePatterns[piece].jumpOffsets?.includes(offset)) {
          newMovePatterns[piece].jumpOffsets?.push(offset);
        }
        return { ...editorState, movePatterns: newMovePatterns };
      });
    }
  
    const removeJumpPattern = (piece:string,offset:number[])=>{
      update((editorState)=>{
        const newMovePatterns = { ...editorState.movePatterns };
        if (newMovePatterns[piece]) {
          const jumpOffsets = newMovePatterns[piece].jumpOffsets;
          if (jumpOffsets) {
            const idx = jumpOffsets.findIndex((o)=> o[0]===offset[0] && o[1]===offset[1]);
            if(idx!==-1){
              jumpOffsets.splice(idx,1)
              if (jumpOffsets.length === 0 && newMovePatterns[piece].slideOffsets.length === 0) {
                delete newMovePatterns[piece];
              }
            }
          }
        } 
        return { ...editorState, movePatterns: newMovePatterns };
      });
    }
  
    // addSlidePattern: arguments include piece character ex. 'p' for pawn along with offset for a direction with 2 values ([0,1] for south)
    const addSlidePattern = (piece:string,offset:number[])=>{
        update((editorState)=>{
            const newMovePatterns = { ...editorState.movePatterns };
            if (!newMovePatterns[piece]) {
                newMovePatterns[piece] = {slideOffsets:[],jumpOffsets:[]}
                newMovePatterns[piece].slideOffsets = [offset];
            } else if(newMovePatterns[piece].slideOffsets && !newMovePatterns[piece].slideOffsets?.includes(offset)) {
                newMovePatterns[piece].slideOffsets?.push(offset);
            }
            return {...editorState, movePatterns: newMovePatterns};
        })
    }

    const removeSlidePattern = (piece:string,offset:number[])=>{
        update((editorState)=>{
            const newMovePatterns = { ...editorState.movePatterns };
            if (newMovePatterns[piece]) {
            
            const slideOffsets = newMovePatterns[piece].slideOffsets;
            if (slideOffsets) {
                const idx = slideOffsets.findIndex((o)=> o===offset);
                if(idx!==-1){
                    slideOffsets.splice(idx,1)
                    if (slideOffsets.length === 0 && newMovePatterns[piece].jumpOffsets.length === 0) {
                        delete newMovePatterns[piece];
                      }
                }
            }
            } 
            return { ...editorState, movePatterns: newMovePatterns };
        })
    }

    const deletePiecePattern = (piece:string)=>{
        update((editorState)=>{
            const newMovePatterns = { ...editorState.movePatterns };
            delete newMovePatterns[piece]
            return {...editorState,movePatterns:newMovePatterns}
        })
    }

    return{
      subscribe,
      set,
      update,
      addJumpPattern,
      removeJumpPattern,
      addSlidePattern,
      removeSlidePattern,
      deletePiecePattern,
    }
  }

function newBoardEditorStore(){
  const {set, update,subscribe}= writable<BoardEditorState>({
    ranks:8,
    files:8,
    theme: "standard"
  })

  return {
    set,
    update,
    subscribe
  }
}




export const editorSettings: Writable<EditorSettings> = writable({
	isWallSelectorOn: false,
	pieceSelection: { pieceType: 'p', color: Color.WHITE, group: 'standard' },
	editorSubTypeSelected: EditorSubType.Board,
});


export const pieceEditor = newPieceEditorStore()
export const boardEditor = newBoardEditorStore()