<template>
    <div class="bg-secondary rounded-lg text-center flex w-6/12">
        <form @submit.prevent="submit" class="py-4 px-10">
        <h1 class="text-2xl font-bold mb-5">Register</h1>

		<router-link :to="{name:'login'}"><button type="button" class="btn-secondary border-none w-full">Already register ? Go sign in.</button></router-link>

        <Input type="text" label="Username" placeholder="Username" name="username" :v$="v$" v-model="form.username" required/>
        
        <div class="flex flex-wrap -mx-2 overflow-hidden">
            <div class="px-2 w-1/2 overflow-hidden">
            <Input type="password" label="Password" placeholder="Password" name="password" :v$="v$" v-model="form.password" required/>
            </div>

            <div class="px-2 w-1/2 overflow-hidden">
            <Input type="password" label="Confirm Password" placeholder="Confirm your Password" name="confirm_password" :v$="v$" v-model="form.confirm_password" required/>
            </div>
        </div>

        <Input type="text" label="Name" placeholder="Name" name="name" :v$="v$" v-model="form.name" required/>
        <Input type="text" label="Email" placeholder="Email" name="email" :v$="v$" v-model="form.email" required/>
    
        <div class="flex flex-wrap -mx-2 overflow-hidden">
            <div class="px-2 w-1/2 overflow-hidden">
            <Input type="password" label="Master Password" placeholder="Master Password" name="masterpassword" :v$="v$" v-model="form.masterpassword" required/>
            </div>

            <div class="px-2 w-1/2 overflow-hidden">
            <Input type="password" label="Confirm Master Password" placeholder="Confirm your Master Password" name="confirm_masterpassword" :v$="v$" v-model="form.confirm_masterpassword" required/>
            </div>
            <p class="px-2 w-full mt-2 text-sm text-left">The master password is used to encrypted your private key. It must be long and robust. Save it into your Keepass. Keys generation can be long. Wait please.</p>
        </div>

        <div class="flex w-full mt-8">
            <button type="submit" :disabled="loading" class="btn-primary w-full flex justify-center items-center">
                <p v-show="!loading">Register</p>
                <svg v-show="loading" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-100" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
            </button>
        </div>
        </form>
    </div>
</template>

<script setup lang="ts">
import { reactive } from "vue";
import useVuelidate from '@vuelidate/core'
import { email as emailValidator, helpers, minLength } from '@vuelidate/validators'

import router from '@routes';
import { useApi } from "@modules/api";
import { Crypto } from "@modules/crypto";
import Input from "@components/ui/Input.vue";
import { showError, showInfo } from "@plugins/modal";

const form = reactive({
    username: "",
    password: "",
    confirm_password: "",
    name: "",
    email: "",
    public_key: "",
    private_key: "",
    masterpassword: "",
    confirm_masterpassword: "",
});

function checkComplexity(password){
    return /[A-Z]/.test(password) && /[a-z]/.test(password) && /\d/.test(password) && /\W/.test(password)
}

const rules = {
    username: {},
    password: {
        minLengthValue: minLength(12),
        checkComplexity: helpers.withMessage("The password must contain at least one capital letter, one number and one special character.",checkComplexity)
    },
    confirm_password: {
        checkComplexity: helpers.withMessage("The password confirmation must match your password.", (value) => value === form.password)
    },
    name: {},
    email: {
        emailValidator
    },
    masterpassword: {
        minLengthValue: minLength(16),
        checkComplexity: helpers.withMessage("The masterpassword must contain at least one capital letter, one number and one special character.",checkComplexity),
        checkPassword: helpers.withMessage("The masterpassword confirmation must not match your account password.", (value) => value != form.password)
    },
    confirm_masterpassword: { checkComplexity: helpers.withMessage("The masterpassword confirmation must match your masterpassword.", (value) => value === form.confirm_masterpassword) },
}

const v$ = useVuelidate(rules, form)

const {
    error,
    loading,
    post,
} = useApi("/api/register");

function register(keys: any){
    form.public_key = keys.public;
    form.private_key = keys.private;
    delete(form.confirm_password)
    delete(form.masterpassword)
    delete(form.confirm_masterpassword)
    console.log("... Keys are generated.");
    post(form).then(() => {
        if(error.value!=undefined){
            showError("Error","There is an error with your informations.")
        } else {
            showInfo("Success","You are registered. Contact your administrator to validate your account.",2000,() => router.push({ name: 'login' }))
        }
    });
}

const submit = () => {
    //@ts-ignore
    v$.value.$touch()
    //@ts-ignore
    if (v$.value.$error) return
    loading.value = true;
    setTimeout(() =>{
        console.log("Start keys generation...");
        Crypto.genKeys(form.masterpassword, register, ()=>{loading.value=false});
    },10)
};
</script>

<style>
</style>