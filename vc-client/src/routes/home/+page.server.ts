import { fail, redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { createWebSocket } from "$lib/websocket";
import { roomId } from "$lib/store/stores";
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
				if(data.roomId) {
					roomId.set(data.roomId)
					throw redirect(302,`/editor`);
				} 			
			
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
			const roomCode = data.get('roomCode') as string;
			if (!roomCode) {
				return fail(400, { roomCode, missing: true });
			}
			throw redirect(302,`/editor`);
		} catch (error) {
			displayAlert('WebSocket connection error:','DANGER',7000);
			return fail(500, { error: 'WebSocket connection failed' });
		}
	}
} 