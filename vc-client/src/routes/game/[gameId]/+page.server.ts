import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, url, locals }) => {
	const session = await locals.auth.validate();
	if (!session) {
		throw redirect(303, '/login');
	}
	const { gameId } = params;

	return {
		gameId
	};
};

export const ssr = false;