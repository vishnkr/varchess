<script>
	import { onMount } from 'svelte';
	import * as echarts from 'echarts';

	let chartContainer;

	onMount(() => {
		const myChart = echarts.init(chartContainer);
		const option = {
            title: {
			    text: 'Games By Result',
                left: 'center',
                textStyle: {
                    color: '#ffffff'
                }
            },
			tooltip: {
				trigger: 'item'
			},
			legend: {
				bottom: 0,
				left: 'center',
				textStyle: {
					color: '#ffffff'
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
						{ value: 100, name: 'Win' ,itemStyle: { color: 'green' }},
                        { value: 5, name: 'Draw' ,itemStyle: { color: 'gray' }},
						{ value: 73, name: 'Loss' ,itemStyle: { color: 'red' }}
					]
				}
			]
		};

		myChart.setOption(option);
		window.addEventListener('resize', () => {
			myChart.resize();
		});

		return () => {
			myChart.dispose();
		};
	});
</script>

<div bind:this={chartContainer} class="chart-container" />

<style>
	.chart-container {
		width: 100%;
		height: 400px;
	}
</style>