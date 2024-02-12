<script setup lang="ts">
import {computed, PropType} from "vue";
import {CardItem, Core} from "@/views/Dashboard/core/core";
import {
  ElButton,
  ElCard,
  ElCol,
  ElCollapse,
  ElCollapseItem,
  ElForm,
  ElFormItem,
  ElInput,
  ElOption,
  ElPopconfirm,
  ElRow,
  ElSelect,
  ElTag
} from 'element-plus'
import {useI18n} from "@/hooks/web/useI18n";
import {ApiEntity} from "@/api/stub";
import {EntitySearch} from "@/components/EntitySearch";
import {KeysSearch} from "@/views/Dashboard/components";
import {CompareProp, comparisonType} from "@/views/Dashboard/core/types";
import {EventStateChange} from "@/api/types";

const {t} = useI18n()

// ---------------------------------
// common
// ---------------------------------
const props = defineProps({
  core: {
    type: Object as PropType<Nullable<Core>>,
    default: () => null
  },
  modelValue: {
    type: Array as PropType<CompareProp[]>,
    default: () => []
  },
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const currentItem = computed(() => props.item as CardItem)

const currentValue = computed({
  get(): CompareProp[] {
    return props.modelValue as CompareProp[]
  },
  set(val: CompareProp[]) {
  }
})

// ---------------------------------
// component methods
// ---------------------------------

const addShowOnProp = () => {

  if (!currentValue.value) {
    currentValue.value = [];
  }

  let counter = 0;
  if (currentValue.value.length) {
    counter = currentValue.value.length;
  }

  currentValue.value.push({
    key: 'key ' + counter,
    value: 'value',
    comparison: comparisonType.EQ
  });
}

const currentCore = computed(() => props.core as Core)

const onEntityChanged = async (entity: ApiEntity, index: number) => {
  currentValue.value[index].key = '';
  if (entity?.id) {
    currentValue.value[index].entity = await currentCore.value.fetchEntity(entity.id);
    currentValue.value[index].entityId = entity.id;
    currentItem.value.lastEvents(entity.id)
  } else {
    currentValue.value[index].entity = undefined;
    currentValue.value[index].entityId = '';
  }
}

const removeShowOnProp = (index: number) => {
  currentValue.value.splice(index, 1);
}

const onChangePropValue = (val, index) => {
  currentValue.value[index].key = val;
}

const lastEvent = (index: number): EventStateChange | undefined => {
  if (currentValue.value[index].entityId) {
    return currentItem.value.lastEvents(currentValue.value[index].entityId)
  } else {
    return currentItem.value.lastEvent
  }
}

</script>

<template>

  <ElRow>
    <ElCol>
      <ElButton class="w-[100%]" @click.prevent.stop="addShowOnProp()">
        <Icon icon="ep:plus" class="mr-5px"/>
        {{ $t('dashboard.editor.addNewProp') }}
      </ElButton>
    </ElCol>
  </ElRow>

  <!-- props -->
  <ElCollapse v-if="currentValue.length">
    <ElCollapseItem
        :name="index"
        :key="index"
        v-for="(prop, index) in currentValue"
    >

      <template #title>
        <ElTag size="small">{{ prop.key }}</ElTag>
        +
        <ElTag size="small">{{ prop.comparison }}</ElTag>
        +
        <ElTag size="small">{{ prop.value }}</ElTag>
      </template>

      <ElCard shadow="never" class="item-card-editor">

        <ElForm
            label-position="top"
            :model="prop"
            style="width: 100%"
            ref="cardItemForm">

          <ElRow>
            <ElCol>
              <ElFormItem :label="$t('dashboard.editor.entity')" prop="entity">
                <EntitySearch v-model="prop.entity" @change="onEntityChanged($event, index)"/>
              </ElFormItem>
            </ElCol>
          </ElRow>

          <ElRow>
            <ElCol>
              <ElFormItem :label="$t('dashboard.editor.attrField')" prop="text">
                <KeysSearch v-model="prop.key" :obj="lastEvent(index)" @change="onChangePropValue($event, index)"/>
              </ElFormItem>

            </ElCol>
          </ElRow>

          <ElRow>
            <ElCol>
              <ElFormItem :label="$t('dashboard.editor.comparison')" prop="comparison">
                <ElSelect
                    v-model="prop.comparison"
                    placeholder="please select type"
                    style="width: 100%"
                >
                  <ElOption label="==" value="eq"/>
                  <ElOption label="<" value="lt"/>
                  <ElOption label="<=" value="le"/>
                  <ElOption label="!=" value="ne"/>
                  <ElOption label=">=" value="ge"/>
                  <ElOption label=">" value="gt"/>
                </ElSelect>
              </ElFormItem>

            </ElCol>
          </ElRow>

          <ElRow>
            <ElCol>

              <ElFormItem :label="$t('dashboard.editor.value')" prop="value">
                <ElInput placeholder="Please input" v-model="prop.value"/>
              </ElFormItem>

            </ElCol>
          </ElRow>

          <ElRow class="mt-10px">
            <ElCol>
              <div style="text-align: right;">
                <ElPopconfirm
                    :confirm-button-text="$t('main.ok')"
                    :cancel-button-text="$t('main.no')"
                    width="250"
                    style="margin-left: 10px;"
                    :title="$t('main.are_you_sure_to_do_want_this?')"
                    @confirm="removeShowOnProp(index)"
                >
                  <template #reference>
                    <ElButton type="danger" plain>
                      <Icon icon="ep:delete" class="mr-5px"/>
                      {{ t('main.remove') }}
                    </ElButton>
                  </template>
                </ElPopconfirm>
              </div>
            </ElCol>
          </ElRow>

        </ElForm>

      </ElCard>

    </ElCollapseItem>
  </ElCollapse>
  <!-- /props -->


</template>

<style lang="less">

</style>
