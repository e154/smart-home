<script setup lang="ts">
import {h, PropType, reactive, unref, watch} from 'vue'
import {Form} from '@/components/Form'
import {useI18n} from '@/hooks/web/useI18n'
import {useForm} from '@/hooks/web/useForm'
import {useValidator} from '@/hooks/web/useValidator'
import {FormSchema} from '@/types/form'
import {Trigger} from "@/views/Automation/components/types";
import {ApiAttribute} from "@/api/stub";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const {t} = useI18n()

const props = defineProps({
  trigger: {
    type: Object as PropType<Nullable<Trigger>>,
    default: () => null
  }
})

const rules = {
  name: [required()],
  pluginName: [required()],
}

const {setValues, setSchema} = methods
const pluginChanged = async (value?: string) => {
  if (!value) {
    value = 'state_change'
  }
  const schema = [
    {field: 'entity', path: 'hidden', value: value !== 'state_change'},
    {field: 'timePluginOptions', path: 'hidden', value: value !== 'time'},
    {field: 'systemPluginOptions', path: 'hidden', value: value !== 'system'},
    {field: 'alexaPluginOptions', path: 'hidden', value: value !== 'alexa'}
  ]
  setSchema(schema)
}

const schema = reactive<FormSchema[]>([
  {
    field: 'name',
    label: t('automation.triggers.name'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('automation.triggers.name')
    }
  },
  {
    field: 'description',
    label: t('automation.description'),
    component: 'Input',
    colProps: {
      span: 24
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
      span: 24
    },
    componentProps: {
      placeholder: t('automation.area'),
    }
  },
  {
    field: 'enabled',
    label: t('automation.enabled'),
    component: 'Switch',
    value: false,
    colProps: {
      span: 24
    },
  },
  {
    field: 'script',
    label: t('automation.triggers.script'),
    component: 'Script',
    colProps: {
      span: 24
    },
    value: null,
    componentProps: {
      placeholder: t('automation.triggers.script')
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
  {
    field: 'pluginName',
    label: t('automation.triggers.pluginName'),
    component: 'Select',
    colProps: {
      span: 24,
    },
    value: 'state_change',
    componentProps: {
      options: [
        {
          label: 'STATE_CHANGE',
          value: 'state_change'
        },
        {
          label: 'TIME',
          value: 'time'
        },
        {
          label: 'SYSTEM',
          value: 'system'
        },
        {
          label: 'ALEXA',
          value: 'alexa'
        }
      ],
      clearable: false,
      onChange: pluginChanged
    }
  },
  {
    field: 'entity',
    label: t('automation.triggers.entity'),
    component: 'Entity',
    colProps: {
      span: 24
    },
    hidden: false,
    componentProps: {
      placeholder: t('automation.triggers.entity')
    }
  },
  {
    hidden: false,
    field: 'timePluginOptions',
    label: t('automation.triggers.timePluginOptions'),
    component: 'Input',
    colProps: {
      span: 24
    }, componentProps: {
      placeholder: t('automation.triggers.timePluginOptions')
    }
  },
  {
    hidden: false,
    field: 'systemPluginOptions',
    label: t('automation.triggers.pluginOptions'),
    component: 'Input',
    colProps: {
      span: 24
    }, componentProps: {
      placeholder: t('automation.triggers.pluginOptions')
    }
  },
  {
    hidden: false,
    field: 'alexaPluginOptions',
    label: t('automation.triggers.alexaSkillId'),
    component: 'InputNumber',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('automation.triggers.alexaSkillId')
    }
  },
])

const prepareAttrs = (attrs?: Record<string, ApiAttribute>) => {
  if (!attrs) {
    return
  }
  if (attrs.hasOwnProperty('cron')) {
    setValues({timePluginOptions: attrs.cron.string})
  }
  if (attrs.hasOwnProperty('skillId')) {
    setValues({alexaPluginOptions: attrs.skillId.int})
  }
}

watch(
    () => props.trigger,
    (val) => {
      if (!val) return
      setValues(val)
      pluginChanged(val.pluginName)
      prepareAttrs(val.attributes)
    },
    {
      deep: true,
      immediate: true
    }
)

pluginChanged(props.trigger?.pluginName)
prepareAttrs(props.trigger?.attributes)

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
