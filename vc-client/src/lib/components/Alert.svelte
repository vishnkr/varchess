<script>
    import { alertMessage, alertType, ALERT_TYPE } from "$lib/store/alert";
    $: color = ()=>{
        if ($alertType === ALERT_TYPE.DANGER){
            return 'bg-red-500'
        } else if ($alertType === ALERT_TYPE.SUCCESS){
            return 'bg-green-500'
        } else{
            return 'bg-gray-500'
        }
    }
</script>
{#if $alertMessage}
<div class="h-10 my-4 mx-6">
    
        <span class={`text-white text-center rounded-md p-2 mb-2 ${color()}`}>
            {$alertMessage}
        </span>
</div>
{/if}

