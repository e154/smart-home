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
  nickname: [required()],
  email: [required()],
  status: [required()],
  role: [required()]
}

const schema = reactive<FormSchema[]>([
  {
    field: 'nickname',
    label: t('users.nickname'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('users.nickname')
    }
  },
  {
    field: 'firstName',
    label: t('users.firstName'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('users.firstName')
    }
  },
  {
    field: 'lastName',
    label: t('users.lastName'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('users.lastName')
    }
  },
  {
    field: 'status',
    label: t('users.status'),
    component: 'Select',
    value: 'active',
    componentProps: {
      options: [
        {
          label: t('main.ACTIVE'),
          value: 'active'
        },
        {
          label: t('main.BLOCKED'),
          value: 'blocked'
        }
      ],
      onChange: (value: string) => {

      }
    }
  },
  {
    field: 'email',
    label: t('users.email'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('users.email')
    }
  },
  {
    field: 'image',
    label: t('users.image'),
    component: 'Image',
    value: null,
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('users.image'),
    }
  },
  {
    field: 'role',
    label: t('users.role'),
    component: 'Role',
    colProps: {
      span: 24
    },
    value: null,
    componentProps: {
      placeholder: t('users.role'),
    }
  },
  {
    field: 'password',
    label: t('users.password'),
    value: '',
    component: 'InputPassword',
    colProps: {
      span: 24
    },
    componentProps: {
      style: {
        width: '100%'
      },
      strength: true,
      placeholder: t('users.passwordPlaceholder')
    }
  },
  {
    field: 'passwordRepeat',
    label: t('users.passwordRepeat'),
    value: '',
    component: 'InputPassword',
    colProps: {
      span: 24
    },
    componentProps: {
      style: {
        width: '100%'
      },
      strength: true,
      placeholder: t('users.passwordPlaceholder')
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
