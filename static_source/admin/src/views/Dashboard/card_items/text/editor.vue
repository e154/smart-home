<script setup lang="ts">
import {computed, onBeforeUnmount, onMounted, PropType, ref, unref, watch} from "vue";
import {Card, CardItem, comparisonType, Core, requestCurrentState, Tab} from "@/views/Dashboard/core";
import {ElDivider, ElCollapse, ElCollapseItem, ElCard, ElForm, ElFormItem, ElPopconfirm, ElSwitch,
  ElRow, ElCol, ElSelect, ElOption, ElInput, ElTag, ElButton } from 'element-plus'
// import { Editor, EditorExpose } from '@/components/Editor'
// import { IDomEditor } from '@wangeditor/editor'
import CommonEditor from "@/views/Dashboard/card_items/common/editor.vue";
import {useI18n} from "@/hooks/web/useI18n";
import {Cache, GetTokens} from "@/views/Dashboard/render";
import {TextProp} from "@/views/Dashboard/card_items/text/types";
import Viewer from "@/components/JsonViewer/JsonViewer.vue";
import { TinycmeEditor } from "@/components/Tinymce";

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
  set(val: CardItem) {}
})

const defaultTextHtml = ref(true);

// ---------------------------------
// component methods
// ---------------------------------

const tokens = ref<string[]>([]);

onMounted(()=>{
  update()
})

onBeforeUnmount(() => {

})

const update = () => {
  updateTokensDefaultText()

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

const updateCurrentState = () => {
  if (currentItem.value.entityId) {
    requestCurrentState(currentItem.value?.entityId)
  }
}
</script>

<template>

  <CommonEditor :item="currentItem" :core="core"/>

  <ElDivider content-position="left">{{ $t('dashboard.editor.textOptions') }}</ElDivider>

  <ElRow style="padding-bottom: 20px">
    <ElCol>
      <div style="padding-bottom: 20px">
        <ElButton type="default" @click.prevent.stop="addProp()">
          <Icon icon="ep:plus" class="mr-5px"/>
          {{ $t('dashboard.editor.addProp') }}
        </ElButton>
      </div>

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

              <ElRow :gutter="24">
                <ElCol
                    :span="8"
                    :xs="8"
                >
                  <ElFormItem :label="$t('dashboard.editor.text')" prop="text">
                    <ElInput
                        placeholder="Please input"
                        v-model="prop.key"/>
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
                    <ElInput
                        placeholder="Please input"
                        v-model="prop.value"/>
                  </ElFormItem>

                </ElCol>
              </ElRow>

              <ElRow>
                <ElCol>
                  <ElFormItem :label="$t('dashboard.editor.html')" prop="enabled">
                    <ElSwitch v-model="defaultTextHtml"/>
                  </ElFormItem>

                  <ElFormItem :label="$t('dashboard.editor.text')" prop="text">
                    <ElInput
                        v-if="!defaultTextHtml"
                        type="textarea"
                        :autosize="{minRows: 10}"
                        placeholder="Please input"
                        v-model="prop.text"
                        @change="propTextUpdated(prop)"
                    />
<!--                    <Editor v-else v-model="prop.text" @change="propTextUpdated(prop)" />-->
                    <TinycmeEditor v-else v-model="prop.text" @text-change="propTextUpdated(prop)"/>
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
                          @confirm="removeProp"
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

  <ElRow>
    <ElCol>
      <ElFormItem :label="$t('dashboard.editor.html')" prop="enabled">
        <ElSwitch v-model="defaultTextHtml"/>
      </ElFormItem>

      <ElFormItem :label="$t('dashboard.editor.textBody')" prop="text">
        <ElInput
            v-if="!defaultTextHtml"
            type="textarea"
            :autosize="{minRows: 10}"
            placeholder="Please input"
            v-model="currentItem.payload.text.default_text"
            @change="defaultTextUpdated"
        />
<!--        <Editor v-else mode="simple" v-model="currentItem.payload.text.default_text" @change="defaultTextUpdated" />-->
        <TinycmeEditor v-else v-model="currentItem.payload.text.default_text" @text-change="defaultTextUpdated"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow>
    <ElCol>
      <ElFormItem :label="$t('dashboard.editor.tokens')">
        <ElTag size="small" v-for="(token, index) in tokens" :key="index" class="mr-10px">{{ token }}</ElTag>
        <div v-if="!tokens.length">{{$t('main.no')}}</div>
      </ElFormItem>
    </ElCol>
  </ElRow>

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
