<script>
	import { goto } from '$app/navigation';
	import { fetchTemplates, deleteTemplate } from '$lib/api/template';
	import Board from '$lib/board/Board.svelte';
	import BoardEditor from '$lib/components/editor/BoardEditor.svelte';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import { toast } from '$lib/store/alert';
	import { BoardType } from '$lib/types';
	import { onMount } from 'svelte';
	let theme = 'light';
	let templates = [];

	const getTemplates = async () => {
		try {
			const templatesResponse = await fetchTemplates(currentPage, itemsPerPage);
			templates = templatesResponse.items;
		} catch (err) {
			console.error(err);
			toast.error('Failed to fetch template. Please try again later.');
		}
	};

	const handleDelete = async (template) => {
		try {
			const id = getTemplateId(template);
			await deleteTemplate(id);
			toast.success('Template deleted successfully');
			await getTemplates();
		} catch (err) {
			console.error(err);
			toast.error('Failed to delete template. Please try again later.');
		}
	};
	onMount(() => {
		getTemplates();
	});

	function nextPage() {
		if (currentPage < totalPages) {
			currentPage++;
			getTemplates();
		}
	}

	function prevPage() {
		if (currentPage > 1) {
			currentPage--;
			getTemplates();
		}
	}

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
	const itemsPerPage = 30;
	const goToEditor = () => goto('/editor');
	const goToTemplateEditor = (template) => {
		const id = getTemplateId(template);
		goto(`/editor?tid=${id}`);
	};

	$: filteredTemplates = templates.filter(
		(template) =>
			template.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
			template.variantType.toLowerCase().includes(searchQuery.toLowerCase())
	);

	$: paginatedTemplates = filteredTemplates.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);

	$: totalPages = Math.ceil(filteredTemplates.length / itemsPerPage);
	const getTemplateId = (template) => template['_id'];
	const getBoardConfig = (template) => {
		return {
			fen: template.fen,
			dimensions: template.dimensions,
			boardType: BoardType.View
		};
	};
</script>

<svelte:head>
	<title>My Templates - Varchess</title>
</svelte:head>
<div class="my-8 relative max-w-7xl mx-auto px-4">
	<h3
		class={`text-center text-2xl font-bold mb-6 ${
			theme === 'dark' ? 'text-white' : 'text-gray-800'
		}`}
	>
		My Templates
	</h3>
	<div class="absolute right-4 top-0 cursor-pointer">
		<a on:click={goToEditor} class={buttonVariants({ variant: 'default' })}>
			<i class="fa-solid fa-plus mr-2" />
			Create Template
		</a>
	</div>

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

	<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
		{#each paginatedTemplates as template}
			<div
				class={`rounded-lg shadow-md p-5 transition 
				${theme === 'dark' ? 'bg-gray-800 text-white' : 'bg-white text-gray-800 border border-gray-200'}
			`}
			>
				<h4 class="text-lg font-semibold mb-3 text-center">{template.name}</h4>
				<div
					class="relative w-full bg-gray-600 rounded-md overflow-hidden mb-4 flex items-center justify-center"
				>
					<Board boardConfig={getBoardConfig(template)} />
				</div>

				<div class="text-sm space-y-1 mb-4">
					<p>
						<strong>Dimensions:</strong>
						{template.dimensions.ranks}x{template.dimensions.files}
					</p>
					<p><strong>Type:</strong> {template.variantType}</p>
				</div>

				<div class="flex flex-wrap gap-2">
					<div class="flex flex-wrap gap-2">
						<Button
							class={`text-sm px-3 py-1.5 rounded border font-medium
								${
									theme === 'dark'
										? 'bg-transparent text-white border-white hover:bg-white/10'
										: 'bg-transparent text-black border-black hover:bg-black/10'
								}
							`}
						>
							<i class="fa-solid fa-play mr-1" /> Play
						</Button>

						<Button
							on:click={goToTemplateEditor(template)}
							class={`text-sm px-3 py-1.5 rounded border font-medium
								${
									theme === 'dark'
										? 'bg-transparent text-white border-white hover:bg-white/10'
										: 'bg-transparent text-black border-black hover:bg-black/10'
								}
							`}
						>
							<i class="fa-solid fa-edit mr-1" /> Edit
						</Button>

						<Button
							on:click={() => handleDelete(template)}
							class={`text-sm px-3 py-1.5 rounded border font-medium
								${
									theme === 'dark'
										? 'bg-transparent text-red-400 border-red-400 hover:bg-red-500/10'
										: 'bg-transparent text-red-600 border-red-400 hover:bg-red-100'
								}
							`}
						>
							<i class="fa-solid fa-trash mr-1" /> Delete
						</Button>
					</div>
				</div>
			</div>
		{/each}
	</div>

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
