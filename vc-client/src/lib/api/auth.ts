import { setAuth } from '$lib/store/auth';

const CORE_URL = import.meta.env.VITE_CORE_URL;

export interface AuthTokens {
	username: string;
	uid: string;
	accessToken: string;
	refreshToken: string;
}

export interface OAuthProviders {
	github: boolean;
	google: boolean;
}

async function readError(res: Response): Promise<string> {
	try {
		const text = await res.text();
		return text?.trim() || res.statusText || 'Request failed';
	} catch {
		return res.statusText || 'Request failed';
	}
}

async function login(identifier: string, password: string): Promise<{ ok: true } | { ok: false; error: string }> {
	const body = identifier.includes('@')
		? { email: identifier, password }
		: { username: identifier, password };

	try {
		const res = await fetch(`${CORE_URL}/login`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(body),
			credentials: 'include'
		});

		if (!res.ok) {
			return { ok: false, error: await readError(res) };
		}

		const data: AuthTokens = await res.json();
		setAuth(data.username, data.uid, data.accessToken, data.refreshToken);
		return { ok: true };
	} catch {
		return { ok: false, error: 'Server is not reachable. Please try again later.' };
	}
}

async function signup(
	email: string,
	username: string,
	password: string
): Promise<{ ok: true } | { ok: false; error: string }> {
	try {
		const res = await fetch(`${CORE_URL}/signup`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ email, username, password }),
			credentials: 'include'
		});
		if (!res.ok) {
			return { ok: false, error: await readError(res) };
		}
		const data: AuthTokens = await res.json();
		if (data.accessToken) {
			setAuth(data.username, data.uid, data.accessToken, data.refreshToken);
		}
		return { ok: true };
	} catch {
		return { ok: false, error: 'Server is not reachable. Please try again later.' };
	}
}

async function fetchOAuthProviders(): Promise<OAuthProviders> {
	try {
		const res = await fetch(`${CORE_URL}/auth/oauth/providers`);
		if (!res.ok) return { github: false, google: false };
		return (await res.json()) as OAuthProviders;
	} catch {
		return { github: false, google: false };
	}
}

function oauthStartUrl(provider: 'github' | 'google'): string {
	return `${CORE_URL}/auth/oauth/${provider}`;
}

export { login, signup, fetchOAuthProviders, oauthStartUrl };
