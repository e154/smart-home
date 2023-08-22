<script setup lang="ts">
import {computed, nextTick, onMounted, PropType, ref, unref, watch} from "vue";
import {ApiEntity, ApiImage, ApiVariable, ApiZigbee2Mqtt} from "@/api/stub";
import {useI18n} from "@/hooks/web/useI18n";
import View from "@/views/Dashboard/view/view.vue";
import api from "@/api/api";
import {isDark} from "@/utils/is";
import {useAppStore} from "@/store/modules/app";

const {t} = useI18n()

const loading = ref(true)
const id = ref<Nullable<number>>(null)

// ---------------------------------
// common
// ---------------------------------

const reloadKey = ref(0);
const reload = () => {
  reloadKey.value += 1
}

const appStore = useAppStore()


const fetchDashboard =  async () => {
  loading.value = true;
  const res = await api.v1.variableServiceGetVariableByName('mainDashboard' + (appStore.isDark? 'Dark' : 'Light'))
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })

  if (!res) {
    return;
  }

  const variable = res.data as ApiVariable

  nextTick(() => {
    id.value = parseInt(variable.value);
    reload()
  })

  loading.value = false;
}


watch(
    () => appStore.isDark,
    (val: boolean) => {
      fetchDashboard()
    },
    {
      immediate: false
    }
)

fetchDashboard()

</script>

<template>
  <View v-if="!loading && id" :id="id" :key="reloadKey"/>
</template>

<style lang="less" scoped>

</style>
