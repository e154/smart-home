<script setup lang="ts">
import {computed} from "vue";
import View from "@/views/Dashboard/view/view.vue";
import {useRoute} from "vue-router";

const route = useRoute();
const dashboardId = computed(() => parseInt(route.params.id as string) as number);

</script>

<template>
  <View v-if="dashboardId" :id="dashboardId"/>
</template>

<style lang="less" scoped>

</style>
