<template>
	<div class="mx-auto bg-secondary rounded-lg text-center flex justify-center items-center w-96">
		<form @submit.prevent="submit" class="p-8 w-full">
			<h1 class="text-2xl font-bold mt-1 mb-4">Login</h1>

			<router-link :to="{name:'register'}"><button type="button" class="btn-secondary border-none w-full">Register</button></router-link>

            <div class="mt-2 text-red-500 text-left overflow-clip overflow-hidden" v-show="hasError">
                Check your information. If the error remains, contact the administrator.
            </div>

			<Input type="text" label="Username" placeholder="Username" name="username" v-model="form.username" required/>
			<Input type="Password" label="Password" placeholder="Password" name="password" v-model="form.password" required/>
		
			<div class="flex w-full mt-8">
                <button type="submit" :disabled="loading" class="btn-primary w-full">
                    <svg v-show="loading" class="animate-spin h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <p v-show="!loading">Sign in</p>
                </button>
			</div>
		</form>
	</div>
</template>

<script setup lang="ts">
import {  reactive, ref } from "vue";

import { useApi } from "@modules/api";
import { Auth } from "@modules/auth";
import Input from "@components/ui/Input.vue";
import { ModalType, showModal } from "@plugins/modal";

const hasError = ref(false)
const form = reactive({
    username: "",
    password: "",
});

const {
    loading,
    post,
    data,
    error,
} = useApi("/api/login");

// showModal(Modal.Error,"test","test",50000)

const submit = () => {
    post(form, { withCredentials: true }).then(() => {
        delete(form.password)
        if(error.value==undefined){
            hasError.value = false
            Auth.setToken(data.value['access_token']);
        } else {
            hasError.value = true
        }
    });
};
</script>

<style>
</style>