<template>
    <div v-if="!docIsLoaded">
    </div>
    <div class="pl-2 flex flex-col h-full " v-else>
        <div class="flex-none bg-secondary rounded text-center flex justify-between p-2 rounded-t-none">
            <form  @submit.prevent="updateDoc" @keydown.ctrl.s.prevent="updateDoc" class="flex-1 flex">
                <input ref="docTitleFormInput" type="text" placeholder="Title" class="flex-1 ml-2 h-10 rounded-lg bg-primary" v-model="doc.title" :disabled="!doc.checkRight(doc.Rights.Writer)"/>
                <button v-if="doc.checkRight(doc.Rights.Writer)" type="submit" class="btn-primary btn-text-animation flex-none ml-2 px-2 h-10 rounded-lg flex justify-center items-center" :disabled="loading">
                    <span>Update title</span>
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4" />
                    </svg>
                </button>
            </form>

            <button 
                v-if="doc.checkRight(doc.Rights.Editor)"
                class="btn-primary btn-text-animation flex-none ml-2 px-2 h-10 rounded-lg flex justify-center items-center"
                :disabled="loading"
                @click="manageAccess = !manageAccess"
                v-show="manageAccess">
                <span>Go to the doc</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
            </button>
            <button class="btn-primary btn-text-animation flex-none ml-2 px-2 h-10 rounded-lg flex justify-center items-center"
                :disabled="loading"
                @click="manageAccess = !manageAccess"
                v-if="doc.checkRight(doc.Rights.Editor)"
                v-show="!manageAccess">
                <span>Manage access</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
            </button>

            <button class="btn-green btn-text-animation flex-none ml-2 px-2 h-10 rounded-lg flex justify-center items-center" :disabled="loading" @click="exportDoc">
                <span>Export</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
            </button>

            <button v-if="doc.checkRight(doc.Rights.Editor)" class="btn-red btn-text-animation flex-none ml-2 px-2 h-10 rounded-lg flex justify-center items-center" :disabled="loading" @click="delIsOpen = true">
                <span>Delete</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
            </button>

            <RawModal v-model:isOpen="delIsOpen">
                <template #title>
                    Delete this document ?
                </template>
                <template #body>       
                    <span>Are you sure you want to delete this document?</span>     
                    <div class="flex justify-between">
                        <button class="btn-primary flex-none mt-5 px-2 h-10 rounded flex justify-center items-center" :disabled="loading" @click="delDoc">
                            <span>Delete it</span>
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 ml-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                            </svg>
                        </button>     
                        <button class="btn-red flex-none mt-5 px-2 h-10 rounded flex justify-center items-center" :disabled="loading" @click="delIsOpen = false">
                            <span>Cancel</span>
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 ml-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
                            </svg>
                        </button>
                    </div>
                </template>
            </RawModal>
		</div>
        <div class="flex-auto min-h-0 flex flex-col h-full">
            <Accesses :doc="doc" v-if="manageAccess"/>
            <div class="flex-auto overflow-y-auto flex flex-row" v-else>
                <div class="flex-initial w-4/12 mt-2 mr-2">
                    <PartList :doc="doc" @select="select"/>
                </div>
                
                <div class="flex-auto mt-2 flex flex-col" v-if="part">
                    <div class="flex-none bg-secondary rounded text-center flex justify-between p-2 rounded-b-none">
                        <form @submit.prevent="partSilentSave(part.titlePayload)" @keydown.ctrl.s.prevent="partSilentSave(part.titlePayload)" class="flex-1 flex">
                            <input type="text" placeholder="Title" class="flex-1 ml-2 h-10 rounded-lg bg-primary" v-model="part.title" :disabled="doc.checkRight(doc.Rights.Writer) || loading"/>
                            <button type="submit" class="btn-primary btn-text-animation flex-none ml-2 px-2 h-10 rounded-lg flex justify-center items-center" :disabled="loading" v-if="doc.checkRight(doc.Rights.Writer)">
                                <span>Update title</span>
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4" />
                                </svg>
                            </button>
                        </form>
                    </div>
                    <div class="flex-grow">
                        <Editor
                            v-if="isWsLoaded"
                            v-model="part.content"
                            @uploadAttachment="uploadAttachment"
                            @save="partSilentSave(part.contentPayload)"
                            @loadImage="loadImage"
                            :cwsID="cwsID"
                            :viewer="!doc.checkRight(doc.Rights.Writer)"/>
                    </div>
                </div>
            </div>
        </div>
    </div>
    

