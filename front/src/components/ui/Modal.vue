<template>
    <TransitionRoot appear :show="isOpen" as="template">
        <Dialog as="div" :open="isOpen" @close="closeModal">
            <div class="fixed inset-0 z-50 overflow-y-auto">
                <div class="min-h-screen px-4 text-center">
                    <TransitionChild
                        as="template"
                        enter="duration-300 ease-out"
                        enter-from="opacity-0"
                        enter-to="opacity-100"
                        leave="duration-200 ease-in"
                        leave-from="opacity-100"
                        leave-to="opacity-0"
                    >
                        <DialogOverlay class="fixed inset-0 bg-smoke-500" />
                    </TransitionChild>

                    <span class="inline-block h-screen align-middle" aria-hidden="true">
                        &#8203;
                    </span>

                    <TransitionChild
                        as="template"
                        enter="duration-300 ease-out"
                        enter-from="opacity-0 scale-95"
                        enter-to="opacity-100 scale-100"
                        leave="duration-200 ease-in"
                        leave-from="opacity-100 scale-100"
                        leave-to="opacity-0 scale-95"
                    >
                    
                        <div :class="divClass" class="inline-block w-full max-w-md p-6 my-8 overflow-hidden text-left align-middle transition-all transform shadow-xl rounded-2xl bg-secondary" >
                            <button class="float-right opacity-100 right-0 p-0" @click="closeModal">
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
                                </svg>
                            </button>
                            <DialogTitle as="h3" class="text-lg font-medium leading-6 text-gray-900">
                                    <div class="flex justify-between">
                                    <div class="">
                                        <slot name="title"></slot>
                                    </div>
                                </div>
                            </DialogTitle>
                            <slot name="body"></slot>
                        </div>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>


<script setup lang="ts">
import {
  TransitionRoot,
  TransitionChild,
  Dialog,
  DialogOverlay,
  DialogTitle,
} from '@headlessui/vue'
import { usePropsWrapper } from "@modules/helpers";

const emit = defineEmits()
const props = defineProps({
    isOpen: Boolean,
    divClass: {type:String, default: ""}
})

const show = usePropsWrapper(props,emit,'isOpen')

function closeModal() {
    show.value = false
}

</script>

<style lang="scss">
</style>