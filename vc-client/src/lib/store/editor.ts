import { EditorSubType } from '$lib/components/types';
import type { BoardEditorState, PieceEditorState, RuleEditorState } from '$lib/types';
import { Color, VariantType } from '$lib/types';
import { writable } from 'svelte/store';

function offsetEq(a: number[], b: number[]) {
	return a[0] === b[0] && a[1] === b[1];
}

function hasOffset(list: number[][] | undefined, offset: number[]) {
	return !!list?.some((o) => offsetEq(o, offset));
}

function newPieceEditorStore() {
	const { subscribe, set, update } = writable<PieceEditorState>({
		movePatterns: {},
		pieceSelection: {
			piece: { pieceType: 'pawn', notation: 'p', color: Color.WHITE },
			group: 'standard'
		}
	});

	/** `piece` is the lowercase notation letter (e.g. "d"). */
	const addJumpPattern = (piece: string, offset: number[]) => {
		const key = piece.toLowerCase();
		update((editorState) => {
			const newMovePatterns = { ...editorState.movePatterns };
			const existing = newMovePatterns[key] ?? { slideDirections: [], jumpOffsets: [] };
			if (hasOffset(existing.jumpOffsets, offset)) {
				return editorState;
			}
			newMovePatterns[key] = {
				slideDirections: [...(existing.slideDirections ?? [])],
				jumpOffsets: [...(existing.jumpOffsets ?? []), offset]
			};
			return { ...editorState, movePatterns: newMovePatterns };
		});
	};

	const removeJumpPattern = (piece: string, offset: number[]) => {
		const key = piece.toLowerCase();
		update((editorState) => {
			const newMovePatterns = { ...editorState.movePatterns };
			const existing = newMovePatterns[key];
			if (!existing) return editorState;
			const jumpOffsets = (existing.jumpOffsets ?? []).filter((o) => !offsetEq(o, offset));
			if (jumpOffsets.length === 0 && (existing.slideDirections ?? []).length === 0) {
				delete newMovePatterns[key];
			} else {
				newMovePatterns[key] = { ...existing, jumpOffsets };
			}
			return { ...editorState, movePatterns: newMovePatterns };
		});
	};

	const addSlidePattern = (piece: string, offset: number[]) => {
		const key = piece.toLowerCase();
		update((editorState) => {
			const newMovePatterns = { ...editorState.movePatterns };
			const existing = newMovePatterns[key] ?? { slideDirections: [], jumpOffsets: [] };
			if (hasOffset(existing.slideDirections, offset)) {
				return editorState;
			}
			newMovePatterns[key] = {
				slideDirections: [...(existing.slideDirections ?? []), offset],
				jumpOffsets: [...(existing.jumpOffsets ?? [])]
			};
			return { ...editorState, movePatterns: newMovePatterns };
		});
	};

	const removeSlidePattern = (piece: string, offset: number[]) => {
		const key = piece.toLowerCase();
		update((editorState) => {
			const newMovePatterns = { ...editorState.movePatterns };
			const existing = newMovePatterns[key];
			if (!existing) return editorState;
			const slideDirections = (existing.slideDirections ?? []).filter((o) => !offsetEq(o, offset));
			if (slideDirections.length === 0 && (existing.jumpOffsets ?? []).length === 0) {
				delete newMovePatterns[key];
			} else {
				newMovePatterns[key] = { ...existing, slideDirections };
			}
			return { ...editorState, movePatterns: newMovePatterns };
		});
	};

	const deletePiecePattern = (piece: string) => {
		const key = piece.toLowerCase();
		update((editorState) => {
			const newMovePatterns = { ...editorState.movePatterns };
			delete newMovePatterns[key];
			return { ...editorState, movePatterns: newMovePatterns };
		});
	};

	const setMovePatterns = (patterns: PieceEditorState['movePatterns']) => {
		update((editorState) => ({ ...editorState, movePatterns: patterns }));
	};

	const updateColor = (color: Color) => {
		update((editorState) => ({
			...editorState,
			pieceSelection: {
				...editorState.pieceSelection,
				piece: { ...editorState.pieceSelection.piece, color }
			}
		}));
	};

	return {
		subscribe,
		set,
		update,
		addJumpPattern,
		removeJumpPattern,
		addSlidePattern,
		removeSlidePattern,
		deletePiecePattern,
		setMovePatterns,
		updateColor
	};
}

function newBoardEditorStore() {
	const { set, update, subscribe } = writable<BoardEditorState>({
		ranks: 8,
		files: 8,
		theme: 'standard',
		isWallSelectorOn: false
	});

	return {
		set,
		update,
		subscribe,
		setDimensions: (ranks: number, files: number) =>
			update((state) => ({ ...state, ranks, files }))
	};
}

function defaultCustomData(
	variantType: VariantType,
	previous: Record<string, unknown> | undefined
): Record<string, unknown> {
	if (variantType === VariantType.NCheck) {
		return { targetChecks: (previous?.targetChecks as number) ?? 3 };
	}
	if (variantType === VariantType.Wormhole) {
		return {
			wormholePairs: Array.isArray(previous?.wormholePairs) ? previous.wormholePairs : [],
			cooldownTurns: 1
		};
	}
	return {};
}

function newRuleStore() {
	const { set, update, subscribe } = writable<RuleEditorState>({
		variantType: VariantType.Checkmate,
		isViewVariantRulesOn: false,
		customData: {}
	});
	function updateVariantType(variantType: VariantType) {
		update((current) => ({
			...current,
			variantType,
			customData: defaultCustomData(variantType, current.customData)
		}));
	}
	function updateCustomData(patch: Record<string, unknown>) {
		update((current) => ({
			...current,
			customData: { ...current.customData, ...patch }
		}));
	}
	return { set, subscribe, update, updateVariantType, updateCustomData };
}

/** Board click pairing for wormhole portals (editor only). */
function newWormholeEditorStore() {
	const { subscribe, set, update } = writable<{ active: boolean; pending: number | null }>({
		active: false,
		pending: null
	});
	return {
		subscribe,
		start: () => set({ active: true, pending: null }),
		stop: () => set({ active: false, pending: null }),
		setPending: (pending: number | null) => update((s) => ({ ...s, pending })),
		resetPending: () => update((s) => ({ ...s, pending: null }))
	};
}

export const editorSubTypeSelected = writable<EditorSubType>(EditorSubType.Board);
/** When true (default in move-pattern mode), board clicks toggle jump offsets. */
export const jumpPatternEditing = writable(true);
/** Request to open a custom piece's move-pattern dialog (notation letter). */
export const openCustomPiecePattern = writable<string | null>(null);
export const pieceEditor = newPieceEditorStore();
export const boardEditor = newBoardEditorStore();
export const ruleEditor = newRuleStore();
export const wormholeEditor = newWormholeEditorStore();

export function resetEditorStores() {
	pieceEditor.set({
		movePatterns: {},
		pieceSelection: {
			piece: {
				pieceType: 'pawn',
				notation: 'p',
				color: Color.WHITE
			},
			group: 'standard'
		}
	});

	boardEditor.set({
		ranks: 8,
		files: 8,
		theme: 'standard',
		isWallSelectorOn: false
	});

	ruleEditor.set({
		variantType: VariantType.Checkmate,
		isViewVariantRulesOn: false,
		customData: {}
	});
	wormholeEditor.stop();
}
