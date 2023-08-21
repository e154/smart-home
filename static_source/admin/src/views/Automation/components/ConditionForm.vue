<script setup lang="ts">
import {PropType, reactive, watch} from 'vue'
import {Form} from '@/components/Form'
import {useI18n} from '@/hooks/web/useI18n'
import {useForm} from '@/hooks/web/useForm'
import {useValidator} from '@/hooks/web/useValidator'
import {FormSchema} from '@/types/form'
import {ApiCondition} from "@/api/stub";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const {t} = useI18n()

const props = defineProps({
  condition: {
    type: Object as PropType<Nullable<ApiCondition>>,
    default: () => null
  }
})

const rules = {
  name: [required()],
  script: [required()],
}

const schema = reactive<FormSchema[]>([
  {
    field: 'name',
    label: t('conditions.name'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('conditions.name')
    }
  },
  {
    field: 'script',
    label: t('conditions.script'),
    component: 'Script',
    colProps: {
      span: 24
    },
    value: null,
    componentProps: {
      placeholder: t('conditions.script')
    }
  },
])

const {setValues} = methods

watch(
    () => props.condition,
    (val) => {
      if (!val) return
      setValues(val)
    },
    {
      deep: true,
      immediate: true
    }
)

defineExpose({
  elFormRef,
  getFormData: methods.getFormData
})
</script>

<template>
  <Form
      :schema="schema"
      :rules="rules"
      label-position="top"
      @register="register"
  />
</template>
