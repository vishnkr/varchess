import type { Handle } from '@sveltejs/kit';
export const handle = (async ({ event, resolve }) => {
	return await resolve(event);
}) satisfies Handle;
