import { writable } from "svelte/store";

export const ALERT_TYPE = {
    DANGER : 'DANGER',
    INFO: 'INFO',
    SUCCESS: 'SUCCESS'
}

export const alertMessage = writable('');
export const alertType = writable('');

export const displayAlert = (message:string, type =ALERT_TYPE.INFO, resetTime?:number) =>{
    alertMessage.set(message);
    alertType.set(type);
    if (resetTime){
        setTimeout(() => {
            alertMessage.set('')
        }, resetTime);
    }
}
export function getErrorMessage(error: unknown) {
    if (error instanceof Error) return error.message
    return String(error)
  }