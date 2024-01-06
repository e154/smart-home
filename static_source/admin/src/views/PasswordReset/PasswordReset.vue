<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {underlineToHump} from '@/utils'
import {useAppStore} from '@/store/modules/app'
import {useDesign} from '@/hooks/web/useDesign'
import {reactive, ref, unref} from 'vue'
import {FormSchema} from "@/types/form";
import {ElButton, ElCheckbox, ElLink, ElInput} from 'element-plus'
import api from "@/api/api";
import {useRouter} from "vue-router";
import {useForm} from "@/hooks/web/useForm";
import {useValidator} from "@/hooks/web/useValidator";
import {Form} from '@/components/Form'

const {getPrefixCls} = useDesign()
const prefixCls = getPrefixCls('login')
const appStore = useAppStore()
const {t} = useI18n()
const {required} = useValidator()
const {register, elFormRef, methods} = useForm()
const {currentRoute, addRoute, push} = useRouter()

const newRequestSchema = reactive<FormSchema[]>([
  {
    field: 'title',
    colProps: {
      span: 24
    }
  },
  {
    field: 'username',
    label: t('login.username'),
    value: '',
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('login.usernamePlaceholder')
    }
  },
  {
    field: 'restore',
    colProps: {
      span: 24
    }
  }
])
const newPasswordSchema = reactive<FormSchema[]>([
  {
    field: 'title',
    colProps: {
      span: 24
    }
  },
  {
    field: 'password',
    label: t('login.password'),
    value: 'admin',
    component: 'InputPassword',
    colProps: {
      span: 24
    },
    componentProps: {
      style: {
        width: '100%'
      },
      placeholder: t('login.passwordPlaceholder')
    }
  },
  {
    field: 'check_password',
    label: t('login.checkPassword'),
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
      placeholder: t('login.passwordPlaceholder')
    }
  },
  {
    field: 'update',
    colProps: {
      span: 24
    }
  }
])
const newRequestRules = {
  username: [required()],
}
const newPasswordRules = {
  password: [required()],
  check_password: [required()],
}

// NEW_REQUEST
// REQUEST_SENDED
// UPDATE_PASSWORD
const state = ref('NEW_REQUEST')
const loading = ref(false)

if (currentRoute.value.query?.t) {
  state.value = 'UPDATE_PASSWORD'
}

const toLogin = () => {
  push('/login')
}

const requestNewPassword = async () => {
  const formRef = unref(elFormRef)
  formRef?.validate(async (valid) => {
    if (valid) {
      try {
        loading.value = true
        const {getFormData} = methods
        const formData = await getFormData()
        const {username} = formData
        const data = await api.v1.authServicePasswordReset({email: username});
        state.value = 'REQUEST_SENDED'
        loading.value = false

      } finally {
        loading.value = false
      }
    }
  })
}
const updatePassword = async () => {
  const formRef = unref(elFormRef)
  formRef?.validate(async (valid) => {
    if (valid) {
      try {
        loading.value = true

        const {getFormData} = methods
        const formData = await getFormData()
        const {password} = formData
        const token = currentRoute.value.query.t as string || '';
        const data = await api.v1.authServicePasswordReset({newPassword: password, token: token});

        toLogin()
      } finally {
        loading.value = false
      }
    }
  })
}

</script>

<template>
  <div
      :class="prefixCls"
      class="h-[100%] relative <xl:bg-v-dark <sm:px-10px <xl:px-10px <md:px-10px"
  >
    <div class="relative h-full flex mx-auto">
      <div
          :class="`${prefixCls}__left flex-1 bg-gray-500 bg-opacity-20 relative p-30px <xl:hidden`"
      >
<!--        <div class="flex items-center relative text-white">-->
<!--          <img src="../../assets/svgs/logo-w.svg" alt="" class="w-48px h-48px mr-10px"/>-->
<!--          <span class="text-20px font-bold">{{ underlineToHump(appStore.getTitle) }}</span>-->
<!--        </div>-->
        <div class="flex justify-center items-center h-[calc(100%-60px)]">
          <TransitionGroup
              appear
              tag="div"
              enter-active-class="animate__animated animate__bounceInLeft"
          >
            <img src="@/assets/svgs/banner-v4.svg" key="1" alt="" class="w-400px" />
