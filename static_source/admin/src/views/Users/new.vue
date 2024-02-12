<script setup lang="ts">
import {ref, unref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {ElButton} from 'element-plus'
import {useForm} from '@/hooks/web/useForm'
import {useCache} from '@/hooks/web/useCache'
import {useAppStore} from '@/store/modules/app'
import {usePermissionStore} from '@/store/modules/permission'
import {useRoute, useRouter} from 'vue-router'
import {useValidator} from '@/hooks/web/useValidator'
import api from "@/api/api";
import Form from './components/Form.vue'
import {ApiNewtUserRequest} from "@/api/stub";
import {ContentWrap} from "@/components/ContentWrap";
import {User} from "@/views/Users/components/Types";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const appStore = useAppStore()
const permissionStore = usePermissionStore()
const {currentRoute, addRoute, push} = useRouter()
const route = useRoute();
const {wsCache} = useCache()
const {t} = useI18n()

const writeRef = ref<ComponentRef<typeof Form>>()
const loading = ref(false)
const currentRow = ref<Nullable<ApiNewtUserRequest>>(null)

const save = async () => {
  const write = unref(writeRef)
  const validate = await write?.elFormRef?.validate()?.catch(() => {
  })
  if (validate) {
    loading.value = true
    const data = (await write?.getFormData()) as User
    const body = {
      nickname: data.nickname,
      firstName: data.firstName,
      lastName: data.lastName,
      password: data.password,
      passwordRepeat: data.passwordRepeat,
      email: data.email,
      status: data.status,
      lang: data.lang,
      imageId: data.image?.id,
      roleName: data.role?.name || null,
    } as ApiNewtUserRequest
    const res = await api.v1.userServiceAddUser(body)
        .catch(() => {
        })
        .finally(() => {
          loading.value = false
        })
    if (res) {
      cancel()
    }
  }
}

const cancel = () => {
  push('/etc/users')
}

</script>

<template>
  <ContentWrap>
    <Form ref="writeRef" :current-row="currentRow"/>

    <div style="text-align: right">

      <ElButton type="primary" @click="save()">
        {{ t('main.save') }}
      </ElButton>

      <ElButton type="default" @click="cancel()">
        {{ t('main.cancel') }}
      </ElButton>

    </div>
  </ContentWrap>

</template>

<style lang="less" scoped>

</style>
