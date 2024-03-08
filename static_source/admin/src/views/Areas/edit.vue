<script setup lang="ts">
import {computed, ref, unref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {ElButton, ElPopconfirm} from 'element-plus'
import {useRoute, useRouter} from 'vue-router'
import api from "@/api/api";
import Form from './Form.vue'
import {ApiArea} from "@/api/stub";
import {ContentWrap} from "@/components/ContentWrap";
import {MapEditor} from "@/components/MapEditor";

const {push} = useRouter()
const route = useRoute();
const {t} = useI18n()

const writeRef = ref<ComponentRef<typeof Form>>()
const loading = ref(false)
const areaId = computed(() => route.params.id as number);
const currentRow = ref<Nullable<ApiArea>>(null)

const fetch = async () => {
  loading.value = true
  const res = await api.v1.areaServiceGetAreaById(areaId.value)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })
  if (res) {
    currentRow.value = res.data
  } else {
    currentRow.value = null
  }
}

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
    const res = await api.v1.areaServiceUpdateArea(areaId.value, body)
        .catch(() => {
        })
        .finally(() => {
          loading.value = false
        })
  }
}

const cancel = () => {
  push('/etc/areas')
}

const remove = async () => {
  loading.value = true
  const res = await api.v1.areaServiceDeleteArea(areaId.value)
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
    <Form ref="writeRef" :current-row="currentRow"/>

    <MapEditor ref="child" :area="currentRow"/>

    <div style="text-align: right" class="mt-20px">

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
