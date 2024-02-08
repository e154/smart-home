<script setup lang="ts">
import {computed, PropType, ref} from "vue";
import {CardItem, Core} from "@/views/Dashboard/core";
import {ElCol, ElDivider, ElFormItem, ElRow, ElSwitch} from 'element-plus'
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {useI18n} from "@/hooks/web/useI18n";
import {EntitiesSearch} from "@/components/EntitiesSearch";

const {t} = useI18n()

const entityIds = ref([])

// ---------------------------------
// common
// ---------------------------------

const props = defineProps({
  core: {
    type: Object as PropType<Nullable<Core>>,
    default: () => null
  },
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const currentItem = computed(() => props.item as CardItem)

// ---------------------------------
// component methods
// ---------------------------------


</script>

<template>
  <div>

    <CommonEditor :item="item" :core="core"/>

    <ElDivider content-position="left">Entity storage options</ElDivider>

    <ElRow :gutter="24">
      <ElCol :span="24">
        <ElFormItem :label="$t('dashboard.editor.entity_storage.entities')" prop="entityIds">
          <EntitiesSearch v-model="currentItem.payload.entityStorage.entityIds"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.entity_storage.showFilter')" prop="filter">
          <ElSwitch v-model="currentItem.payload.entityStorage.filter"/>
        </ElFormItem>
      </ElCol>
      <ElCol :span="12" :xs="12"/>
    </ElRow>

  </div>
</template>

<style lang="less">

</style>
