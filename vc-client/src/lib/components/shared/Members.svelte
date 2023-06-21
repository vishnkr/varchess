<script lang="ts">
	import { generateUsername } from '$lib/utils';

	let members = [{ id: 1 }, { id: 2 }, { id: 3 }, { id: 4 }, { id: 5 }, { id: 6 }];
	const handleAction = (memberId: number, action: string) => {
		// Handle the action for the selected member
	};
	export let actions: { type: string; handler: () => void }[] = [];
	let openDropdownId: number | null = null;

	const toggleDropdown = (memberId: number) => {
		openDropdownId = openDropdownId === memberId ? null : memberId;
	};
</script>

<div class="bg-zinc-900 flex flex-col rounded py-2 my-2">
	<h4 class="font-semibold text-center text-xl text-white">Members</h4>
	<div class="max-h-60 overflow-y-auto bg-white">
		{#each members as member (member.id)}
			<div class="flex justify-around items-center py-1 px-2 border-b border-gray-300">
				<span>{generateUsername()}</span>
				<div class="dropdown inline-block relative">
					<button
						class="bg-gray-200 text-gray-600 py-1 px-2 rounded-md hover:bg-gray-300 focus:outline-none"
						on:click={() => toggleDropdown(member.id)}
					>
						Actions
					</button>
					{#if openDropdownId === member.id}
						<ul class="absolute z-10 bg-white border border-gray-300 rounded-md shadow mt-2 w-40">
							{#each actions as action}
								<li>
									<button
										class="hover:bg-gray-100 px-4 py-2 w-full text-left"
										on:click={() => action.handler}
									>
										{action.type}
									</button>
								</li>
							{/each}
						</ul>
					{/if}
				</div>
			</div>
		{/each}
	</div>
</div>
