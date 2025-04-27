import type { Game, Template } from '$lib/types';
import { CORE_URL_PROTECTED } from './config';
import { customFetch } from './fetch';

export async function fetchGames(): Promise<Game[]> {
    const res = await customFetch(`${CORE_URL_PROTECTED}/games`);
    if (!res.ok) throw new Error('Failed to fetch games');
    return res.json();
}

export async function createGame(data: { gc: Template|null; templateId: string|null }) {
    const res = await customFetch(`${CORE_URL_PROTECTED}/games`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });
    console.log('res',res)
    if (!res.ok) throw new Error('Failed to create game');
    return res.json();
}