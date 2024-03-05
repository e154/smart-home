<script setup lang="ts">
import {onMounted, onUnmounted, ref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {ElButton} from 'element-plus'
import {useRouter} from 'vue-router'
import api from "@/api/api";
import {ApiVariable} from "@/api/stub";
import {ContentWrap} from "@/components/ContentWrap";
import VariableForm from "@/views/Variables/components/VariableForm.vue";

const {push} = useRouter()
const {t} = useI18n()

const currentRow = ref<Nullable<ApiVariable>>(null)

onMounted(() => {
  currentRow.value = {
    name: '',
    value: '',
    tags: []
  } as ApiVariable
})

onUnmounted(() => {

})

const save = async () => {
  const res = await api.v1.variableServiceAddVariable(currentRow.value)
      .catch(() => {
      })
      .finally(() => {
      })
  if (res) {
    push(`/etc/variables/edit/${currentRow.value.name}`)
  }
}

const cancel = () => {
  push('/etc/variables')
}

</script>

<template>
  <ContentWrap>
    <VariableForm v-if="currentRow" v-model="currentRow"/>

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
