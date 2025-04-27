// stores/game.ts
import { writable } from 'svelte/store';

export enum Status {
	Waiting,
	InProgress,
	Completed
}

export const gameId = writable<string | null>(null);

export interface GameState {
	status: Status;
	turn: 'w' | 'b';
	players: { playerWhite: string; playerBlack: string };
}

function createGameState() {
	const { subscribe, update, set } = writable<GameState>({
		status: Status.Waiting,
		turn: 'w',
		players: { playerWhite: '', playerBlack: '' }
	});

	return {
		subscribe,
		startGame: (players: { white: string; black: string }) =>
			update((g) => ({
				...g,
				status: Status.InProgress,
				players: {
					playerWhite: players.white,
					playerBlack: players.black
				}
			})),
		changeTurn: () =>
			update((g) => ({
				...g,
				turn: g.turn === 'w' ? 'b' : 'w'
			})),
		reset: () => set({ status: Status.Waiting, turn: 'w', players: { playerWhite: '', playerBlack: '' } })
	};
}

export const gameState = createGameState();
