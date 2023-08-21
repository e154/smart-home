<script setup lang="ts">
import { PropType, reactive, watch} from 'vue'
import {Form} from '@/components/Form'
import {useI18n} from '@/hooks/web/useI18n'
import {useForm} from '@/hooks/web/useForm'
import {useValidator} from '@/hooks/web/useValidator'
import {FormSchema} from '@/types/form'
import {ApiArea} from "@/api/stub";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const {t} = useI18n()

const props = defineProps({
  currentRow: {
    type: Object as PropType<Nullable<ApiArea>>,
    default: () => null
  }
})

const rules = {
  name: [required()],
  description: []
}

const schema = reactive<FormSchema[]>([
  {
    field: 'name',
    label: t('areas.name'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('areas.name')
    }
  },
  {
    field: 'description',
    label: t('areas.description'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('areas.description')
    }
  },
])

watch(
    () => props.currentRow,
    (currentRow) => {
      if (!currentRow) return
      const { setValues, setSchema } = methods
      setValues(currentRow)
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
