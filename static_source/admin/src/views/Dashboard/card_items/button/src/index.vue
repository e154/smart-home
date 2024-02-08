<script setup lang="ts">
import {onMounted, PropType, ref} from "vue";
import {ElButton, ElMessage} from 'element-plus'
import {CardItem, requestCurrentState} from "@/views/Dashboard/core/core";
import api from "@/api/api";
import {propTypes} from "@/utils/propTypes";
import {useI18n} from "@/hooks/web/useI18n";

const {t} = useI18n()

// ---------------------------------
// common
// ---------------------------------
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
    >
      <Icon v-if="item.payload.button.icon" :icon="item.payload.button.icon"/>
      {{ item.payload.button.text }}
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
