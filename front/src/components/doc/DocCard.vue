<template>
    <div class="rounded-xl flex flex-col bg-secondary w-40 p-3 h-full pb-6 cursor-pointer hover:bg-gray-400">
        <svg xmlns="http://www.w3.org/2000/svg" class="rounded-xl bg-gray-100 w-16 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <div class="mx-auto mt-3 text-center">
            <span>{{ doc.title }}</span>
        </div>
    </div>
</template>

<script setup lang="ts">
import { Doc } from "@modules/models";

const props = defineProps({
    doc: Object as () => Doc
})

</script>

<style>
</style>

script