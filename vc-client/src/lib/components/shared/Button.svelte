<script lang="ts">
	export let transition = false
</script>

<button
	on:click
	on:mouseover
	on:focus
	on:mouseenter
	on:mouseleave
	class="btn"
	class:hover={transition}
>
	<slot />
</button>

<style>
	/*button {
		border: 0;
		padding: 8px 12px;
		box-shadow: 1px 1px 2px rgba(255, 255, 255, 0.8);
		height: 3rem;
		border-radius: 6px;
		cursor: pointer;
		transition: all 0.3s;
		max-width:max-content;
	}

	.hover:hover {
		transform: translateY(-10px);
		background-color: var(--primary-color);
	}

	.dark {
		background-color: var(--dark);
	}

	.dark.inverse {
		color: var(--dark);
	}

	.positive {
		background-color: var(--positive);
		color: white;
	}
	.negative {
		background-color: var(--negative);
		color: white;
	}

	.flat {
		box-shadow: none;
	}

	.primary {
		background-color: var(--primary-color);
		color: black;
	}
	.secondary {
		background-color: var(--secondary-color);
		color: black;
	}

	.primary.inverse {
		background-color: var(--dark);
		color: var(--primary-color);
		border: 2px solid var(--primary-color);
	}

	.negative.inverse {
		background-color: #1a1a2b;
		color: var(--negative);
		border: 2px solid var(--negative);
	}
	.positive.inverse {
		background-color: var(--dark);
		color: var(--positive);
		border: 2px solid var(--positive);
	}
	.secondary.inverse {
		background-color: var(--dark);
		color: var(--secondary-color);
		border: 2px solid var(--secondary-color);
	}*/
</style>
