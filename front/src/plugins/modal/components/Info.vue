<template>
    <Modal v-model:isOpen="isOpen">
        <template #title>
            <span class="text-green-500">
                <svg xmlns="http://www.w3.org/2000/svg" class="inline-flex self-center h-5 w-5 " fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {{ props.title }}
            </span>
        </template>
        <template #body>         
            {{ props.body }}
        </template>
    </Modal>
</template>

<script setup lang="ts">
import Modal from "@components/ui/Modal.vue";
import { ref, watch } from "vue";

const isOpen = ref(true)

const props = defineProps({
    title: String,
    body: String,
    time: {type: Number, default: 2000},
    close: Function,
})

setTimeout(() => isOpen.value = false, props.time)

watch(isOpen, (isOpen) => {if(!isOpen)props.close()})

</script>