<template>
    <Listbox v-model="model">
        <div class="relative">
            <ListboxButton class="relative z-0 flex flex-row items-center text-right bg-primary rounded min-h-0 h-6 mx-2 w-28">
                <span class="flex-1 no-selection mr-1">{{button}}</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 bg-gray-0" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M10 3a1 1 0 01.707.293l3 3a1 1 0 01-1.414 1.414L10 5.414 7.707 7.707a1 1 0 01-1.414-1.414l3-3A1 1 0 0110 3zm-3.707 9.293a1 1 0 011.414 0L10 14.586l2.293-2.293a1 1 0 011.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clip-rule="evenodd" />
                </svg>
            </ListboxButton>

            <transition
                leave-active-class="transition duration-100 ease-in"
                leave-from-class="opacity-100"
                leave-to-class="opacity-0"
            >
            <ListboxOptions class="absolute mt-1 overflow-auto text-base text-primary bg-primary border-2 border-gray-600 rounded-md shadow-lg max-h-60 z-50">
                <ListboxOption 
                    v-slot="{ active, selected }"
                    v-for="option in options"
                    :key="option.val"
                    :value="option.val"
                    as="template"
                >
                <li
                    :class="[
                    active ? 'bg-gray-300': '',
                    'cursor-default select-none relative py-2 pl-10 pr-4',
                    ]"
                >
                    <span :class="[selected ? 'font-medium' : 'font-normal','block truncate w-100']">
                        {{ option.str }}
                    </span>
                    <span v-if="selected" class="absolute inset-y-0 left-0 flex items-center pl-3 text-amber-600">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                            <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                        </svg>
                    </span>
                </li>
                </ListboxOption>
            </ListboxOptions>
            </transition>
        </div>
    </Listbox>
</template>

<script setup lang="ts">
import {
    Listbox,
    ListboxButton,
    ListboxOptions,
    ListboxOption,
} from '@headlessui/vue'
import { usePropsWrapper } from '@modules/helpers';
import type { ListOptionsMap } from '@modules/helpers';
import type { PropType } from 'vue';

const emit = defineEmits()
const props = defineProps({
    modelValue: [String, Number, Object],
    button: String,
    options: Object as PropType<ListOptionsMap>,
})
const model = usePropsWrapper(props, emit, "modelValue")
</script>

<style>
</style>