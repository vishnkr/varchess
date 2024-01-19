import { fail, redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { createWebSocket } from "$lib/websocket";
import { gameId } from "$lib/store/stores";
import { goto } from "$app/navigation";
import { apiServerUrl, wsServerUrl } from "$lib/server/private";
import { displayAlert } from "$lib/store/alert";

export const load: PageServerLoad = async ({ locals }) => {
	const session = await locals.auth.validate();
	if (!session) throw redirect(302, "/login");
	return {
		userId: session.user.userId,
		username: session.user.username
	};
};

export const actions = {
	createRoom: async({locals})=>{
			const session = await locals.auth.validate();
			if (!session) throw redirect(302, "/login");
			const username = session.user.username;
			const response = await fetch(`http://${apiServerUrl}/rooms`, {
				method: 'POST',
				headers: {
				'Content-Type': 'application/json'
				},
				body: JSON.stringify({username}) 
			});
			if (response.ok){
				const data = await response.json();
			} else {
				displayAlert('Unable to create room. Please try again later.','DANGER',6000)
				return fail(500, { error: 'Connection failed' });
			}
	},

	joinRoom: async({locals,request})=>{
		try {
			const data = await request.formData();
			const session = await locals.auth.validate();
			if (!session) throw redirect(302, "/login");
			const gameIdInp = data.get('gameId') as string;
			if (!gameIdInp) {
				return fail(400, { gameId, missing: true });
			}
			const response = await fetch(`http://${apiServerUrl}/join`, {
				method: 'POST',
				headers: {
				'Content-Type': 'application/json'
				},
				body: JSON.stringify({username: session.user.username, gameId: gameIdInp}) 
			});
			if (response.ok){
				const data = await response.json();
				gameId.set(gameIdInp)
				return {
					game_state : data.state,
					game_config: data.config
				}
			}
			//throw redirect(302,`/game/${gameId}`);
			return fail(500, { error: 'Unable to join game' });
		} catch (error) {
			//displayAlert('WebSocket connection error:','DANGER',7000);
			//return fail(500, { error: 'WebSocket connection failed' });
		}
	}
} 