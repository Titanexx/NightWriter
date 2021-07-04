<template>
    <Modal v-model:isOpen="isOpen">
        <template #title>
            {{ props.title }}
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