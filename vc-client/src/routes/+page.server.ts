import { fail, redirect } from '@sveltejs/kit'

import type { Actions } from './$types'
const baseUrl:string = import.meta.env.VITE_VARCHESS_SERVER_BASE;
export const actions: Actions = {
	createRoom: async({request})=>{
		const formData = await request.formData();
		const username = formData.get('username');
		if(username?.length ==0){
			fail(400,{
				error:true,
				message:'Username cannot be empty'
			})
		}
		console.log('we here',username,baseUrl);
		const response = await fetch(`${baseUrl}/create-room`,{
			method:'POST',
			headers:{
				'Content-Type':'application/json',
			},
			body: JSON.stringify({username})
		})
		return { success: true };
		
	},
}
