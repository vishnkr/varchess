<script lang="ts">
	import BoardEditor from '$lib/components/editor/BoardEditor.svelte';
	import pieceSvg from '$lib/assets/svg/piece.svg';
	import boardSvg from '$lib/assets/svg/board.svg';
	import Tabs from '$lib/components/shared/Tabs.svelte';
	import type { BoardConfig } from '$lib/board/types';
	import EditableBoard from '$lib/board/EditableBoard.svelte';
	import PieceEditor from '$lib/components/editor/PieceEditor.svelte';
	import GameSettings from '$lib/components/editor/GameSettings.svelte';
	import ExpandableCard from '$lib/components/ExpandableCard.svelte';
	let items = ['Custom', 'Predefined'];
	let activeItem = 'Custom';
	let inputValue = 'localhost:5sfds137';
	const tabChange = (e: CustomEvent<string>) => (activeItem = e.detail);

	export let boardConfig: BoardConfig = {
		fen: 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR',
		dimensions: { ranks: 8, files: 8 },
		editable: true,
		interactive: false,
		isFlipped: false
	};

	let editorSettings = {
		curPieceSelected: 'p',
		isDisableToggled: false
	};

	const copyToClipboard = () => {
		navigator.clipboard.writeText(inputValue);
	};

	let clearBoard: () => void;
	let shiftBoard: (direction: string) => void;

	let members = [
		{ id: 1, username: 'User1' },
		{ id: 2, username: 'User2' },
		{ id: 3, username: 'User3' },
		{ id: 4, username: 'User4' },
		{ id: 5, username: 'User5' },
		{ id: 6, username: 'User6' }
	];
	const handleAction = (memberId: number, action: string) => {
		// Handle the action for the selected member
	};

	let openDropdownId: number | null = null;

	const toggleDropdown = (memberId: number) => {
		openDropdownId = openDropdownId === memberId ? null : memberId;
	};
</script>

<svelte:head>
	<title>Editor - Varchess</title>
</svelte:head>
<div class="font-inter text-zinc-90 flex-grow">
	<div class="flex-1 flex m-4 lg:flex-row flex-col">
		<div class="bg-zinc-700 rounded-md lg:w-3/12 mx-3 p-3 max-h-[45rem] overflow-y-auto">
			<div class="border-b border-gray-200 dark:border-gray-700 flex flex-col text-center">
				<div class="flex justify-center">
					<Tabs {activeItem} {items} on:tabChange={tabChange} />
				</div>
				{#if activeItem === 'Custom'}
					<ExpandableCard svg={boardSvg} title="Board Editor">
						<BoardEditor
							bind:dimensions={boardConfig.dimensions}
							on:shift={(e) => shiftBoard(e.detail)}
							on:clear={() => clearBoard()}
						/>
					</ExpandableCard>
					<ExpandableCard svg={pieceSvg} title="Piece Editor">
						<PieceEditor />
					</ExpandableCard>

					<ExpandableCard iconClass="fa-solid fa-gear fa-lg" title="Game Settings">
						<GameSettings />
					</ExpandableCard>
				{:else}
					<ExpandableCard title="Duck Chess">
						<div>
							<p>Chess with ducks!</p>
							<button class="rounded bg-blue-400 p-1 my-2">View Rules</button>
							<button class="rounded bg-orange-400 p-1 my-2">Select Variant</button>
						</div>
					</ExpandableCard>
					<ExpandableCard title="Wormhole">
						<div>
							<p>Teleporation!</p>
							<button class="rounded bg-blue-400 p-1 my-2">View Rules</button>
							<button class="rounded bg-orange-400 p-1 my-2">Select Variant</button>
						</div>
					</ExpandableCard>
					<ExpandableCard title="Sniper Chess">
						<div>
							<p>Pieces can make long range attacks</p>
							<button class="rounded bg-blue-400 p-1 my-2">View Rules</button>
							<button class="rounded bg-orange-400 p-1 my-2">Select Variant</button>
						</div>
					</ExpandableCard>
				{/if}
			</div>
		</div>
		<div class="bg-zinc-700 rounded-md lg:w-6/12 mx-3 my-3 p-3">
			<EditableBoard {boardConfig} bind:shift={shiftBoard} bind:clear={clearBoard} />
		</div>
		<div class="flex flex-col bg-zinc-700 rounded-md lg:w-3/12 mx-3 p-3">
			<div class="flex mb-5">
				<input
					class="border border-gray-300 px-4 py-2 text-white rounded-l w-64"
					bind:value={inputValue}
					disabled
				/>
				<button
					class="bg-blue-500 hover:bg-blue-700 text-white font-semibold py-2 px-2 rounded-r"
					on:click={copyToClipboard}
				>
					Share <i class="fa-solid fa-link" />
				</button>
			</div>
			<div class="bg-zinc-900 flex flex-col rounded">
				<h4 class="font-semibold text-center text-xl text-white">Members</h4>
				<div class="max-h-60 overflow-y-auto bg-white">
					{#each members as member (member.id)}
						<div class="flex justify-between items-center py-2 px-4 border-b border-gray-300">
							<span>{member.username}</span>
							<div class="dropdown inline-block relative">
								<button
									class="bg-gray-200 text-gray-600 py-1 px-2 rounded-md hover:bg-gray-300 focus:outline-none"
									on:click={() => toggleDropdown(member.id)}
								>
									Actions
								</button>
								{#if openDropdownId === member.id}
									<ul
										class="absolute z-10 bg-white border border-gray-300 rounded-md shadow mt-2 w-40"
									>
										<li>
											<button
												class="hover:bg-gray-100 px-4 py-2 w-full text-left"
												on:click={() => handleAction(member.id, 'Kick')}
											>
												Remove
											</button>
										</li>
										<li>
											<button
												class="hover:bg-gray-100 px-4 py-2 w-full text-left"
												on:click={() => handleAction(member.id, 'Ban')}
											>
												Set as White
											</button>
										</li>
										<li>
											<button
												class="hover:bg-gray-100 px-4 py-2 w-full text-left"
												on:click={() => handleAction(member.id, 'Ban')}
											>
												Set as Black
											</button>
										</li>
									</ul>
								{/if}
							</div>
						</div>
					{/each}
				</div>
			</div>
		</div>
	</div>
</div>

<style>
	::-webkit-scrollbar {
		width: 10px;
	}
</style>
