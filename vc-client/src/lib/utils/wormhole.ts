import { VariantType } from '$lib/types';
import { largeToSmallIndex, smallToLargeIndex, smallIndexToAlgebraic } from '$lib/utils/index';

export type WormholePair = [number, number];

export type WormholePlayState = {
	/** Dense FE board indices. */
	pairs: WormholePair[];
	/** Remaining half-moves of cooldown per pair (0 = open). */
	cooldown: number[];
};

export function parseWormholePairs(raw: unknown): WormholePair[] {
	if (!Array.isArray(raw)) return [];
	const out: WormholePair[] = [];
	for (const item of raw) {
		if (!Array.isArray(item) || item.length < 2) continue;
		const a = Number(item[0]);
		const b = Number(item[1]);
		if (!Number.isFinite(a) || !Number.isFinite(b) || a === b) continue;
		out.push([a, b]);
	}
	return out;
}

export function parseWormholeCooldown(raw: unknown, pairCount: number): number[] {
	const cooldown = Array.from({ length: pairCount }, () => 0);
	if (!Array.isArray(raw)) return cooldown;
	for (let i = 0; i < pairCount; i++) {
		const n = Number(raw[i]);
		cooldown[i] = Number.isFinite(n) && n > 0 ? Math.floor(n) : 0;
	}
	return cooldown;
}

/** Build play state from gameConfig.customData + optional live variantState. */
export function wormholePlayStateFromConfig(
	customData: Record<string, unknown> | null | undefined,
	variantState: Record<string, unknown> | null | undefined,
	files: number,
	ranks: number
): WormholePlayState | null {
	const wirePairs = parseWormholePairs(
		variantState?.wormholePairs ?? customData?.wormholePairs
	);
	if (wirePairs.length === 0) return null;
	const pairs = wormholePairsFromWire(wirePairs, files, ranks);
	const cooldown = parseWormholeCooldown(variantState?.pairCooldown, pairs.length);
	return { pairs, cooldown };
}

export function applyWormholeVariantState(
	current: WormholePlayState | null,
	variantState: Record<string, unknown> | null | undefined,
	files: number,
	ranks: number
): WormholePlayState | null {
	if (!variantState) return current;
	const wirePairs = parseWormholePairs(variantState.wormholePairs);
	if (wirePairs.length === 0) {
		if (!current) return null;
		return {
			pairs: current.pairs,
			cooldown: parseWormholeCooldown(variantState.pairCooldown, current.pairs.length)
		};
	}
	const pairs = wormholePairsFromWire(wirePairs, files, ranks);
	return {
		pairs,
		cooldown: parseWormholeCooldown(variantState.pairCooldown, pairs.length)
	};
}

/** Editor / FE use dense indices; engine uses padded large-board indices. */
export function wormholePairsToWire(
	pairs: WormholePair[],
	files: number,
	ranks: number
): WormholePair[] {
	return pairs.map(
		([a, b]) =>
			[smallToLargeIndex(a, files, ranks), smallToLargeIndex(b, files, ranks)] as WormholePair
	);
}

export function wormholePairsFromWire(
	pairs: WormholePair[],
	files: number,
	ranks: number
): WormholePair[] {
	return pairs.map(
		([a, b]) =>
			[largeToSmallIndex(a, files, ranks), largeToSmallIndex(b, files, ranks)] as WormholePair
	);
}

export function wormholePairLabel(pair: WormholePair, files: number, ranks: number): string {
	return `${smallIndexToAlgebraic(pair[0], files, ranks)} ↔ ${smallIndexToAlgebraic(pair[1], files, ranks)}`;
}

export function validateVariantConfig(
	variantType: string,
	customData: Record<string, unknown> | null | undefined
): string | null {
	if (variantType === VariantType.NCheck) {
		const n = Number(customData?.targetChecks);
		if (!Number.isFinite(n) || n < 1 || n > 10) {
			return 'Set checks to win between 1 and 10 for n-Check.';
		}
		return null;
	}
	if (variantType === VariantType.Wormhole) {
		const pairs = parseWormholePairs(customData?.wormholePairs);
		if (pairs.length === 0) {
			return 'Wormhole requires at least one portal pair. Use Rule Editor → Add pair on board.';
		}
		const seen = new Set<number>();
		for (const [a, b] of pairs) {
			if (a === b) return 'Each wormhole pair needs two different squares.';
			if (seen.has(a) || seen.has(b)) {
				return 'A square can only belong to one wormhole pair.';
			}
			seen.add(a);
			seen.add(b);
		}
		return null;
	}
	return null;
}
