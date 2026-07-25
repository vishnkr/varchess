<script lang="ts">
	import './board-styles.css';
	import {
		BoardType,
		doesSupportDragDrop,
		isGameBoard,
		type IPiece,
		type SquareColor,
		type SquareInfo,
		isEditor
	} from './types';
	import { editorMaxBoard } from './board';
	import { pieceEditor, boardEditor, jumpPatternEditing, ruleEditor, wormholeEditor } from '$lib/store/editor';
	import { wormholePlayState } from '$lib/store/stores';
	import wallSvg from '$lib/assets/svg/wall.svg';
	import { src, dest, moveSelectorPiece, recentMove } from '$lib/store/stores';
	import { settings } from '$lib/store/settings';
	import { Color, VariantType } from '$lib/types';
	import { parseWormholePairs } from '$lib/utils/wormhole';
	import { get } from 'svelte/store';

	export let squareData: SquareInfo;

	export let color: SquareColor;
	export let piece: IPiece | null = null;
	export let wall: boolean = false;
	export let boardId: string = 'board';
	export let boardType: BoardType = BoardType.GameBoard;

	let isDraggable: boolean = doesSupportDragDrop(boardType);
	export let isMarkedTarget = false;
	export let isPremoveSrc = false;
	export let isPremoveDest = false;
	export let isInCheck = false;
	/** Move-pattern editor: blue slide ray / red jump target. */
	export let mpHighlight: 'slide' | 'jump' | null = null;
	export let showCoordinates = false;
	export let boardRanks = 8;
	/** Hide piece briefly while a move animation flies onto this square. */
	export let animHidePiece = false;

	export let nonPieceSvg: string | null = null;
	/**
	 * When set (game play), only this color's pieces can be selected / dragged.
	 * Opponent pieces are only valid as capture destinations.
	 */
	export let selectableColor: Color | null = null;

	const DRAG_MIME = 'application/x-vc-piece';

	function getPieceClass(p: IPiece) {
		return p.color.charAt(0).toLowerCase() + p.notation.toLowerCase();
	}

	let pieceEl: HTMLElement;
	let squareEl: HTMLElement;
	let hover = false;
	let dropHover = false;
	/** True if a drag started from this square — skip the trailing click. */
	let didDrag = false;

	function isOwnPiece(p: IPiece | null | undefined): boolean {
		if (!p) return false;
		if (selectableColor == null) return true;
		return p.color === selectableColor;
	}

	$: isMoveSrc = $src != null && $src === squareData.squareIndex && isGameBoard(boardType);
	$: showMarkedTarget = isGameBoard(boardType) && isMarkedTarget;
	$: isRecentSrc =
		isGameBoard(boardType) &&
		$settings.highlightLastMove &&
		!!$recentMove &&
		$recentMove.from === squareData.squareIndex;
	$: isRecentDest =
		isGameBoard(boardType) &&
		$settings.highlightLastMove &&
		!!$recentMove &&
		$recentMove.to === squareData.squareIndex;

	$: notation = squareData.squareNotation ?? '';
	$: fileLabel = notation.slice(0, 1);
	$: rankLabel = notation.slice(1);
	$: showFileCoord = showCoordinates && squareData.gridX === boardRanks;
	$: showRankCoord = showCoordinates && squareData.gridY === 1;

	$: gameDraggable = isGameBoard(boardType) && isOwnPiece(piece);
	$: ownPieceHere = isGameBoard(boardType) && isOwnPiece(piece);
	/** Hover chrome: own pieces anytime; any square once a move is started / marked. */
	$: allowHoverChrome =
		!isGameBoard(boardType) || ownPieceHere || $src != null || isMarkedTarget || !piece;

	$: wormholePairs = parseWormholePairs($ruleEditor.customData?.wormholePairs);
	$: wormholePairIndex = (() => {
		const idx = squareData.squareIndex;
		if (isEditor(boardType)) {
			return wormholePairs.findIndex((p) => p[0] === idx || p[1] === idx);
		}
		const play = $wormholePlayState;
		if (!play) return -1;
		return play.pairs.findIndex((p) => p[0] === idx || p[1] === idx);
	})();
	$: isWormholeSquare = wormholePairIndex >= 0;
	$: isWormholePending = isEditor(boardType) && $wormholeEditor.pending === squareData.squareIndex;
	$: isWormholeClosed =
		!isEditor(boardType) &&
		isWormholeSquare &&
		(($wormholePlayState?.cooldown[wormholePairIndex] ?? 0) > 0);

	function selectPiece(p: IPiece) {
		if (!isOwnPiece(p)) return;
		src.set(squareData.squareIndex);
		dest.set(null);
		moveSelectorPiece.set(p);
	}

	function clearSelection() {
		src.set(null);
		dest.set(null);
		moveSelectorPiece.set(null);
	}

	function sameSide(a: IPiece, b: IPiece) {
		return a.color === b.color;
	}

	function handleGameClick() {
		if (didDrag) {
			didDrag = false;
			return;
		}

		const selected = $moveSelectorPiece;
		const from = $src;

		if (from == null) {
			// Only pick up your own pieces
			if (piece && isOwnPiece(piece)) selectPiece(piece);
			return;
		}

		// Clicked the already-selected square → deselect
		if (from === squareData.squareIndex) {
			clearSelection();
			return;
		}

		// Clicked another own piece → switch selection (lichess/chess.com behavior)
		if (piece && selected && sameSide(piece, selected) && isOwnPiece(piece) && !isMarkedTarget) {
			selectPiece(piece);
			return;
		}

		// Empty square, capture, or legal target → attempt move
		dest.set(squareData.squareIndex);
	}

	function handleGameDragStart(e: DragEvent) {
		if (!isGameBoard(boardType) || !piece || !isOwnPiece(piece)) {
			e.preventDefault();
			return;
		}
		didDrag = true;
		selectPiece(piece);
		e.dataTransfer?.setData(
			DRAG_MIME,
			JSON.stringify({ idx: squareData.squareIndex, color: piece.color, notation: piece.notation })
		);
		e.dataTransfer!.effectAllowed = 'move';
		if (pieceEl) {
			pieceEl.classList.add('dragging');
		}
	}

	function handleGameDragEnd() {
		pieceEl?.classList.remove('dragging');
		dropHover = false;
		// Keep selection after a cancelled drag so user can click a destination.
		setTimeout(() => {
			didDrag = false;
		}, 0);
	}

	function handleGameDragOver(e: DragEvent) {
		if (!isGameBoard(boardType)) return;
		e.preventDefault();
		if (e.dataTransfer) e.dataTransfer.dropEffect = 'move';
		dropHover = true;
	}

	function handleGameDragLeave() {
		dropHover = false;
	}

	function handleGameDrop(e: DragEvent) {
		if (!isGameBoard(boardType)) return;
		e.preventDefault();
		dropHover = false;
		const raw = e.dataTransfer?.getData(DRAG_MIME);
		if (!raw) return;

		let payload: { idx: number; color: string; notation: string };
		try {
			payload = JSON.parse(raw);
		} catch {
			return;
		}

		if (payload.idx === squareData.squareIndex) {
			return;
		}

		const draggedColor = payload.color as Color;
		// Dropping on another own piece that isn't a legal capture target → reselect
		if (
			piece &&
			piece.color === draggedColor &&
			isOwnPiece(piece) &&
			!isMarkedTarget
		) {
			selectPiece(piece);
			return;
		}

		src.set(payload.idx);
		dest.set(squareData.squareIndex);
	}

	/* ── Editor handlers (unchanged behavior) ── */

	function handleEditorDragStart(e: DragEvent) {
		if (!isEditor(boardType) || !piece) return;
		e.dataTransfer?.setData('dragInfo', JSON.stringify({ idx: squareData.squareIndex, piece }));
		if (pieceEl) pieceEl.style.opacity = '0.4';
	}

	function handleEditorDragEnd() {
		if (pieceEl) pieceEl.style.opacity = '1';
	}

	function handleEditorDragOver(e: DragEvent) {
		if (!isEditor(boardType)) return;
		e.preventDefault();
		const sel = $pieceEditor.pieceSelection?.piece;
		if (sel) piece = sel;
		if ($boardEditor.isWallSelectorOn) {
			editorMaxBoard.updatePieceInfo(squareData.row, squareData.column, {
				isPiecePresent: false,
				wall: true
			});
		} else if (sel) {
			editorMaxBoard.updatePieceInfo(squareData.row, squareData.column, {
				isPiecePresent: true,
				piece: sel
			});
		}
	}

	function handleEditorDrop(e: DragEvent) {
		e.preventDefault();
		const data = e.dataTransfer?.getData('dragInfo');
		if (data && !wall && isDraggable) {
			const obj = JSON.parse(data);
			piece = obj.piece;
			editorMaxBoard.updatePieceInfo(squareData.row, squareData.column, {
				isPiecePresent: true,
				piece
			});
		}
	}

	function handleEditorClick() {
		if ($wormholeEditor.active && $ruleEditor.variantType === VariantType.Wormhole) {
			handleWormholePairClick();
			return;
		}
		if ($boardEditor.isWallSelectorOn) {
			wall = !wall;
			editorMaxBoard.updatePieceInfo(squareData.row, squareData.column, {
				isPiecePresent: false,
				wall,
				piece: null
			});
		} else {
			if (wall) {
				editorMaxBoard.updatePieceInfo(squareData.row, squareData.column, {
					isPiecePresent: false,
					wall: false
				});
				return;
			}
			editorMaxBoard.updatePieceInfo(squareData.row, squareData.column, {
				isPiecePresent: piece ? false : true,
				piece: piece ? null : $pieceEditor.pieceSelection.piece
			});
		}
	}

	function handleWormholePairClick() {
		const idx = squareData.squareIndex;
		const pairs = parseWormholePairs(get(ruleEditor).customData?.wormholePairs);
		const used = new Set<number>();
		for (const [a, b] of pairs) {
			used.add(a);
			used.add(b);
		}
		const pending = get(wormholeEditor).pending;
		if (pending == null) {
			if (used.has(idx)) return;
			wormholeEditor.setPending(idx);
			return;
		}
		if (pending === idx) {
			wormholeEditor.resetPending();
			return;
		}
		if (used.has(idx)) return;
		ruleEditor.updateCustomData({
			wormholePairs: [...pairs, [pending, idx]]
		});
		wormholeEditor.start();
	}

	function handleMPEditorClick() {
		if (piece) return;
		if (!$jumpPatternEditing) return;
		const selectedPiece = $pieceEditor.pieceSelection;
		if (!selectedPiece?.piece?.notation) return;
		const notation = selectedPiece.piece.notation.toLowerCase();
		const jumpOffset = [squareData.row - 4, squareData.column - 4];
		if (jumpOffset[0] === 0 && jumpOffset[1] === 0) return;

		const jumps = $pieceEditor.movePatterns[notation]?.jumpOffsets ?? [];
		const isJumpOffsetPresent = jumps.some(
			(o) => o[0] === jumpOffset[0] && o[1] === jumpOffset[1]
		);
		if (isJumpOffsetPresent) {
			pieceEditor.removeJumpPattern(notation, jumpOffset);
		} else {
			pieceEditor.addJumpPattern(notation, jumpOffset);
		}
	}

	function handleClick(e: MouseEvent | KeyboardEvent) {
		e.preventDefault();
		switch (boardType) {
			case BoardType.Editor:
				handleEditorClick();
				break;
			case BoardType.MovePatternEditor:
				handleMPEditorClick();
				break;
			case BoardType.GameBoard:
				handleGameClick();
				break;
			default:
				break;
		}
	}

	function onDragStart(e: DragEvent) {
		if (isGameBoard(boardType)) handleGameDragStart(e);
		else if (isEditor(boardType)) handleEditorDragStart(e);
	}

	function onDragEnd(e: DragEvent) {
		e.preventDefault();
		if (isGameBoard(boardType)) handleGameDragEnd();
		else if (isEditor(boardType)) handleEditorDragEnd();
	}

	function onDragOver(e: DragEvent) {
		if (isGameBoard(boardType)) handleGameDragOver(e);
		else if (isEditor(boardType)) handleEditorDragOver(e);
	}

	function onDragLeave(e: DragEvent) {
		e.preventDefault();
		if (isGameBoard(boardType)) handleGameDragLeave();
		dropHover = false;
		hover = false;
	}

	function onDrop(e: DragEvent) {
		if (isGameBoard(boardType)) handleGameDrop(e);
		else if (isEditor(boardType)) handleEditorDrop(e);
	}
