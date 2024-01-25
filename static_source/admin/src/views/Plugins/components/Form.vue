<script setup lang="ts">
import {PropType, reactive, ref, watch} from 'vue'
import {Form} from '@/components/Form'
import {useI18n} from '@/hooks/web/useI18n'
import {useForm} from '@/hooks/web/useForm'
import {useValidator} from '@/hooks/web/useValidator'
import {FormSchema} from '@/types/form'
import {Plugin} from './Types.ts'

const {register, elFormRef, methods} = useForm()
const {t} = useI18n()

const form = ref<Nullable<Plugin>>(null)
const props = defineProps({
  currentRow: {
    type: Object as PropType<Nullable<Plugin>>,
    default: () => null
  }
})

const schema = reactive<FormSchema[]>([
  {
    field: 'name',
    label: t('plugins.name'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      disabled: true,
      placeholder: t('plugins.name')
    }
  },
  {
    field: 'version',
    label: t('plugins.version'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      disabled: true,
      placeholder: t('plugins.version')
    }
  },
  {
    field: 'enabled',
    label: t('plugins.enabled'),
    component: 'Switch',
    colProps: {
      span: 24
    },
  },
  {
    field: 'system',
    label: t('plugins.system'),
    component: 'Switch',
    componentProps: {
      disabled: true,
    },
    colProps: {
      span: 24
    },
  },
  {
    field: 'triggers',
    label: t('plugins.trigger'),
    component: 'Switch',
    componentProps: {
      disabled: true,
    },
    colProps: {
      span: 24
    },
  },
  {
    field: 'actors',
    label: t('actors'),
    component: 'Switch',
    componentProps: {
      disabled: true,
    },
    colProps: {
      span: 24
    },
  },
  {
    field: 'actorCustomAttrs',
    label: t('plugins.actorCustomAttrs'),
    component: 'Switch',
    componentProps: {
      disabled: true,
    },
    colProps: {
      span: 24
    },
  },
  {
    field: 'actorCustomActions',
    label: t('plugins.actorCustomActions'),
    component: 'Switch',
    componentProps: {
      disabled: true,
    },
    colProps: {
      span: 24
    },
  },
  {
    field: 'actorCustomStates',
    label: t('plugins.actorCustomStates'),
    component: 'Switch',
    componentProps: {
      disabled: true,
    },
    colProps: {
      span: 24
    },
  },
  {
    field: 'actorCustomSetts',
    label: t('plugins.actorCustomSettings'),
    component: 'Switch',
    componentProps: {
      disabled: true,
    },
    colProps: {
      span: 24
    },
  },
])

watch(
    () => props.currentRow,
    (currentRow) => {
      if (!currentRow) return
      const {setValues, setSchema} = methods
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
      label-position="top"
      @register="register"
  />
</template>
