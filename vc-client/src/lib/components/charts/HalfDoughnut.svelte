<script>
	import { onMount, onDestroy } from 'svelte';
	import * as echarts from 'echarts';
	import { browser } from '$app/environment';

	export let theme;
	let chartContainer;
	let chartInstance;

	function getChartOptions(theme) {
		return {
			title: {
				text: 'Games By Result',
				left: 'center',
				textStyle: {
					color: theme === 'dark' ? '#ffffff' : '#000000'
				}
			},
			tooltip: {
				trigger: 'item'
			},
			legend: {
				bottom: 0,
				left: 'center',
				textStyle: {
					color: theme === 'dark' ? '#ffffff' : '#000000'
				}
			},
			series: [
				{
					name: 'Result',
					type: 'pie',
					radius: ['40%', '70%'],
					center: ['50%', '70%'],
					labelLine: {
						show: false
					},
					label: {
						show: false,
						position: 'center'
					},
					startAngle: 180,
					endAngle: 360,
					data: [
						{ value: 100, name: 'Win', itemStyle: { color: 'green' }},
						{ value: 5, name: 'Draw', itemStyle: { color: 'gray' }},
						{ value: 73, name: 'Loss', itemStyle: { color: 'red' }}
					]
				}
			]
		};
	}

	function resizeChart() {
		if (chartInstance) {
			chartInstance.resize();
		}
	}

	onMount(() => {
		chartInstance = echarts.init(chartContainer);
		chartInstance.setOption(getChartOptions(theme));
		if (browser) {
			window.addEventListener('resize', resizeChart);
		}
	});

	$: if (chartInstance) {
		chartInstance.setOption(getChartOptions(theme), true);
	}

	onDestroy(() => {
		if (chartInstance) chartInstance.dispose();
		if (browser) window.removeEventListener('resize', resizeChart);
	});
</script>

<div bind:this={chartContainer} class="chart-container" />

<style>
	.chart-container {
		width: 100%;
		height: 400px;
	}
</style>
