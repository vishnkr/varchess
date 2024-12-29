<script lang="ts">
    import VariantRulesModal from "./VariantRulesModal.svelte";
</script>

<VariantRulesModal>
    <div class="p-4">
      <h2 class="text-xl font-semibold mb-2">Antichess Rules</h2>
      <p class="text-gray-800">
        Antichess, also known as Losing Chess, is a variant where the goal is to lose all of your pieces or to be stalemated.
      </p>
      <h2 class="text-lg font-semibold mt-4 mb-2">Key Differences:</h2>
      <ul class="list-disc pl-6">
        <li class="text-gray-800">The player who loses all their pieces or is stalemated wins the game.</li>
        <li class="text-gray-800">Capturing is mandatory if possible.</li>
        <li class="text-gray-800">Pieces move as in standard chess, but the strategy is reversed—players aim to lose their pieces instead of protecting them.</li>
        <li class="text-gray-800">The king can be captured like any other piece, and there is no check or checkmate.</li>
      </ul>
    </div>
</VariantRulesModal>