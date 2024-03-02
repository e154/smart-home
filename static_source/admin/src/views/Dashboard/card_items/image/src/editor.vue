<script setup lang="ts">
import {computed, PropType} from "vue";
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
  ElOption,
  ElPopconfirm,
  ElRow,
  ElSelect,
  ElSwitch,
  ElTag
} from 'element-plus'
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {Cache} from "@/views/Dashboard/core/cache";
import {ApiImage} from "@/api/stub";
import {ImageSearch} from "@/components/ImageSearch";
import {KeysSearch} from "@/views/Dashboard/components";
import {comparisonType} from "@/views/Dashboard/core/types";
import {ItemPayloadImage} from "@/views/Dashboard/card_items/image";
import {useI18n} from "@/hooks/web/useI18n";

const {t} = useI18n()

// ---------------------------------
// common
// ---------------------------------

const _cache: Cache = new Cache();

const props = defineProps({
  core: {
    type: Object as PropType<Core>,
  },
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const currentItem = computed({
  get(): CardItem {
    return props.item as CardItem
  },
  set(val: CardItem) {
  }
})

// ---------------------------------
// component methods
// ---------------------------------

const onSelectImage = (image: ApiImage) => {
  currentItem.value.payload.image.image = image || undefined;
}

const onChangePropValue = (val: string) => {
  currentItem.value.payload.image.attrField = val;
}

const initDefaultValue = () => {
  currentItem.value.payload.image = {
    items: [],
    image: undefined,
    background: false,
    attrField: undefined,
  } as ItemPayloadImage;
}

const onChangeValue = (val) => {
  currentItem.value.payload.image.attrField = val;
}

const addProp = () => {
  // console.log('add prop');

  if (!props.item.payload.image?.items) {
    currentItem.value.payload.image = {
      items: []
    };
  }

  let counter = 0;
  if (props.item.payload.image.items.length) {
    counter = props.item.payload.image.items.length;
  }

  currentItem.value.payload.image.items.push({

    key: 'new proper ' + counter,
    value: '',
    comparison: comparisonType.EQ,
    image: undefined
  });
  props.item.update();
}

const removeProp = (index: number) => {
  if (!props.item.payload.image) {
    initDefaultValue();
  }

  props.item.payload.image.items!.splice(index, 1);
  props.item.update();
}

const onChangePropKey = (index, key) => {
  currentItem.value.payload.image.items[index].key = key;
}


const onChangePropImage = (index, image) => {
  currentItem.value.payload.image.items[index].image = image || null;
}

</script>

<template>

  <CommonEditor :item="currentItem" :core="core"/>

  <ElRow class="mb-10px mt-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.editor.imageOptions') }}</ElDivider>
    </ElCol>
  </ElRow>

  <ElRow>
    <ElCol>
      <div class="mb-10px">
        <ElButton class="w-[100%]" @click.prevent.stop="addProp()">
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
      v-for="(prop, index) in item.payload.image.items"
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
              <ElFormItem :label="$t('dashboard.editor.attrField')" prop="text">
                <KeysSearch v-model="prop.key" :obj="currentItem.lastEvent" @change="onChangePropKey(index, $event)"/>
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

          <ElRow>
            <ElCol>
              <ElDivider content-position="left">{{ $t('dashboard.editor.appearanceOptions') }}</ElDivider>
            </ElCol>
          </ElRow>


          <ElRow>
            <ElCol>
              <ElFormItem :label="$t('dashboard.editor.background')" prop="background">
                <ElSwitch v-model="prop.background"/>
              </ElFormItem>
            </ElCol>
          </ElRow>

          <ElFormItem :label="$t('dashboard.editor.image')" prop="image">
            <ImageSearch v-model="prop.image" @change="onChangePropImage(index, $event)"/>
          </ElFormItem>

          <ElRow class="mb-10px mt-10px">
            <ElCol>
              <ElDivider content-position="left">{{ $t('main.or') }}</ElDivider>
            </ElCol>
          </ElRow>

          <ElFormItem :label="$t('dashboard.editor.attrField')" prop="text">
            <KeysSearch v-model="prop.attrField" :obj="currentItem.lastEvent" @change="onChangePropValue"/>
          </ElFormItem>

          <ElRow>
            <ElCol>
              <div class="mb-10px">
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
                      <ElButton type="danger" plain>
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

  <ElRow class="mb-10px mt-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.editor.defaultImage') }}</ElDivider>
    </ElCol>
  </ElRow>

  <ElRow>
    <ElCol>
      <ElFormItem :label="$t('dashboard.editor.background')" prop="background">
        <ElSwitch v-model="currentItem.payload.image.background"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElFormItem :label="$t('dashboard.editor.image')" prop="image">
    <ImageSearch v-model="currentItem.payload.image.image" @change="onSelectImage"/>
  </ElFormItem>

  <ElRow class="mb-10px mt-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('main.or') }}</ElDivider>
    </ElCol>
  </ElRow>

  <ElFormItem :label="$t('dashboard.editor.attrField')" prop="text">
    <KeysSearch v-model="currentItem.payload.image.attrField" :obj="currentItem.lastEvent" @change="onChangeValue"/>
  </ElFormItem>

</template>

<style lang="less">

</style>
