<script setup lang="ts">
import {PropType} from "vue";
import {CardItem, Core} from "@/views/Dashboard/core";
import {ElDivider} from 'element-plus'
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {useI18n} from "@/hooks/web/useI18n";

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

<!--    <ElDivider content-position="left">Logs options</ElDivider>-->

    <!--    todo: add pagination limits-->
  </div>
</template>

<style lang="less">

</style>
