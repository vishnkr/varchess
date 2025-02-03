<script lang="ts">
    import { goto } from '$app/navigation';
    import { type ConnectParams, gameId } from '$lib/store/stores';
    import { wsStore } from '$lib/websocket';
    import PieChart from '$lib/components/charts/PieChart.svelte';
	import HalfDoughnut from '$lib/components/charts/HalfDoughnut.svelte';
	import { authStore } from '$lib/store/auth';
	import { onDestroy } from 'svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	
    let username: string = '';
    const unsubscribe = authStore.subscribe(state => {
        username = state.username || '';
    });

    onDestroy(() => {
        unsubscribe(); 
    });

    //export let data;
    //$: ({ username } = data);
    let gameIdInput: string;
    
    const createGame = () =>{
        goto('/editor')
    }
    /*export async function fetchDashboardData() {
    try {
        const response = await fetch('/api/dashboard');
        if (!response.ok) throw new Error('Failed to fetch dashboard data');
        return await response.json();
    } catch (error) {
        console.error('Error fetching dashboard data:', error);
        return null;
    }
    }

    const joinGame = () => {
        if (gameIdInput) {
            const url = `ws://${import.meta.env.VITE_WS_HOST}/ws`;
            const params: ConnectParams = {
                sessionId: 'sdf',
                gameId: gameIdInput,
                username: username
            };
            wsStore.newWebSocketConnection(url, params);
        }
    };
    $: {
        if ($gameId !== null && $wsStore) {
            goto(`/game/${$gameId}/waiting`);
        }
    }*/
</script>

<svelte:head>
    <title>Dashboard - Varchess</title>
</svelte:head>

<div class="flex flex-col md:flex-row h-screen">
    <!-- Left Section - Dashboard -->
    <div class="w-full md:w-2/3 p-4">
        <div class="bg-gray-900 p-6 rounded-lg shadow-lg">
            <h2 class="text-xl font-bold text-center text-white mb-4">Dashboard</h2>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div class="bg-gray-800 p-4 rounded-lg">
                    <PieChart />
                </div>
                <div class="bg-gray-800 p-4 rounded-lg">
                    <HalfDoughnut />
                </div>
            </div>
        </div>
    </div>

    <div class="w-full md:w-1/3 p-4">
        <div class="bg-gray-800 text-white p-6 rounded-lg flex flex-col items-center space-y-4 shadow-lg">
            <h2 class="text-xl font-bold mb-4">Play a Game</h2>
            <a href="/editor" class="w-full">
                <Button on:click={createGame} class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-2xl flex items-center w-full justify-center">
                    <i class="fa-solid fa-plus mr-2"></i> Create New Game
                </Button>
            </a>
            <span class="text-white">OR</span>
            <input
                type="text"
                name="gameId"
                bind:value={gameIdInput}
                placeholder="Enter Room Code"
                class="border border-gray-600 bg-gray-700 text-white px-4 py-2 rounded-xl w-full"
            />
            <Button
                on:click={()=>{}}
                class="bg-blue-600 hover:bg-blue-800 text-white font-bold py-2 px-4 rounded-2xl w-full"
            >
                Join Room
            </Button>
        </div>
    </div>
</div>
