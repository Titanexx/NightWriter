<template>
    <div class="mt-4">
        <p>{{ user }}</p>
    </div>
</template>

<script setup lang="ts">
import {
    defineProps,
} from "vue";

defineProps(['user'])
</script>

<style>
</style>