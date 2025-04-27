import { writable } from 'svelte/store';
import { browser } from '$app/environment';

const STORAGE_KEY = 'auth';

interface AuthState {
	username: string | null;
	userId: string | null;
	accessToken: string | null;
	refreshToken: string | null;
}

const defaultState: AuthState = {
	username: null,
	userId:null,
	accessToken: null,
	refreshToken: null
};

const storedAuth = browser ? localStorage.getItem(STORAGE_KEY) : null;
const initialAuth = storedAuth ? JSON.parse(storedAuth) : defaultState;

const authStore = writable<AuthState>(initialAuth);

authStore.subscribe((state) => {
	if(browser){
		localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
	}
});

function setAuth(username: string,userId:string, accessToken: string, refreshToken: string) {
	authStore.set({ username,userId, accessToken, refreshToken });
	console.log('setting auth',username,accessToken,refreshToken)
}

function logout() {
	authStore.set(defaultState);
	if(browser){
		localStorage.removeItem(STORAGE_KEY);
	}
	
}

export { authStore, setAuth, logout };	export type { AuthState };

