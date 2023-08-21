<script setup lang="ts">
import {computed, defineEmits, onMounted, PropType, reactive, ref, watch} from 'vue'
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
import {User} from "@/views/Users/components/types";
import api from "@/api/api";
import {ApiImage, ApiPluginOptionsResultEntityState, ApiRole} from "@/api/stub";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const appStore = useAppStore()
const permissionStore = usePermissionStore()
const {currentRoute, addRoute, push} = useRouter()
const route = useRoute();
const {wsCache} = useCache()
const {t} = useI18n()

const props = defineProps({
  currentRow: {
    type: Object as PropType<Nullable<User>>,
    default: () => null
  }
})

const rules = {
  name: [required()],
  lang: [required()],
}

const schema = reactive<FormSchema[]>([
  {
    field: 'name',
    label: t('scripts.name'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('scripts.name')
    }
  },
  {
    field: 'lang',
    label: t('scripts.lang'),
    component: 'Select',
    value: 'coffeescript',
    componentProps: {
      options: [
        {
          label: 'coffeescript',
          value: 'coffeescript'
        },
        {
          label: 'javascript',
          value: 'javascript'
        },
        {
          label: 'typescript',
          value: 'ts'
        }
      ],
      onChange: (value: string) => {

      }
    }
  },
  {
    field: 'description',
    label: t('scripts.description'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('scripts.description')
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
