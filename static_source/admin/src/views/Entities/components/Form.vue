<script setup lang="ts">
import { PropType, reactive, watch} from 'vue'
import {Form} from '@/components/Form'
import {useI18n} from '@/hooks/web/useI18n'
import {useForm} from '@/hooks/web/useForm'
import {useValidator} from '@/hooks/web/useValidator'
import {FormSchema} from '@/types/form'
import {Entity} from "@/views/Entities/components/types";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const {t} = useI18n()

const props = defineProps({
  currentRow: {
    type: Object as PropType<Nullable<Entity>>,
    default: () => null
  }
})

const rules = {
  name: [required()],
  plugin: [required()],
  description: []
}

const schema = reactive<FormSchema[]>([
  {
    field: 'name',
    label: t('entities.name'),
    component: 'Input',
    colProps: {
      span: 12
    },
    componentProps: {
      placeholder: t('entities.name')
    }
  },
  {
    field: 'id',
    label: t('entities.id'),
    component: 'Input',
    colProps: {
      span: 12
    },
    hidden: true,
    componentProps: {
      placeholder: t('entities.id'),
      disabled: true
    }
  },
  {
    field: 'plugin',
    label: t('entities.pluginName'),
    component: 'Plugin',
    value: null,
    colProps: {
      span: 12
    },
    hidden: false,
    componentProps: {
      placeholder: t('entities.pluginName')
    }
  },
  {
    field: 'scriptIds',
    label: t('entities.scripts'),
    component: 'Scripts',
    colProps: {
      span: 24
    },
    value: [],
    hidden: false,
    componentProps: {
      placeholder: t('entities.scripts')
    }
  },
  {
    field: 'description',
    label: t('entities.description'),
    component: 'Input',
    colProps: {
      span: 12
    },
    componentProps: {
      placeholder: t('entities.description')
    }
  },
  {
    field: 'icon',
    label: t('entities.icon'),
    component: 'Input',
    colProps: {
      span: 12
    },
    componentProps: {
      placeholder: t('entities.icon')
    }
  },
  {
    field: 'image',
    label: t('entities.image'),
    component: 'Image',
    value: null,
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('entities.image'),
    }
  },
  {
    field: 'autoLoad',
    label: t('entities.autoLoad'),
    component: 'Switch',
    value: false,
    colProps: {
      span: 12
    },
  },
  {
    field: 'restoreState',
    label: t('entities.restoreState'),
    component: 'Switch',
    value: true,
    colProps: {
      span: 12
    },
  },
  {
    field: 'parent',
    label: t('entities.parent'),
    value: null,
    component: 'Entity',
    colProps: {
      span: 12
    },
    componentProps: {
      placeholder: t('entities.parent'),
    }
  },
  {
    field: 'area',
    label: t('entities.area'),
    value: null,
    component: 'Area',
    colProps: {
      span: 12
    },
    componentProps: {
      placeholder: t('entities.area'),
    }
  },
])

const {setValues, setSchema} = methods
watch(
    () => props.currentRow,
    (currentRow) => {
      if (!currentRow) return

      let exist: boolean
      if (currentRow.id) {
        exist = true
      }
      const schema = [
        {field: 'name', path: 'hidden', value: exist},
        {field: 'id', path: 'hidden', value: !exist},
        {field: 'plugin', path: 'componentProps.disabled', value: exist},
      ]
      setSchema(schema)
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
