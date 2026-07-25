import { authStore, setAuth, logout } from '$lib/store/auth';
import { get } from 'svelte/store';
import { CORE_URL } from './config';

async function refreshAccessToken(): Promise<boolean> {
	const state = get(authStore);
	const refreshToken = state.refreshToken;
	const username = state.username;
	const userId = state.userId;

	if (!refreshToken) {
		logout();
		return false;
	}

	const res = await fetch(`${CORE_URL}/api/refresh`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ refreshToken }),
		credentials: 'include'
	});

	if (!res.ok) {
		logout();
		return false;
	}

	const data = await res.json();
	const accessToken = data.accessToken as string | undefined;
	if (!accessToken || !username || !userId) {
		logout();
		return false;
	}
	setAuth(username, userId, accessToken, refreshToken);
	return true;
}

async function customFetch(url: string, options: RequestInit = {}, retry = true) {
	const { accessToken } = get(authStore);

	const headers = new Headers(options.headers || {});
	if (accessToken) {
		headers.set('Authorization', `Bearer ${accessToken}`);
	}

	const res = await fetch(url, { ...options, headers, credentials: 'include' });

	if (res.status === 401 && retry) {
		const refreshed = await refreshAccessToken();
		if (refreshed) {
			return customFetch(url, options, false);
		}
	}

	return res;
}

export { customFetch, refreshAccessToken };
