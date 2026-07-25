import { CORE_URL_PROTECTED } from './config';
import { customFetch } from './fetch';

/** Permanently delete the current account. `confirm` must equal the username. */
export async function deleteAccount(confirm: string): Promise<void> {
	const res = await customFetch(`${CORE_URL_PROTECTED}/account`, {
		method: 'DELETE',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ confirm })
	});
	if (!res.ok) {
		const text = await res.text().catch(() => '');
		throw new Error(text?.trim() || 'Failed to delete account');
	}
}
