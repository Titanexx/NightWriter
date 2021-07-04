<template>
	<div class="w-full flex justify-center content-center items-center h-screen">
		<div class="mx-auto bg-secondary w-6/12 rounded-lg text-center flex justify-center items-center">
			<form @submit.prevent="submit" class="w-10/12 p-8">
				<h1 class="text-2xl font-bold mt-1 mb-4">Create a new document.</h1>

				<Input type="text" label="Title" placeholder="Title" name="title" v-model="title" required/>
			
				<div class="flex w-full mt-8">
					<button type="submit" :disabled="loadingDoc" class="btn-primary w-full">
						Create
					</button>
				</div>
			</form>
		</div>

	</div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

import router from '@routes';

import { useApiToken } from "@modules/api";
import { Crypto } from "@modules/crypto";
import { Auth } from "@modules/auth";

import Input from "@components/ui/Input.vue";
import { showError, showInfo } from '@plugins/modal';
const title = ref()

const {
    error: errorDoc,
    loading: loadingDoc,
    put: putDoc,
} = useApiToken("/api/docs/");

function submit(){
    const encTitle = Crypto.encryptAES(title.value);
   
    const userWithKeys = {ID: Auth.user.value.ID, key: Crypto.encryptRSA(encTitle.key),right:3}

    putDoc({
        title: {
            content: encTitle.data,
            iv: encTitle.iv
        },
        user: userWithKeys
    }).then(()=>{
        if(errorDoc.value!=undefined){
            if(errorDoc.value.response.status != 401){
                console.log(errorDoc.value)
                showError('Error','Contact your administrator')
            } else {
                showError('Unauthorized', "You can't create an document." )
            }
        } else {
            showInfo('Success','The document is created', 1500,()=>{router.push({name:'docs'})})
        }
    })
}
</script>

<style>
</style>