<!--            <div class="text-3xl text-white" key="2">{{ t('login.welcome') }}</div>-->
<!--            <div class="mt-5 font-normal text-white text-14px" key="3">-->
<!--              {{ t('login.message') }}-->
<!--            </div>-->
          </TransitionGroup>
        </div>
      </div>
      <div class="flex-1 p-30px <sm:p-10px dark:bg-v-dark relative">
        <div class="flex justify-between items-center text-white @2xl:justify-end @xl:justify-end">
          <div class="flex items-center @2xl:hidden @xl:hidden">
            <img src="../../assets/svgs/logo-w.svg" alt="" class="w-48px h-48px mr-10px"/>
            <span class="text-20px font-bold">{{ underlineToHump(appStore.getTitle) }}</span>
          </div>

          <div class="flex justify-end items-center space-x-10px">
            <ThemeSwitch/>
            <LocaleDropdown class="<xl:text-white dark:text-white"/>
          </div>
        </div>

        <Transition appear enter-active-class="animate__animated animate__bounceInRight">
        <div
            class="h-full flex items-center m-auto w-[100%] @2xl:max-w-500px @xl:max-w-500px @md:max-w-500px @lg:max-w-500px">

          <div v-if="state==='REQUEST_SENDED'" class="dark:(border-1 border-[var(--el-border-color)] border-solid)" style="padding: 20px">
            <p style="padding-bottom: 20px">
              Check your email for a link to reset your password. If it doesn’t appear within a few minutes, check your
              spam
              folder.
            </p>
            <el-button
                type="primary"
                class="w-[100%]"
                @click="toLogin"
            >
              {{ $t('login.returnToSignIn') }}
            </el-button>
          </div>
          <div v-else-if="state==='NEW_REQUEST'" class="dark:(border-1 border-[var(--el-border-color)] border-solid)">
            <Form
                :schema="newRequestSchema"
                :rules="newRequestRules"
                label-position="top"
                hide-required-asterisk
                size="large"
                @register="register"
                style="margin: 20px"
            >

              <template #title>
                <h2 class="text-2xl font-bold text-center w-[100%]">{{ t('login.restorePassword') }}</h2>
              </template>

              <template #code="form">
                <div class="w-[100%] flex">
                  <ElInput v-model="form['code']" :placeholder="t('login.codePlaceholder')"/>
                </div>
              </template>

              <template #restore>
                <div class="w-[100%]">
                  <ElButton type="primary" class="w-[100%]" :loading="loading" @click="requestNewPassword">
                    {{ t('login.restorePassword') }}
                  </ElButton>
                </div>
                <div class="w-[100%] mt-15px">
                  <ElButton class="w-[100%]" @click="toLogin">
                    {{ t('login.hasUser') }}
                  </ElButton>
                </div>
              </template>
            </Form>

          </div>
          <div v-else-if="state==='UPDATE_PASSWORD'"
               class="dark:(border-1 border-[var(--el-border-color)] border-solid)">

            <Form
                :schema="newPasswordSchema"
                :rules="newPasswordRules"
                label-position="top"
                hide-required-asterisk
                size="large"
                @register="register"
                style="margin: 20px"
            >

              <template #title>
                <h2 class="text-2xl font-bold text-center w-[100%]">{{ t('login.updatePassword') }}</h2>
              </template>

              <template #code="form">
                <div class="w-[100%] flex">
                  <ElInput v-model="form['code']" :placeholder="t('login.codePlaceholder')"/>
                </div>
              </template>

              <template #update>
                <div class="w-[100%]">
                  <ElButton type="primary" class="w-[100%]" :loading="loading" @click="updatePassword">
                    {{ t('main.update') }}
                  </ElButton>
                </div>
              </template>
            </Form>

          </div>

        </div>
        </Transition>

      </div>
    </div>
  </div>
</template>

<style lang="less" scoped>
@prefix-cls: ~'@{namespace}-login';

.@{prefix-cls} {
  &__left {
    &::before {
      position: absolute;
      top: 0;
      left: 0;
      z-index: -1;
      width: 100%;
      height: 100%;
      background-image: url('@/assets/svgs/login-bg.svg');
      background-position: center;
      background-repeat: no-repeat;
      content: '';
    }
  }
}
</style>
