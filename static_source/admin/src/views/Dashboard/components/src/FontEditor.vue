<script setup lang="ts">

import {computed, PropType} from "vue";
import {Core, Tab} from "@/views/Dashboard/core";
import {ElCol, ElDivider, ElRow} from "element-plus";
import {useI18n} from "@/hooks/web/useI18n";
import VariablesSearch from "@/components/VariablesSearch/src/VariablesSearch.vue";

const {t} = useI18n()

// ---------------------------------
// common
// ---------------------------------

const props = defineProps({
  core: {
    type: Object as PropType<Core>,
  },
  tab: {
    type: Object as PropType<Nullable<Tab>>,
    default: () => null
  },
})

const currentTab = computed({
  get(): Tab {
    return props.tab as Tab
  },
  set(val: Tab) {
  }
})

// ---------------------------------
// component methods
// ---------------------------------


</script>

<template>
  <div v-if="currentTab">
    <ElRow class="mt-10px mb-10px">
      <ElCol>
        <ElDivider content-position="left">{{ $t('dashboard.editor.FontsOptions') }}</ElDivider>
      </ElCol>
    </ElRow>

    <ElRow class="mb-10px">
      <ElCol>
        <VariablesSearch v-model="currentTab.fonts"/>
      </ElCol>
    </ElRow>
  </div>


</template>

<style scoped lang="less">

</style>
