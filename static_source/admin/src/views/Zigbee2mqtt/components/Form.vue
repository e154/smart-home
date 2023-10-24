<script setup lang="ts">
import { PropType, reactive, watch} from 'vue'
import {Form} from '@/components/Form'
import {useI18n} from '@/hooks/web/useI18n'
import {useForm} from '@/hooks/web/useForm'
import {useValidator} from '@/hooks/web/useValidator'
import {FormSchema} from '@/types/form'
import {ApiZigbee2Mqtt} from "@/api/stub";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const {t} = useI18n()

const props = defineProps({
  currentRow: {
    type: Object as PropType<Nullable<ApiZigbee2Mqtt>>,
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
    label: t('zigbee2mqtt.name'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('zigbee2mqtt.name')
    }
  },
  {
    field: 'baseTopic',
    label: t('zigbee2mqtt.baseTopic'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('zigbee2mqtt.baseTopic')
    }
  },
  {
    field: 'login',
    label: t('zigbee2mqtt.login'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('zigbee2mqtt.login')
    }
  },
  {
    field: 'password',
    label: t('zigbee2mqtt.password'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('zigbee2mqtt.password')
    }
  },
  {
    field: 'permitJoin',
    label: t('zigbee2mqtt.permitJoin'),
    component: 'Switch',
    value: false,
    colProps: {
      span: 24
    },
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
