<script setup lang="ts">

import {onMounted, PropType, ref, watch} from "vue";
import {EntitiesActionOptions} from "./types";
import {ElCol, ElForm, ElFormItem, ElInput, ElOption, ElRow, ElSelect} from "element-plus";
import {EntitySearch} from "@/components/EntitySearch";
import {ApiArea, ApiEntity} from "@/api/stub";
import api from "@/api/api";
import {TagsSearch} from "@/components/TagsSearch";
import {AreaSearch} from "@/components/AreaSearch";

const emit = defineEmits(['change', 'update:options'])

const props = defineProps({
  options: {
    type: Object as PropType<Nullable<EntitiesActionOptions>>,
    default: () => null
  },
  entity: {
    type: Object as PropType<Nullable<ApiEntity>>,
    default: () => null
  }
})

const action = ref<EntitiesActionOptions>()
const once = ref(false)

let currentEntity = ref<ApiEntity>();
let currentArea = ref<ApiArea>();

onMounted(() => {

})

watch(
    () => props.options,
    (val?: EntitiesActionOptions) => {
      if (!val) return;
      // if (val == action.value) return;
      action.value = {
        entityId: val.entityId || val.entity?.id || '',
        action: val.action || val.actionName || '',
        tags: val.tags,
        areaId: val.areaId,
      }

      if (action.value.areaId) {
        api.v1.areaServiceGetAreaById(action.value.areaId)
            .then(({data}) => {
              currentArea.value = data;
            })
      }

      if (action.value.entityId) {
        api.v1.entityServiceGetEntity(action.value.entityId)
            .then(({data}) => {
              currentEntity.value = data;
            })
      } else if (props.entity && !once.value) {
        once.value = true
        currentEntity.value = props.entity
      } else {
        currentEntity.value = undefined
      }
    },
    {
      immediate: true,
    }
)

const changedEntity = async (entity: ApiEntity) => {
  if (entity && entity?.id) {
    const {data} = await api.v1.entityServiceGetEntity(entity.id);
    currentEntity.value = data;
    action.value.entityId = entity.id;
  } else {
    currentEntity.value = undefined;
    action.value.entityId = '';
  }
  handleSelect()
}

const getActionList = () => {
  if (!currentEntity.value) {
    return [];
  }
  return currentEntity.value.actions;
}

const changedTags = async (tags: string[]) => {
  // console.log(tags)
  action.value.tags = tags
  handleSelect()
}

const changedArea = async (area: ApiArea) => {
  // console.log(area)
  if (area?.id) {
    currentArea.value = area
    action.value.areaId = area.id
  } else {
    currentArea.value = undefined
    action.value.areaId = undefined
  }
  handleSelect()
}

const changedAction = (action: any) => {
  // console.log(action)
  handleSelect()
}

const handleSelect = () => {
  // console.log(action.value)
  emit('change', {
    entityId: action.value.entityId || undefined,
    action: action.value.action || '',
    tags: action.value.tags || [],
    areaId: action.value.areaId || undefined,
  })
}

</script>

<template>
  <ElForm v-if="action"
          label-position="top"
          :model="action"
          style="width: 100%"
          ref="cardItemForm">

    <ElRow :gutter="24">
      <ElCol :span="12">
        <ElFormItem :label="$t('entityAction.entity')" prop="entity">
          <EntitySearch v-model="currentEntity" @change="changedEntity($event)"/>
        </ElFormItem>
      </ElCol>

      <ElCol :span="12" v-if="currentEntity">
        <ElFormItem :label="$t('entityAction.action')" prop="action">
          <ElSelect
              v-model="action.action"
              clearable
              :placeholder="$t('entityAction.selectAction')"
              style="width: 100%"
              @change="changedAction($event)"
          >
            <ElOption
                v-for="item in getActionList()"
                :key="item.name"
                :label="item.name"
                :value="item.name"/>
          </ElSelect>
        </ElFormItem>
      </ElCol>

      <ElCol :span="12" v-if="!currentEntity">
        <ElFormItem :label="$t('entityAction.action')" prop="action">
          <ElInput v-model="action.action" clearable @change="changedAction($event)"
                   :placeholder=" $t('common.inputText')"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow :gutter="24">
      <ElCol :span="12">
        <ElFormItem :label="$t('entityAction.tags')" prop="tags">
          <TagsSearch v-model="action.tags" @change="changedTags($event)"/>
        </ElFormItem>
      </ElCol>
      <ElCol :span="12">
        <ElFormItem :label="$t('entityAction.area')" prop="area">
          <AreaSearch v-model="currentArea" @change="changedArea($event)"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

  </ElForm>

</template>

<style lang="less">

</style>
