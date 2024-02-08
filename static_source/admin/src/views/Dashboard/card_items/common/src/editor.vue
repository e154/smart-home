<script setup lang="ts">
import {computed, onMounted, PropType} from "vue";
import {CardItem, Core} from "@/views/Dashboard/core/core";
import {
  ElButton,
  ElCard,
  ElCol,
  ElCollapse,
  ElCollapseItem,
  ElColorPicker,
  ElDivider,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElOption,
  ElPopconfirm,
  ElRow,
  ElSelect,
  ElSwitch
} from 'element-plus'
import {ApiEntity, ApiImage} from "@/api/stub";
import {EntitySearch} from "@/components/EntitySearch";
import ShowOn from "./show-on.vue";
import {useI18n} from "@/hooks/web/useI18n";
import {ImageSearch} from "@/components/ImageSearch";

const {t} = useI18n()

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

const currentCore = computed(() => props.core as Core)

const currentItem = computed(() => props.item as CardItem)

// ---------------------------------
// component methods
// ---------------------------------

onMounted(() => {
  setTimeout(() => {
    if (props.item?.entityId) {
      fetchEntity(props.item.entityId);
    }
  }, 1000);
})

const fetchEntity = async (id: string) => {
  const entity = await currentCore.value.fetchEntity(id);
  currentItem.value.entity = entity;
}

const changedEntity = (entity: ApiEntity, event?: any) => {
  if (!entity?.id) {
    currentItem.value.entity = undefined;
    return;
  }
  fetchEntity(entity.id);
}

const changedForActionButton = async (entity: ApiEntity, index: number) => {
  if (entity?.id) {
    currentItem.value.buttonActions[index].entity = await currentCore.value.fetchEntity(entity.id);
    currentItem.value.buttonActions[index].entityId = entity.id;
  } else {
    currentItem.value.buttonActions[index].entity = undefined;
    currentItem.value.buttonActions[index].entityId = '';
    currentItem.value.buttonActions[index].action = '';
  }
}

const updateButtonActions = () => {
  for (const index in currentItem.value.buttonActions) {
    changedForActionButton(currentItem.value.buttonActions[index].entity, index)
  }
}

updateButtonActions()

const getActionList = (entity?: ApiEntity) => {
  if (!entity) {
    return [];
  }
  return entity.actions;
}

const onSelectImageForAction = (index: number, image: ApiImage) => {
  // console.log('select image', index, image);
  if (!currentItem.value.buttonActions[index]) {
    return;
  }

  currentItem.value.buttonActions[index].image = image as ApiImage || undefined;
}

const addAction = () => {
  currentItem.value.buttonActions.push({
    entity: undefined,
    entityId: currentItem.value.entityId,
    action: '',
    image: null,
  });
}

const removeAction = (index: number) => {
  if (!currentItem.value.buttonActions) {
    return;
  }

  currentItem.value.buttonActions.splice(index, 1);
}
</script>

<template>

  <!-- main options -->
  <ElDivider content-position="left">{{ $t('dashboard.editor.mainOptions') }}</ElDivider>

  <ElRow :gutter="24">
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.entity')" prop="entity">
        <EntitySearch v-model="currentItem.entity" @change="changedEntity"/>
      </ElFormItem>

      <ElFormItem :label="$t('dashboard.editor.frozen')" prop="frozen">
        <ElSwitch v-model="currentItem.frozen"/>
      </ElFormItem>
    </ElCol>

    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.enabled')" prop="enabled">
        <ElSwitch v-model="currentItem.enabled"/>
      </ElFormItem>
      <ElFormItem :label="$t('dashboard.editor.hidden')" prop="hidden">
        <ElSwitch v-model="currentItem.hidden"/>
      </ElFormItem>
    </ElCol>
  </ElRow>
  <!-- /main options -->

  <!-- show on -->
  <ElDivider content-position="left">{{ $t('dashboard.editor.showOn') }}</ElDivider>
  <ShowOn v-model="currentItem.showOn" :item="currentItem" :core="core"/>
  <!-- /show on -->

  <!-- hide on-->
  <ElDivider content-position="left">{{ $t('dashboard.editor.hideOn') }}</ElDivider>
  <ShowOn v-model="currentItem.hideOn" :item="currentItem" :core="core"/>
  <!-- /hide on-->

  <!-- button options -->
  <div
      v-if="!['button', 'chart', 'chart_custom', 'chartCustom', 'map', 'slider',
      'streamPlayer', 'tiles', 'grid', 'progress', 'colorPicker'].includes(item.type)">
    <ElDivider content-position="left">{{ $t('dashboard.editor.buttonOptions') }}</ElDivider>
    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.asButton')" prop="enabled">
          <ElSwitch v-model="currentItem.asButton"/>
        </ElFormItem>
      </ElCol>
    </ElRow>
    <ElRow v-if="currentItem.asButton">
      <ElCol>
        <div style="padding-bottom: 20px">
          <ElButton type="default" @click.prevent.stop="addAction()">
            <Icon icon="ep:plus" class="mr-5px"/>
            {{ $t('dashboard.editor.addAction') }}
          </ElButton>
        </div>

        <!-- props -->
        <ElCollapse>
          <ElCollapseItem
              v-for="(prop, index) in item.buttonActions"
              :name="index"
              :key="index"
          >

            <template #title>
              {{ prop.entityId }} - {{ prop.action }}
            </template>

            <ElCard shadow="never" class="item-card-editor">

              <ElForm
                  label-position="top"
                  :model="prop"
                  style="width: 100%"
                  ref="cardItemForm">

                <ElRow :gutter="24">
                  <ElCol :span="12" :xs="12">
                    <ElFormItem :label="$t('dashboard.editor.entity')" prop="entity">
                      <EntitySearch v-model="prop.entity" @change="changedForActionButton($event, index)"/>
                    </ElFormItem>
                  </ElCol>

                  <ElCol :span="12" :xs="12">
                    <ElFormItem :label="$t('dashboard.editor.action')" prop="action" :aria-disabled="!item.entity">
                      <ElSelect
                          v-model="prop.action"
                          clearable
                          :placeholder="$t('dashboard.editor.selectAction')"
                          style="width: 100%"
                      >
                        <ElOption
                            v-for="item in getActionList(prop.entity)"
                            :key="item.name"
                            :label="item.name"
                            :value="item.name"/>
                      </ElSelect>
                    </ElFormItem>
                  </ElCol>
                </ElRow>

                <ElFormItem :label="$t('dashboard.editor.image')" prop="image">
                  <ImageSearch v-model="prop.image" @change="onSelectImageForAction(index, ...arguments)"/>
                </ElFormItem>

                <ElRow :gutter="24">
                  <ElCol :span="8" :xs="8">
                    <ElFormItem :label="$t('dashboard.editor.icon')" prop="icon">
                      <ElInput v-model="prop.icon"/>
                    </ElFormItem>
                  </ElCol>
                  <ElCol :span="8" :xs="8">
                    <ElFormItem :label="$t('dashboard.editor.iconColor')" prop="iconColor">
                      <ElColorPicker show-alpha v-model="prop.iconColor"/>
                    </ElFormItem>
                  </ElCol>
                  <ElCol :span="8" :xs="8">
                    <ElFormItem :label="$t('dashboard.editor.iconSize')" prop="iconSize">
                      <ElInputNumber v-model="prop.iconSize" :min="1" :value-on-clear="12"/>
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
                            @confirm="removeAction(index)"
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
  </div>
  <!-- /button options -->


</template>

<style lang="less">

</style>
