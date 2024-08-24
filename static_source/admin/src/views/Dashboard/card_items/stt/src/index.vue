<script setup lang="ts">
import {onMounted, onUnmounted, PropType, ref, watch} from "vue";
import {CardItem} from "@/views/Dashboard/core";
import {ElButton} from "element-plus";
import {startRecording, stopRecording} from "@/components/Stt";

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

})

onUnmounted(() => {

})

// ---------------------------------
// component methods
// ---------------------------------

watch(
  () => props.item,
  (val?: CardItem) => {
    if (!val) return;

  },
  {
    deep: true,
    immediate: true
  }
)

const setIsRecording = ref(false)

const start = () => {
  setIsRecording.value = true;
  startRecording()
}

const stop = () => {
  setIsRecording.value = false;
  stopRecording()
};

const toggleRecording = () => {
  if (!setIsRecording.value) {
    start()
  } else {
    stop()
  }
}


</script>

<template>
  <div ref="el" :class="[{'hidden': item.hidden}]" class="w-[100%] h-[100%]">
    <ElButton type="primary" @click="toggleRecording()">
      {{ !setIsRecording ? 'recording start' : 'recording stop' }}
    </ElButton>
  </div>

</template>

<style lang="less" scoped>

</style>
