<script setup lang="ts">
import {computed, ref, unref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {ElButton, ElMessage, ElPopconfirm} from 'element-plus'
import {useRoute, useRouter} from 'vue-router'
import api from "@/api/api";
import Form from './components/Form.vue'
import {ApiAction} from "@/api/stub";
import {ContentWrap} from "@/components/ContentWrap";
import ActionForm from "@/views/Automation/components/ActionForm.vue";

const {push} = useRouter()
const route = useRoute();
const {t} = useI18n()

const writeRef = ref<ComponentRef<typeof Form>>()
const loading = ref(true)
const actionId = computed(() => route.params.id as number);
const currentRow = ref<Nullable<ApiAction>>(null)

const fetch = async () => {
  loading.value = true
  const res = await api.v1.actionServiceGetActionById(actionId.value)
    .catch(() => {
    })
    .finally(() => {
      loading.value = false
    })
  if (res) {
    currentRow.value = res.data;
  } else {
    currentRow.value = null
  }
}

const save = async () => {
  const write = unref(writeRef)
  const validate = await write?.elFormRef?.validate()?.catch(() => {
  })
  if (validate) {
    loading.value = true
    const act = (await write?.getFormData()) as ApiAction;
    const data = {
      name: act.name,
      description: act.description,
      scriptId: act.script?.id,
      entityId: act.entity?.id,
      areaId: act.area?.id,
      entityActionName: act.entityActionName,
    }
    const res = await api.v1.actionServiceUpdateAction(actionId.value, data)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
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
}

const cancel = () => {
  push('/automation/actions')
}

const remove = async () => {
  loading.value = true
  const res = await api.v1.actionServiceDeleteAction(actionId.value)
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

    <ActionForm ref="writeRef" :action="currentRow"/>

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
