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
import {ApiNewScriptRequest} from "@/api/stub";
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";
import {useEmitt} from "@/hooks/web/useEmitt";

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
const currentScript = ref<Nullable<ApiNewScriptRequest>>(null)
const sourceScript = ref('')

const save = async () => {
  const write = unref(writeRef)
  const validate = await write?.elFormRef?.validate()?.catch(() => {
  })
  if (validate) {
    loading.value = true
    const data = (await write?.getFormData())
    const body = {
      lang: data.lang,
      name: data.name,
      source: sourceScript.value,
      description: data.description,
    } as ApiNewScriptRequest
    const res = await api.v1.scriptServiceAddScript(body)
        .catch(() => {
        })
        .finally(() => {
          loading.value = false
        })
    if (res) {
      const {id} = res.data;
      push(`/scripts/edit/${id}`)
    }
  }
}

const cancel = () => {
  push('/scripts')
}

useEmitt({
  name: 'updateSource',
  callback: (val: string) => {
    if (sourceScript.value == val) {
      return
    }
    sourceScript.value = val
  }
})

</script>

<template>
  <ContentWrap>
    <Form ref="writeRef" :current-row="currentScript"/>

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
