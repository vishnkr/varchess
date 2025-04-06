import { setAuth } from '$lib/store/auth';
const CORE_URL = import.meta.env.VITE_CORE_URL;

interface LoginResponse {
	username: string;
	accessToken: string;
	refreshToken: string;
}

async function login(username: string, password: string): Promise<boolean> {
	const res = await fetch(`${CORE_URL}/login`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ username, password }),
		credentials: 'include'
	});

	if (!res.ok) return false;

	const data: LoginResponse = await res.json();
	setAuth(data.username, data.accessToken, data.refreshToken);
	return true;
}

async function signup(email: string, username: string, password: string): Promise<boolean> {
	try {
		const res = await fetch(`${CORE_URL}/signup`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ email, username, password })
		});
		return res.ok;
	} catch (err) {
		return false;
	}
}

export { login, signup };
