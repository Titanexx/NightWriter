<template>
    <div class="flex flex-col h-full">
        <div class="mb-2 flex flex-row">
            <div class="flex-1">
                <SearchInput :data="doc.parts" @found="(partsFound) => filteredParts = partsFound"></SearchInput>
            </div>
            <button v-if="doc.checkRight(doc.Rights.Writer)" class="btn-primary btn-text-animation flex-none ml-2 px-2 h-10 rounded flex justify-center items-center" :disabled="loading.value" @click="openModal">
                <span>Add Part</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v3m0 0v3m0-3h3m-3 0H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
            </button>

            <button class="btn-primary btn-text-animation flex-none ml-2 px-2 h-10 rounded flex justify-center items-center" :disabled="loading.value" @click="refreshParts() && emit('select',null)">
                <span>Refresh</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6"  :class="{'animate-spin':loading.value}" style="animation-direction: reverse;" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
            </button>
        </div>

        <div class="flex-auto overflow-y-auto min-h-0 rounded rounded-b-none bg-secondary p-2  -pt-1 -pb-2 ">
            <draggable 
                v-model="filteredParts"
                :disabled="isDraggable"
                group="part"
                @end="checkEndPart" 
                item-key="ID">
                <template #item="{element}">
                    <PartItem @click="emit('select',element)" :is-draggable="isDraggable" :part="element"/>
                </template>
            </draggable>
        </div>
    </div>
    

    <Modal v-if="doc.checkRight(doc.Rights.Writer)" v-model:isOpen="isOpen">
        <template #title>
            Create part
        </template>
        <template #body>            
            <form @submit.prevent="CreatePart" class="">
                <Input type="text" label="Title" v-model="partTitle" />
                <button class="btn-primary flex-none mt-5 px-2 h-10 rounded flex justify-center items-center" :disabled="loadingParts">
                    <span>Add Part</span>
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 ml-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v3m0 0v3m0-3h3m-3 0H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                </button>
            </form>
        </template>
    </Modal>
</template>


<script setup lang="ts">
import { ref, computed } from "vue";

import { Doc, Part } from "@modules/models";
import { useApiToken } from '@modules/api';

import Input from '@components/ui/Input.vue';
import Modal from '@components/ui/Modal.vue';
import SearchInput from "@components/ui/SearchBar.vue";
import PartItem from "@components/part/PartItem.vue";
import draggable from 'vuedraggable'

const props = defineProps({
    doc: Doc,
})
const emit = defineEmits<{
  (event: 'select', part: Part): void
}>()

const {
    loading: loadingParts,
    put: putParts,
} = useApiToken("/api/docs/"+props.doc.ID+"/parts")

const {
    loading: loadingRefreshParts,
    refreshParts : refreshParts
} = props.doc.refreshParts()

const partTitle = ref<string>("");
const isOpen = ref(false);
const filteredParts = ref<Part[]>([])
const draggingPart = ref<Part>()
const loading = computed(()=> loadingParts || loadingRefreshParts)
const isDraggable = computed(() => filteredParts.value.length != props.doc.parts.length)

function checkEndPart(event){
    if(event.oldIndex != event.newIndex){
        props.doc.movePart(event.oldIndex, event.newIndex)
    }
}

function closeModal() {
    isOpen.value = false
}
function openModal() {
    isOpen.value = true
}

function CreatePart(){
    let newPart = new Part(props.doc)
    newPart.title = partTitle.value
    props.doc.addPart(newPart)

    putParts(newPart.payload).then((partID)=>{
        newPart.ID = partID
        closeModal()
    })
}

document.addEventListener('keyup', (event) => {
  if (document.activeElement == document.body) {
    if (event.key == 'a'){
        openModal()
    } else if (event.key == 'r'){
        refreshParts()
    }
  }
}, false);


</script>

<style lang="scss">

</style>