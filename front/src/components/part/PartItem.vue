<template>
    <div class="flex flex-row bg-primary p-2 rounded items-center my-1" :style="{'marginLeft': `${10*part.level}px`}">
        <span class="flex-1">{{ part.title }}</span>
        <div v-if="part.doc.checkRight(part.doc.Rights.Writer)" class="flex-none flex flex-row">
            <button class="flex-1 btn-primary ml-2 px-2 h-6 rounded flex justify-center items-center" :disabled="isDraggable">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l4-4 4 4m0 6l-4 4-4-4" />
                </svg>
            </button>
            <button class="flex-1 btn-primary ml-2 px-2 h-6 rounded flex justify-center items-center" :disabled="loading || part.level == 0" @click.stop="part.levelDown()">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 17l-5-5m0 0l5-5m-5 5h12" />
                </svg>
            </button>
            <button class="flex-1 btn-primary ml-2 px-2 h-6 rounded flex justify-center items-center" :disabled="loading || part.level == 3" @click.stop="part.levelUp()">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" /> 
                </svg>
            </button>
    
            <button class="flex-1 btn-red ml-2 px-2 h-6 rounded flex justify-center items-center" :disabled="loading" @click.stop="part.silentDelete()">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
            </button>
        </div>
    </div>
</template>


<script setup lang="ts">
import { Part } from "@modules/models";

const props = defineProps({
    loading: Boolean,
    isDraggable: {type: Boolean, default: true},
    part: Part,
})
</script>

<style lang="scss">
</style>