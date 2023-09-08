<script setup lang="ts">
import {onMounted, ref, unref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {ElButton, ElCol, ElRow} from 'element-plus'
import {useForm} from '@/hooks/web/useForm'
import {useCache} from '@/hooks/web/useCache'
import {useAppStore} from '@/store/modules/app'
import {usePermissionStore} from '@/store/modules/permission'
import {useRoute, useRouter} from 'vue-router'
import {useValidator} from '@/hooks/web/useValidator'
import api from "@/api/api";
import Form from './components/Form.vue'
import {ApiArea} from "@/api/stub";
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";
import MapEditor from "@/views/Areas/components/MapEditor.vue";

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
const currentRow = ref<Nullable<ApiArea>>(null)

onMounted(() => {
  currentRow.value = {
    name: '',
    description: '',
    polygon: [],
    center: {lat: 0, lon: 0},
    zoom: 0,
    resolution: 0,
  }
})

const child = ref(null)
const save = async () => {
  const {polygon, zoom, center, resolution} = child.value.save()

  const write = unref(writeRef)
  const validate = await write?.elFormRef?.validate()?.catch(() => {
  })
  if (validate) {
    loading.value = true
    const data = (await write?.getFormData()) as ApiArea
    const body = {
      name: data.name,
      description: data.description,
      polygon: polygon,
      center: center,
      zoom: zoom,
      resolution: resolution,
    }
    const res = await api.v1.areaServiceAddArea(body)
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
  push('/etc/areas')
}

</script>

<template>
  <ContentWrap>
    <Form ref="writeRef" :current-row="currentRow"/>

    <MapEditor ref="child" :area="currentRow"/>

    <div style="text-align: right" class="mt-20px">

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
