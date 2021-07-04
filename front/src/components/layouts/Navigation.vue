<template>
  <div id="main-nav">
    <!-- <div class="w-full h-20 border-b flex px-4 items-center mb-8">
        <p class="font-semibold text-3xl text-blue-400 pl-4">LOGO</p>
    </div> -->

    <div class="mb-4 px-4 w-full">
        <p class="headers">Reports</p>
        <router-link class="items" :to="{name:'home'}">Dashboard</router-link>
        <router-link class="items" :to="{name:'docs'}">Documents</router-link>
    </div>

    <div class="mb-4 px-4 w-full">
        <p class="headers">Settings</p>
        <router-link class="items" :to="{name:'users-id',params:{id:'me'}}">Me</router-link>
        <router-link class="items" :to="{name:'users'}">Users</router-link>
    </div>

    <div class="mb-4 px-4 w-full">
        <div class="items" @click="Auth.logout()">Logout</div>
    </div>

    <div class="bottom-6 w-full">
      <style-switcher />
    </div>
  </div>
</template>

<script setup lang="ts">
import { Auth } from "../../modules/auth";
import StyleSwitcher from "./StyleSwitcher.vue";
</script>

<style lang="postcss" scoped>
.items {
  @apply w-full flex items-center h-10 pl-2 rounded-lg cursor-pointer mb-1;

  &:hover {
    background-color: var(--color-gray-400);
  }

  &.router-link-active {
    background-color: var(--color-gray-400);
  }
}

.headers {
  @apply text-sm font-semibold mb-1;
}
</style>