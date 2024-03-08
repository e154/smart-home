<script setup lang="ts">
import {computed, ref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {ElButton, ElPopconfirm} from 'element-plus'
import {useRoute, useRouter} from 'vue-router'
import api from "@/api/api";
import {ApiVariable} from "@/api/stub";
import {ContentWrap} from "@/components/ContentWrap";
import VariableForm from "@/views/Variables/components/VariableForm.vue";

const {push} = useRouter()
const route = useRoute();
const {t} = useI18n()

const variableName = computed(() => route.params.name);
const currentRow = ref<Nullable<ApiVariable>>(null)

const fetch = async () => {
  const res = await api.v1.variableServiceGetVariableByName(variableName.value as string)
      .catch(() => {
      })
      .finally(() => {
      })
  if (res) {
    currentRow.value = res.data
  } else {
    currentRow.value = null
  }
}

const download = async () => {
  // window.location.href = 'data:application/octet-stream;base64,' + currentRow.value.value;
  const a = document.createElement("a"); //Create <a>
  a.href = 'data:application/octet-stream;base64,' + currentRow.value.value;
  a.download = currentRow.value.name; //File name Here
  a.click();
}

const save = async () => {
  const data = {
    value: currentRow.value.value,
    tags: currentRow.value.tags,
  }
  const res = await api.v1.variableServiceUpdateVariable(variableName.value as string, data)
      .catch(() => {
      })
      .finally(() => {
      })
 
}

const cancel = () => {
  push('/etc/variables')
}

const remove = async () => {
  const res = await api.v1.variableServiceDeleteVariable(variableName.value as string)
      .catch(() => {
      })
      .finally(() => {
      })
  if (res) {
    cancel()
  }
}
fetch()

</script>

<template>
  <ContentWrap>
    <VariableForm v-if="currentRow" v-model="currentRow" :edit="true"/>

    <div style="text-align: right">

      <ElButton type="primary" @click="download()">
        {{ t('main.download') }}
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
