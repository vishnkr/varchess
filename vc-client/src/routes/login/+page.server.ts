import { error, redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { displayAlert, getErrorMessage } from '$lib/store/alert';
import { invalidateAll } from '$app/navigation';

export const load = async ({ locals }) => {
	const session = await locals.auth.validate();
	if (session) throw redirect(302, '/home');
	return {};
};