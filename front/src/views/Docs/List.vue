<template>
    <div class="w-full px-2 h-full flex flex-col">
        <div class="bg-secondary rounded-xl p-2 rounded-t-none mb-2">
            <SearchBar :data="docs" @found="(data) => filteredDocs = data" />
        </div>
        <div class="flex flex-row flex-wrap gap-2">
            <router-link class="rounded-xl flex flex-col bg-secondary w-40 p-3 pb-6 cursor-pointer hover:bg-gray-400" :to="{name:'docs-new'}">
                <svg xmlns="http://www.w3.org/2000/svg" class="rounded-full bg-gray-100 w-16 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                </svg>
                <div class="flex-1 mx-auto mt-3 text-center align-middle">
                    <span>Create a new document</span>
                </div>
            </router-link>
            
            <div class="rounded-xl flex flex-col bg-secondary w-40 p-3 pb-6" v-if="loading">
                <svg xmlns="http://www.w3.org/2000/svg" class="bg-gray-100 w-16 mx-auto rounded-full animate-spin p-1" style="animation-direction: reverse;" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                <div class="mx-auto mt-3 text-center items-center">
                    <span>Loading documents</span>
                </div>
            </div> 

            <router-link v-for="doc in filteredDocs" :to="{name:'docs-id',params:{id:doc.ID}}" >
                <DocCard :doc="doc"/>
            </router-link>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

import { useApiToken } from '@modules/api';
import { Doc } from '@modules/models';
import DocCard from '@components/doc/DocCard.vue';
import SearchBar from "@components/ui/SearchBar.vue"
import { showError } from '@plugins/modal';

const {
    loading,
    get,
    data,
    error,
} = useApiToken("/api/docs/");

const docs = ref<Doc[]>()
const filteredDocs = ref<Doc[]>([])

get().then(() => {
    docs.value = []
    if(error.value!=undefined){
        console.log(error)
        showError('Error','Contact your administrator');
    } else {
        data.value.forEach(element => {
            docs.value.push(new Doc(element))
        });
    }
});
</script>

<style lang="scss" scoped>
</style>