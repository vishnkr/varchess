import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ url, locals }) => {
	const session = await locals.auth.validate();
	if (!session) {
		throw redirect(303, '/login');
	}
	return {
		userId: session.user.userId,
		username: session.user.username
	};
};
