<script setup lang="ts">
import {computed, onMounted, PropType, ref, unref, watch} from "vue";
import {ApiEntity, ApiImage, ApiVariable} from "@/api/stub";
import {useI18n} from "@/hooks/web/useI18n";
import View from "@/views/Dashboard/view/view.vue";
import api from "@/api/api";

const {t} = useI18n()

const loading = ref(true)
const id = ref<Nullable<number>>(null)

// ---------------------------------
// common
// ---------------------------------

const fetchDashboard =  async () => {
  loading.value = true;

  const res = await api.v1.variableServiceGetVariableByName('mainDashboard')
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })

  if (!res) {
    loading.value = false;
    return;
  }

  const variable = res.data as ApiVariable

  id.value = parseInt(variable.value);

  loading.value = false;
}

fetchDashboard()

</script>

<template>
  <View v-if="!loading && id" :id="id" />
</template>

<style lang="less" scoped>

</style>
