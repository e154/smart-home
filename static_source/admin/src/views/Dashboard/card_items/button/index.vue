<script setup lang="ts">
import {computed, onMounted, PropType, ref, unref, watch} from "vue";
import {ElButton, ElMessage, ElText} from 'element-plus'
import {Card, CardItem, Core, requestCurrentState, Tab} from "@/views/Dashboard/core";
import api from "@/api/api";
import {propTypes} from "@/utils/propTypes";
import {useI18n} from "@/hooks/web/useI18n";

const {t} = useI18n()

// ---------------------------------
// common
// ---------------------------------
const item = ref<CardItem>({} as CardItem)

const props = defineProps({
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
  disabled: propTypes.bool.def(false),
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

const callAction = async () => {
  await api.v1.interactServiceEntityCallAction({
    id: props.item?.entityId,
    name: props.item?.payload.button?.action || ''
  })
  ElMessage({
    title: t('Success'),
    message: t('message.callSuccessful'),
    type: 'success',
    duration: 2000
  })
}

const onClick = () => {
  callAction()
}

requestCurrentState(props.item?.entityId);
</script>

<template>
  <div ref="el" class="w-[100%] h-[100%]">
    <ElButton
        style="width: 100%; height: 100%"
        v-if="item.enabled" v-show="!item.hidden"
        :size="item.payload.button.size"
        :type="item.payload.button.type"
        :round="item.payload.button.round"
        @click.prevent.stop="onClick"
        :disabled="props.disabled"
    >{{ item.payload.button.text }}
    </ElButton>
  </div>
</template>

<style lang="less" scoped>

.clearfix:before,
.clearfix:after {
  display: table;
  content: "";
}

.clearfix:after {
  clear: both
}
</style>
