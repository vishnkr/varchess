
<script>
	import { createWebSocket, me, roomId } from "$lib/store/stores";
    import { goto } from '$app/navigation';
    let username;
    export let data;
    async function handleSubmit() {
		if (username?.length == 0) {
            return;
		}
        try{
		    let ws = await createWebSocket(data.roomId,username);
            me.set({id:0,isHost:true,role:0,username:username})
			roomId.set(data.roomId)
			goto(`/editor/${data.roomId}`)
        } catch (error){
            console.error('WebSocket connection failed:', error);
        }
	}
</script>

<div class="flex items-center justify-center">
    <form on:submit|preventDefault={handleSubmit}>
        <div class="grid grid-cols-1 grid-rows-2">
            <div class="flex justify-center items-center">
                <div class="p-1 items-center">
                    <label for="username">
                        Your username
                        <input
                            type="text"
                            bind:value={username}
                            class="rounded-md border bg-slate-600 text-white border-gray-300 px-4 py-2 focus:border-blue-300 outline-none"
                            name="username"
                        />
                    </label>
                </div>
            </div>
            <div class="flex justify-center items-center">
                <div>
                    <button type="submit" class="flex items-center justify-center gap-x-6">
                        <!-- svelte-ignore a11y-click-events-have-key-events -->
                        <span
                            class="btn-custom-1"
                        >
                            Create Room
                        </span>
                    </button>
                </div>
            </div>
        </div>
    </form>
</div>