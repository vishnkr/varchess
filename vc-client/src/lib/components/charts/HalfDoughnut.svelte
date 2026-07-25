<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import * as echarts from 'echarts';
	import { browser } from '$app/environment';

	export let theme = 'light';
	export let data: { name: string; value: number; color?: string }[] = [];
	export let title = 'Games by result';

	let chartContainer: HTMLDivElement;
	let chartInstance: echarts.ECharts | null = null;
	let lastOptionKey = '';

	function getChartOptions(
		currentTheme: string,
		seriesData: { name: string; value: number; color?: string }[]
	) {
		const textColor = currentTheme === 'dark' ? '#ffffff' : '#000000';
		const empty = !seriesData.length || seriesData.every((d) => d.value === 0);
		const chartData = empty
			? [{ value: 1, name: 'No games yet', itemStyle: { color: '#64748b' } }]
			: seriesData.map((d) => ({
					name: d.name,
					value: d.value,
					itemStyle: d.color ? { color: d.color } : undefined
				}));

		return {
			title: {
				text: title,
				left: 'center',
				textStyle: { color: textColor, fontSize: 14 }
			},
			tooltip: { trigger: 'item' },
			legend: {
				bottom: 0,
				left: 'center',
				textStyle: { color: textColor }
			},
			series: [
				{
					name: 'Result',
					type: 'pie',
					radius: ['40%', '70%'],
					center: ['50%', '65%'],
					labelLine: { show: false },
					label: { show: false },
					startAngle: 180,
					endAngle: 360,
					data: chartData
				}
			]
		};
	}

	function applyOptions() {
		if (!chartInstance) return;
		const key = `${theme}|${title}|${JSON.stringify(data)}`;
		if (key === lastOptionKey) return;
		lastOptionKey = key;
		chartInstance.setOption(getChartOptions(theme, data), true);
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
