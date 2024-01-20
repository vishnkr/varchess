<script lang="ts">
	import { chats, gameId, wsStore, type ChatMessage, MessageType } from '$lib/store/stores';
	import { EventChatMessage } from '$lib/store/types';
	import { camelToSnake } from '$lib/utils';

	let chatMessages: ChatMessage[] = [];
	let inputMessage = '';
	
	function sendMessage() {
		if (inputMessage.trim() !== '') {
			const wsMessage = {
				event: EventChatMessage,
				params:{
					gameId: $gameId,
					message: inputMessage.trim()
				}
			}
			$wsStore?.send(JSON.stringify(camelToSnake(wsMessage)));
			inputMessage = '';
		}
	}

	$: {
		chatMessages = $chats;
		console.log(chatMessages)
	}

</script>

<div class="m-2 bg-white">
	<div class="max-h-60 overflow-y-auto">
		{#each chatMessages as message}
			{#if message.messageType === MessageType.ChatMessage}
				<div class="flex pt-2">
					<p class="pl-4 font-bold text-red-700">{message.username}:</p>
					<p class="ml-2">{message.content}</p>
				</div>
			{:else}
				<div class="flex justify-center">
					<p class="text-blue-500">{message.content}</p>
				</div>
			{/if}
		{/each}
	</div>

	<div class="p-2 w-full flex">
		<input
			class="border border-gray-300 px-4 py-2 rounded-l flex-grow"
			type="text"
			bind:value={inputMessage}
			placeholder="Type a message..."
		/>
		<button
			class="bg-blue-600 hover:bg-blue-800 text-white font-semibold py-2 px-4 rounded-r"
			on:click={sendMessage}>Send</button
		>
	</div>
</div>
