<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import * as echarts from 'echarts';
	import { browser } from '$app/environment';

	export let theme = 'light';
	export let data: { name: string; value: number }[] = [];
	export let title = 'Games by variant';

	let chartContainer: HTMLDivElement;
	let chartInstance: echarts.ECharts | null = null;
	let lastOptionKey = '';

	function getPieChartOptions(currentTheme: string, seriesData: { name: string; value: number }[]) {
		const isDark = currentTheme === 'dark';
		const textColor = isDark ? '#ffffff' : '#1f2937';
		const empty = !seriesData.length || seriesData.every((d) => d.value === 0);

		return {
			title: {
				text: title,
				left: 'center',
				textStyle: { color: textColor, fontSize: 14 }
			},
			tooltip: { trigger: 'item' },
			legend: {
				orient: 'horizontal',
				left: 'center',
				bottom: 0,
				type: 'scroll',
				textStyle: { color: textColor },
				pageTextStyle: { color: textColor }
			},
			series: [
				{
					name: 'Variant',
					type: 'pie',
					radius: ['45%', '68%'],
					center: ['50%', '46%'],
					labelLine: { show: false },
					label: { show: false },
					data: empty
						? [{ value: 1, name: 'No games yet', itemStyle: { color: '#64748b' } }]
						: seriesData
				}
			]
		};
	}

	function applyOptions() {
		if (!chartInstance) return;
		const key = `${theme}|${title}|${JSON.stringify(data)}`;
		if (key === lastOptionKey) return;
		lastOptionKey = key;
		chartInstance.setOption(getPieChartOptions(theme, data), true);
	}

	function resizeChart() {
		chartInstance?.resize();
	}

	onMount(() => {
		chartInstance = echarts.init(chartContainer);
		applyOptions();
		if (browser) window.addEventListener('resize', resizeChart);
	});

	$: theme, data, title, applyOptions();

	onDestroy(() => {
		chartInstance?.dispose();
		chartInstance = null;
		if (browser) window.removeEventListener('resize', resizeChart);
	});
</script>

<div bind:this={chartContainer} class="chart-container" />

<style>
	.chart-container {
		width: 100%;
		height: 280px;
	}
</style>
