<script lang="ts">
	import {
		Status,
		clearMoveSelectorStores,
		gameId,
		gameState,
		templateStore,
		gameResult,
		resetGameHistory,
		seekEnd,
		seekNext,
		seekPrev,
		seekStart,
		viewPly,
		moveHistory,
		lastGameConfig,
		chats
	} from '$lib/store/stores';
	import { wsStore, wsStatus, incomingDrawOffer } from '$lib/websocket';
	import { beforeNavigate, goto } from '$app/navigation';
	import { browser } from '$app/environment';
	import { onMount } from 'svelte';
	import { authStore } from '$lib/store/auth.js';
	import { page } from '$app/stores';
	import { get } from 'svelte/store';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import GameBoard from '$lib/board/GameBoard.svelte';
	import MoveList from '$lib/components/play/MoveList.svelte';
	import CustomPiecesPanel from '$lib/components/play/CustomPiecesPanel.svelte';
	import WormholesPanel from '$lib/components/play/WormholesPanel.svelte';
	import GameOverModal from '$lib/components/play/GameOverModal.svelte';
	import Chat from '$lib/components/Chat.svelte';
	import { createGame } from '$lib/api/games';
	import { toast } from '$lib/store/alert';
	import { settings } from '$lib/store/settings';

	let isFlipped = false;
	let username = '';
	let copiedToClipboard = false;
	let currentGameId: string | null = null;
	let showGameOverModal = false;
	let gameOverPrompted = false;
	let replayingSetup = false;
	let showResignConfirm = false;

	const goHome = () => {
		if (!browser) return;
		goto('/home');
	};

	$: {
		if ($authStore?.username && $gameState?.players) {
			username = $authStore.username;
			const { playerBlack } = $gameState.players;
			isFlipped = username === playerBlack?.name;
		}
	}

	$: if ($gameState.status === Status.Completed && !gameOverPrompted) {
		gameOverPrompted = true;
		showGameOverModal = true;
	}

	$: isReviewing = $viewPly !== null;

	$: resultBanner = (() => {
		const r = $gameResult;
		if (!r) return 'Game over';
		const w = (r.winner || '').toLowerCase();
		if (w === 'draw') return `Draw${r.reason ? ` — ${r.reason}` : ''}`;
		if (w === 'white' || w === 'w') {
			const name = $gameState.players?.playerWhite?.name ?? 'White';
			return `${name} wins${r.reason ? ` — ${r.reason}` : ''}`;
		}
		if (w === 'black' || w === 'b') {
			const name = $gameState.players?.playerBlack?.name ?? 'Black';
			return `${name} wins${r.reason ? ` — ${r.reason}` : ''}`;
		}
		return r.reason && r.reason.toLowerCase() !== 'game over' ? `Game over — ${r.reason}` : 'Game over';
	})();

	const copyToClipboard = () => {
		copiedToClipboard = true;
		setTimeout(() => (copiedToClipboard = false), 2000);
		if (currentGameId) navigator.clipboard.writeText(currentGameId);
	};

	onMount(() => {
		currentGameId = $page.params.gameId;
		const auth = get(authStore);
		if (!auth.accessToken || !auth.userId || !currentGameId) {
			goto('/home');
			return;
		}
		clearMoveSelectorStores();
		gameId.set(currentGameId);

		// Ensure we have a live socket for this room (home join may have started one already).
		if (!$wsStore) {
			const resumed = wsStore.tryResume();
			if (!resumed) {
				gameState.updateStatus(Status.Waiting);
				wsStore.newWebSocketConnection(
					`ws://${import.meta.env.VITE_WS_HOST}/play/${currentGameId}`,
					{
						token: auth.accessToken,
						userId: auth.userId
					}
				);
			}
		}
	});

	function clearStores() {
		wsStore.close({ clearGame: true });
		templateStore.removeTemplate();
		gameId.set(null);
		gameState.updateStatus(Status.None);
		resetGameHistory();
		gameResult.set(null);
		chats.clear();
	}

	beforeNavigate(({ to }) => {
		const next = to?.url?.pathname ?? '';
		if (next.startsWith('/play/')) return;
		clearStores();
	});

	const handleDraw = () => wsStore.sendDrawOffer();
	const handleAcceptDraw = () => {
		wsStore.sendDrawAccept();
		incomingDrawOffer.set(false);
	};
	const handleRejectDraw = () => {
		wsStore.sendDrawReject();
		incomingDrawOffer.set(false);
	};
	const handleResignClick = () => {
		if ($settings.confirmResign) {
			showResignConfirm = true;
		} else {
			wsStore.sendResign();
		}
	};
	const confirmResign = () => {
		showResignConfirm = false;
		wsStore.sendResign();
	};

	async function handleReplaySetup() {
		const cfg = get(lastGameConfig);
		const auth = get(authStore);
		if (!cfg || !auth.accessToken || !auth.userId) {
			toast.error('Cannot replay — missing game setup');
			return;
		}
		replayingSetup = true;
		try {
			const colorPref = isFlipped ? 'b' : 'w';
			const { gameId: newGameId } = await createGame({ gc: cfg, templateId: null });
			wsStore.close({ clearGame: false });
			resetGameHistory();
			gameResult.set(null);
			chats.clear();
			gameState.updateStatus(Status.Waiting);
			gameId.set(newGameId);
			wsStore.newWebSocketConnection(`ws://${import.meta.env.VITE_WS_HOST}/play/${newGameId}`, {
				token: auth.accessToken,
				userId: auth.userId,
				colorPref
			});
			showGameOverModal = false;
			gameOverPrompted = false;
			goto(`/play/${newGameId}`);
		} catch (e) {
			console.error(e);
			toast.error('Failed to start rematch with same setup');
		} finally {
			replayingSetup = false;
		}
	}
