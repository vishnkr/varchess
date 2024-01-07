// routes/login/github/callback/+server.ts
import { auth, githubAuth } from "$lib/server/lucia.js";
import { OAuthRequestError } from "@lucia-auth/oauth";
import { redirect } from "@sveltejs/kit";

export const GET = async ({ url, cookies, locals }) => {
	const storedState = cookies.get("github_oauth_state");
	const state = url.searchParams.get("state");
	const code = url.searchParams.get("code");
    console.log('here4')
	// validate state
	if (!storedState || !state || storedState !== state || !code) {
        console.log('here3')
		return new Response(null, {
			status: 400
		});
	}
	try {
		const { getExistingUser, githubUser, createUser } =
			await githubAuth.validateCallback(code);

		const getUser = async () => {
			const existingUser = await getExistingUser();
			if (existingUser) return existingUser;
			const user = await createUser({
				attributes: {
					username: githubUser.login
				}
			});
			return user;
		};

		const user = await getUser();
		const session = await auth.createSession({
			userId: user.userId,
			attributes: {}
		});
		locals.auth.setSession(session);
		return new Response(null, {
			status: 302,
			headers: {
				Location: "/home"
			}
		});
	} catch (e) {

		if (e instanceof OAuthRequestError) {
			// invalid code
            console.log('here1')
			return new Response(null, {
				status: 400
			});
		}
        console.log('here2',e)
		return new Response(null, {
			status: 500
		});
	}
};