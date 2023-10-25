import { error, redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { displayAlert, getErrorMessage } from '$lib/store/alert';
import { invalidateAll } from '$app/navigation';

export type OutputType = { 
	githubAuthProviderRedirect: string; 
	githubAuthProviderState: string ;
	googleAuthProviderRedirect: string;
	googleAuthProviderState: string;
};
export const load: PageServerLoad<OutputType> = async ({ locals, url }) =>{
	const authMethods = await locals.pb?.collection('users').listAuthMethods();
	if (!authMethods) {
		return {
			githubAuthProviderRedirect: '',
			githubAuthProviderState: '',
			googleAuthProviderRedirect: '',
			googleAuthProviderState: '',
		};
	}

	const redirectURL = `${url.origin}/auth/callback`;
	const githubAuthProvider = authMethods.authProviders.find(provider => provider.name === 'github');
    const googleAuthProvider = authMethods.authProviders.find(provider => provider.name === 'google');
	const githubAuthProviderRedirect = `${githubAuthProvider?.authUrl}${redirectURL}`;
	const googleAuthProviderRedirect = `${googleAuthProvider?.authUrl}${redirectURL}`;
	return {
		githubAuthProviderRedirect: githubAuthProviderRedirect,
		githubAuthProviderState: githubAuthProvider?.state,
		googleAuthProviderRedirect: googleAuthProviderRedirect,
		googleAuthProviderState: googleAuthProvider?.state
	};
}
export const actions = {
	login: async ({ locals, request }) => {
		const body = Object.fromEntries(await request.formData());
		console.log(body)
		try {
			await locals.pb?.collection('users').authWithPassword(body.username, body.password);
			
		} catch (err) {
			console.log('Error: ', err);
			displayAlert('Authentication failed','DANGER',9000)
			
			return { status: 400, body: 'Authentication failed' };
		}

		throw redirect(303, '/home');
	}
};