<script lang="ts">
	import { page } from '$app/stores';
	import { gameState } from '$lib/store/stores';
	import { goto } from '$app/navigation';
	let copiedToClipboard = false;
	let currentGameId = $page.data.gameId;

	const copyToClipboard = () => {
		copiedToClipboard = true;
		setTimeout(() => {
			copiedToClipboard = false;
		}, 2000);
		navigator.clipboard.writeText(currentGameId);
	};
	$: {
		if ($gameState.status === 'InProgress') {
			goto(`/game/${currentGameId}`);
		}
	}
</script>

<div class="font-inter text-zinc-90 flex-grow">
	<!-- Modal -->
	<div class="fixed inset-0 z-10 overflow-y-auto">
		<div class="flex items-center justify-center min-h-screen">
			<div class="bg-white p-8 rounded-md shadow-md">
				<h1 class="text-xl text-gray-800 mb-4">Waiting for opponent...</h1>

				<!-- Shareable input and Copy button -->
				<div class="mb-4 flex flex-col items-center">
					<label for="shareableUrl" class="text-gray-800 mb-2">Share Game ID</label>
					<div class="flex items-center">
						<input
							type="text"
							id="shareableUrl"
							value={currentGameId}
							class="w-48 border p-2 rounded-md mr-2"
							readonly
						/>
						<button on:click={copyToClipboard} class="bg-blue-500 text-white px-2 py-1 rounded-md"
							>Copy</button
						>
					</div>
					{#if copiedToClipboard}
						<p class="text-green-500 mt-2">Copied to clipboard</p>
					{/if}
				</div>
			</div>
		</div>
	</div>
</div>
