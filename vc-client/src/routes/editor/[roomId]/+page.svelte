<script lang="ts">
	import BoardEditor from "$lib/components/editor/BoardEditor.svelte";
	import pieceSvg from "$lib/assets/svg/piece.svg";
    import boardSvg from "$lib/assets/svg/board.svg";
    import Tabs from "$lib/components/shared/Tabs.svelte";
	import type { BoardConfig } from "$lib/board/types";
	import EditableBoard from "$lib/board/EditableBoard.svelte";
	import PieceEditor from "$lib/components/editor/PieceEditor.svelte";
	import GameSettings from "$lib/components/editor/GameSettings.svelte";
    let items = ['Custom','Predefined'];
    let activeItem = 'Custom';
    let inputValue = "localhost:5sfds137";
    const tabChange=(e:CustomEvent<string>)=> activeItem = e.detail;
    
    export let boardConfig:BoardConfig = {
		fen: "rnbqkbnr/pnpppppp/p6p/5.../7P/P6P/PPPPPPPP/RNBQKBNR",
		dimensions: { ranks: 8, files: 8 },
		editable: false,
		interactive: false,
        isFlipped:false,
	}

    const copyToClipboard = () => {
        navigator.clipboard.writeText(inputValue);
    };
  
    let shiftBoard:(direction:string)=>void;
    function handleShift(event:CustomEvent<string>){
        const direction = event.detail;
        shiftBoard(direction);
    }

    let members = [
    { id: 1, username: "User1" },
    { id: 2, username: "User2" },
    { id: 3, username: "User3" },
    { id: 4, username: "User4" },
    { id: 5, username: "User5" },
    { id: 6, username: "User6" },
  ];
  const handleAction = (memberId:number, action:string) => {
    // Handle the action for the selected member
    console.log(`Performing ${action} on member with ID ${memberId}`);
  };

  let openDropdownId:number|null = null;

  const toggleDropdown = (memberId:number) => {
    openDropdownId = openDropdownId === memberId ? null : memberId;
  };
</script>

<svelte:head>
	<title> Editor - Varchess</title>
</svelte:head>
<div class="font-inter grid text-zinc-900">
    <div class="flex-1 flex m-4 lg:flex-row flex-col"> 
        <div class="bg-zinc-700 rounded-md lg:w-3/12 mx-3 p-3">
            <div class="border-b border-gray-200 dark:border-gray-700 flex flex-col text-center">
                <div class="flex justify-center">
                    <Tabs {activeItem} {items} on:tabChange={tabChange}/>
                </div>
                {#if activeItem==="Custom"}
                    <details class="bg-white shadow rounded group mb-4">
                        <summary class="list-none flat flex flex-wrap items-center cursor-pointer">
                            <h3 class="flex flex-1 p-4 lg:text-xl text-lg font-semibold">Board Editor  </h3>
                            <img class="w-g h-6"  src={boardSvg} alt="board editor" />
                        <div class="w-8 flex items-center justify-center">
                            <div class="border-8 border-transparent border-l-gray-600 ml-2
                            group-open:rotate-90 transition-transform origin-left" />
                        </div>
                    </summary>
                        <div><BoardEditor 
                            bind:dimensions={boardConfig.dimensions}
                            on:shift={handleShift}
                            /></div>
                    </details>
                    <details class="bg-white shadow rounded group mb-4">
                        <summary class="list-none flat flex flex-wrap items-center cursor-pointer">
                            <h3 class="flex flex-1 p-4 lg:text-xl text-lg font-semibold">Piece Editor </h3>
                            <img class="w-g h-6" src={pieceSvg} alt="piece editor" />
                        <div class="w-8 flex items-center justify-center">
                            <div class="border-8 border-transparent border-l-gray-600 ml-2
                            group-open:rotate-90 transition-transform origin-left" />
                        </div>
                    </summary>
                    <div><PieceEditor /></div>
                    </details>

                    <details class="bg-white shadow rounded group mb-4">
                        <summary class="list-none flat flex flex-wrap items-center cursor-pointer">
                            <h3 class="flex flex-1 p-4 lg:text-xl text-lg font-semibold">Game Settings </h3>
                            <i class="fa-solid fa-gear fa-lg"></i>
                        <div class="w-8 flex items-center justify-center">
                            <div class="border-8 border-transparent border-l-gray-600 ml-2
                            group-open:rotate-90 transition-transform origin-left" />
                        </div>
                    </summary>
                    <div> <GameSettings /></div>
                    </details>
                {:else}
                <details class="bg-white shadow rounded group mb-4">
                    <summary class="list-none flat flex flex-wrap items-center cursor-pointer">
                        <h3 class="flex flex-1 p-4 font-semibold">Duck Chess</h3>
                    <div class="w-8 flex items-center justify-center">
                        <div class="border-8 border-transparent border-l-gray-600 ml-2
                        group-open:rotate-90 transition-transform origin-left"></div>
                    </div>
                </summary>
                <div><p>sfdg</p></div>
                </details>
                {/if}
            </div>
        </div>
        <div class="bg-zinc-700 rounded-md lg:w-6/12 mx-3 my-3 p-3">
            
            <EditableBoard 
            boardConfig={boardConfig}
            bind:shift={shiftBoard}
            />
        </div>
        <div class="flex flex-col bg-zinc-700 rounded-md lg:w-3/12 mx-3 p-3">
            <div class="flex mb-5">
                <input
                  class="border border-gray-300 px-4 py-2 text-white rounded-l w-64"
                  bind:value={inputValue}
                  disabled
                />
                <button
                  class="bg-blue-500 hover:bg-blue-700 text-white font-semibold py-2 px-2 rounded-r"
                  on:click={copyToClipboard}
                >
                Share <i class="fa-solid fa-link"></i>
                </button>
            </div>
            <div class="bg-zinc-900 flex flex-col rounded">
                <h4 class="font-semibold text-center text-xl text-white">Members</h4>
                <div class="max-h-60 overflow-y-auto bg-white">
                    {#each members as member (member.id)}
                      <div class="flex justify-between items-center py-2 px-4 border-b border-gray-300">
                        <span>{member.username}</span>
                        <div class="dropdown inline-block relative">
                          <button class="bg-gray-200 text-gray-600 py-1 px-2 rounded-md hover:bg-gray-300 focus:outline-none" on:click={() => toggleDropdown(member.id)}>
                            Actions
                          </button>
                          {#if openDropdownId === member.id}
                            <ul class="absolute z-10 bg-white border border-gray-300 rounded-md shadow mt-2 w-40">
                              <li>
                                <button class="hover:bg-gray-100 px-4 py-2 w-full text-left" on:click={() => handleAction(member.id, 'Kick')}>
                                  Remove
                                </button>
                              </li>
                              <li>
                                <button class="hover:bg-gray-100 px-4 py-2 w-full text-left" on:click={() => handleAction(member.id, 'Ban')}>
                                  Set as White
                                </button>
                              </li>
                              <li>
                                <button class="hover:bg-gray-100 px-4 py-2 w-full text-left" on:click={() => handleAction(member.id, 'Ban')}>
                                  Set as Black
                                </button>
                              </li>
                            </ul>
                          {/if}
                        </div>
                      </div>
                    {/each}
                  </div>
            </div>
        </div>
    </div>
    
</div>