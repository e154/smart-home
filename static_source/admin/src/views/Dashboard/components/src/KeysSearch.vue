<script setup lang="ts">

import {PropType, ref, watch} from "vue";
import {ElAutocomplete} from 'element-plus'
import {getAllKeys, getFilteredKeys} from "@/views/Dashboard/core";
import {propTypes} from "@/utils/propTypes";

const currentValue = ref<Nullable<string>>(null)
const AllKeys = ref<{ value: string }[]>([])

const emit = defineEmits(['change', 'update:modelValue'])

const props = defineProps({
  modelValue: propTypes.string.def(''),
  allKeys: propTypes.bool.def(false),
  obj: {
    type: Object as PropType<Nullable<Object>>,
    default: () => undefined
  }

})

watch(
    () => props.obj,
    (val: any) => {
      if (!props.obj) {
        return
      }
      if (props.allKeys) {
        const keys = getAllKeys(val)
        AllKeys.value = keys.map(value => {
          return {value: value}
        })
      } else {
        const keys = getFilteredKeys(val)
        AllKeys.value = keys.map(value => {
          return {value: value}
        })
      }
    },
    {
      immediate: true
    }
)

watch(
    () => props.modelValue,
    (val: any) => {
      if (!props.modelValue) {
        currentValue.value = ''
        return
      }
      currentValue.value = val
    },
    {
      immediate: true
    }
)


const querySearch = (queryString: string, cb: any) => {
  const results = queryString
      ? AllKeys.value.filter(createFilter(queryString))
      : AllKeys.value
  // call callback function to return suggestions
  cb(results)
}
const createFilter = (queryString: string) => {
  return (item: { value: string }) => {
    return (
        item.value.toLowerCase().indexOf(queryString.toLowerCase()) === 0
    )
  }
}

const handleChange = (val) => {
  if (val == '') {
    currentValue.value = null
    emit('change', '')
  }
}
const handleSelect = (val: { value: string }) => {
  currentValue.value = val.value
  emit('change', val.value)
}

</script>

<template>
  <ElAutocomplete
      class="w-[100%]"
      v-model="currentValue"
      :fetch-suggestions="querySearch"
      placeholder="Please input"
      @select="handleSelect"
      @change="handleChange"
      clearable
  />
</template>

<style lang="less">

</style>
