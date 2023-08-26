import { writable } from "svelte/store";

const serverUrl =
	import.meta.env.VITE_ENVIRONMENT === 'production'
		? import.meta.env.VITE_SERVER_BASE
		: 'localhost:5000';

function createWebSocketStore() {
	const { subscribe, set, update } = writable<WebSocket | null>(null);
	let ws: WebSocket | null = null;
		  
	function connect(){
		if (!ws || ws.readyState === WebSocket.CLOSED) {
			ws = new WebSocket(`ws://${serverUrl}/ws`);
			ws.addEventListener('message', (event) => {
				const data = JSON.parse(event.data);
				update((currentWs) => {
				if (currentWs) {
					currentWs.dispatchEvent(new CustomEvent('message', { detail: data }));
				}
				return currentWs;
				});
			});
		}
		set(ws);
		console.log('connected ws')
	}
	return {
		subscribe,
		connect,
	};
}
		  
export const webSocketStore = createWebSocketStore();
		  