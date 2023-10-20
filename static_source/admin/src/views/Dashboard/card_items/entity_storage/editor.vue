<script setup lang="ts">
import {computed, PropType, ref, unref, watch} from "vue";
import {Card, CardItem, Core, Tab} from "@/views/Dashboard/core";
import {ElDivider, ElRow, ElCol, ElFormItem} from 'element-plus'
import {Vuuri} from "@/views/Dashboard/Vuuri"
import {useBus} from "@/views/Dashboard/bus";
import ViewCard from "@/views/Dashboard/editor/ViewCard.vue";
import CommonEditor from "@/views/Dashboard/card_items/common/editor.vue";
import {useI18n} from "@/hooks/web/useI18n";
import EntitiesSearch from "@/views/Entities/components/EntitiesSearch.vue";

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
        <ElFormItem :label="$t('dashboard.editor.entity_storage.entities')" prop="background">
          <EntitiesSearch v-model="currentItem.payload.entityStorage.entityIds"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

  </div>
</template>

<style lang="less" >

</style>
