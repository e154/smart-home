<script setup lang="ts">
import {computed, defineEmits, ref, unref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {ElButton, ElMessage, ElPopconfirm} from 'element-plus'
import {useForm} from '@/hooks/web/useForm'
import {useRoute, useRouter} from 'vue-router'
import {useValidator} from '@/hooks/web/useValidator'
import api from "@/api/api";
import Form from './components/Form.vue'
import {ApiCondition, ApiScript} from "@/api/stub";
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";
import ConditionForm from "@/views/Automation/components/ConditionForm.vue";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const {currentRoute, addRoute, push} = useRouter()
const route = useRoute();
const {t} = useI18n()

const writeRef = ref<ComponentRef<typeof Form>>()
const loading = ref(true)
const conditionId = computed(() => route.params.id as number);
const currentRow = ref<Nullable<ApiCondition>>(null)

const fetch = async () => {
  loading.value = true
  const res = await api.v1.conditionServiceGetConditionById(conditionId.value)
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
    const con = (await write?.getFormData()) as ApiCondition;
    const data = {
      name: con.name,
      description: con.description,
      scriptId: con.script?.id,
      areaId: con.area?.id,
    }
    const res = await api.v1.conditionServiceUpdateCondition(conditionId.value, data)
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
  push('/automation/conditions')
}

const remove = async () => {
  loading.value = true
  const res = await api.v1.conditionServiceDeleteCondition(conditionId.value)
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

    <ConditionForm ref="writeRef" :condition="currentRow"/>

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
