<script setup lang="ts">
import {computed, onMounted, PropType, ref, unref, watch} from "vue";
import {ApiVariable} from "@/api/stub";
import View from "@/views/Dashboard/view/view.vue";
import api from "@/api/api";

const loading = ref(true)
const id = ref<Nullable<number>>(null)

// ---------------------------------
// common
// ---------------------------------

const fetchDashboard =  async () => {
  loading.value = true;

  const res = await api.v1.variableServiceGetVariableByName('devDashboard')
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
  <div class="ml-20px mt-20px">
    <View v-if="!loading && id" :id="id" />
  </div>
</template>

<style lang="less" scoped>

</style>
