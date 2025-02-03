import { writable } from 'svelte/store';

const STORAGE_KEY = 'auth';

interface AuthState {
	username: string | null;
	accessToken: string | null;
	refreshToken: string | null;
}

const defaultState: AuthState = {
	username: null,
	accessToken: null,
	refreshToken: null
};

const storedAuth = localStorage.getItem(STORAGE_KEY);
const initialAuth = storedAuth ? JSON.parse(storedAuth) : defaultState;

const authStore = writable<AuthState>(initialAuth);

authStore.subscribe((state) => {
	localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
});

function setAuth(username: string, accessToken: string, refreshToken: string) {
	authStore.set({ username, accessToken, refreshToken });
	console.log('setting auth',username,accessToken,refreshToken)
}

function logout() {
	authStore.set(defaultState);
	localStorage.removeItem(STORAGE_KEY);
}

export { authStore, setAuth, logout };	export type { AuthState };

