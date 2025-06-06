<script setup lang="ts">
import { PropType } from 'vue'
import { Highlight } from '@/components/Highlight'
import { useDesign } from '@/hooks/web/useDesign'
import { propTypes } from '@/utils/propTypes'
import { TipSchema } from '@/types/infoTip'

const { getPrefixCls } = useDesign()

const prefixCls = getPrefixCls('infotip')

defineProps({
  title: propTypes.string.def(''),
  type: propTypes.string.def(''),
  schema: {
    type: Array as PropType<Array<string | TipSchema>>,
    required: true,
    default: () => []
  },
  showIndex: propTypes.bool.def(true),
  highlightColor: propTypes.string.def('var(--el-color-primary)')
})

const emit = defineEmits(['click'])

const keyClick = (key: string) => {
  emit('click', key)
}
</script>

<template>
  <div v-if="type != 'warning'"
      style="border-left: 5px solid var(--el-color-primary); border-radius: 4px;"
      :class="[
        prefixCls,
        'p-20px mb-20px border-solid border-[var(--el-color-primary)] bg-[var(--el-color-primary-light-9)]'
      ]"
  >
    <div v-if="title" :class="[`${prefixCls}__header`, 'flex items-center']">
      <span :class="[`${prefixCls}__title`, 'text-14px font-bold']">{{ title }}</span>
    </div>
    <div :class="`${prefixCls}__content`">
      <p v-for="(item, $index) in schema" :key="$index" class="text-12px mt-10px">
        <Highlight
          :keys="typeof item === 'string' ? [] : item.keys"
          :color="highlightColor"
          @click="keyClick"
        >
          {{ showIndex ? `${$index + 1}、` : '' }}{{ typeof item === 'string' ? item : item.label }}
        </Highlight>
      </p>
    </div>
  </div>
  <div v-if="type == 'warning'"
      style="border-left: 5px solid var(--el-color-warning); border-radius: 4px;"
      :class="[
        prefixCls,
        'p-20px mb-20px border-solid border-[var(--el-color-warning)] bg-[var(--el-color-warning-light-9)]'
      ]"
  >
    <div v-if="title" :class="[`${prefixCls}__header`, 'flex items-center']">
      <span :class="[`${prefixCls}__title`, 'text-14px font-bold']">{{ title }}</span>
    </div>
    <div :class="`${prefixCls}__content`">
      <p v-for="(item, $index) in schema" :key="$index" class="text-12px mt-10px">
        <Highlight
          :keys="typeof item === 'string' ? [] : item.keys"
          :color="highlightColor"
          @click="keyClick"
        >
          {{ showIndex ? `${$index + 1}、` : '' }}{{ typeof item === 'string' ? item : item.label }}
        </Highlight>
      </p>
    </div>
  </div>
</template>
