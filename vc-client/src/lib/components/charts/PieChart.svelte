<script>
	import { onMount, onDestroy } from 'svelte';
	import * as echarts from 'echarts';
	import { browser } from '$app/environment';

	export let theme = 'light'; // 👈 allow parent to pass theme

	let chartContainer;
	let chartInstance;

	function getPieChartOptions(theme) {
		const isDark = theme === 'dark';
		const textColor = isDark ? '#ffffff' : '#1f2937';

		return {
			title: {
				text: 'Games By Variant',
				left: 'center',
				textStyle: { color: textColor }
			},
			tooltip: {
				trigger: 'item'
			},
			legend: {
				orient: 'horizontal',
				left: 'left',
				padding: 0,
				bottom: 0,
				type: 'scroll',
				textStyle: { color: textColor },
				pageTextStyle: { color: textColor }
			},
			series: [
				{
					name: 'Variant',
					type: 'pie',
					radius: ['50%', '70%'],
					labelLine: { show: false },
					data: [
						{ value: 108, name: 'Checkmate' },
						{ value: 75, name: 'Antichess' },
						{ value: 50, name: 'Wormhole' },
						{ value: 44, name: 'Duck chess' },
						{ value: 30, name: 'wwwwwwwwwwww' },
						{ value: 50, name: 'Woole' },
						{ value: 44, name: 'Dchess' },
						{ value: 30, name: 'wwwwwww' },
						{ value: 50, name: 'Wle' },
						{ value: 44, name: 'Duck ch1ess' }
					],
					label: {
						show: false,
						position: 'center'
					}
				}
			]
		};
	}

	function resizeChart() {
		if (chartInstance) chartInstance.resize();
	}

	onMount(() => {
		chartInstance = echarts.init(chartContainer);
		chartInstance.setOption(getPieChartOptions(theme));
		if (browser) {
			window.addEventListener('resize', resizeChart);
		}
	});

	$: if (chartInstance) {
		chartInstance.setOption(getPieChartOptions(theme), true);
	}

	onDestroy(() => {
		if (chartInstance) chartInstance.dispose();
		if (browser) {
			window.removeEventListener('resize', resizeChart);
		}
	});
</script>


<div bind:this={chartContainer} class="chart-container" />

<style>
	.chart-container {
		width: 100%;
		height: 400px;
	}
</style>
