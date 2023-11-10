<script setup lang="ts">
import {computed, onMounted, PropType, ref, unref, watch} from "vue";
import {Card, CardItem, CompareProp, comparisonType, Core, Tab} from "@/views/Dashboard/core";
import {ElDivider, ElOption, ElCollapse, ElFormItem, ElSwitch, ElCol, ElRow, ElButton, ElSelect,
  ElPopconfirm, ElForm, ElCard, ElCollapseItem, ElInput, ElTag} from 'element-plus'
import {useI18n} from "@/hooks/web/useI18n";
import {ApiEntity} from "@/api/stub";
import EntitySearch from "@/views/Entities/components/EntitySearch.vue";

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
})

const currentValue = computed({
  get(): CompareProp[] {
    return props.modelValue as CompareProp[]
  },
  set(val: CompareProp[]) {}
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
  if (entity?.id) {
    currentValue.value[index].entity = await currentCore.value.fetchEntity(entity.id);
    currentValue.value[index].entityId = entity.id;
  } else {
    currentValue.value[index].entity = undefined;
    currentValue.value[index].entityId = '';
  }
}

const removeShowOnProp = (index: number) => {
  currentValue.value.splice(index, 1);
}
</script>

<template>

  <ElRow>
    <ElCol>
      <div style="padding-bottom: 20px">
        <ElButton type="default" @click.prevent.stop="addShowOnProp()">
          <Icon icon="ep:plus" class="mr-5px"/>
          {{ $t('dashboard.editor.addNewProp') }}
        </ElButton>
      </div>

      <!-- props -->
      <ElCollapse>
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

              <ElRow :gutter="24">
                <ElCol
                    :span="8"
                    :xs="8"
                >
                  <ElFormItem :label="$t('dashboard.editor.text')" prop="text">
                    <ElInput placeholder="Please input" v-model="prop.key"/>
                  </ElFormItem>

                </ElCol>

                <ElCol
                    :span="8"
                    :xs="8"
                >
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

                <ElCol
                    :span="8"
                    :xs="8"
                >

                  <ElFormItem :label="$t('dashboard.editor.value')" prop="value">
                    <ElInput placeholder="Please input" v-model="prop.value"/>
                  </ElFormItem>

                </ElCol>
              </ElRow>

              <ElRow>
                <ElCol :span="12" :xs="12">
                  <ElFormItem :label="$t('dashboard.editor.entity')" prop="entity">
                    <EntitySearch v-model="prop.entity" @change="onEntityChanged($event, index)"/>
                  </ElFormItem>
                </ElCol>
              </ElRow>

              <ElRow>
                <ElCol>
                  <div style="padding-bottom: 20px">
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
                          <ElButton class="mr-10px" type="danger" plain>
                            <Icon icon="ep:delete" class="mr-5px"/>
                            {{ t('main.remove') }}
                          </ElButton>
                        </template>
                      </ElPopconfirm>
                    </div>
                  </div>
                </ElCol>
              </ElRow>

            </ElForm>

          </ElCard>

        </ElCollapseItem>
      </ElCollapse>
      <!-- /props -->

    </ElCol>
  </ElRow>

</template>

<style lang="less" >

</style>
