export function camelToSnake(obj: any): any {
	if (obj !== null && typeof obj === 'object') {
		if (Array.isArray(obj)) {
			return obj.map((item) => camelToSnake(item));
		} else {
			return Object.fromEntries(
				Object.entries(obj).map(([key, value]) => [
					key.replace(/[A-Z]/g, (match) => `_${match.toLowerCase()}`),
					camelToSnake(value)
				])
			);
		}
	} else {
		return obj;
	}
}

/**
 * Backend chess engine indexes squares on a padded "large" board
 * (width 8 if max(ranks,files) ≤ 8, else 16). The FE uses dense
 * rank*files+file indices. Convert before sending / after receiving moves.
 */
export function smallToLargeIndex(sq: number, files: number, ranks: number): number {
	const largestDim = Math.max(files, ranks);
	const largeWidth = largestDim > 8 ? 16 : 8;
	if (files === largeWidth) return sq;
	const rank = Math.floor(sq / files);
	const file = sq % files;
	return rank * largeWidth + file;
}

export function largeToSmallIndex(sq: number, files: number, ranks: number): number {
	const largestDim = Math.max(files, ranks);
	const largeWidth = largestDim > 8 ? 16 : 8;
	if (files === largeWidth) return sq;
	const rank = Math.floor(sq / largeWidth);
	const file = sq % largeWidth;
	return rank * files + file;
}

/** Dense FE square index → algebraic (a1…), matching engine convention (rank 0 = top). */
export function smallIndexToAlgebraic(sq: number, files: number, ranks: number): string {
	const rank = Math.floor(sq / files);
	const file = sq % files;
	if (file < 0 || file >= files || rank < 0 || rank >= ranks) return '??';
	const chessRank = ranks - rank;
	return `${String.fromCharCode(97 + file)}${chessRank}`;
}

/**
 * Human move notation: e4, exd5, Nf3, Nxe5, e8=Q, O-O.
 * Uses FE dense square indices.
 */
export function formatMoveNotation(
	move: { from: number; to: number; piece: string; capture?: boolean; promotion?: string },
	files: number,
	ranks: number
): string {
	const from = smallIndexToAlgebraic(move.from, files, ranks);
	const to = smallIndexToAlgebraic(move.to, files, ranks);
	const pieceChar = (move.piece || '').toString();
	const lower = pieceChar.toLowerCase();
	const isPawn = lower === 'p' || pieceChar === '';

	if (lower === 'k') {
		const fileDiff = (move.to % files) - (move.from % files);
		if (Math.abs(fileDiff) === 2) {
			return fileDiff > 0 ? 'O-O' : 'O-O-O';
		}
	}

	const promo = move.promotion ? `=${move.promotion.toUpperCase()}` : '';
	if (isPawn) {
		if (move.capture) return `${from[0]}x${to}${promo}`;
		return `${to}${promo}`;
	}
	const letter = pieceChar.toUpperCase();
	return move.capture ? `${letter}x${to}${promo}` : `${letter}${to}${promo}`;
}

export const COLOR_THEMES: Record<string, { lightColor: string; darkColor: string }> = {
	Default: { lightColor: 'hsl(51deg 24% 84%)', darkColor: 'hsl(145deg 32% 44%)' },
	Brown: { lightColor: 'hsl(36, 81%, 84%)', darkColor: 'hsl(25, 31%, 51%)' },
	Aqua: { lightColor: 'hsl(197, 34%, 83%)', darkColor: 'hsl(217, 68%, 52%)' },
	Classic: { lightColor: 'hsl(0, 0%, 100%)', darkColor: 'hsl(0, 0%, 45%)' },
	Candy: { lightColor: 'hsl(314, 100%, 90%)', darkColor: 'hsl(328, 100%, 55%)' }
};