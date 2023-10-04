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
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";
import {Entity} from "@/views/Entities/components/types";
import {ApiNewEntityRequest} from "@/api/stub";

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
const currentRow = ref<Nullable<Entity>>(null)

const save = async () => {
  const write = unref(writeRef)
  const validate = await write?.elFormRef?.validate()?.catch(() => {
  })
  if (validate) {
    loading.value = true
    const data = (await write?.getFormData()) as Entity
    const body = {
      name: data.name,
      pluginName: data.plugin?.name,
      description: data.description,
      areaId: data.area?.id,
      icon: data.icon,
      imageId: data.image?.id,
      autoLoad: data.autoLoad,
      scriptIds: data.scriptIds,
      parentId: data.parent?.id,
    } as ApiNewEntityRequest
    const res = await api.v1.entityServiceAddEntity(body)
        .catch(() => {
        })
        .finally(() => {
          loading.value = false
        })
    if (res) {
      push(`/entities/edit/${res.data.id}`)
    }
  }
}

const cancel = () => {
  push('/entities')
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
