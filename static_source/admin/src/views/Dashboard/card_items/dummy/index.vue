<script setup lang="ts">
import {computed, onMounted, PropType, ref, unref, watch} from "vue";
import {Card, CardItem, Core, Tab} from "@/views/Dashboard/core";

// ---------------------------------
// common
// ---------------------------------
const item = ref<CardItem>({} as CardItem)

const props = defineProps({
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const el = ref(null)
onMounted(() => {
  // store dom element moveable
  props.item.setTarget(el.value)
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
  <div ref="el"></div>
</template>

<style lang="less" >

</style>
