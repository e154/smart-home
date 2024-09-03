<script setup lang="ts">
import {computed, PropType} from "vue";
import {CardItem, Core} from "@/views/Dashboard/core";
import {ElCol, ElDivider, ElFormItem, ElInput, ElOption, ElRow, ElSelect, ElSwitch} from 'element-plus'
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {EntitiesAction, EntitiesActionOptions} from "@/components/EntitiesAction";

// ---------------------------------
// common
// ---------------------------------

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

const currentItem = computed(() => props.item as CardItem)

// ---------------------------------
// component methods
// ---------------------------------

const changedForActionButton = async (options: EntitiesActionOptions) => {
  currentItem.value.payload.button.entityId = options.entityId
  currentItem.value.payload.button.action = options.action
  currentItem.value.payload.button.tags = options.tags
  currentItem.value.payload.button.areaId = options.areaId
  currentItem.value.payload.button.eventName = options.eventName
  currentItem.value.payload.button.eventArgs = options.eventArgs
}

</script>

<template>
  <div>

    <CommonEditor :item="currentItem" :core="core"/>

    <ElRow class="mb-10px mt-10px">
      <ElCol>
        <ElDivider content-position="left">{{ $t('dashboard.editor.buttonOptions') }}</ElDivider>
      </ElCol>
    </ElRow>

    <ElRow>
      <ElCol>
        <ElFormItem :label="$t('dashboard.editor.icon')" prop="icon">
          <ElInput v-model="currentItem.payload.button.icon"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow>
      <ElCol>
        <ElFormItem :label="$t('dashboard.editor.text')" prop="text">
          <ElInput v-model="currentItem.payload.button.text"/>
        </ElFormItem>
      </ElCol>
    </ElRow>


    <ElRow>
      <ElCol>
        <ElFormItem :label="$t('dashboard.editor.type')" prop="type">
          <ElSelect
              v-model="currentItem.payload.button.type"
              placeholder="please select type"
              style="width: 100%"
          >
            <ElOption label="Default" value=""/>
            <ElOption label="Primary" value="primary"/>
            <ElOption label="Success" value="success"/>
            <ElOption label="Info" value="info"/>
            <ElOption label="Warning" value="warning"/>
            <ElOption label="Danger" value="danger"/>
          </ElSelect>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow>
      <ElCol>

        <ElFormItem :label="$t('dashboard.editor.round')" prop="round">
          <ElSwitch v-model="currentItem.payload.button.round"/>
        </ElFormItem>

        <ElFormItem :label="$t('dashboard.editor.text')" prop="text">
          <ElSwitch v-model="currentItem.payload.button.asText"/>
        </ElFormItem>

      </ElCol>
    </ElRow>

    <ElDivider content-position="left">{{ $t('dashboard.editor.actionOptions') }}</ElDivider>

    <EntitiesAction :options="currentItem.payload.button" :entity="currentItem.entity"
                    @change="changedForActionButton($event)"/>

  </div>
</template>

<style lang="less">

</style>
