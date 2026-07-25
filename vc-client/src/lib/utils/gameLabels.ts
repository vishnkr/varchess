import type { Game } from '$lib/types';

const VARIANT_LABELS: Record<string, string> = {
	checkmate: 'Checkmate',
	antichess: 'Antichess',
	ncheck: 'N-Check',
	archerchess: 'Archer Chess',
	wormhole: 'Wormhole',
	poisonedpawn: 'Poisoned Pawn',
	hiddenqueen: 'Hidden Queen',
	footballchess: 'Football Chess',
	spychess: 'Spy Chess'
};

export function variantLabel(variantType?: string): string {
	if (!variantType) return 'Custom';
	return VARIANT_LABELS[variantType.toLowerCase()] ?? variantType;
}

export function gameIdOf(game: Game): string {
	const raw = game._id ?? game.id ?? '';
	if (typeof raw === 'string') return raw;
	if (raw && typeof raw === 'object' && '$oid' in (raw as object)) {
		return String((raw as { $oid: string }).$oid);
	}
	return String(raw);
}

export function boardSizeLabel(game: Game): string | null {
	const d = game.config?.dimensions;
	if (!d?.ranks || !d?.files) return null;
	return `${d.files}×${d.ranks}`;
}

/** Template name if present, else "Variant · WxH". */
export function gameTitle(game: Game): string {
	const name = game.config?.name?.trim();
	if (name) return name;
	const variant = variantLabel(game.config?.variantType);
	const size = boardSizeLabel(game);
	if (size) return `${variant} · ${size}`;
	return variant !== 'Custom' ? variant : 'Past game';
}

export function formatGameResult(game: Game): string {
	const winner = (game.result?.winner || '').toLowerCase();
	const reason = game.result?.reason?.trim();
	const names = game.playerNames ?? {};

	if (winner === 'draw') {
		return reason ? `Draw — ${reason}` : 'Draw';
	}
	if (winner === 'white' || winner === 'w') {
		const who = names.w || names.white || 'White';
		return reason ? `${who} won — ${reason}` : `${who} won`;
	}
	if (winner === 'black' || winner === 'b') {
		const who = names.b || names.black || 'Black';
		return reason ? `${who} won — ${reason}` : `${who} won`;
	}
	if (reason && reason.toLowerCase() !== 'unknown' && reason.toLowerCase() !== 'game over') {
		return reason;
	}
	return 'Finished';
}

function playerId(players: Record<string, string> | undefined, ...keys: string[]): string {
	if (!players) return '';
	for (const k of keys) {
		const v = players[k];
		if (v != null && v !== '') return String(v);
	}
	return '';
}

/** Outcome of a finished game from the given user's perspective. */
export function myGameOutcome(
	game: Game,
	userId: string | null | undefined
): 'win' | 'loss' | 'draw' | 'unknown' {
	if (!userId) return 'unknown';
	const winner = (game.result?.winner || '').toLowerCase();
	if (!winner) return 'unknown';
	if (winner === 'draw') return 'draw';

	const uid = String(userId);
	const whiteId = playerId(game.players, 'w', 'white');
	const blackId = playerId(game.players, 'b', 'black');
	const isWhite = whiteId === uid;
	const isBlack = blackId === uid;
	if (!isWhite && !isBlack) return 'unknown';

	if (winner === 'white' || winner === 'w') return isWhite ? 'win' : 'loss';
	if (winner === 'black' || winner === 'b') return isBlack ? 'win' : 'loss';
	return 'unknown';
}

export type DashboardStats = {
	total: number;
	byVariant: { name: string; value: number }[];
	byResult: { name: string; value: number; color: string }[];
	recent: Game[];
};

export function buildDashboardStats(
	games: Game[],
	userId: string | null | undefined
): DashboardStats {
	const variantCounts = new Map<string, number>();
	let wins = 0;
	let losses = 0;
	let draws = 0;

	for (const g of games) {
		const label = variantLabel(g.config?.variantType);
		variantCounts.set(label, (variantCounts.get(label) ?? 0) + 1);
		const outcome = myGameOutcome(g, userId);
		if (outcome === 'win') wins += 1;
		else if (outcome === 'loss') losses += 1;
		else if (outcome === 'draw') draws += 1;
	}

	const byVariant = [...variantCounts.entries()]
		.map(([name, value]) => ({ name, value }))
		.sort((a, b) => b.value - a.value);

	const byResult = [
		{ name: 'Win', value: wins, color: '#22c55e' },
		{ name: 'Draw', value: draws, color: '#94a3b8' },
		{ name: 'Loss', value: losses, color: '#ef4444' }
	];

	return {
		total: games.length,
		byVariant,
		byResult,
		recent: games.slice(0, 5)
	};
}
