import type { LayoutServerLoad } from './$types';
export type OutputType = { user: object; isLoggedIn: boolean };

export const load: LayoutServerLoad = async ({ locals }) => {
	const session = await locals.auth.validate();
	if (session) {
		return {
			userId: session.user.userId,
			username: session.user.username
		};
	}
};
