<script setup lang="ts">
import {computed, onMounted, PropType, ref, unref, watch} from "vue";
import View from "@/views/Dashboard/view/view.vue";
import {useRoute} from "vue-router";

const route = useRoute();
const dashboardId = computed(() => parseInt(route.params.id as string));

</script>

<template>
  <View v-if="dashboardId" :id="dashboardId" class="ml-20px mt-20px"/>
</template>

<style lang="less" scoped>

</style>
