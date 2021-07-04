<template>
    <div v-if='Auth.isAuthenticated.value' class="h-full">
        <navigation class="w-32 h-full fixed top-0 left-0 bg-secondary z-10 text-center flex flex-wrap content-center"/>
        <div class="pl-32 w-full h-full flex flex-col">
            <router-view />
        </div>
    </div>
    <div v-else class="h-screen flex justify-center items-center">
        <div class="flex justify-center items-center pb-7">
            <router-view />
        </div>
        <div class="absolute bottom-4">
            <StyleSwitcher />
        </div>
    </div>
</template>

<script setup lang="ts">
import Navigation from "@layouts/Navigation.vue";
import { Auth } from "@modules/auth";
import StyleSwitcher from "@components/layouts/StyleSwitcher.vue";

</script>

<style>
</style>