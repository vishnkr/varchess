<script lang="ts">
    import { webSocketStore } from '$lib/utils/websocket';
    import { onDestroy, onMount } from 'svelte';
	  import type { ChatMessage } from './types';
  
    let chatMessages: ChatMessage[] = []; 
    let inputMessage = ''; 

    function sendMessage() {
      if (inputMessage.trim() !== '') {
        const message = {
          type: 'chat',
          user: 'UserA', 
          text: inputMessage.trim(),
        };
        chatMessages = [...chatMessages, message];
        inputMessage = ''; 
      }
    }
  
    const unsubscribe = webSocketStore.subscribe((ws) => {
      if (ws) {
        ws.addEventListener('message', handleWebSocketMessage);
      }
    });

    function handleWebSocketMessage(event:any) {
      const message = event.detail;
      chatMessages = [...chatMessages, message];
    
    }

    onMount(() => {
      chatMessages = [
        { type: 'roomUpdate', text: 'UserA has joined the room' },
        { type: 'chat', user: 'UserA', text: 'Hello, how are you?' },
        { type: 'chat', user: 'UserA', text: 'I\'m doing great, thanks!' },
        { type: 'roomUpdate', text: 'UserC has joined the room' },
        { type: 'chat', user: 'UserC', text: 'Hello, how are you?' },
        { type: 'chat', user: 'UserA', text: 'I\'m doing great, thanks!' },
        { type: 'roomUpdate', text: 'UserB has joined the room' },
        { type: 'chat', user: 'UserB', text: 'Hello, how are you?' },
        { type: 'chat', user: 'UserA', text: 'I\'m doing great, thanks!' },
        { type: 'roomUpdate', text: 'UserA has left the room' },
        { type: 'chat', user: 'UserB', text: 'Hello, how are you?' },
        { type: 'chat', user: 'UserA', text: 'I\'m doing great, thanks!' },
      ];
    });
    onDestroy(() => {
      unsubscribe();
    });
  </script>
  
  <div class="m-2 bg-white">
    <div class="max-h-60 overflow-y-auto">
      {#each chatMessages as message}
        {#if message.type === 'roomUpdate'}
          <div class="flex justify-center">
            <p class="text-blue-500">{message.text}</p>
          </div>
        {:else if message.type === 'chat'}
          <div class="flex pt-2">
            <p class="pl-4 font-bold text-red-700">{message.user}: </p>
            <p class="ml-2">{message.text}</p>
          </div>
        {/if}
      {/each}
    </div>
  
    <div class="p-2 w-full flex">
        <input class="border border-gray-300 px-4 py-2 text-white rounded-l flex-grow" type="text" bind:value={inputMessage} placeholder="Type a message..." />
        <button class="bg-blue-600 hover:bg-blue-800 text-white font-semibold py-2 px-4 rounded-r" on:click={sendMessage}>Send</button>
      </div>
      
  </div>
  