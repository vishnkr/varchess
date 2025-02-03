import { authStore, setAuth, logout } from '$lib/store/auth';

async function refreshAccessToken(): Promise<boolean> {
	let refreshToken: string | null = null;
	let username: string | null = null;

	authStore.subscribe((state) => {
		refreshToken = state.refreshToken;
		username = state.username;
	})();

	if (!refreshToken) {
		logout();
		return false;
	}

	const res = await fetch('/api/refresh', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ refreshToken }),
		credentials: 'include'
	});

	if (!res.ok) {
		logout();
		return false;
	}

	const { accessToken } = await res.json();
	if (username) {
		setAuth(username, accessToken, refreshToken);
	}

	return true;
}

async function customFetch(url: string, options: RequestInit = {}) {
	let accessToken: string | null = null;

	authStore.subscribe((state) => {
		accessToken = state.accessToken;
	})();

	const headers = new Headers(options.headers || {});
	if (accessToken) {
		headers.set('Authorization', `Bearer ${accessToken}`);
	}

	const res = await fetch(url, { ...options, headers, credentials: 'include' });

	if (res.status === 401) {
		const refreshed = await refreshAccessToken();
		if (refreshed) {
			return customFetch(url, options);
		}
	}

	return res;
}

export { customFetch };
