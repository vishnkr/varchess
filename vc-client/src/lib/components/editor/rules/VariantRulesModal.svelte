<script lang="ts">
	import { ruleEditor } from '$lib/store/editor';
  let isOpen: boolean;
  $: {
    isOpen = $ruleEditor.isViewVariantRulesOn;
  }
  
  function closeModal() {
    $ruleEditor.isViewVariantRulesOn = false;
    
  }

  </script>
  
  {#if isOpen}
  <div class="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50 modal-overlay">
    <div class="bg-white rounded-lg p-6 w-full sm:w-96 relative">
      <button class="absolute top-0 right-0 mt-2 mr-2 text-gray-700" on:click={closeModal}>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
      <div class="overflow-y-auto max-h-96">
        <slot />
      </div>
    </div>
  </div>
{/if}
  