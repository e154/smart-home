<script setup lang="ts">
import {PropType} from "vue";
import {ElCol, ElDivider, ElRow} from 'element-plus'
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {useI18n} from "@/hooks/web/useI18n";
import {CardItem, Core} from "@/views/Dashboard/core";

const {t} = useI18n()

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


// ---------------------------------
// component methods
// ---------------------------------


</script>

<template>
  <div>

    <CommonEditor :item="item" :core="core"/>

    <ElRow class="mb-10px mt-10px">
      <ElCol>
        <ElDivider content-position="left">Dummy options</ElDivider>
      </ElCol>
    </ElRow>

    [Unknown item type "{{ item.type }}"]

  </div>
</template>

<style lang="less">

</style>
