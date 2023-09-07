<script setup lang="ts">
import {computed, defineEmits, PropType, reactive, ref, watch} from 'vue'
import {Form} from '@/components/Form'
import {useI18n} from '@/hooks/web/useI18n'
import {ElButton, ElPopconfirm} from 'element-plus'
import {useForm} from '@/hooks/web/useForm'
import {useCache} from '@/hooks/web/useCache'
import {useAppStore} from '@/store/modules/app'
import {usePermissionStore} from '@/store/modules/permission'
import {useRoute, useRouter} from 'vue-router'
import {useValidator} from '@/hooks/web/useValidator'
import {FormSchema} from '@/types/form'
import {ApiTask} from "@/api/stub";
import api from "@/api/api";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const appStore = useAppStore()
const permissionStore = usePermissionStore()
const {currentRoute, addRoute, push} = useRouter()
const route = useRoute();
const {wsCache} = useCache()
const {t} = useI18n()

const props = defineProps({
  task: {
    type: Object as PropType<Nullable<ApiTask>>,
    default: () => null
  }
})

const rules = {
  name: [required()],
  condition: [required()],
}

const schema = reactive<FormSchema[]>([
  {
    field: 'name',
    label: t('automation.name'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('automation.name')
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
    field: 'enabled',
    label: t('automation.enabled'),
    component: 'Switch',
    value: false,
    colProps: {
      span: 24
    },
  },
  {
    field: 'condition',
    label: t('automation.condition'),
    component: 'Select',
    value: 'and',
    componentProps: {
      options: [
        {
          label: 'AND',
          value: 'and'
        },
        {
          label: 'OR',
          value: 'or'
        }
      ],
      onChange: (value: string) => {

      }
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
])

watch(
    () => props.task,
    (task) => {
      if (!task) return
      const { setValues, setSchema } = methods
      setValues(task)
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
