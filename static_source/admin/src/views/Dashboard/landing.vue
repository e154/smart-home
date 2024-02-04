<script setup lang="ts">
import View from "@/views/Dashboard/view/view.vue";
import {useRoute} from "vue-router";
import {ElEmpty} from "element-plus";
import {useAppStore} from "@/store/modules/app";

const appStore = useAppStore()
const route = useRoute();

const dashboardId = parseInt(route.params.id as string) as number
const accessToken = route.query.access_token as string
const serverId = route.query.serverId as string

if (accessToken && !appStore.getToken) {
  appStore.SetToken(accessToken)
}
if (appStore.getIsGate && serverId) {
  appStore.setServerId(serverId)
}

</script>

<template>
  <View v-if="dashboardId" :id="dashboardId"/>
  <ElEmpty v-if="!dashboardId" :rows="5"/>
</template>

<style lang="less" scoped>

</style>