</script>

<div
	class={`relative
	${showMarkedTarget ? 'marked-target' : ''} 
	${piece ? 'has-piece' : ''}
	${isMoveSrc ? 'move-src' : ''}
	${isRecentSrc ? 'recent-src' : ''}
	${isRecentDest ? 'recent-dest' : ''}
	${isPremoveSrc ? 'premove-src' : ''}
	${isPremoveDest ? 'premove-dest' : ''}
	${isInCheck ? 'in-check' : ''}
	${mpHighlight === 'slide' ? 'mp-slide' : ''}
	${mpHighlight === 'jump' ? 'mp-jump' : ''}
	${isWormholeSquare ? 'wormhole' : ''}
	${isWormholeClosed ? 'wormhole-closed' : ''}
	${isWormholePending ? 'wormhole-pending' : ''}
	${ownPieceHere ? 'own-piece' : ''}
	${dropHover ? 'drop-hover' : ''}`}
	style="--x:{squareData.gridX}; --y:{squareData.gridY};"
	data-square-color={color}
	data-square
	id={`${boardId}-s-${squareData.squareIndex}`}
	class:hover={hover && allowHoverChrome}
	role="button"
	tabindex="0"
	aria-label={`Square ${squareData.squareNotation ?? squareData.squareIndex}`}
	on:keydown={(e) => {
		if (e.key === 'Enter' || e.key === ' ') handleClick(e);
	}}
	bind:this={squareEl}
	on:mouseenter={() => {
		if (allowHoverChrome) hover = true;
	}}
	on:mouseleave={() => {
		hover = false;
		dropHover = false;
	}}
	on:click={handleClick}
	on:dragover={onDragOver}
	on:dragleave={onDragLeave}
	on:drop={onDrop}
