import { writable } from 'svelte/store';

export type ToastType = 'success' | 'error' | 'warning';

export interface ToastMessage {
	id: number;
	type: ToastType;
	message: string;
}

function createToastStore() {
	const { subscribe, update } = writable<ToastMessage[]>([]);

	let id = 0;

	function show(type: ToastType, message: string, duration = 3000) {
		const toast = { id: ++id, type, message };
		update((toasts) => [...toasts, toast]);

		setTimeout(() => {
			update((toasts) => toasts.filter((t) => t.id !== toast.id));
		}, duration);
	}

	return {
		subscribe,
		success: (msg: string,duration?:number) => show('success', msg,duration),
		error: (msg: string,duration?:number) => show('error', msg,duration),
		warning: (msg: string,duration?:number) => show('warning', msg,duration),
	};
}

export const toast = createToastStore();
