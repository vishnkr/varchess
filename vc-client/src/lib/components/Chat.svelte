<script lang="ts">
	import { chats, MessageType } from '$lib/store/stores';
	import { wsStore } from '$lib/websocket';
	import { tick } from 'svelte';

	let inputMessage = '';
	let listEl: HTMLDivElement | null = null;
	let lastChatLen = 0;

	async function sendMessage() {
		const text = inputMessage.trim();
		if (!text) return;
		wsStore.sendChat(text);
		inputMessage = '';
		await tick();
		scrollToBottom();
	}

	function scrollToBottom() {
		if (listEl) listEl.scrollTop = listEl.scrollHeight;
	}

	function onKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter' && !e.shiftKey) {
			e.preventDefault();
			sendMessage();
		}
	}

	$: {
		const len = $chats.length;
		if (listEl && len !== lastChatLen) {
			lastChatLen = len;
			tick().then(scrollToBottom);
		}
	}
</script>

<div class="flex flex-col h-full min-h-0 rounded-md border border-gray-300 dark:border-gray-600 bg-white/60 dark:bg-black/30">
	<h3
		class="text-sm font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400 px-3 py-2 border-b border-gray-200 dark:border-gray-700"
	>
		Chat
	</h3>
	<div bind:this={listEl} class="flex-1 min-h-0 overflow-y-auto px-3 py-2 space-y-1.5">
		{#if $chats.length === 0}
			<p class="text-xs text-gray-500 dark:text-gray-400">No messages yet</p>
		{:else}
			{#each $chats as message}
				{#if message.messageType === MessageType.ChatMessage}
					<p class="text-sm break-words">
						<span class="font-semibold text-emerald-700 dark:text-emerald-400"
							>{message.username}:</span
						>
						<span class="text-gray-800 dark:text-gray-100 ml-1">{message.content}</span>
					</p>
				{:else if message.messageType === MessageType.UserJoin}
					<p class="text-xs text-center text-emerald-600 dark:text-emerald-400 italic">
						{message.content}
					</p>
				{:else if message.messageType === MessageType.UserLeave}
					<p class="text-xs text-center text-amber-600 dark:text-amber-400 italic">
						{message.content}
					</p>
				{:else}
					<p class="text-xs text-center text-sky-600 dark:text-sky-400 italic">{message.content}</p>
				{/if}
			{/each}
		{/if}
	</div>
	<div class="flex gap-1 p-2 border-t border-gray-200 dark:border-gray-700">
		<input
			class="flex-1 min-w-0 rounded border border-gray-300 dark:border-gray-600 bg-transparent px-2 py-1.5 text-sm dark:text-white"
			type="text"
			bind:value={inputMessage}
			placeholder="Type a message…"
			on:keydown={onKeydown}
		/>
		<button
			type="button"
			class="shrink-0 rounded bg-gray-900 dark:bg-white text-white dark:text-black px-3 py-1.5 text-sm font-medium"
			on:click={sendMessage}>Send</button
		>
	</div>
</div>
