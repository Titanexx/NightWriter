<template>
    <div class="flex items-center rounded w-full bg-secondary p-2">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-gray-600 mr-2" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd" />
        </svg>
        <input ref="searchBarInput" class="w-full bg-secondary focus-visible:outline-none" placeholder="Search" v-model="searchValue">
    </div>
</template>

<script lang="ts">
import { defineComponent, ref, watch } from "vue";

export default defineComponent({
    props: ['data'],
    setup(props, { emit }) {
        const searchValue = ref<string>("")
        const searchBarInput = ref();
        watch(searchValue, ()=>{
            if(searchValue.value != ""){
                emit("found",props.data.filter((e) => e.search(searchValue.value.toLowerCase())))
            } else {
                emit("found",props.data)
            }
        })

        watch(()=>props.data,(data)=>{
            searchValue.value = ""
            emit("found",data)
        })

        emit("found",props.data)

        document.addEventListener('keyup', (event) => {
        if (document.activeElement == document.body) {
            if (event.key == 's'){
                searchBarInput.value.focus()
            }
        }
        }, false);


        return {
            searchValue,
            searchBarInput,
        };
    },
});
</script>

<style lang="postcss" scoped>
</style>