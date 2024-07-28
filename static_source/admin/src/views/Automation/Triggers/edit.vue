<script setup lang="ts">
import {computed, ref, unref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {ElButton, ElMessage, ElPopconfirm} from 'element-plus'
import {useRoute, useRouter} from 'vue-router'
import api from "@/api/api";
import Form from './components/Form.vue'
import {ApiAttribute, ApiTrigger} from "@/api/stub";
import {ContentWrap} from "@/components/ContentWrap";
import TriggerForm from "@/views/Automation/components/TriggerForm.vue";

const {currentRoute, addRoute, push} = useRouter()
const route = useRoute();
const {t} = useI18n()

const writeRef = ref<ComponentRef<typeof Form>>()
const loading = ref(true)
const triggerId = computed(() => +route.params.id);
const currentRow = ref<Nullable<ApiTrigger>>(null)

const fetch = async () => {
  loading.value = true
  const res = await api.v1.triggerServiceGetTriggerById(triggerId.value)
    .catch(() => {
    })
    .finally(() => {
      loading.value = false
    })
  if (res) {
    currentRow.value = res.data;
    if (currentRow.value.attributes) {
      for (const key in currentRow.value.attributes) {
        if (key == 'notice') {
          continue
        }
        currentRow.value.attributes[key].name = key
      }
    }
  } else {
    currentRow.value = null
  }
}

const save = async () => {
  const write = unref(writeRef)
  const validate = await write?.form?.validate()?.catch(() => {
  })
  if (!validate) {
    return
  }

  const tr = unref(currentRow)
  let data = {
    name: tr.name,
    description: tr.description,
    entityIds: tr.entityIds || [],
    scriptId: tr.script?.id || null,
    areaId: tr.area?.id || null,
    pluginName: tr.pluginName,
    attributes: {},
    enabled: tr.enabled,
  }

  let attributes: { [key: string]: ApiAttribute } = {};
  for (const index in tr.attributes) {
    attributes[tr.attributes[index].name] = tr.attributes[index];
  }
  data.attributes = attributes

  const res = await api.v1.triggerServiceUpdateTrigger(triggerId.value, data)
    .catch(() => {
    })
    .finally(() => {
    })
  if (res) {
    ElMessage({
      title: t('Success'),
      message: t('message.updatedSuccessfully'),
      type: 'success',
      duration: 2000
    })
  }

}

const cancel = () => {
  push('/automation/triggers')
}

const remove = async () => {
  loading.value = true
  const res = await api.v1.triggerServiceDeleteTrigger(triggerId.value)
    .catch(() => {
    })
    .finally(() => {
      loading.value = false
    })
  if (res) {
    cancel()
  }
}

const callTrigger = async () => {
  await api.v1.developerToolsServiceCallTrigger({id: triggerId.value})
    .catch(() => {
    })
    .finally(() => {
      ElMessage({
        title: t('Success'),
        message: t('message.callSuccessful'),
        type: 'success',
        duration: 2000
      })
    })
}

fetch()

</script>

<template>
  <ContentWrap>

    <TriggerForm ref="writeRef" :trigger="currentRow"/>

    <div style="text-align: right">

      <ElButton type="success" @click="callTrigger()">
        {{ t('main.call') }}
      </ElButton>

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
