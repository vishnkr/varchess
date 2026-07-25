<script lang="ts">
	import {
		clearMoveSelectorStores,
		clearPremove,
		checkedKingSquare,
		dest,
		gameState,
		legalMoves,
		lastHintsFrom,
		moveSelectorPiece,
		playerPresence,
		premove,
		src,
		viewPly,
		Status,
		chats
	} from '$lib/store/stores';
	import { dimensions, fen, positionStore } from '$lib/store/stores';
	import { authStore } from '$lib/store/auth';
	import { settings } from '$lib/store/settings';
	import { BoardType, Color, type BoardConfig, type Move, type Position } from '$lib/types';
	import Board from './Board.svelte';
	import { wsStore } from '$lib/websocket';
	import { get } from 'svelte/store';
	import { onDestroy } from 'svelte';
	import { playSfx } from '$lib/utils/sfx';

	export let isFlipped: boolean = false;
	/** When false, clicks do not send moves (replay / game over). */
	export let interactive: boolean = true;

	let boardConfig: BoardConfig | null = null;
	let position: Position | null = null;
	let markedTargets: number[] = [];
	let lastHintSrc: number | null = null;
	let lastPremoveSentKey: string | null = null;
	let processingSelection = false;

	$: username = $authStore.username;
	$: players = $gameState.players;
	$: turn = $gameState.turn;
	$: reviewing = $viewPly !== null;
	$: status = $gameState.status;

	$: mySide =
		players?.playerWhite?.name === username
			? ('w' as const)
			: players?.playerBlack?.name === username
				? ('b' as const)
				: null;

	$: canMove =
		interactive &&
		!reviewing &&
		status === Status.InProgress &&
		mySide != null &&
		turn === mySide;

	$: canPremove =
		interactive &&
		!reviewing &&
		status === Status.InProgress &&
		!!$settings.enablePremove &&
		mySide != null &&
		turn !== mySide;

	/** Live game board (not review / completed): own pieces selectable. */
	$: boardPlayable = interactive && !reviewing && status === Status.InProgress;

	$: selectableColor =
		boardPlayable && mySide === 'w'
			? Color.WHITE
			: boardPlayable && mySide === 'b'
				? Color.BLACK
				: null;

	$: {
		const f = $fen;
		const d = $dimensions;
		if (f && d) {
			const nextType = boardPlayable ? BoardType.GameBoard : BoardType.View;
			if (
				!boardConfig ||
				boardConfig.fen !== f ||
				boardConfig.dimensions !== d ||
				boardConfig.boardType !== nextType ||
				boardConfig.isFlipped !== isFlipped
			) {
				boardConfig = {
					fen: f,
					dimensions: d,
					boardType: nextType,
					isFlipped
				};
			}
		} else if (boardConfig) {
			boardConfig = null;
		}
	}

	$: position = $positionStore;

	$: {
		const from = $src;
		const pm = $premove;
		let next: number[] = [];
		if (canMove && $settings.showPossibleMoves && from != null) {
			next = $legalMoves.filter((m) => m.from === from).map((m) => m.to);
		} else if (pm) {
			next = [pm.to];
		}
		if (
			next.length !== markedTargets.length ||
			next.some((v, i) => v !== markedTargets[i])
		) {
			markedTargets = next;
		}
	}

	function isOwnPiece(piece: { color: Color } | null | undefined): boolean {
		if (!piece || !mySide) return false;
		if (mySide === 'w') return piece.color === Color.WHITE;
		return piece.color === Color.BLACK;
	}

	function requestHintsFor(from: number) {
		if (lastHintSrc === from) return;
		lastHintSrc = from;
		legalMoves.set([]);
		lastHintsFrom.set(null);
		wsStore.requestHints(from);
	}

	function clearHints() {
		if (lastHintSrc == null && get(legalMoves).length === 0) return;
		lastHintSrc = null;
		legalMoves.set([]);
		lastHintsFrom.set(null);
	}

	// Subscribe instead of $: store writes — avoids Svelte reactive cycles.
	const unsubSrc = src.subscribe((from) => {
		const piece = get(moveSelectorPiece);
		if (
			from != null &&
			get(gameState).status === Status.InProgress &&
			piece &&
			isOwnPiece(piece)
		) {
			const gs = get(gameState);
			const name = get(authStore).username;
			const mine =
				(gs.players?.playerWhite?.name === name && gs.turn === 'w') ||
				(gs.players?.playerBlack?.name === name && gs.turn === 'b');
			if (mine) requestHintsFor(from);
			else clearHints();
		} else {
			clearHints();
		}
	});

	const unsubDest = dest.subscribe((to) => {
		if (to == null || processingSelection) return;
		const from = get(src);
		const piece = get(moveSelectorPiece);
		if (from == null || !piece) return;

		processingSelection = true;
		try {
			const pieceNotation =
				piece.color === Color.BLACK ? piece.notation : piece.notation.toUpperCase();
			const pos = get(positionStore);
			const target = pos?.piecePositions[to];
			const isCapture = !!(target && target.color !== piece.color);
			const move: Move = { from, to, piece: pieceNotation, capture: isCapture };

			const gs = get(gameState);
			const name = get(authStore).username;
			const reviewingNow = get(viewPly) !== null;
			const my =
				gs.players?.playerWhite?.name === name
					? 'w'
					: gs.players?.playerBlack?.name === name
						? 'b'
						: null;
			const ourTurn =
				interactive &&
				!reviewingNow &&
				gs.status === Status.InProgress &&
				my != null &&
				gs.turn === my;
			const premoveOk =
				interactive &&
				!reviewingNow &&
				gs.status === Status.InProgress &&
				get(settings).enablePremove &&
				my != null &&
				gs.turn !== my;

			if (ourTurn && isOwnPiece(piece)) {
				// If hints for this piece have arrived, reject illegal destinations locally.
				if (get(lastHintsFrom) === from) {
					const legal = get(legalMoves).some((m) => m.from === from && m.to === to);
					if (!legal) {
						if (get(settings).soundEnabled) playSfx('illegal');
						clearMoveSelectorStores();
						clearHints();
						return;
					}
				}
				wsStore.sendMove(move);
				clearMoveSelectorStores();
				clearPremove();
				lastPremoveSentKey = null;
				clearHints();
			} else if (premoveOk && isOwnPiece(piece)) {
				premove.set(move);
				clearMoveSelectorStores();
				lastPremoveSentKey = null;
				clearHints();
			} else {
				clearMoveSelectorStores();
				clearHints();
			}
		} finally {
			processingSelection = false;
		}
	});

	const unsubTurn = gameState.subscribe((gs) => {
		const pm = get(premove);
		if (!pm || gs.status !== Status.InProgress) return;

		const name = get(authStore).username;
		const my =
			gs.players?.playerWhite?.name === name
				? 'w'
				: gs.players?.playerBlack?.name === name
					? 'b'
					: null;
		if (my == null || gs.turn !== my) return;

		const key = `${pm.from}-${pm.to}`;
		if (key === lastPremoveSentKey) return;

		const pos = get(positionStore);
		const pieceAt = pos?.piecePositions[pm.from];
		const own =
			pieceAt &&
			((my === 'w' && pieceAt.color === Color.WHITE) ||
				(my === 'b' && pieceAt.color === Color.BLACK));
		if (own) {
			lastPremoveSentKey = key;
			wsStore.sendMove(pm);
			clearPremove();
		} else {
			clearPremove();
			lastPremoveSentKey = null;
			chats.pushSystem('Premove canceled');
		}
	});

	onDestroy(() => {
		unsubSrc();
		unsubDest();
		unsubTurn();
	});

	function turnClass(side: 'w' | 'b') {
		const active = turn === side && status === Status.InProgress && !reviewing;
		return active
			? 'ring-2 ring-amber-400 ring-offset-2 ring-offset-transparent shadow-md'
			: 'opacity-70';
	}

	function isOnline(userId: string | undefined): boolean {
		if (!userId) return false;
		return $playerPresence[userId] === true;
	}

	function presenceDot(online: boolean) {
		return online
			? 'h-2.5 w-2.5 rounded-full bg-emerald-400 shadow-[0_0_8px_2px_rgba(52,211,153,0.7)]'
			: 'h-2.5 w-2.5 rounded-full bg-gray-500';
	}
