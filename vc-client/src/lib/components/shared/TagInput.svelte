<script lang="ts">
	import { pieceEditor } from '$lib/store/editor';

	let selectedOption = '';
	let tagList: string[] = [];

	export let dropDownText: string;
	export let slideDirections: Record<string, number[]>;
	export let initialOptions: string[] = Object.keys(slideDirections);
	let options = [...initialOptions];

	const addSlide = (offset: number[]) => {
		const pieceType = $pieceEditor.pieceSelection?.pieceType;
		if (pieceType) {
			pieceEditor.addSlidePattern(pieceType, offset);
		}
	};
	const removeSlide = (offset: number[]) => {
		const pieceType = $pieceEditor.pieceSelection?.pieceType;
		if (pieceType) {
			pieceEditor.removeSlidePattern(pieceType, offset);
		}
	};

	function addTag() {
		if (selectedOption !== '' && !tagList.includes(selectedOption)) {
			tagList = [...tagList, selectedOption];
			options = initialOptions.filter((option) => !tagList.includes(option));
			addSlide(slideDirections[selectedOption]);
			selectedOption = '';
		}
	}

	function removeTag(tag: string) {
		tagList = tagList.filter((item) => item !== tag);
		options = initialOptions.filter((option) => !tagList.includes(option));
		removeSlide(slideDirections[tag]);
	}
</script>

<div class="flex space-x-4 mb-4">
	<select class="p-2 border rounded" bind:value={selectedOption}>
		<option value="" disabled>{dropDownText}</option>
		{#each options as option (option)}
			<option value={option}>{option}</option>
		{/each}
	</select>
	<button
		class="p-2 bg-blue-500 text-white rounded"
		on:click={addTag}
		disabled={selectedOption === ''}
	>
		Add
	</button>
</div>

<div class="flex flex-wrap space-x-2">
	{#each tagList as tag (tag)}
		<div class="bg-blue-500 text-white p-2 rounded flex items-center space-x-2 mb-2">
			<span>{tag}</span>
			<button class="text-white" on:click={() => removeTag(tag)}>x</button>
		</div>
	{/each}
</div>