</script>

<svelte:head>
	<title>Play - Varchess</title>
</svelte:head>

<div class="font-inter h-full max-h-full overflow-hidden flex flex-col">
	{#if $wsStatus === 'reconnecting'}
		<div class="shrink-0 z-20 bg-amber-500/90 text-black text-center text-sm py-1">
			Reconnecting to game…
		</div>
	{/if}

	{#if $gameState.status === Status.Waiting}
		<div class="flex-1 flex items-center justify-center overflow-auto p-4">
			<div class="dark:bg-darkbg bg-lightbg2 p-8 rounded-md shadow-md">
				<h1 class="text-xl dark:text-white text-gray-800 mb-4">Waiting for opponent...</h1>
				<div class="mb-4 flex flex-col items-center">
					<label for="gameId" class="dark:text-white text-gray-800 mb-2">Share Game ID</label>
					<div class="flex items-center gap-2">
						<Input id="gameId" bind:value={currentGameId} type="text" readonly />
						<Button
							on:click={copyToClipboard}
							class="text-sm px-3 py-1.5 rounded border font-medium
								bg-transparent dark:text-white dark:border-white dark:hover:bg-white/10 text-black border-black hover:bg-black/10"
							>Copy</Button
						>
					</div>
					{#if copiedToClipboard}
						<p class="dark:text-white text-black mt-2">Copied to clipboard</p>
					{/if}
				</div>
			</div>
		</div>
	{:else if $gameState.status === Status.InProgress || $gameState.status === Status.Completed}
		<div class="flex-1 min-h-0 flex flex-col lg:flex-row gap-2 p-2">
			<div class="flex-1 min-h-0 min-w-0 flex flex-col gap-2 overflow-hidden">
				{#if $gameState.status === Status.InProgress && isReviewing}
					<div
						class="shrink-0 self-start inline-flex items-center gap-2 rounded px-3 py-1.5 text-sm font-medium
						bg-amber-100 text-amber-900 dark:bg-amber-500/20 dark:text-amber-200"
					>
						Reviewing move {$viewPly} / {$moveHistory.length}
					</div>
				{:else if $gameState.status === Status.Completed}
					<div
						class="shrink-0 self-start inline-flex items-center gap-2 rounded px-3 py-1.5 text-sm font-medium
						bg-emerald-100 text-emerald-900 dark:bg-emerald-500/20 dark:text-emerald-200"
					>
						{resultBanner}
					</div>
				{/if}

				<div class="flex-1 min-h-0 overflow-hidden">
					<GameBoard {isFlipped} interactive={$gameState.status === Status.InProgress} />
				</div>
			</div>

			<aside
				class="lg:w-80 xl:w-96 shrink-0 min-h-0 max-h-[40vh] lg:max-h-none flex flex-col gap-2
					bg-lightbg dark:bg-darkbg2 border border-black dark:border-lightbg rounded-md p-2 overflow-hidden"
			>
				{#if $incomingDrawOffer && $gameState.status === Status.InProgress}
					<div class="flex gap-2 shrink-0">
						<Button on:click={handleAcceptDraw} class="flex-1">Accept draw</Button>
						<Button on:click={handleRejectDraw} variant="outline" class="flex-1">Decline</Button>
					</div>
				{/if}

				{#if $gameState.status === Status.InProgress}
					<div class="flex flex-wrap gap-2 justify-center shrink-0">
						<Button
							on:click={handleDraw}
							class="w-[calc(50%-0.5rem)] text-sm px-2 py-1.5 rounded border font-medium
					bg-transparent text-blue-600 border-blue-600 hover:bg-blue-500/10
					dark:text-blue-400 dark:border-blue-400 dark:hover:bg-blue-500/10"
						>
							<i class="fa-solid fa-handshake-simple mr-1" />
							Draw
						</Button>

						<Button
							on:click={handleResignClick}
							class="w-[calc(50%-0.5rem)] text-sm px-2 py-1.5 rounded border font-medium
					bg-transparent text-red-600 border-red-600 hover:bg-red-600/10
					dark:text-red-400 dark:border-red-400 dark:hover:bg-red-600/10"
						>
							<i class="fa-solid fa-flag mr-1" />
							Resign
						</Button>

						<Button
							on:click={() => (isFlipped = !isFlipped)}
							class="w-[calc(50%-0.5rem)] text-sm px-2 py-1.5 rounded border font-medium
					bg-transparent border-black text-black hover:bg-gray-100
					dark:border-white dark:text-white dark:hover:bg-white/10"
						>
							<i class="fa-solid fa-repeat mr-1" />
							Flip
						</Button>
					</div>
				{:else}
					<div class="flex flex-wrap gap-2 justify-center shrink-0">
						<Button variant="outline" size="sm" on:click={() => (showGameOverModal = true)}
							>Result</Button
						>
						<Button variant="outline" size="sm" on:click={() => (isFlipped = !isFlipped)}>Flip</Button>
						<Button size="sm" on:click={goHome}>Home</Button>
					</div>
				{/if}

				<div class="flex gap-1 justify-center shrink-0">
					<Button variant="outline" size="sm" on:click={seekStart} aria-label="First move">
						<i class="fa-solid fa-backward-fast" />
					</Button>
					<Button variant="outline" size="sm" on:click={seekPrev} aria-label="Previous move">
						<i class="fa-solid fa-backward-step" />
					</Button>
					<Button variant="outline" size="sm" on:click={seekNext} aria-label="Next move">
						<i class="fa-solid fa-forward-step" />
					</Button>
					<Button variant="outline" size="sm" on:click={seekEnd} aria-label="Last move">
						<i class="fa-solid fa-forward-fast" />
					</Button>
				</div>

				<div class="min-h-0 flex-1 overflow-y-auto flex flex-col gap-3 rounded-md bg-black/5 dark:bg-black/40 p-2">
					<CustomPiecesPanel />
					<WormholesPanel />
					<div class="min-h-[8rem] flex-1 overflow-hidden">
						<MoveList />
					</div>
				</div>

				{#if !$settings.hideChat && ($gameState.status === Status.InProgress || $gameState.status === Status.Completed)}
					<div class="h-40 shrink-0 min-h-0">
						<Chat />
					</div>
				{/if}
			</aside>
		</div>

		<GameOverModal
			bind:open={showGameOverModal}
			onHome={goHome}
			onClose={() => (showGameOverModal = false)}
			onReplaySetup={handleReplaySetup}
			replaying={replayingSetup}
		/>

		<Dialog.Root bind:open={showResignConfirm}>
			<Dialog.Content class="max-w-sm dark:bg-darkbg bg-white text-gray-900 dark:text-white">
				<Dialog.Header>
					<Dialog.Title>Resign?</Dialog.Title>
					<Dialog.Description class="text-gray-600 dark:text-gray-300">
						Are you sure you want to resign this game?
					</Dialog.Description>
				</Dialog.Header>
				<Dialog.Footer class="flex gap-2 justify-end">
					<Button variant="outline" on:click={() => (showResignConfirm = false)}>Cancel</Button>
					<Button variant="destructive" on:click={confirmResign}>Resign</Button>
				</Dialog.Footer>
			</Dialog.Content>
		</Dialog.Root>
	{:else}
		<div class="text-center dark:text-white p-8">Connecting…</div>
	{/if}
</div>
