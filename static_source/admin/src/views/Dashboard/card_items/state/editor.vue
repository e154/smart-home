<script setup lang="ts">
import {computed, PropType, ref, unref, watch} from "vue";
import {Card, CardItem, comparisonType, Core, requestCurrentState, Tab} from "@/views/Dashboard/core";
import {ElDivider, ElCollapse, ElCollapseItem, ElCard, ElForm, ElFormItem, ElPopconfirm, ElSwitch,
  ElRow, ElCol, ElSelect, ElOption, ElInput, ElTag, ElButton } from 'element-plus'
import CommonEditor from "@/views/Dashboard/card_items/common/editor.vue";
import {useI18n} from "@/hooks/web/useI18n";
import {Cache} from "@/views/Dashboard/render";
import {ApiImage} from "@/api/stub";
import {ItemPayloadState} from "@/views/Dashboard/card_items/state/types";
import Viewer from "@/components/JsonViewer/JsonViewer.vue";
import ImageSearch from "@/views/Images/components/ImageSearch.vue";

const {t} = useI18n()

// ---------------------------------
// common
// ---------------------------------
const _cache: Cache = new Cache();

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

watch(
    () => props.item,
    (val?: CardItem) => {
      if (!val) return;

    },
    {
      deep: true,
      immediate: true
    }
)

const currentItem = computed({
  get(): CardItem {
    return props.item as CardItem
  },
  set(val: CardItem) {}
})

const addProp = () => {
  // console.log('add prop');

  if (!props.item.payload.state?.items) {
    currentItem.value.payload.state = {
      items: []
    };
  }

  let counter = 0;
  if (props.item.payload.state.items.length) {
    counter = props.item.payload.state.items.length;
  }

  currentItem.value.payload.state.items.push({

    key: 'new proper ' + counter,
    value: '',
    comparison: comparisonType.EQ,
    image: undefined
  });
  props.item.update();
}

const removeProp = (index: number) => {
  if (!props.item.payload.state) {
    currentItem.value.payload.state = {
      items: [],
      default_image: undefined
    } as ItemPayloadState;
  }

  props.item.payload.state.items!.splice(index, 1);
  props.item.update();
}

const onSelectImageForState = (index: number, image: ApiImage) => {
  console.log('select image', index, image);

  if (!props.item.payload.state) {
    currentItem.value.payload.state = {
      items: [],
      default_image: undefined
    } as ItemPayloadState;
  }

  currentItem.value.payload.state.items[index].image = image as ApiImage || undefined;
  props.item.update();
}

const reloadKeyDefaultImage = 0;

const onSelectDefaultImage = (image: ApiImage) => {
  console.log('select image', image);

  if (!props.item.payload.state) {
    currentItem.value.payload.state = {
      items: [],
      default_image: undefined
    } as ItemPayloadState;
  }

  currentItem.value.payload.state.default_image = image as ApiImage || undefined;
  // this.reloadKeyDefaultImage += 1
  props.item.update();
}
// ---------------------------------
// component methods
// ---------------------------------

const updateCurrentState = () => {
  if (currentItem.value.entityId) {
    requestCurrentState(currentItem.value?.entityId)
  }
}
</script>

<template>

  <CommonEditor :item="item" :core="core"/>

  <ElDivider content-position="left">{{ $t('dashboard.editor.stateOptions') }}</ElDivider>

  <ElRow>
    <ElCol>
      <div style="padding-bottom: 20px">
        <ElButton type="default" @click.prevent.stop="addProp()">
          <Icon icon="ep:plus" class="mr-5px"/>
          {{ $t('dashboard.editor.addNewProp') }}
        </ElButton>
      </div>
    </ElCol>
  </ElRow>

  <!-- props -->
  <ElCollapse>
    <ElCollapseItem
        :name="index"
        :key="index"
        v-for="(prop, index) in item.payload.state.items"
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
            <ElCol :span="8" :xs="8">
              <ElFormItem :label="$t('dashboard.editor.text')" prop="text">
                <ElInput placeholder="Please input" v-model="prop.key"/>
              </ElFormItem>
            </ElCol>

            <ElCol :span="8" :xs="8">
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

            <ElCol :span="8" :xs="8">
              <ElFormItem :label="$t('dashboard.editor.value')" prop="value">
                <ElInput placeholder="Please input" v-model="prop.value"/>
              </ElFormItem>
            </ElCol>
          </ElRow>

          <ElRow>
            <ElCol>
              <ElFormItem :label="$t('dashboard.editor.image')" prop="image">
                <ImageSearch v-model="prop.image" @change="onSelectImageForState(index, ...arguments)"/>
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
                      @confirm="removeProp(index)"
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
    <ElCollapseItem>
      <template #title>
        {{ $t('dashboard.editor.defaultImage') }}
      </template>
      <ElRow>
        <ElCol>
          <ElCard shadow="never" class="item-card-editor">
            <ImageSearch :key="reloadKeyDefaultImage" v-model="currentItem.payload.state.default_image" @change="onSelectDefaultImage"/>
          </ElCard>
        </ElCol>
      </ElRow>
    </ElCollapseItem>
  </ElCollapse>
  <!-- /props -->

  <ElRow style="padding-bottom: 20px" v-if="currentItem.entity">
    <ElCol>
      <ElCollapse>
        <ElCollapseItem :title="$t('dashboard.editor.eventstateJSONobject')">
          <ElButton type="default" @click.prevent.stop="updateCurrentState()" style="margin-bottom: 20px">
            <Icon icon="ep:refresh" class="mr-5px"/>
            {{ $t('dashboard.editor.getEvent') }}
          </ElButton>

          <Viewer v-model="currentItem.lastEvent"/>

        </ElCollapseItem>
      </ElCollapse>
    </ElCol>
  </ElRow>

</template>

<style lang="less" >

</style>
