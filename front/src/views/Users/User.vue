<template>
    <div>
        <p>
            User {{ $route.params.id }}
        </p>
    </div>
    <div class="mt-4">
        <p>{{ Auth.user }}</p>
    </div>
</template>

<script setup lang="ts">
import { Auth } from '@modules/auth';

</script>