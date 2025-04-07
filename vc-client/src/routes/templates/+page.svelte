<script>
	import { Button, buttonVariants } from '$lib/components/ui/button';
	let theme = 'light';
	let allTemplates = [
		{
			num: 0,
			name: 'Template 0',
			boardDimensions: '14x13',
			variantType: 'Classic',
			objective: 'Capture the King'
		},
		{
			num: 2,
			name: 'quack quack',
			boardDimensions: '8x8',
			variantType: 'Fast',
			objective: 'Checkmate the opponent'
		},
		{
			num: 1,
			name: 'archer',
			boardDimensions: '10x10',
			variantType: 'Tactical',
			objective: 'Complete the challenge'
		},
		{
			num: 3,
			name: 'no mans land',
			boardDimensions: '15x15',
			variantType: 'War',
			objective: 'Conquer the opponent'
		},
		{
			num: 4,
			name: 'varchess temp',
			boardDimensions: '16x16',
			variantType: 'Strategy',
			objective: 'Outsmart your opponent'
		},
		{
			num: 12,
			name: 'Template er',
			boardDimensions: '12x12',
			variantType: 'Advanced',
			objective: 'Survive the longest'
		},

		...Array.from({ length: 20 }, (_, i) => ({
			num: i + 20,
			name: `Extra Template ${i + 1}`,
			boardDimensions: '10x10',
			variantType: 'Fun',
			objective: 'Dominate'
		}))
	];

	if (typeof window !== 'undefined') {
		const updateTheme = () => {
			theme = document.documentElement.classList.contains('dark') ? 'dark' : 'light';
		};

		updateTheme();
		const observer = new MutationObserver(updateTheme);
		observer.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] });
	}
	let searchQuery = '';
	let currentPage = 1;
	const itemsPerPage = 6;

	$: filteredTemplates = allTemplates.filter(
		(template) =>
			template.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
			template.variantType.toLowerCase().includes(searchQuery.toLowerCase()) ||
			template.objective.toLowerCase().includes(searchQuery.toLowerCase())
	);

	$: paginatedTemplates = filteredTemplates.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);

	$: totalPages = Math.ceil(filteredTemplates.length / itemsPerPage);

	function nextPage() {
		if (currentPage < totalPages) currentPage++;
	}

	function prevPage() {
		if (currentPage > 1) currentPage--;
	}
</script>

<div class="my-8 relative max-w-7xl mx-auto px-4">
	<h3
		class={`text-center text-2xl font-bold mb-6 ${
			theme === 'dark' ? 'text-white' : 'text-gray-800'
		}`}
	>
		My Templates
	</h3>

	<div class="mb-6 flex justify-center">
		<input
			type="text"
			placeholder="Search Templates..."
			class={`w-full sm:w-80 md:w-96 lg:w-1/2 p-3 rounded-lg 
				${
					theme === 'dark'
						? 'bg-gray-700 text-white placeholder-gray-400 focus:ring-blue-500'
						: 'bg-gray-100 text-black placeholder-gray-500 focus:ring-blue-400'
				}
				focus:outline-none focus:ring-2`}
			bind:value={searchQuery}
		/>
	</div>

	<!-- Cards Layout -->
	<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
		{#each paginatedTemplates as template}
			<div
				class={`rounded-lg shadow-md p-5 transition 
				${theme === 'dark' ? 'bg-gray-800 text-white' : 'bg-white text-gray-800 border border-gray-200'}
			`}
			>
				<div
					class="relative w-full aspect-[1/1] bg-gray-600 rounded-md overflow-hidden mb-4 flex items-center justify-center"
				>
					<span class="text-white font-semibold text-lg">{template.name}</span>
				</div>

				<div class="text-sm space-y-1 mb-4">
					<p><strong>Board:</strong> {template.boardDimensions}</p>
					<p><strong>Type:</strong> {template.variantType}</p>
					<p><strong>Objective:</strong> {template.objective}</p>
				</div>

				<div class="flex flex-wrap gap-2">
					<div class="flex flex-wrap gap-2">
						<!-- Play Button -->
						<Button
							class={`text-sm px-3 py-1.5 rounded border font-medium
								${theme === 'dark' 
									? 'bg-transparent text-white border-white hover:bg-white/10' 
									: 'bg-transparent text-black border-black hover:bg-black/10'}
							`}
						>
							<i class="fa-solid fa-play mr-1" /> Play
						</Button>
					
						<!-- Edit Button -->
						<Button
							class={`text-sm px-3 py-1.5 rounded border font-medium
								${theme === 'dark' 
									? 'bg-transparent text-white border-white hover:bg-white/10' 
									: 'bg-transparent text-black border-black hover:bg-black/10'}
							`}
						>
							<i class="fa-solid fa-edit mr-1" /> Edit
						</Button>
					
						<Button
							class={`text-sm px-3 py-1.5 rounded border font-medium
								${theme === 'dark' 
									? 'bg-transparent text-red-400 border-red-400 hover:bg-red-500/10' 
									: 'bg-transparent text-red-600 border-red-400 hover:bg-red-100'}
							`}
						>
							<i class="fa-solid fa-trash mr-1" /> Delete
						</Button>
					</div>
					
				</div>
			</div>
		{/each}
	</div>

	<!-- Pagination -->
	{#if totalPages > 1}
		<div class="flex justify-center items-center gap-4 mt-8">
			<Button on:click={prevPage} disabled={currentPage === 1} class="text-sm">Previous</Button>

			<span class={theme === 'dark' ? 'text-white' : 'text-gray-800'}>
				Page {currentPage} of {totalPages}
			</span>

			<Button on:click={nextPage} disabled={currentPage === totalPages} class="text-sm">
				Next
			</Button>
		</div>
	{/if}
</div>
