<script setup lang="ts">
import {computed, PropType, watch} from "vue";
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
  ElInputNumber,
  ElOption,
  ElPopconfirm,
  ElRow,
  ElSelect,
  ElTag
} from 'element-plus'
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {useI18n} from "@/hooks/web/useI18n";
import {Cache} from "@/views/Dashboard/core/cache";
import {ApiImage} from "@/api/stub";
import {ItemPayloadState} from "./types";
import {ImageSearch} from "@/components/ImageSearch";
import {KeysSearch} from "@/views/Dashboard/components";
import {comparisonType} from "@/views/Dashboard/core/types";
import {ColorPicker} from "@/components/ColorPicker";

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

const currentItem = computed(() => props.item as CardItem)

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
    initDefaultValue();
  }

  props.item.payload.state.items!.splice(index, 1);
  props.item.update();
}

const onSelectImageForState = (index: number, image: ApiImage) => {
  // console.log('select image', index, image);

  if (!props.item.payload.state) {
    initDefaultValue();
  }

  currentItem.value.payload.state.items[index].image = image as ApiImage || undefined;
  props.item.update();
}

const onSelectIconForState = (index: number, icon: string) => {
  // console.log('select icon', index, icon);

  if (!props.item.payload.state) {
    initDefaultValue();
  }

  currentItem.value.payload.state.items[index].icon = icon || undefined;
  props.item.update();
}

const reloadKeyDefaultImage = 0;

const onSelectDefaultImage = (image: ApiImage) => {
  // console.log('select image', image);

  if (!props.item.payload.state) {
    initDefaultValue();
  }

  currentItem.value.payload.state.defaultImage = image as ApiImage || undefined;
  // this.reloadKeyDefaultImage += 1
  props.item.update();
}

const onSelectDefaultIcon = (icon?: string) => {
  // console.log('select icon', icon);

  if (!props.item.payload.state) {
    initDefaultValue();
  }

  currentItem.value.payload.state.defaultIcon = icon || undefined;
  // this.reloadKeyDefaultImage += 1
  props.item.update();
}

const initDefaultValue = () => {
  currentItem.value.payload.state = {
    items: [],
    defaultImage: undefined,
    defaultIcon: undefined,
    defaultIconColor: "#000000",
    defaultIconSize: 12,
  } as ItemPayloadState;
}
// ---------------------------------
// component methods
// ---------------------------------


const onChangePropKey = (val, index) => {
  currentItem.value.payload.state.items[index].key = val.key;
}

</script>

<template>

  <CommonEditor :item="item" :core="core"/>

  <ElRow class="mt-10px mb-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.editor.stateOptions') }}</ElDivider>
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

          <ElRow>
            <ElCol>
              <ElFormItem :label="$t('dashboard.editor.attrField')" prop="text">
                <KeysSearch v-model="prop.key" :obj="currentItem.lastEvent" @change="onChangePropKey(prop, index)"/>
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
              <ElFormItem :label="$t('dashboard.editor.image')" prop="image">
                <ImageSearch v-model="prop.image" @change="onSelectImageForState(index, $event)"/>
              </ElFormItem>
            </ElCol>
          </ElRow>

          <ElDivider content-position="left">{{ $t('dashboard.editor.iconOptions') }}</ElDivider>

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
              <ElFormItem :label="$t('dashboard.editor.iconSize')" prop="iconSize">
                <ElInputNumber v-model="prop.iconSize" :min="1" :value-on-clear="12"/>
              </ElFormItem>
            </ElCol>
          </ElRow>

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


  <ElRow>
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.editor.defaultImage') }}</ElDivider>
    </ElCol>
  </ElRow>
  <ElRow>
    <ElCol>
      <ElFormItem :label="$t('dashboard.editor.image')" prop="image">
        <ImageSearch :key="reloadKeyDefaultImage" v-model="currentItem.payload.state.defaultImage"
                     @change="onSelectDefaultImage"/>
      </ElFormItem>
    </ElCol>
  </ElRow>
  <ElRow>
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.editor.defaultIcon') }}</ElDivider>
    </ElCol>
  </ElRow>
  <ElRow>
    <ElCol>
      <ElFormItem :label="$t('dashboard.editor.icon')" prop="icon">
        <ElInput v-model="currentItem.payload.state.defaultIcon" @change="onSelectDefaultIcon"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow>
    <ElCol>
      <ElFormItem :label="$t('dashboard.editor.iconColor')" prop="iconColor">
        <ColorPicker show-alpha v-model="currentItem.payload.state.defaultIconColor"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow class="mb-10px">
    <ElCol>
      <ElFormItem :label="$t('dashboard.editor.iconSize')" prop="iconSize">
        <ElInputNumber v-model="currentItem.payload.state.defaultIconSize" :min="1" :value-on-clear="12"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

</template>

<style lang="less">

</style>
