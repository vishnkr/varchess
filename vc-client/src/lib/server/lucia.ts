import { lucia } from 'lucia';
import { sveltekit } from 'lucia/middleware';
import { pg } from '@lucia-auth/adapter-postgresql';
import postgres from 'pg';
import { github } from '@lucia-auth/oauth/providers';
import {
	GITHUB_CLIENT_ID,
	GITHUB_CLIENT_SECRET,
	ENVIRONMENT,
	POSTGRES_USER,
	POSTGRES_PASSWORD,
	POSTGRES_HOST,
	POSTGRES_PORT,
	POSTGRES_DB
} from '$env/static/private';

const pool = new postgres.Pool({
	user: POSTGRES_USER,
	password: POSTGRES_PASSWORD,
	host: POSTGRES_HOST,
	port: parseInt(POSTGRES_PORT),
	database: POSTGRES_DB
});

export const auth = lucia({
	env: ENVIRONMENT === 'development' ? 'DEV' : 'PROD',
	middleware: sveltekit(),
	adapter: pg(pool, {
		user: 'auth_user',
		key: 'user_key',
		session: 'user_session'
	}),
	getUserAttributes: (data) => {
		return {
			username: data.username
		};
	}
});

export const githubAuth = github(auth, {
	clientId: GITHUB_CLIENT_ID,
	clientSecret: GITHUB_CLIENT_SECRET,
	redirectUri: 'http://localhost:5173/login/github/callback'
});

export type Auth = typeof auth;
