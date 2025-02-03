<script>
	let templates = [
        { num: 0, name: 'Template 0', boardDimensions: '14x13', variantType: 'Classic', objective: 'Capture the King' },
        { num: 2, name: 'quack quack', boardDimensions: '8x8', variantType: 'Fast', objective: 'Checkmate the opponent' },
        { num: 1, name: 'archer', boardDimensions: '10x10', variantType: 'Tactical', objective: 'Complete the challenge' },
        { num: 3, name: 'no mans land', boardDimensions: '15x15', variantType: 'War', objective: 'Conquer the opponent' },
        { num: 4, name: 'varchess temp', boardDimensions: '16x16', variantType: 'Strategy', objective: 'Outsmart your opponent' },
        { num: 12, name: 'Template er', boardDimensions: '12x12', variantType: 'Advanced', objective: 'Survive the longest' },
	];

	let activeTooltipIndex = null;
	let searchQuery = "";

	function showTooltip(index) {
		activeTooltipIndex = index;
	}

	function hideTooltip() {
		activeTooltipIndex = null;
	}

	// Filter templates based on search query
	$: filteredTemplates = templates.filter(template => 
		template.name.toLowerCase().includes(searchQuery.toLowerCase())
	);
</script>

<div class="my-8 relative max-w-7xl mx-auto px-4">
	<h3 class="text-center text-white text-2xl font-bold mb-6">My Templates</h3>
	
	<!-- Search Bar -->
	<div class="mb-6 flex justify-center">
		<input 
			type="text" 
			placeholder="Search Templates..." 
			class="w-full sm:w-80 md:w-96 lg:w-1/2 p-3 rounded-lg bg-gray-700 text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
			bind:value={searchQuery}
		/>
	</div>
	
	<!-- Cards Layout (2 cards per row) -->
	<div class="grid grid-cols-1 sm:grid-cols-2 gap-8">
		{#each filteredTemplates as template, index}
			<div class="bg-gray-800 rounded-lg shadow-lg p-6 flex flex-col sm:flex-row items-center">
				<!-- Chessboard preview (left side) -->
				<div class="flex-none w-full sm:w-40 h-40 bg-gray-600 rounded-md mb-4 sm:mb-0 relative">
					<p class="text-center text-white absolute inset-0 flex items-center justify-center">{template.name}</p>
				</div>

				<!-- Template Info and Action Buttons (right side) -->
				<div class="flex-1 sm:ml-4">
					<!-- Template Name -->
					<h4 class="text-white font-bold text-lg mb-4">{template.name}</h4>

					<!-- Additional Info -->
					<div class="text-gray-400 text-sm mb-4">
						<p><strong>Board Dimensions:</strong> {template.boardDimensions}</p>
						<p><strong>Variant Type:</strong> {template.variantType}</p>
						<p><strong>Objective:</strong> {template.objective}</p>
					</div>
					
					<!-- Action Buttons -->
					<div class="flex space-x-4">
						<button class="bg-green-500 text-white hover:bg-green-600 font-semibold py-2 px-4 rounded-lg transition-colors w-auto flex items-center">
							<i class="fa-solid fa-play mr-2"></i> <!-- Play Icon -->
							<span class="text-sm">Play</span> <!-- Small Text -->
						</button>
						<button class="bg-blue-500 text-white hover:bg-blue-600 font-semibold py-2 px-4 rounded-lg transition-colors w-auto flex items-center">
							<i class="fa-solid fa-edit mr-2"></i> <!-- Edit Icon -->
							<span class="text-sm">Edit</span> <!-- Small Text -->
						</button>
						<button class="bg-red-500 text-white hover:bg-red-600 font-semibold py-2 px-4 rounded-lg transition-colors w-auto flex items-center">
							<i class="fa-solid fa-trash mr-2"></i> <!-- Trash Icon -->
							<span class="text-sm">Delete</span> <!-- Small Text -->
						</button>
					</div>
				</div>
			</div>
		{/each}
	</div>

	<!-- Pagination Buttons (Optional) -->
	<div class="flex justify-center mt-8 space-x-4">
		<button class="bg-blue-500 text-white py-2 px-4 rounded-lg hover:bg-blue-700 transition-colors">Previous</button>
		<button class="bg-blue-500 text-white py-2 px-4 rounded-lg hover:bg-blue-700 transition-colors">Next</button>
	</div>
</div>
