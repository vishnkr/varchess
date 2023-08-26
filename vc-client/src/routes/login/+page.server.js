import { redirect } from '@sveltejs/kit'
export const actions = {
    default: async ({cookies, request,url}) => {
        const formData = await request.formData()
        const email = formData.get("email")
        const password = formData.get("password")
        if(email==="a@email.com" && password==="123"){
            cookies.set("logged_in","true")
            throw redirect(303, url.searchParams.get('redirectTo') ?? '/home');
        }
        return {
            email,
            message: "Login info not valid"
        }
        
    }
}