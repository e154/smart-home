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
import {ApiVariable} from "@/api/stub";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const emit = defineEmits(['to-restore'])
const appStore = useAppStore()
const permissionStore = usePermissionStore()
const {currentRoute, addRoute, push} = useRouter()
const route = useRoute();
const {wsCache} = useCache()
const {t} = useI18n()

const props = defineProps({
  currentRow: {
    type: Object as PropType<Nullable<ApiVariable>>,
    default: () => null
  }
})

const rules = {
  name: [required()],
  value: [required()]
}

const variableName = computed(() => route.params.name);

const schema = reactive<FormSchema[]>([
  {
    field: 'name',
    label: t('variables.name'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('variables.name')
    }
  },
  {
    field: 'value',
    label: t('variables.value'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('variables.value')
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
