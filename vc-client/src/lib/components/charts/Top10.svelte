<script lang="ts">
  import * as Card from "$lib/components/ui/card/index.js";
  import * as Chart from "$lib/components/ui/chart/index.js";
  import { scaleBand } from "d3-scale";
  import { BarChart, type ChartContextValue } from "layerchart";
  import { cubicInOut } from "svelte/easing";
  import CrownIcon from "@lucide/svelte/icons/crown";
  import SkullIcon from "@lucide/svelte/icons/skull";

  const chartData = [
    { variant: "Standard", wins: 120, losses: 60 },
    { variant: "Antichess", wins: 90, losses: 30 },
    { variant: "KnightStorm", wins: 55, losses: 25 },
    { variant: "Minichess", wins: 42, losses: 18 },
    { variant: "Chaos960", wins: 34, losses: 21 },
    { variant: "Others", wins: 67, losses: 48 },
  ];

  const chartConfig = {
    wins: { label: "Wins", color: "var(--chart-1)", icon: CrownIcon },
    losses: { label: "Losses", color: "var(--chart-2)", icon: SkullIcon },
  } satisfies Chart.ChartConfig;

  let context = $state<ChartContextValue>();
</script>

<Card.Root>
  <Card.Header>
    <Card.Title>Games by Variant</Card.Title>
    <Card.Description>Stacked wins and losses for your top variants.</Card.Description>
  </Card.Header>
  <Card.Content>
    <Chart.Container config={chartConfig}>
      <BarChart
        bind:context
        data={chartData}
        xScale={scaleBand().padding(0.25)}
        x="variant"
        axis="x"
        rule={false}
        seriesLayout="stack"
        grid={false}
        highlight={false}
        series={[
          {
            key: "wins",
            label: "Wins",
            color: chartConfig.wins.color,
            props: { rounded: "bottom" },
          },
          {
            key: "losses",
            label: "Losses",
            color: chartConfig.losses.color,
          },
        ]}
        props={{
          bars: {
            stroke: "none",
            initialY: context?.height,
            initialHeight: 0,
            motion: {
              y: { type: "tween", duration: 500, easing: cubicInOut },
              height: { type: "tween", duration: 500, easing: cubicInOut },
            },
          },
          xAxis: {
            format: (d) => d,
            tickLabelProps: {
              svgProps: {
                y: 13,
              },
            },
          },
        }}
      >
          <Chart.Tooltip hideLabel />
      </BarChart>
    </Chart.Container>
  </Card.Content>
</Card.Root>
