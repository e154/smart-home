<!--
  - This file is part of the Smart Home
  - Program complex distribution https://github.com/e154/smart-home
  - Copyright (C) 2024, Filippov Alex
  -
  - This library is free software: you can redistribute it and/or
  - modify it under the terms of the GNU Lesser General Public
  - License as published by the Free Software Foundation; either
  - version 3 of the License, or (at your option) any later version.
  -
  - This library is distributed in the hope that it will be useful,
  - but WITHOUT ANY WARRANTY; without even the implied warranty of
  - MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
  - Library General Public License for more details.
  -
  - You should have received a copy of the GNU Lesser General Public
  - License along with this library.  If not, see
  - <https://www.gnu.org/licenses/>.
  -->

<script setup lang="ts">

import {PropType, ref, unref, watch} from "vue";
import {ApiEntityShort} from "@/api/stub";
import {ElAutocomplete} from 'element-plus'
import api from "@/api/api";

const currentEntity = ref<Nullable<ApiEntityShort>>(null)
const entityId = ref<Nullable<string>>(null)
const emit = defineEmits(['change', 'update:modelValue'])

const props = defineProps({
  modelValue: {
    type: Object as PropType<Nullable<ApiEntityShort>>,
    default: () => undefined
  }
})

watch(
    () => props.modelValue,
    (val?: ApiEntityShort) => {
      if (val === unref(currentEntity)) return
      entityId.value = val?.id || null;
      currentEntity.value = val || null;
    },
    {
      immediate: true
    }
)

// 监听
watch(
    () => currentEntity.value,
    (val?: ApiEntityShort) => {
      if (props.modelValue == unref(val)) return;
      emit('update:modelValue', val)
      if (!val) {
        emit('change', val)
      }
    },

)

const querySearchAsync = async (queryString: string, cb: Fn) => {
  if (queryString == null || queryString == 'null') {
    queryString = ''
  }
  const params = {query: queryString, limit: 25, offset: 0}
  const {data} = await api.v1.entityServiceSearchEntity(params)
  const {items} = data
  cb(items)
}

const handleChange = (val) => {
  if (val == '') {
    currentEntity.value = null
  }
}
const handleSelect = (val: ApiEntityShort) => {
  currentEntity.value = val
  emit('change', val)
}

</script>

<template>
  <ElAutocomplete
      class="w-[100%]"
      v-model="entityId"
      :fetch-suggestions="querySearchAsync"
      placeholder="Please input"
      value-key="id"
      @select="handleSelect"
      @change="handleChange"
      clearable
  />
</template>

<style lang="less">

</style>
