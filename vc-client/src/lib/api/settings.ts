import type { UserSettings } from '$lib/store/settings';
import { CORE_URL_PROTECTED } from './config';
import { customFetch } from './fetch';

export async function fetchSettings(): Promise<UserSettings | null> {
	const res = await customFetch(`${CORE_URL_PROTECTED}/settings`);
	if (res.status === 404) return null;
	if (!res.ok) throw new Error('Failed to fetch settings');
	return res.json();
}

export async function saveSettings(settings: UserSettings): Promise<UserSettings> {
	const res = await customFetch(`${CORE_URL_PROTECTED}/settings`, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(settings)
	});
	if (!res.ok) throw new Error('Failed to save settings');
	return res.json();
}
