import { auth } from '$lib/server/lucia';
import { redirect } from '@sveltejs/kit';

export const GET = async ({ locals }) => {
	const session = await locals.auth.validate();
	if (session) {
		await auth.invalidateSession(session.sessionId);
		locals.auth.setSession(null);
	}
	throw redirect(303, '/');
};
