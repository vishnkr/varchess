import { writable, type Writable } from 'svelte/store';
import type { userInfo } from './types';


function createUserStore(){
    const {subscribe, set, update} = writable<userInfo>({
        username: ''
    })
    return {
        subscribe,
        set,
        updateRoomId:(roomId:string) =>{
            update((userInfo: userInfo)=>{
                userInfo.currentRoomId = roomId;
                return userInfo
            })
        }
    }
}

export const userStore = createUserStore();
