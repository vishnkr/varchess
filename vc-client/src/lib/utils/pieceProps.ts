import type { MovePattern, Offset } from '$lib/types';

/** Wire shape expected by vc-server chess.PieceProps. */
export type WirePieceProps = {
	slideOffsets: Offset[];
	jumpOffsets: Offset[];
};

/**
 * Editor stores offsets as [Δrow, Δcol] (= [Δrank, Δfile] on the 9×9 pattern board).
 * Engine MoveOffset uses { x: Δfile, y: Δrank }.
 */
function toWireOffset(rowCol: number[]): Offset {
	return { x: rowCol[1] ?? 0, y: rowCol[0] ?? 0 };
}

function fromWireOffset(o: Offset): number[] {
	return [o.y, o.x];
}

/** Convert editor movePatterns (keyed by notation) to API pieceProps. */
export function toWirePieceProps(
	patterns: Record<string, MovePattern>
): Record<string, WirePieceProps> {
	const out: Record<string, WirePieceProps> = {};
	for (const [key, pat] of Object.entries(patterns ?? {})) {
		if (!key || key.length !== 1) continue;
		const notation = key.toLowerCase();
		const slideOffsets = (pat.slideDirections ?? []).map(toWireOffset);
		const jumpOffsets = (pat.jumpOffsets ?? []).map(toWireOffset);
		if (slideOffsets.length === 0 && jumpOffsets.length === 0) continue;
		out[notation] = { slideOffsets, jumpOffsets };
	}
	return out;
}

/** Restore editor movePatterns from API/DB pieceProps. */
export function fromWirePieceProps(
	props: Record<string, WirePieceProps | { slideOffsets?: Offset[]; jumpOffsets?: Offset[]; jumpProps?: Offset[] }> | null | undefined
): Record<string, MovePattern> {
	const out: Record<string, MovePattern> = {};
	if (!props || typeof props !== 'object') return out;
	for (const [key, raw] of Object.entries(props)) {
		if (!key || key.length !== 1 || !raw || typeof raw !== 'object') continue;
		const notation = key.toLowerCase();
		const slides = Array.isArray(raw.slideOffsets) ? raw.slideOffsets : [];
		const jumps = Array.isArray(raw.jumpOffsets)
			? raw.jumpOffsets
			: Array.isArray((raw as { jumpProps?: Offset[] }).jumpProps)
				? (raw as { jumpProps: Offset[] }).jumpProps
				: [];
		out[notation] = {
			slideDirections: slides.filter((o) => o && typeof o === 'object').map(fromWireOffset),
			jumpOffsets: jumps.filter((o) => o && typeof o === 'object').map(fromWireOffset)
		};
	}
	return out;
}
