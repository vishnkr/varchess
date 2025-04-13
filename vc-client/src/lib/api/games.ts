import type { Game } from '$lib/types';
import { CORE_URL } from './config';

export async function fetchGames(): Promise<Game[]> {
    const res = await fetch(`${CORE_URL}/games`);
    if (!res.ok) throw new Error('Failed to fetch games');
    return res.json();
}

export async function createGame(data: { name: string; type: string }) {
    const res = await fetch(`${CORE_URL}/games`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });
    if (!res.ok) throw new Error('Failed to create game');
    return res.json();
}