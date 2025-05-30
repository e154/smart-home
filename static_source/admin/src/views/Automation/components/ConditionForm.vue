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
    label: t('automation.conditions.name'),
    component: 'Input',
    colProps: {
      span: 12
    },
    componentProps: {
      placeholder: t('automation.conditions.name')
    }
  },
  {
    field: 'description',
    label: t('automation.description'),
    component: 'Input',
    colProps: {
      span: 12
    },
    componentProps: {
      placeholder: t('automation.description')
    }
  },
  {
    field: 'area',
    label: t('automation.area'),
    value: null,
    component: 'Area',
    colProps: {
      span: 12
    },
    componentProps: {
      placeholder: t('automation.area'),
    }
  },
  {
    field: 'script',
    label: t('automation.conditions.script'),
    component: 'Script',
    colProps: {
      span: 24
    },
    value: null,
    componentProps: {
      placeholder: t('automation.conditions.script')
    }
  },
  {
    field: 'script',
    component: 'ScriptHelper',
    colProps: {
      span: 24
    },
    value: null,
    componentProps: {
      placeholder: t('automation.triggers.script'),
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
