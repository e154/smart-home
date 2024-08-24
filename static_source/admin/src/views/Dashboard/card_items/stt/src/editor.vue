<script setup lang="ts">
import {computed, onBeforeUnmount, onMounted, PropType, ref, watch} from "vue";
import {Cache, CardItem, comparisonType, Core, GetTokens} from "@/views/Dashboard/core";
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
import {useI18n} from "@/hooks/web/useI18n";
import {TextProp} from "./types";
import {TinycmeEditor} from "@/components/Tinymce";
import {KeysSearch} from "@/views/Dashboard/components";

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

const currentItem = computed(() => props.item as CardItem)

const defaultTextHtml = ref(false);

// ---------------------------------
// component methods
// ---------------------------------

const tokens = ref<string[]>([]);

onMounted(() => {
  update()
})

onBeforeUnmount(() => {

})

const update = () => {
  // updateTokensDefaultText()

  if (currentItem.value?.payload?.text?.default_text) {
    for (const prop of currentItem.value.payload.text.items) {
      updateTokensPropText(prop)
    }
  }
}

const updateTokensDefaultText = () => {
  if (!currentItem.value?.payload?.text?.default_text) {
    tokens.value = []
    return
  }

  const _tokens = GetTokens(currentItem.value.payload.text.default_text, _cache)
  tokens.value = _tokens || []
}

const defaultTextUpdated = () => {
  updateTokensDefaultText()
}

watch(
  () => props.item,
  () => {
    updateTokensDefaultText()
  },
  {
    immediate: true
  }
)


const updateTokensPropText = (prop: TextProp) => {
  prop.tokens = GetTokens(prop.text, _cache) || []
}

const propTextUpdated = (prop: TextProp) => {
  updateTokensPropText(prop)
}

const addProp = () => {
  // console.log('addProp')

  if (!currentItem.value?.payload.text?.items) {
    currentItem.value.payload.text.items = []
  }

  let counter = 0
  if (currentItem.value?.payload.text!.items.length) {
    counter = currentItem.value?.payload.text!.items.length
  }

  currentItem.value?.payload.text!.items.push({
    key: 'new proper ' + counter,
    value: '',
    comparison: comparisonType.EQ,
    text: ''
  })
}

const removeProp = (index: number) => {
  if (!currentItem.value?.payload.text?.items) {
    return
  }

  currentItem.value?.payload.text?.items.splice(index, 1)
}

const onChangePropValue = (val, index) => {
  currentItem.value.payload.text.items[index].key = val;
}

const getFonts = () : string[] => {
  if (!props.core?.getActiveTab) {
    return []
  }
  return props.core.getActiveTab?.fonts || []
}

</script>

<template>

  <CommonEditor :item="currentItem" :core="core"/>

  <!-- text options -->
  <ElRow class="mb-10px mt-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.editor.textOptions') }}</ElDivider>
    </ElCol>
  </ElRow>

  <ElRow>
    <ElCol>
      <ElButton class="w-[100%]" @click.prevent.stop="addProp()">
        <Icon icon="ep:plus" class="mr-5px"/>
        {{ $t('dashboard.editor.addProp') }}
      </ElButton>
    </ElCol>
  </ElRow>

  <ElRow>
    <ElCol>
      <!-- props -->
      <ElCollapse>
        <ElCollapseItem
          :name="index"
          :key="index"
          v-for="(prop, index) in currentItem.payload.text.items"
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
              ref="cardItemForm"
            >

              <ElRow>
                <ElCol>
                  <ElFormItem :label="$t('dashboard.editor.attrField')" prop="text">
                    <KeysSearch v-model="prop.key" :obj="currentItem.lastEvent"
                                @change="onChangePropValue($event, index)"/>
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
                    <ElInput
                      placeholder="Please input"
                      v-model="prop.value"/>
                  </ElFormItem>

                </ElCol>
              </ElRow>

              <ElRow>
                <ElCol>
                  <ElFormItem :label="$t('dashboard.editor.html')" prop="enabled">
                    <ElSwitch v-model="prop.defaultTextHtml"/>
                  </ElFormItem>
                </ElCol>
              </ElRow>

              <ElRow v-if="!prop.defaultTextHtml">
                <ElCol>
                  <ElFormItem :label="$t('dashboard.editor.text')" prop="text">
                    <ElInput
                      type="textarea"
                      :autosize="{minRows: 10}"
                      placeholder="Please input"
                      v-model="prop.text"
                      @update:modelValue="propTextUpdated(prop)"
                    />
                  </ElFormItem>
                </ElCol>
              </ElRow>

              <ElRow v-else>
                <ElCol>
                  <ElFormItem :label="$t('dashboard.editor.text')" prop="text">
                    <TinycmeEditor v-model="prop.text" @update:modelValue="propTextUpdated(prop)" :fonts="getFonts()"/>
                  </ElFormItem>
                </ElCol>
              </ElRow>

              <ElRow>
                <ElCol>
                  <ElFormItem :label="$t('dashboard.editor.tokens')">
                    <ElTag size="small" v-for="(token, idx) in prop.tokens" :key="idx">{{ token }}</ElTag>
                  </ElFormItem>
                </ElCol>
              </ElRow>


              <div style="text-align: right;">
                <ElPopconfirm
                  :confirm-button-text="$t('main.ok')"
                  :cancel-button-text="$t('main.no')"
                  width="250"
                  style="margin-left: 10px;"
                  :title="$t('main.are_you_sure_to_do_want_this?')"
                  @confirm="removeProp"
                >
                  <template #reference>
                    <ElButton type="danger" plain>
                      <Icon icon="ep:delete" class="mr-5px"/>
                      {{ t('main.remove') }}
                    </ElButton>
                  </template>
                </ElPopconfirm>
              </div>


            </ElForm>

          </ElCard>

        </ElCollapseItem>
      </ElCollapse>
      <!-- /props -->

    </ElCol>
  </ElRow>

  <ElRow>
    <ElCol>
      <ElFormItem :label="$t('dashboard.editor.html')" prop="enabled">
        <ElSwitch v-model="defaultTextHtml"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow v-if="!defaultTextHtml">
    <ElCol>
      <ElFormItem :label="$t('dashboard.editor.textBody')" prop="text">
        <ElInput
          type="textarea"
          :autosize="{minRows: 10}"
          placeholder="Please input"
          v-model="currentItem.payload.text.default_text"
          @update:modelValue="defaultTextUpdated"
        />
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow v-else>
    <ElCol>
      <ElFormItem :label="$t('dashboard.editor.textBody')" prop="text">
        <TinycmeEditor v-model="currentItem.payload.text.default_text" @update:modelValue="defaultTextUpdated" :fonts="getFonts()"/>
      </ElFormItem>
    </ElCol>
  </ElRow>
  <!-- /text options -->

  <ElRow>
    <ElCol class="tag-list">
      <ElFormItem :label="$t('dashboard.editor.tokens')">
        <ElTag size="small" v-for="(token, index) in tokens" :key="index" class="mr-10px">{{ token }}</ElTag>
        <div v-if="!tokens.length">{{ $t('main.no') }}</div>
      </ElFormItem>
    </ElCol>
  </ElRow>

</template>

<style lang="less" scoped>
:deep(.tag-list .el-tag--small) {
  margin: 0 7px 7px 0;
}
</style>
