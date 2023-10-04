<script setup lang="ts">
import {h, PropType, reactive, unref, watch} from 'vue'
import {Form} from '@/components/Form'
import {useI18n} from '@/hooks/web/useI18n'
import {useForm} from '@/hooks/web/useForm'
import {useValidator} from '@/hooks/web/useValidator'
import {FormSchema} from '@/types/form'
import {EntityState} from "@/views/Entities/components/types";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const {t} = useI18n()

const props = defineProps({
  state: {
    type: Object as PropType<Nullable<EntityState>>,
    default: () => null
  }
})

const rules = {
  name: [required()],
}

const {setValues, setSchema} = methods

const schema = reactive<FormSchema[]>([
  {
    field: 'name',
    label: t('entities.name'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('entities.name')
    }
  },
  {
    field: 'description',
    label: t('entities.description'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('entities.description')
    }
  },
  {
    field: 'image',
    label: t('entities.image'),
    component: 'Image',
    colProps: {
      span: 24
    },
    value: null,
    componentProps: {
      placeholder: t('entities.image')
    }
  },
])

watch(
    () => props.state,
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
