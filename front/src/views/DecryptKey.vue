<template>
	<div class="w-full flex justify-center content-center items-center h-screen">
		<div class="mx-auto bg-secondary rounded-lg text-center flex justify-center items-center min-w-max w-3/12">
			<form @submit.prevent="submit" class="p-8 w-full">
				<h1 class="text-2xl font-bold mt-1 mb-4">Decrypt your key</h1>

				<Input type="Password" label="Master Password" placeholder="Password" name="p" v-model="form.password" required/>

                <div class="block w-full pt-2 text-red-500 text-left overflow-clip overflow-hidden" v-show="hasError">
                    The masterpassword is wrong.
                </div>
                
                <Input type="checkbox" label="Remember it" v-model="form.remember" :oneline="true" />

				<div class="flex w-full mt-8">
                    <button type="submit" :disabled="loading" class="btn-primary w-full flex justify-center">
                        <svg v-show="loading" class="animate-spin h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                        <p v-show="!loading">Decrypt</p>
                    </button>
				</div>
			</form>
		</div>
	</div>
</template>

<script setup lang="ts">
import { reactive, ref } from "vue";

import router from '@routes';
import { Auth } from "@modules/auth";
import Input from "@components/ui/Input.vue";

const form = reactive({
    password: "",
    remember: false
})
const loading = ref(false);
const hasError = ref(false)

const submit = () => {
    loading.value = true;
    setTimeout(() => {
        const isLoaded = Auth.loadKeys(form.password, form.remember);
        delete form.remember;
        if(isLoaded){
            if(["/login","/login/","/decrypt-key/","/decrypt-key",""].includes(router.options.history.state.back?.toString())){
                router.push({ name: 'home' })
            } else {
                router.back()
            }
        } else {
            hasError.value = true
        }
        loading.value = false;
    }, 10);
};
</script>

<style>
</style>