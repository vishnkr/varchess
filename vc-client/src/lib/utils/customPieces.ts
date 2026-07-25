import type { MovePattern } from '$lib/types';
import { MoveType } from '$lib/types';

export const CUSTOM_PIECE_NAMES: Record<string, string> = {
	d: 'Dolphin',
	i: 'Ninja',
	u: 'Unicorn',
	g: 'Giraffe',
	j: 'Juicer',
	a: 'Astronaut',
	v: 'Phage',
	z: 'Zebra'
};

export const STANDARD_NOTATIONS = new Set(['p', 'k', 'q', 'b', 'n', 'r']);

export function customPieceName(notation: string): string {
	const n = notation.toLowerCase();
	return CUSTOM_PIECE_NAMES[n] ?? n.toUpperCase();
}

export function isCustomNotation(notation: string): boolean {
	const n = notation.toLowerCase();
	return n.length === 1 && !STANDARD_NOTATIONS.has(n);
}

const CENTER = 4;

function inBounds(row: number, col: number) {
	return row >= 0 && row < 9 && col >= 0 && col < 9;
}

/** Build square → Jump/Slide map for a 9×9 pattern board centered at (4,4). */
export function buildMpSquares(
	pattern: MovePattern | null | undefined
): Record<number, MoveType> {
	const squares: Record<number, MoveType> = {};
	if (!pattern) return squares;

	pattern.jumpOffsets?.forEach((offset) => {
		const row = CENTER + offset[0];
		const col = CENTER + offset[1];
		if (!inBounds(row, col)) return;
		squares[row * 9 + col] = MoveType.Jump;
	});

	pattern.slideDirections?.forEach((offset) => {
		let x = offset[0];
		let y = offset[1];
		while (inBounds(CENTER + x, CENTER + y)) {
			squares[(CENTER + x) * 9 + (CENTER + y)] = MoveType.Slide;
			x += offset[0];
			y += offset[1];
		}
	});
	return squares;
}
