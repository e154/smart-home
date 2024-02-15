<script setup lang="ts">
import {computed, onMounted, PropType} from "vue";
import {CardItem, Core} from "@/views/Dashboard/core/core";
import {
  ElButton,
  ElCard,
  ElCol,
  ElCollapse,
  ElCollapseItem,
  ElDivider,
  ElForm,
  ElFormItem,
  ElInput,
  ElPopconfirm,
  ElRow,
  ElSwitch
} from 'element-plus'
import {ApiEntity, ApiImage} from "@/api/stub";
import {EntitySearch} from "@/components/EntitySearch";
import ShowOn from "./show-on.vue";
import {useI18n} from "@/hooks/web/useI18n";
import {ImageSearch} from "@/components/ImageSearch";
import {EntitiesAction, EntitiesActionOptions} from "@/components/EntitiesAction";
import {ColorPicker} from "@/components/ColorPicker";

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

const changedForActionButton = async (options: EntitiesActionOptions, index: number) => {
  currentItem.value.buttonActions[index] = {
    entityId: options.entityId,
    action: options.action,
    tags: options.tags,
    areaId: options.areaId,
    image: currentItem.value.buttonActions[index].image || undefined,
  }
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
    action: 'ACTION NAME ' + currentItem.value.buttonActions.length || '',
    image: null,
    tags: [],
    areaId: undefined,
    icon: 'icomoon-free:switch'
  });
}

const removeAction = (index: number) => {
  if (!currentItem.value.buttonActions) {
    return;
  }

  currentItem.value.buttonActions.splice(index, 1);
}

const asButtonHandler = () => {
  if (currentItem.value.asButton) {

    if (!currentItem.value.buttonActions) {
      currentItem.value.buttonActions = []
    }

    if (currentItem.value.buttonActions.length == 0)
    currentItem.value.buttonActions.push({
      entity: undefined,
      entityId: currentItem.value.entityId,
      action: 'ACTION NAME ' + currentItem.value.buttonActions.length || '',
      tags: [],
      areaId: undefined,
    });
  }
}

</script>

<template>

  <!-- main options -->
  <ElRow class="mb-10px mt-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.editor.mainOptions') }}</ElDivider>
    </ElCol>
  </ElRow>

  <ElRow>
    <ElCol>
      <ElFormItem :label="$t('dashboard.editor.entity')" prop="entity">
        <EntitySearch v-model="currentItem.entity" @change="changedEntity"/>
      </ElFormItem>

      <ElFormItem :label="$t('dashboard.editor.frozen')" prop="frozen">
        <ElSwitch v-model="currentItem.frozen"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow>
    <ElCol>
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
  <ElRow class="mb-10px mt-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.editor.showOn') }}</ElDivider>
    </ElCol>
  </ElRow>

  <ShowOn v-model="currentItem.showOn" :item="currentItem" :core="core"/>
  <!-- /show on -->

  <!-- hide on-->
  <ElRow class="mb-10px mt-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.editor.hideOn') }}</ElDivider>
    </ElCol>
  </ElRow>

  <ShowOn v-model="currentItem.hideOn" :item="currentItem" :core="core"/>
  <!-- /hide on-->

  <!-- menu options -->
  <div
    v-if="['icon', 'image'].includes(item.type)">
    <ElDivider content-position="left">{{ $t('dashboard.editor.buttonOptions') }}</ElDivider>
    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.asButton')" prop="enabled">
          <ElSwitch v-model="currentItem.asButton"/>
        </ElFormItem>
      </ElCol>
      <ElCol :span="12" :xs="12" v-if="currentItem.asButton">
        <ElFormItem label="&nbsp;" prop="addAction">
          <ElButton class="w-[100%]" @click.prevent.stop="addAction()">
            <Icon icon="ep:plus" class="mr-5px"/>
            {{ $t('dashboard.editor.addAction') }}
          </ElButton>
        </ElFormItem>

      </ElCol>
    </ElRow>

    <ElRow v-if="currentItem.asButton">
      <ElCol>
        <!-- props -->
        <ElCollapse>
          <ElCollapseItem
            v-for="(prop, index) in item.buttonActions"
            :name="index"
            :key="index"
          >

            <template #title>
              {{ prop.action }}
            </template>

            <ElCard shadow="never" class="item-card-editor">

              <ElForm
                label-position="top"
                :model="prop"
                style="width: 100%"
                ref="cardItemForm">

                <ElRow class="mb-10px mt-10px">
                  <ElCol>
                    <ElDivider content-position="left">{{ $t('dashboard.editor.actionOptions') }}</ElDivider>
                  </ElCol>
                </ElRow>

                <EntitiesAction :options="prop" :entity="currentItem.entity"
                                @change="changedForActionButton($event, index)"/>

                <ElDivider content-position="left">{{ $t('dashboard.editor.appearanceOptions') }}</ElDivider>

                <ElFormItem :label="$t('dashboard.editor.image')" prop="image">
                  <ImageSearch v-model="prop.image" @change="onSelectImageForAction(index, ...arguments)"/>
                </ElFormItem>

                <ElDivider content-position="left">{{ $t('main.or') }}</ElDivider>

                <ElRow>
                  <ElCol>
                    <ElFormItem :label="$t('dashboard.editor.icon')" prop="icon">
                      <ElInput v-model="prop.icon"/>
                    </ElFormItem>
                  </ElCol>
                </ElRow>

                <ElRow>
                  <ElCol>
                    <ElFormItem :label="$t('dashboard.editor.iconColor')" prop="iconColor">
                      <ColorPicker show-alpha v-model="prop.iconColor"/>
                    </ElFormItem>
                  </ElCol>
                </ElRow>

                <ElRow>
                  <ElCol>
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

      </ElCol>
    </ElRow>
  </div>
  <!-- /menu options -->

  <!-- button options -->
  <div
      v-if="['text'].includes(item.type)">
    <ElDivider content-position="left">{{ $t('dashboard.editor.buttonOptions') }}</ElDivider>
    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.asButton')" prop="enabled">
          <ElSwitch v-model="currentItem.asButton" @click="asButtonHandler"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow v-if="currentItem.asButton">
      <ElCol>
        <!-- props -->
                <ElRow class="mb-10px mt-10px">
                  <ElCol>
                    <ElDivider content-position="left">{{ $t('dashboard.editor.actionOptions') }}</ElDivider>
                  </ElCol>
                </ElRow>

                <EntitiesAction :options="item.buttonActions[0]" :entity="currentItem.entity"
                                @change="changedForActionButton($event, 0)"/>
        <!-- /props -->

      </ElCol>
    </ElRow>
  </div>
  <!-- /button options -->


</template>

<style lang="less">

</style>
