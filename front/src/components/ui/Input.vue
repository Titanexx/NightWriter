<template>
    <div class="mt-2">
        <div class="flex justify-between" v-if="oneline">
            <label class="font-semibold my-auto">
                {{ label }}
            </label>
            <input
                v-model="data"
                class="p-4 rounded leading-tight bg-primary active:ring-0 focus:outline-none focus:ring-0 focus:shadow-outline"
                :class="{'border-red-500 border-2' : hasError}"
                :type="type"
                :placeholder="placeholder"
                :required="required"
                @blur="v$ ? v$[name].$touch() : () => {}"
            />
        </div>
        <div class="flex flex-col" v-else >
            <label class="block font-semibold mb-2 text-left">
                {{ label }}
            </label>
            <input
                v-model="data"
                class="w-full appearance-none rounded bg-primary leading-tight active:ring-0 focus:outline-none focus:ring-0 focus:shadow-outline"
                :class="{'border-red-500 border-2' : hasError}"
                :type="type"
                :placeholder="placeholder"
                :required="required"
                @blur="v$ ? v$[name].$touch() : () => {}"
            />
        </div>
        <div class="block w-full" v-if="v$">
            <div class="pt-2" v-show="v$[name].$errors.length > 0">
                <div class="text-red-500 text-left overflow-clip overflow-hidden" v-for="error of v$[name].$errors" :key="error.$uid" >{{ error.$message }}</div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { usePropsWrapper } from "@modules/helpers";
import { ref, watch } from "vue";

const emit = defineEmits()
const props = defineProps(['name','type','label','placeholder','required','modelValue','oneline','v$'])
const data = usePropsWrapper(props, emit, "modelValue")
const hasError = ref(false);

if(props.v$){
    watch(()=>props.v$[props.name].$errors, errors => hasError.value = errors.length > 0)
}

</script>

<style lang="postcss" scoped>
</style>