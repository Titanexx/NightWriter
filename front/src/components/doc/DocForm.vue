<template>
    <form @submit.prevent="emit('submit')" class="w-10/12 p-8">
        <slot name="header"></slot>

        <Input type="text" label="Title" placeholder="Title" name="title" v-model="title" required/>
    
        <div class="flex w-full mt-8">
            <button type="submit" :disabled="props.loading" class="btn-primary w-full">
                <slot>Submit</slot>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { usePropsWrapper } from "@modules/helpers";
import Input from "@components/ui/Input.vue";

const props = defineProps(['loading','title'])
const emit = defineEmits(['submit','update:title'])

const title = usePropsWrapper(props, emit, "title")
</script>

<style lang="postcss" scoped>
</style>