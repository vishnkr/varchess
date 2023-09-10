<script lang="ts">
	import Board from '$lib/board/Board.svelte';
import { me, members, type Member, Role  } from '$lib/store/stores';
	import { onDestroy } from 'svelte';

	let membersList: Member[];
	let unsubscribe = members.subscribe((newVal)=> membersList = newVal );

	const handleAction = (memberId: number, action: string) => {
    // Handle the selected action here
    if (action === 'SetBlack') {
      // Perform the "Set Black" action
    } else if (action === 'SetWhite') {
      // Perform the "Set White" action
    } else if (action === 'RemoveUser') {
      // Perform the "Remove User from Room" action
    }
  };

  let openDropdownId: number | null = null;

  const toggleDropdown = (memberId: number) => {
    openDropdownId = openDropdownId === memberId ? null : memberId;
  };
  onDestroy(() => unsubscribe());
</script>

<div class="bg-zinc-900 flex flex-col rounded py-2 my-2">
  <h4 class="font-semibold text-center text-xl text-white">Members</h4>
  <div class="min-h-[20em] overflow-y-auto bg-white">
    {#each membersList as member (member.id)}
      <div class="flex justify-between items-center py-1 px-2 border-b border-gray-300">
		{#if member.isHost}
		<div class="bg-red-900 text-white text-sm p-1 mx-1 rounded-sm">Host</div>
		{/if}
		{#if member.role == Role.Black}
		<div class="bg-black text-white text-sm p-1 mx-1 rounded-sm">Black</div>
		{/if}
		
		{#if member.role == Role.White}
		<div class="bg-white text-sm p-1 mx-1 rounded-sm">White</div>
		{/if}
        <span>{member.username}</span>
        <div class="flex space-x-2">
          {#if member.role!=Role.Black}
		  <div
            class="bg-black text-white py-1 px-2 rounded-md cursor-pointer"
            on:click={() => handleAction(member.id, 'SetBlack')}
          >
            Set Black
          </div>
		  {/if}
		  {#if member.role!=Role.White}
          <div
            class="bg-white text-black py-1 px-2 rounded-md cursor-pointer"
            on:click={() => handleAction(member.id, 'SetWhite')}
          >
            Set White
          </div>
		  {/if}
		  {#if $me.isHost}
          <div
            class="bg-red-500 text-white py-1 px-2 rounded-md cursor-pointer"
            on:click={() => handleAction(member.id, 'RemoveUser')}
          >
            Remove
          </div>
		  {/if}
        </div>
      </div>
    {/each}
  </div>
</div>