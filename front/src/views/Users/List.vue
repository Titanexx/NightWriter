<template>
    <span>Todo</span>
</template>

<script setup lang="ts">
import { useApiToken } from '../../modules/api';

const {
    loading,
    get,
    data,
    error,
} = useApiToken("/api/users/");

get().then(() => {
    if(error.value!=undefined && error.value.response.status != 401){
        console.log('there is an error')
    }
});
</script>