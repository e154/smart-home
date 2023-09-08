<script setup lang="ts">

const src: string = (import.meta.env.VITE_API_BASEPATH as string) + '/swagger-ui'
</script>

<template>
  <iframe
      :src="src"
      style="position: absolute; border: none; width: 100%; height: 100%"
      width="100%"
      height="100%"
      frameborder="0"
  >
  </iframe>
</template>

<style lang="less">

</style>
