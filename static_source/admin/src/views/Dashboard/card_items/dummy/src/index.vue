<script setup lang="ts">
import {onMounted, PropType, ref} from "vue";
import {CardItem} from "@/views/Dashboard/core";

// ---------------------------------
// common
// ---------------------------------

const props = defineProps({
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const el = ref<ElRef>(null)
onMounted(() => {
  // store dom element moveable
  props.item.setTarget(el.value)
})

// ---------------------------------
// component methods
// ---------------------------------


</script>

<template>
  <div ref="el"></div>
</template>

<style lang="less">

</style>
