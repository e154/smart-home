<script setup lang="ts">
import {computed, PropType} from "vue";
import {CardItem, Core, Cache, comparisonType} from "@/views/Dashboard/core";
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
  ElTag,
} from 'element-plus'
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {useI18n} from "@/hooks/web/useI18n";
import {KeysSearch} from "@/views/Dashboard/components";
import {ColorPicker} from "@/components/ColorPicker";
import {ItemPayloadIcon} from "@/views/Dashboard/card_items/icon";

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

const initDefaultValue = () => {
  currentItem.value.payload.icon = {
    items: [],
    value: undefined,
    attrField: undefined,
    iconColor: "#eee",
  } as ItemPayloadIcon;
}

const onChangeValue = (val) => {
  currentItem.value.payload.icon.attrField = val;
}

const addProp = () => {
  // console.log('add prop');

  if (!props.item.payload.icon?.items) {
    currentItem.value.payload.icon = {
      items: []
    };
  }

  let counter = 0;
  if (props.item.payload.icon.items.length) {
    counter = props.item.payload.icon.items.length;
  }

  currentItem.value.payload.icon.items.push({

    key: 'new proper ' + counter,
    value: '',
    comparison: comparisonType.EQ,
    image: undefined
  });
  props.item.update();
}

const removeProp = (index: number) => {
  if (!props.item.payload.icon) {
    initDefaultValue();
  }

  props.item.payload.icon.items!.splice(index, 1);
  props.item.update();
}

const onChangePropKey = (val, index, event) => {
  currentItem.value.payload.icon.items[index].key = event;
}

</script>

<template>

  <CommonEditor :item="currentItem" :core="core"/>

  <ElRow class="mb-10px mt-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.editor.iconOptions') }}</ElDivider>
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
      v-for="(prop, index) in item.payload.icon.items"
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
                <KeysSearch v-model="prop.key" :obj="currentItem.lastEvent"
                            @change="onChangePropKey(prop, index, $event)"/>
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
              <ElFormItem :label="$t('dashboard.editor.icon')" prop="icon">
                <ElInput v-model="prop.icon"/>
              </ElFormItem>
            </ElCol>
          </ElRow>

          <ElRow>
            <ElCol>
              <ElFormItem :label="$t('dashboard.editor.attrField')" prop="text">
                <KeysSearch v-model="prop.attrField" :obj="currentItem.lastEvent" @change="onChangeValue"/>
              </ElFormItem>
            </ElCol>
          </ElRow>

          <ElRow class="mb-10px">
            <ElCol>
              <ElFormItem :label="$t('dashboard.editor.iconColor')" prop="iconColor">
                <ColorPicker show-alpha v-model="prop.iconColor"/>
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

  <ElRow class="mb-10px mt-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.editor.defaultIcon') }}</ElDivider>
    </ElCol>
  </ElRow>

  <ElRow>
    <ElCol>
      <ElFormItem :label="$t('dashboard.editor.icon')" prop="icon">
        <ElInput v-model="currentItem.payload.icon.value"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElFormItem :label="$t('dashboard.editor.attrField')" prop="text">
    <KeysSearch v-model="currentItem.payload.icon.attrField" :obj="currentItem.lastEvent" @change="onChangeValue"/>
  </ElFormItem>

  <ElRow :gutter="24">
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.iconColor')" prop="iconColor">
        <ColorPicker show-alpha v-model="currentItem.payload.icon.iconColor"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

</template>

<style lang="less">

</style>
