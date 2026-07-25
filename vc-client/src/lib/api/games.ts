import type { Game } from '$lib/types';
import { CORE_URL_PROTECTED } from './config';
import { customFetch } from './fetch';

export async function fetchGames(opts?: { page?: number; pageSize?: number }): Promise<Game[]> {
	const page = opts?.page ?? 1;
	const pageSize = opts?.pageSize ?? 50;
	const res = await customFetch(`${CORE_URL_PROTECTED}/games?page=${page}&pageSize=${pageSize}`);
	if (!res.ok) throw new Error('Failed to fetch games');
	const data = await res.json();
	// BE may return a bare array or a paginated envelope.
	if (Array.isArray(data)) return data;
	if (data?.items && Array.isArray(data.items)) return data.items;
	return [];
}

export async function fetchGame(id: string): Promise<Game> {
	const res = await customFetch(`${CORE_URL_PROTECTED}/games/${id}`);
	if (!res.ok) throw new Error('Failed to fetch game');
	return res.json();
}

export async function createGame(data: { gc: unknown; templateId: string | null }) {
	const body =
		data.templateId != null && data.templateId !== ''
			? { templateId: data.templateId }
			: { gc: data.gc };
	const res = await customFetch(`${CORE_URL_PROTECTED}/games`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(body)
	});
	if (!res.ok) throw new Error('Failed to create game');
	return res.json();
}
