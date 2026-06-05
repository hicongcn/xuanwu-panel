<script lang="ts" setup>
import type { ToasterProps } from "vue-sonner"
import { Toaster as Sonner } from "vue-sonner"

const props = withDefaults(defineProps<ToasterProps>(), {
  theme: 'light',
  richColors: true
})
</script>

<template>
  <Sonner
    class="toaster group"
    v-bind="props"
  />
</template>