</template>

<script setup lang="ts">

import { ref, watch } from 'vue';
import { useRoute } from 'vue-router';

import router from '@routes';
import { useApiToken } from '@modules/api';
import { Doc, Part } from '@modules/models';
import { CollabWS } from '@modules/collaboration';
import RawModal from '@components/ui/Modal.vue';
import Accesses from '@components/doc/DocAccesses.vue';
import PartList from '@components/part/PartList.vue'; 
import Editor from '@components/ui/Editor.vue'; 
import { showError } from '@plugins/modal';

const route = useRoute()
const doc = ref<Doc>()
const part = ref<Part>()
const docTitleFormInput = ref<HTMLElement>()
const manageAccess = ref(false)
const docIsLoaded = ref(false)
const isWsLoaded = ref(false)
const delIsOpen = ref(false)
var cws: CollabWS
const cwsID = ref(-1)

function select(e){
    part.value = e
}

const {
    loading,
    data,
    get, del, post,
    error
} = useApiToken("/api/docs/"+route.params.id)

get().then(() => {
    if(error.value!=undefined){
        console.log(error)
        showError("Unauthorized","You are not authorized to delete this document.")
    } else {
        doc.value = new Doc(data.value)
        docIsLoaded.value = true
    }
})

function delDoc(){
    if(doc.value.checkRight(Doc.Rights.Editor)){
        del().then(()=>{
            if(error.value!=undefined){
                showError("Unauthorized","You are not authorized to delete this document.")
            } else {
                delIsOpen.value = false
                router.push({'name':'docs'})
            }
        })
    } else {
        showError("Unauthorized","You are not authorized to delete this document.")
    }
}

function updateDoc(){
    if(doc.value.checkRight(Doc.Rights.Writer)){
        post(doc.value.payload).then(()=>{
            if(error.value!=undefined){
                if(error.value.response.status == 401){
                    showError("Unauthorized","You are not authorized to update this document.")
                } else {
                    showError("Error","Your data isn't valid.")
                }
            }
        })
    } else {
        showError("Unauthorized","You are not authorized to update this document.")
    }
}

function uploadAttachment(attachmentBlob:File, callback:Function){
    if(doc.value.right >= Doc.Rights.Writer){
        const {addAttachment} = part.value.addAttachment(attachmentBlob)
        addAttachment().then((data)=>{
            console.log(data)
            callback(data,attachmentBlob.name)
        })
    } else {
        showError("Unauthorized","You are not authorized to upload an attachment.")
    }
}

function loadImage(url: string, node: HTMLImageElement){
    if(node.src == ""){
        part.value.getAttachment(url, (data)=>{
            var urlCreator = window.URL || window.webkitURL;
            var imageUrl = urlCreator.createObjectURL(data);
            node.src = imageUrl;
            data.arrayBuffer().then(d => {node.src = `data:img;base64,${btoa(String.fromCharCode(...new Uint8Array(d)))}`})
        })
    }
}

function exportDoc(){
    doc.value.export()
}

function partSilentSave(payload){
    if(doc.value.right >= Doc.Rights.Writer){
        part.value.silentSave(payload)
    }
}

watch(part, (_part)=>{
    if(doc.value.checkRight(Doc.Rights.Writer)){
        cws?.close()
        isWsLoaded.value=false
        if(_part != null){
            cws = new CollabWS(_part.wsUrl, (u)=>doc.value.encryptUpdate(u), (u)=>doc.value.decryptUpdate(u),()=>{cwsID.value = cws.id;isWsLoaded.value=true})
            cws.init()
            cwsID.value = cws.id
        }
    }
})

document.addEventListener('keyup', (event) => {
  if (document.activeElement == document.body) {
    if (event.key == 'm'){
        manageAccess.value = !manageAccess.value
    } else if (event.key == 't'){
        docTitleFormInput.value.focus()
    }
  }
}, false);
</script>

<style lang="scss" scoped>
</style>