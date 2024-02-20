<script setup lang="ts">
import {defineEmits, reactive, ref, unref, watch} from 'vue'
import {Form} from '@/components/Form'
import {useI18n} from '@/hooks/web/useI18n'
import {ElButton, ElCheckbox, ElLink} from 'element-plus'
import {useForm} from '@/hooks/web/useForm'
import {useCache} from '@/hooks/web/useCache'
import {useAppStore} from '@/store/modules/app'
import {usePermissionStore} from '@/store/modules/permission'
import type {RouteLocationNormalizedLoaded, RouteRecordRaw} from 'vue-router'
import {useRouter} from 'vue-router'
import {UserType} from '@/api/login/types'
import {useValidator} from '@/hooks/web/useValidator'
import {FormSchema} from '@/types/form'
import api from "@/api/api";
import stream from "@/api/stream";
import {ApiSigninResponse} from "@/api/stub";
import {prepareUrl} from "@/utils/serverId";
import pushService from "@/api/pushService";

const {required} = useValidator()
const emit = defineEmits(['to-restore'])
const appStore = useAppStore()
const permissionStore = usePermissionStore()
const {currentRoute, addRoute, push} = useRouter()
const {wsCache} = useCache()
const {t} = useI18n()

const rules = {
  username: [required()],
  password: [required()],
  server: []
}

const schema = reactive<FormSchema[]>([
  {
    field: 'title',
    colProps: {
      span: 24
    }
  },
  {
    field: 'username',
    label: t('login.username'),
    value: wsCache.get("username"),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('login.usernamePlaceholder')
    }
  },
  {
    field: 'password',
    label: t('login.password'),
    value: wsCache.get("password"),
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
    field: 'server',
    label: t('login.server'),
    value: wsCache.get("serverId"),
    component: 'Input',
    colProps: {
      span: 24
    },
    hidden: !appStore.getIsGate,
    componentProps: {
      style: {
        width: '100%'
      },
      placeholder: t('login.serverPlaceholder')
    }
  },
  {
    field: 'tool',
    colProps: {
      span: 24
    }
  },
  {
    field: 'login',
    colProps: {
      span: 24
    }
  }
])
const iconSize = 30
const remember = ref(wsCache.get("remember"))
const {register, elFormRef, methods} = useForm()
const loading = ref(false)
const iconColor = '#999'
const redirect = ref<string>('')
const ipAddress = ref();

watch(
    () => currentRoute.value,
    (route: RouteLocationNormalizedLoaded) => {
      redirect.value = route?.query?.redirect as string
    },
    {
      immediate: true
    }
)

const getIpAddress = async () => {
  fetch('https://api.ipify.org?format=json')
      .then(x => x.json())
      .then(({ ip }) => {
        ipAddress.value = ip;
      });
}

getIpAddress()

// 登录
const signIn = async () => {
  const formRef = unref(elFormRef)
  await formRef?.validate(async (isValid) => {
    if (isValid) {
      loading.value = true

      const {getFormData} = methods
      const formData = await getFormData<UserType>()
      let {username, password, server} = formData;
      username = username.trim();
      if (server) {
        server = server.trim();
      }

      appStore.setServerId(server);

      try {
        const resp = await api.v1.authServiceSignin({
          headers: {Authorization: 'Basic ' + btoa(username + ':' + password), ip: ipAddress.value}
        });

        const {accessToken, currentUser} = resp.data as ApiSigninResponse;
        if (accessToken) {
          wsCache.set("accessToken", accessToken)
          wsCache.set("currentUser", currentUser)

          if (remember.value) {
            wsCache.set("username", username)
            wsCache.set("password", password)
            wsCache.set("remember", true)
          } else {
            wsCache.delete("username")
            wsCache.delete("password")
            wsCache.delete("remember")
          }

          appStore.SetToken(accessToken);
          appStore.SetUser(currentUser);

          if (currentUser?.image) {
            appStore.SetAvatar(prepareUrl(import.meta.env.VITE_API_BASEPATH as string + currentUser.image.url));
          } else {
            appStore.SetAvatar('');
          }

          // ws
          stream.connect(import.meta.env.VITE_API_BASEPATH as string || window.location.origin, accessToken);
          // push
          pushService.start()

          await permissionStore.generateRoutes('none').catch(() => {
          })
          permissionStore.getAddRouters.forEach((route) => {
            addRoute(route as RouteRecordRaw) // 动态添加可访问路由表
          })
          permissionStore.setIsAddRouters(true)
          push({path: redirect.value || permissionStore.addRouters[0].path})
        }
      } finally {
        loading.value = false
      }
    }
  })
}

// 获取角色信息
// const getRole = async () => {
//   const { getFormData } = methods
//   const formData = await getFormData<UserType>()
//   const params = {
//     roleName: formData.username
//   }
//   // admin - 模拟后端过滤菜单
//   // test - 模拟前端过滤菜单
//   const res =
//     formData.username === 'admin' ? await getAdminRoleApi(params) : await getTestRoleApi(params)
//   if (res) {
//     const { wsCache } = useCache()
//     const routers = res.data || []
//     wsCache.set('roleRouters', routers)
//
//     formData.username === 'admin'
//       ? await permissionStore.generateRoutes('admin', routers).catch(() => {})
//       : await permissionStore.generateRoutes('test', routers).catch(() => {})
//
//     permissionStore.getAddRouters.forEach((route) => {
//       addRoute(route as RouteRecordRaw) // 动态添加可访问路由表
//     })
//     permissionStore.setIsAddRouters(true)
//     push({ path: redirect.value || permissionStore.addRouters[0].path })
//   }
// }

// 去注册页面
const toRestore = () => {
  push('/password_reset')
}
</script>

<template>
  <Form
      :schema="schema"
      :rules="rules"
      label-position="top"
      hide-required-asterisk
      size="large"
      class="dark:(border-1 border-[var(--el-border-color)] border-solid)"
      @register="register"
  >
    <template #title>
      <h2 class="text-2xl font-bold text-center w-[100%]">{{ t('login.login') }}</h2>
    </template>

    <template #tool>
      <div class="flex justify-between items-center w-[100%]">
        <ElCheckbox v-model="remember" :label="t('login.remember')"/>
        <ElLink  type="primary" :underline="false" @click="toRestore()">{{ t('login.forgetPassword') }}</ElLink>
      </div>
    </template>

    <template #login>
      <div class="w-[100%]">
        <ElButton :loading="loading" type="primary" class="w-[100%]" @click="signIn">
          {{ t('login.login') }}
        </ElButton>
      </div>
    </template>
  </Form>
</template>

<style lang="less" scoped>
:deep(.anticon) {
  &:hover {
    color: var(--el-color-primary) !important;
  }
}
</style>
