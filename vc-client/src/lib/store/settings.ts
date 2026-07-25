import { browser } from '$app/environment';
import { get, writable } from 'svelte/store';
import { COLOR_THEMES } from '$lib/utils/index';
import { theme as appThemeStore } from '$lib/store/theme';
import { fetchSettings, saveSettings as apiSaveSettings } from '$lib/api/settings';

export interface UserSettings {
	boardTheme: string;
	showPossibleMoves: boolean;
	enablePremove: boolean;
	highlightLastMove: boolean;
	showCoordinates: boolean;
	confirmResign: boolean;
	hideChat: boolean;
	soundEnabled: boolean;
	appTheme: 'light' | 'dark';
}

export const DEFAULT_SETTINGS: UserSettings = {
	boardTheme: 'Default',
	showPossibleMoves: true,
	enablePremove: false,
	highlightLastMove: true,
	showCoordinates: false,
	confirmResign: true,
	hideChat: false,
	soundEnabled: false,
	appTheme: 'dark'
};

const STORAGE_KEY = 'vc-user-settings';

function normalize(raw: Partial<UserSettings> | null | undefined): UserSettings {
	const merged = { ...DEFAULT_SETTINGS, ...(raw ?? {}) };
	if (!COLOR_THEMES[merged.boardTheme]) merged.boardTheme = DEFAULT_SETTINGS.boardTheme;
	if (merged.appTheme !== 'light' && merged.appTheme !== 'dark') {
		merged.appTheme = DEFAULT_SETTINGS.appTheme;
	}
	return merged;
}

function readLocal(): UserSettings {
	if (!browser) return { ...DEFAULT_SETTINGS };
	try {
		const raw = localStorage.getItem(STORAGE_KEY);
		if (raw) return normalize(JSON.parse(raw));
		// Migrate legacy keys
		const legacy: Partial<UserSettings> = {};
		const bt = localStorage.getItem('board-theme');
		if (bt) legacy.boardTheme = bt;
		const sm = localStorage.getItem('show-possible-moves');
		if (sm != null) legacy.showPossibleMoves = JSON.parse(sm);
		const ep = localStorage.getItem('enable-premove');
		if (ep != null) legacy.enablePremove = JSON.parse(ep);
		const app = localStorage.getItem('theme') as 'light' | 'dark' | null;
		if (app === 'light' || app === 'dark') legacy.appTheme = app;
		return normalize(legacy);
	} catch {
		return { ...DEFAULT_SETTINGS };
	}
}

function writeLocal(s: UserSettings) {
	if (!browser) return;
	localStorage.setItem(STORAGE_KEY, JSON.stringify(s));
	localStorage.setItem('board-theme', s.boardTheme);
	localStorage.setItem('show-possible-moves', JSON.stringify(s.showPossibleMoves));
	localStorage.setItem('enable-premove', JSON.stringify(s.enablePremove));
	localStorage.setItem('theme', s.appTheme);
}

export function applyBoardTheme(boardTheme: string) {
	if (!browser) return;
	const theme = COLOR_THEMES[boardTheme] ?? COLOR_THEMES.Default;
	document.documentElement.style.setProperty('--default-light-square', theme.lightColor);
	document.documentElement.style.setProperty('--default-dark-square', theme.darkColor);
}

export function applyAppTheme(mode: 'light' | 'dark') {
	if (!browser) return;
	const cur = get(appThemeStore);
	if (cur !== mode) appThemeStore.set(mode);
	document.documentElement.classList.toggle('dark', mode === 'dark');
	localStorage.setItem('theme', mode);
}

function createSettingsStore() {
	const { subscribe, set, update } = writable<UserSettings>(DEFAULT_SETTINGS);

	function apply(s: UserSettings) {
		applyBoardTheme(s.boardTheme);
		applyAppTheme(s.appTheme);
	}

	return {
		subscribe,
		set,
		update,
		/** Load from localStorage immediately; optionally merge from API when authenticated. */
		async load(opts?: { syncRemote?: boolean }) {
			const local = readLocal();
			set(local);
			apply(local);
			if (!opts?.syncRemote) return local;
			try {
				const remote = await fetchSettings();
				if (remote) {
					const merged = normalize(remote);
					set(merged);
					writeLocal(merged);
					apply(merged);
					return merged;
				}
			} catch (err) {
				console.warn('settings remote load failed', err);
			}
			return local;
		},
		patch(partial: Partial<UserSettings>) {
			update((cur) => {
				const next = normalize({ ...cur, ...partial });
				writeLocal(next);
				if (partial.boardTheme != null) applyBoardTheme(next.boardTheme);
				if (partial.appTheme != null) applyAppTheme(next.appTheme);
				return next;
			});
		},
		async save(next?: UserSettings) {
			const value = normalize(next ?? get({ subscribe }));
			set(value);
			writeLocal(value);
			apply(value);
			try {
				await apiSaveSettings(value);
			} catch (err) {
				console.warn('settings remote save failed', err);
				throw err;
			}
			return value;
		}
	};
}

export const settings = createSettingsStore();
