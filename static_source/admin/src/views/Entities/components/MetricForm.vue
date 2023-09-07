<script setup lang="ts">
import api from "@/api/api";
import {h, PropType, reactive, ref, unref, watch} from 'vue'
import {ElButton, ElRow, ElCol, ElCard, ElForm,
  ElFormItem,  ElDivider, ElCollapse, ElCollapseItem,
  ElColorPicker, ElPopconfirm, ElInput, ElSkeleton} from 'element-plus'
import {useI18n} from '@/hooks/web/useI18n'
import {useForm} from '@/hooks/web/useForm'
import {useValidator} from '@/hooks/web/useValidator'
import {ApiMetric} from "@/api/stub";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const {t} = useI18n()

interface Current {
  item: ApiMetric
}

const current  = reactive<Current>({
  item: {} as ApiMetric
})

const props = defineProps({
  modelValue: {
    type: Object as PropType<Nullable<ApiMetric>>,
    default: () => null
  }
})

watch(
    () => props.modelValue,
    (val) => {
     if (val == unref(current)) return;
      current.item = val
    },
    {
      deep: true,
      immediate: true
    }
)

const rules = {
  name: [
    {required: true, trigger: 'blur'},
    {min: 4, max: 255, trigger: 'blur'}
  ],
  description: [
    {required: false, trigger: 'blur'},
    {max: 255, trigger: 'blur'}
  ],
  type: [
    {required: false, trigger: 'blur'},
    {max: 255, trigger: 'blur'}
  ]
};

const rules2 = {
  name: [
    {required: true, trigger: 'blur'},
    {min: 4, max: 255, trigger: 'blur'}
  ],
  description: [
    {required: false, trigger: 'blur'},
    {max: 255, trigger: 'blur'}
  ],
  // color: [
  //   {required: false, trigger: 'blur'},
  //   {max: 255, trigger: 'blur'}
  // ],
  translate: [
    {required: false, trigger: 'blur'},
    {max: 255, trigger: 'blur'}
  ],
  label: [
    {required: false, trigger: 'blur'},
    {max: 255, trigger: 'blur'}
  ]
};

const addProp = () => {

  current.item.options!.items!.push({
    name: 'new label',
    description: '',
    // color: '#FF0000',
    translate: '',
    label: ''
  });
}

const removeProp = (index: number) => {
  current.item.options!.items!.splice(index, 1);
}

</script>

<template>

  <ElForm v-if="current.item" label-position="top" label-width="100px" ref="current.item" :model="current.item" :rules="rules2" style="width: 100%">
    <ElFormItem :label="$t('metrics.name')" prop="name">
      <ElInput size="small" v-model="current.item.name"/>
    </ElFormItem>
    <ElFormItem :label="$t('metrics.description')" prop="description">
      <ElInput size="small" v-model="current.item.description"/>
    </ElFormItem>

  </ElForm>

  <ElDivider v-if="current.item" content-position="left">{{ $t('metrics.properties') }}</ElDivider>

  <div v-if="current.item" style="padding-bottom: 20px" >
    <ElButton type="default" @click.prevent.stop="addProp()">
      <Icon icon="ep:plus" class="mr-5px"/>
      {{ $t('metrics.addProperty') }}
    </ElButton>
  </div>

  <!-- props -->
  <ElSkeleton v-if="current.item && !current.item?.options?.items.length" :rows="5" />

  <ElCollapse v-if="current.item && current.item?.options?.items.length">
    <ElCollapseItem :title="prop.name" :name="index" :key="index" v-for="(prop, index) in current.item.options?.items">

<!--      <ElCard shadow="never" class="item-card-editor">-->
        <ElForm label-position="top" :model="prop" style="width: 100%" ref="cardItemForm">

          <ElRow :gutter="20">
            <ElCol>
              <ElFormItem :label="$t('metrics.name')" prop="name">
                <ElInput size="small" v-model="prop.name"/>
              </ElFormItem>

              <ElFormItem :label="$t('metrics.description')" prop="description">
                <ElInput size="small" v-model="prop.description"/>
              </ElFormItem>

<!--              <ElFormItem :label="$t('metrics.color')" prop="background">-->
<!--                <ElColorPicker show-alpha v-model="prop.color"/>-->
<!--              </ElFormItem>-->

              <ElFormItem :label="$t('metrics.translate')" prop="translate">
                <ElInput size="small" v-model="prop.translate"/>
              </ElFormItem>

              <ElFormItem :label="$t('metrics.label')" prop="label">
                <ElInput size="small" v-model="prop.label"/>
              </ElFormItem>

            </ElCol>
          </ElRow>

        </ElForm>
<!--      </ElCard>-->

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
            {{ t('metrics.removeProp') }}
          </ElButton>
        </template>
      </ElPopconfirm>

    </ElCollapseItem>
  </ElCollapse>
  <!-- /props -->

</template>
