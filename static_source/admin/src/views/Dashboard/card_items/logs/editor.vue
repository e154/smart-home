<script setup lang="ts">
import {computed, PropType, ref, unref, watch} from "vue";
import {Card, CardItem, Core, Tab} from "@/views/Dashboard/core";
import {ElDivider} from 'element-plus'
import {Vuuri} from "@/views/Dashboard/Vuuri"
import {useBus} from "@/views/Dashboard/bus";
import ViewCard from "@/views/Dashboard/editor/ViewCard.vue";
import CommonEditor from "@/views/Dashboard/card_items/common/editor.vue";
import {useI18n} from "@/hooks/web/useI18n";

const {t} = useI18n()

// ---------------------------------
// common
// ---------------------------------
const item = ref<CardItem>({} as CardItem)

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

watch(
    () => props.item,
    (val?: CardItem) => {
      if (!val) return;
      item.value = val;
    },
    {
      deep: true,
      immediate: true
    }
)

// ---------------------------------
// component methods
// ---------------------------------


</script>

<template>
  <div>

    <CommonEditor :item="item" :core="core"/>

    <ElDivider content-position="left">Logs options</ElDivider>

<!--    todo: add pagination limits-->
  </div>
</template>

<style lang="less" >

</style>
