import { writable, type Writable } from 'svelte/store';
import { Color, type EditorSettings, type PiecePresentInfo } from './types';

export const editorSettings: Writable<EditorSettings> = writable({
	disableSelected: false,
	pieceSelection: { pieceType: 'p', color: Color.WHITE }
});

function createEditorMaxBoard() {
	const { subscribe, set, update } = writable<PiecePresentInfo[][]>([[]]);

	return {
		subscribe,
		set,
		updatePieceInfo: (rowIndex: number, colIndex: number, newValue: PiecePresentInfo) => {
			update((maxBoard: PiecePresentInfo[][]) => {
				const updatedMaxBoard = [...maxBoard];
				updatedMaxBoard[rowIndex][colIndex] = newValue;
				return updatedMaxBoard;
			});
		}
	};
}

export const editorMaxBoard = createEditorMaxBoard();