>
	{#if showFileCoord}
		<span class="coord coord-file">{fileLabel}</span>
	{/if}
	{#if showRankCoord}
		<span class="coord coord-rank">{rankLabel}</span>
	{/if}
	{#if piece}
		<div
			class={`absolute bg-piece ${gameDraggable || isDraggable ? 'draggable' : ''} w-full h-full ${getPieceClass(
				piece
			)} ${animHidePiece ? 'anim-hidden' : ''}`}
			draggable={gameDraggable || (isEditor(boardType) && isDraggable)}
			role="img"
			aria-label="Game piece"
			id={`${boardId}-p-${squareData.squareIndex}`}
			bind:this={pieceEl}
			on:dragstart={onDragStart}
			on:dragend={onDragEnd}
		/>
	{:else if wall}
		<div class="absolute inset-0 flex items-center justify-center bg-red-400">
			<img alt="wall" draggable={isEditor(boardType)} src={wallSvg} class="w-full h-full" />
		</div>
	{:else if nonPieceSvg}
		<div class="absolute inset-0 flex items-center justify-center">
			<img alt="nonpiece" draggable={false} src={nonPieceSvg} class="w-full h-full" />
		</div>
		<slot />
	{:else}
		<slot />
	{/if}
	{#if isWormholeSquare || isWormholePending}
		<span class="wormhole-badge" aria-hidden="true">
			{isWormholePending ? '·' : isWormholeClosed ? '✕' : wormholePairIndex + 1}
		</span>
	{/if}
</div>

<style>
	[data-square-color] {
		width: 100%;
		height: 100%;
		min-width: 0;
		min-height: 0;
		position: relative;
		grid-column: var(--y);
		grid-row: var(--x);
		background-color: var(--square-color);
	}

	.hover {
		background-color: var(--default-hover-square);
	}

	[data-square].drop-hover {
		outline: 2px solid rgba(59, 130, 246, 0.85);
		outline-offset: -2px;
	}

	[data-square-color='dark'] {
		--square-color: var(--default-dark-square);
		--label-color: var(--default-light-square);
		--square-color-hover: var(--square-color-dark-hover);
		--move-target-marker-color: var(--move-target-marker-color-dark-square);
		--square-color-active: var(--square-color-dark-active);
		--outline-color-active: var(--outline-color-dark-active);
	}

	[data-square-color='light'] {
		--square-color: var(--default-light-square);
		--label-color: var(--default-dark-square);
		--square-color-hover: var(--square-color-light-hover);
		--move-target-marker-color: var(--move-target-marker-color-light-square);
		--square-color-active: var(--square-color-light-active);
		--outline-color-active: var(--outline-color-light-active);
	}
	[data-square].marked-target {
		background: radial-gradient(
			var(--move-target-marker-color) var(--move-target-marker-radius),
			var(--square-color) calc(var(--move-target-marker-radius) + 1px)
		);
	}

	[data-square].move-src {
		--square-color: var(--square-color-active);
	}
	[data-square].recent-src {
		background-image: linear-gradient(
			to bottom,
			var(--recent-move-src-color),
			var(--recent-move-src-color)
		) !important;
		background-color: var(--square-color);
	}

	[data-square].recent-dest {
		background-image: linear-gradient(
			to bottom,
			var(--recent-move-dest-color),
			var(--recent-move-dest-color)
		) !important;
		background-color: var(--square-color);
	}

	[data-square].recent-src:hover,
	[data-square].recent-dest:hover {
		background-image: linear-gradient(
			to bottom,
			var(--recent-move-src-color),
			var(--recent-move-src-color)
		);
	}

	[data-square].recent-dest:hover {
		background-image: linear-gradient(
			to bottom,
			var(--recent-move-dest-color),
			var(--recent-move-dest-color)
		);
	}

	[data-square]:hover:not(.move-src):not(.recent-src):not(.recent-dest) {
		background-image: linear-gradient(
			to bottom,
			var(--square-color-dark-hover),
			var(--square-color-dark-hover)
		);
		background-color: var(--square-color);
	}

	[data-square].has-piece.marked-target {
		background: radial-gradient(
			var(--square-color) var(--move-target-marker-radius-occupied),
			var(--move-target-marker-color) calc(var(--move-target-marker-radius-occupied) + 1px)
		);
	}

	[data-square].premove-src,
	[data-square].premove-dest {
		background-image: linear-gradient(
			to bottom,
			rgba(96, 165, 250, 0.45),
			rgba(96, 165, 250, 0.45)
		) !important;
		background-color: var(--square-color);
	}

	[data-square].in-check {
		background-image: radial-gradient(
			circle at center,
			hsla(0, 95%, 52%, 0.85) 0%,
			hsla(0, 90%, 45%, 0.55) 42%,
			hsla(0, 80%, 40%, 0.2) 68%,
			transparent 72%
		) !important;
		background-color: var(--square-color);
		animation: check-pulse 1.1s ease-in-out infinite;
	}

	[data-square].mp-slide {
		background-image: linear-gradient(
			to bottom,
			rgba(37, 99, 235, 0.55),
			rgba(37, 99, 235, 0.55)
		) !important;
		background-color: var(--square-color);
	}

	[data-square].mp-jump {
		background-image: linear-gradient(
			to bottom,
			rgba(220, 38, 38, 0.55),
			rgba(220, 38, 38, 0.55)
		) !important;
		background-color: var(--square-color);
	}

	[data-square].wormhole {
		box-shadow: inset 0 0 0 3px rgba(139, 92, 246, 0.85);
	}

	[data-square].wormhole-closed {
		box-shadow: inset 0 0 0 3px rgba(217, 119, 6, 0.9);
		opacity: 0.92;
	}

	[data-square].wormhole-pending {
		box-shadow: inset 0 0 0 3px rgba(16, 185, 129, 0.9);
	}

	.wormhole-badge {
		position: absolute;
		top: 2px;
		right: 2px;
		z-index: 3;
		min-width: 1rem;
		height: 1rem;
		padding: 0 3px;
		border-radius: 999px;
		background: rgba(91, 33, 182, 0.92);
		color: #fff;
		font-size: 0.65rem;
		font-weight: 700;
		line-height: 1rem;
		text-align: center;
		pointer-events: none;
	}

	[data-square].wormhole-closed .wormhole-badge {
		background: rgba(180, 83, 9, 0.95);
	}

	@keyframes check-pulse {
		0%,
		100% {
			filter: brightness(1);
		}
		50% {
			filter: brightness(1.12);
		}
	}

	.coord {
		position: absolute;
		z-index: 2;
		pointer-events: none;
		font-size: var(--coords-font-size, 0.65rem);
		font-family: var(--coords-font-family, sans-serif);
		line-height: 1;
		color: var(--label-color);
		opacity: 0.85;
		font-weight: 600;
	}
	.coord-file {
		right: 0.15rem;
		bottom: 0.1rem;
	}
	.coord-rank {
		left: 0.15rem;
		top: 0.1rem;
	}

	.draggable {
		cursor: grab;
		touch-action: none;
		z-index: 1;
	}
	.draggable:active,
	.draggable.dragging {
		cursor: grabbing;
		opacity: 0.35;
	}
	[data-square].own-piece {
		cursor: grab;
	}
	[data-square].has-piece:not(.own-piece) {
		cursor: default;
	}
	.anim-hidden {
		opacity: 0 !important;
	}
</style>
