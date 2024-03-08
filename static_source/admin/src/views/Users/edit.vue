<script setup lang="ts">
import {computed, defineEmits, ref, unref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {ElButton, ElMessage, ElPopconfirm} from 'element-plus'
import {useForm} from '@/hooks/web/useForm'
import {useCache} from '@/hooks/web/useCache'
import {useAppStore} from '@/store/modules/app'
import {usePermissionStore} from '@/store/modules/permission'
import {useRoute, useRouter} from 'vue-router'
import {useValidator} from '@/hooks/web/useValidator'
import api from "@/api/api";
import Form from './components/Form.vue'
import {ApiUserFull, ApiUserMeta} from "@/api/stub";
import {ContentWrap} from "@/components/ContentWrap";
import {User} from "@/views/Users/components/Types";
import {prepareUrl} from "@/utils/serverId";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const emit = defineEmits(['to-restore'])
const appStore = useAppStore()
const permissionStore = usePermissionStore()
const {currentRoute, addRoute, push} = useRouter()
const route = useRoute();
const {wsCache} = useCache()
const {t} = useI18n()

const writeRef = ref<ComponentRef<typeof Form>>()
const loading = ref(false)
const userId = computed(() => route.params.id as number);
const currentUser = ref<Nullable<User>>(null)

const fetch = async () => {
  loading.value = true
  const res = await api.v1.userServiceGetUserById(userId.value)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })
  if (res) {
    const user = res.data as ApiUserFull
    currentUser.value = {
      nickname: user.nickname,
      firstName: user.firstName,
      lastName: user.lastName,
      email: user.email,
      status: user.status,
      lang: user.lang,
      image: user.image,
      roleName: user.roleName,
      role: user.role,
      meta: user.meta,
    } as User
  } else {
    currentUser.value = null
  }
}

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
    }

    if (data?.image && data.image.url) {
      appStore.SetAvatar(prepareUrl(import.meta.env.VITE_API_BASEPATH as string + data.image.url));
    } else {
      appStore.SetAvatar('')
    }

    const res = await api.v1.userServiceUpdateUserById(userId.value, body)
        .catch(() => {
        })
        .finally(() => {
          loading.value = false
        })
    if (res) {
      fetch()
      ElMessage({
        title: t('Success'),
        message: t('message.updatedSuccessfully'),
        type: 'success',
        duration: 2000
      })
    }
  }
}

const cancel = () => {
  push('/etc/users')
}

const remove = async () => {
  loading.value = true
  const res = await api.v1.userServiceDeleteUserById(userId.value)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })
  if (res) {
    cancel()
  }
}

fetch()

</script>

<template>
  <ContentWrap>
    <Form ref="writeRef" :current-row="currentUser"/>

    <div style="text-align: right">

      <ElButton type="primary" @click="save()">
        {{ t('main.save') }}
      </ElButton>

      <ElButton type="default" @click="cancel()">
        {{ t('main.return') }}
      </ElButton>

      <ElPopconfirm
          :confirm-button-text="$t('main.ok')"
          :cancel-button-text="$t('main.no')"
          width="250"
          style="margin-left: 10px;"
          :title="$t('main.are_you_sure_to_do_want_this?')"
          @confirm="remove"
      >
        <template #reference>
          <ElButton class="mr-10px" type="danger" plain>
            <Icon icon="ep:delete" class="mr-5px"/>
            {{ t('main.remove') }}
          </ElButton>
        </template>
      </ElPopconfirm>

    </div>
  </ContentWrap>
</template>

<style lang="less" scoped>

</style>