</script>

<div class="grid grid-rows-[auto_minmax(0,1fr)_auto] w-full h-full min-h-0 gap-2">
	<div class="flex justify-center sm:justify-start items-center gap-2 shrink-0 px-1">
		{#if players?.playerWhite?.name && players?.playerBlack?.name}
			{@const side = isFlipped ? 'w' : 'b'}
			{@const name = isFlipped ? players.playerWhite.name : players.playerBlack.name}
			{@const uid = isFlipped ? players.playerWhite.userId : players.playerBlack.userId}
			<div
				class="inline-flex items-center gap-2 min-w-[10rem] justify-center px-3 py-1.5 rounded font-medium text-sm
    {isFlipped
					? 'bg-white text-black'
					: 'bg-black text-white border border-white'} {turnClass(side)}"
			>
				<span class={presenceDot(isOnline(uid))} title={isOnline(uid) ? 'Online' : 'Offline'} />
				{name}
			</div>
		{/if}
	</div>

	<div class="w-full h-full min-h-0 min-w-0 flex items-center justify-center overflow-hidden">
		{#if boardConfig && position}
			<Board
				{boardConfig}
				{position}
				{isFlipped}
				{markedTargets}
				fitContainer
				animateMoves
				premoveFrom={$premove?.from ?? null}
				premoveTo={$premove?.to ?? null}
				showCoordinates={$settings.showCoordinates}
				checkedSquare={$checkedKingSquare}
				{selectableColor}
			/>
		{:else}
			<p class="text-center text-gray-500 mt-4">Waiting for game to start...</p>
		{/if}
	</div>

	<div class="flex justify-center sm:justify-start items-center gap-2 shrink-0 px-1">
		{#if players?.playerWhite?.name && players?.playerBlack?.name}
			{@const side = isFlipped ? 'b' : 'w'}
			{@const name = isFlipped ? players.playerBlack.name : players.playerWhite.name}
			{@const uid = isFlipped ? players.playerBlack.userId : players.playerWhite.userId}
			<div
				class="inline-flex items-center gap-2 min-w-[10rem] justify-center px-3 py-1.5 rounded font-medium text-sm
    {isFlipped
					? 'bg-black text-white border border-white'
					: 'bg-white text-black border border-gray-800'} {turnClass(side)}"
			>
				<span class={presenceDot(isOnline(uid))} title={isOnline(uid) ? 'Online' : 'Offline'} />
				{name}
			</div>
		{/if}
	</div>
</div>
