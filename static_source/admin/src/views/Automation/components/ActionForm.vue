<script setup lang="ts">
import {PropType, reactive, ref, unref, watch} from 'vue'
import {Form, FormExpose} from '@/components/Form'
import {useI18n} from '@/hooks/web/useI18n'
import {useForm} from '@/hooks/web/useForm'
import {useValidator} from '@/hooks/web/useValidator'
import {FormSchema} from '@/types/form'
import {Trigger} from "@/views/Automation/components/types";
import {ApiAction, ApiEntityShort} from "@/api/stub";
import api from "@/api/api";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const {t} = useI18n()

const props = defineProps({
  action: {
    type: Object as PropType<Nullable<ApiAction>>,
    default: () => null
  }
})

const rules = {
  name: [required()],
  pluginName: [required()],
}

const getEntity = async (entityId?: string) => {
  if (!entityId) {
    return null
  }
  const res = await api.v1.entityServiceGetEntity(entityId)
  return res.data
}

const entityActionOptions = (options) => {
  unref(formRef)?.setSchema([
    {
      field: 'entityActionName',
      path: 'componentProps.options',
      value: options,
    }
  ])
}

const formRef = ref<FormExpose>()
const updateEntityActions = async (entity?: ApiEntityShort) => {
  if (!entity) {
    entityActionOptions([])
    return
  }
  const res = await getEntity(entity.id)
  let options = []
  if (res) {
    for (const action of res?.actions) {
      options.push({
        key: action.name,
        value: action.name,
      })
    }
  }
  entityActionOptions(options)
}

updateEntityActions(props.action?.entity)

const schema = reactive<FormSchema[]>([
  {
    field: 'name',
    label: t('automation.actions.name'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('automation.actions.name')
    }
  },
  {
    field: 'script',
    label: t('automation.actions.script'),
    component: 'Script',
    colProps: {
      span: 24
    },
    value: null,
    componentProps: {
      placeholder: t('automation.actions.script')
    }
  },
  {
    field: 'entity',
    label: t('automation.actions.entity'),
    component: 'Entity',
    colProps: {
      span: 24
    },
    hidden: false,
    componentProps: {
      placeholder: t('automation.actions.entity'),
      onChange: (entity?: ApiEntityShort) => {
        updateEntityActions(entity)
      }
    }
  },
  {
    field: 'entityActionName',
    label: t('automation.actions.entityActionName'),
    component: 'Select',
    colProps: {
      span: 24,
    },
    componentProps: {
      options: []
    }
  },
])

const {setValues, setSchema} = methods

watch(
    () => props.action,
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
      ref="formRef"
      :schema="schema"
      :rules="rules"
      label-position="top"
      @register="register"

  />
</template>
