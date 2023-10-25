import type { LayoutServerLoad } from './$types';

export type OutputType = { user: object; isLoggedIn: boolean };

export const load: LayoutServerLoad = async ({ locals }) => {
	const user = locals.user
	if (user){
		return {user, isLoggedIn:true};
	}
	return {
		user: undefined,
		isLoggedIn:false
	}
}