<script lang="ts">
    import VariantRulesModal from "./VariantRulesModal.svelte";
</script>

<VariantRulesModal>
    <div class="p-4">
        <h2 class="text-xl font-semibold mb-2">N-Check Rules</h2>
        <p class="text-gray-800">
            In N-Check, the objective is to check your opponent's king N times.
        </p>
        <ul class="list-disc pl-6">
            <li class="text-gray-800">Each time a player checks their opponent's king, it counts towards the N checks required to win.</li>
            <li class="text-gray-800">Players can still win by checkmating their opponent's king, but the primary objective is to accumulate three checks.</li>
            <li class="text-gray-800">All other standard chess rules apply.</li>
        </ul>
    </div>
</VariantRulesModal>
 