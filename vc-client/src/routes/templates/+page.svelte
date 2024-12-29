<script>
	let templates = [
		{ num: 0, name: 'Template 0' },
        { num: 2, name: 'quack quack' },
        { num: 1, name: 'archer' },
        { num: 3, name: 'no mans land' },
        { num: 4, name: 'varchess temp' },
        { num: 12, name: 'Template er' },
	];

	let activeTooltipIndex = null;

	function showTooltip(index) {
		activeTooltipIndex = index;
	}

	function hideTooltip() {
		activeTooltipIndex = null;
	}
</script>

<div class="my-8 relative">
	<h3 class="text-center text-white text-2xl font-bold mb-6">My Templates</h3>
	<div class="flex items-center justify-center">
		<table class="w-full max-w-4xl bg-transparent rounded-lg overflow-hidden shadow-md">
			<thead class="bg-gray-900 bg-opacity-75">
				<tr>
					<th class="text-white px-6 py-3 border-r border-gray-700">Variant Template</th>
					<th class="text-white px-6 py-3">Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each templates as template, index}
					<tr class="bg-gray-800 bg-opacity-50 hover:bg-gray-700 transition-colors relative">
						<td class="text-gray-200 text-center px-6 py-4 border-r border-gray-700">
							<div class="relative flex items-center justify-center">
								{template.name}
								<div 
									class="relative"
									on:mouseenter={() => showTooltip(index)}
									on:mouseleave={hideTooltip}>
									<i class="fa-solid fa-chess-board ml-2 cursor-pointer" title="View"></i>
									{#if activeTooltipIndex === index}
										<div class="absolute left-full top-1/2 transform -translate-y-1/2 ml-3 w-40 p-2 bg-gray-900 text-white text-sm rounded-lg shadow-lg opacity-100 transition-opacity duration-200">
											<div class="absolute left-[-6px] top-1/2 transform -translate-y-1/2 w-3 h-3 bg-gray-900 rotate-45"></div>
											Preview of template goes here.
										</div>
									{/if}
								</div>
							</div>
						</td>
						<td class="text-gray-200 flex justify-center px-6 py-4">
							<button class="bg-transparent border border-green-500 text-green-500 hover:bg-green-700 hover:text-white font-semibold py-2 px-4 rounded-lg mr-2 transition-colors">Play</button>
							<button class="bg-transparent border border-blue-500 text-blue-500 hover:bg-blue-700 hover:text-white font-semibold py-2 px-4 rounded-lg mr-2 transition-colors">Edit</button>
							<button class="bg-transparent border border-red-500 text-red-500 hover:bg-red-700 hover:text-white font-semibold py-2 px-4 rounded-lg transition-colors">Delete</button>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
