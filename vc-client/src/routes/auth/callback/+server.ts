import { redirect } from '@sveltejs/kit'

import type { RequestEvent, RequestHandler } from './$types';
import { me } from '$lib/store/stores';
import { ALERT_TYPE, alertMessage, alertType, displayAlert, getErrorMessage } from '$lib/store/alert';


export const GET:RequestHandler = async ({ url, locals,cookies }:RequestEvent) => {
  const redirectURL = `${url.origin}/auth/callback`;
	const expectedState = cookies.get('state');

  const query = new URLSearchParams(url.search);
	const state = query.get('state');
	const code = query.get('code');
  
  const authProviderName = cookies.get('provider');
  const authMethods = await locals.pb?.collection('users').listAuthMethods();

	if (!authMethods?.authProviders) {
		throw redirect(303, '/login');
	}

	const provider = authMethods.authProviders.find(provider => provider.name === authProviderName);;

	if (!provider) {
		console.log('Provider not found');
		throw redirect(303, '/login');
	}

	if (expectedState !== state) {
		console.log('state does not match expected', expectedState, state);
		throw redirect(303, '/login');
	}
  try{
    const userdata = await locals.pb?.collection('users').authWithOAuth2Code(provider.name,code||'',provider.codeVerifier,redirectURL);
	console.log('hy',userdata?.meta?.username,userdata)
  } catch(err){
	console.log(err)
	alertMessage.set('ERROR');
    alertType.set('DANGER');
	setTimeout(() => {
        alertMessage.set('');
        alertType.set('');
    }, 3000);

  }

  throw redirect(303, '/home')
}