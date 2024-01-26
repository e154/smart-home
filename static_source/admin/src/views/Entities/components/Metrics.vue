<script setup lang="ts">

import {useI18n} from '@/hooks/web/useI18n'
import {PropType, reactive, ref, unref, watch} from 'vue'
import {ElButton, ElRow, ElCol, ElCard, ElPopconfirm, ElSkeleton, ElEmpty} from 'element-plus'
import {ApiMetric} from "@/api/stub";
import MetricForm from "@/views/Entities/components/MetricForm.vue";

const {t} = useI18n()
const props = defineProps({
  metrics: {
    type: Array as PropType<ApiMetric[]>,
    default: () => []
  }
})

interface Current {
  metrics: ApiMetric[];
  item?: ApiMetric;
  index: number;
}

const current = reactive<Current>({
  metrics: [],
  index: -1,
})

watch(
    () => props.metrics,
    (val: ApiMetric[]) => {
      if (val == unref(current.metrics)) return;
      current.metrics = val
      if (val.length) {
        current.index = 0
        current.item = val[0]
      }
    },
    {
      immediate: true
    }
)

const addNew = () => {
  current.metrics.push({
    description: undefined,
    name: `new metric ${current.metrics.length}`,
    ranges: [],
    type: 'LINE',
    options: {
      items: []
    }
  } as ApiMetric);
  current.index = current.metrics.length - 1 || 0;
  current.item = Object.assign({}, current.metrics[current.index])
}

const edit = (val: ApiMetric, $index) => {
  current.item = val
  current.index = $index
}

const addProperty = () => {

}
const remove = () => {
  if (!current.metrics || !current.metrics.length || current.index < 0) {
    return;
  }
  current.metrics.splice(current.index, 1);
  current.index = current.metrics.length - 1;

  if (current.metrics.length) {
    current.index = current.metrics.length - 1
    current.item = current.metrics[current.index]
  } else {
    current.index = -1
    current.item = undefined
  }
}


</script>

<template>
  <ElRow :gutter="20" >
    <ElCol :span="6" :xs="12">
      <ElCard class="box-card">
        <template #header>
          <div class="card-header">
            <span>{{$t('metrics.list')}}</span>
            <ElButton @click="addNew()" text>
              {{ t('entities.addNewMetric') }}
            </ElButton>
          </div>
        </template>
        <div @click="edit(metric, index)" v-for="(metric, index) in metrics" :key="metric" class="text item cursor-pointer">{{ metric.name }}</div>
          <div v-if="!metrics.length" class="text item">{{$t('metrics.noMetrics')}}</div>
      </ElCard>
    </ElCol>

    <ElCol :span="14" :xs="12">
      <ElCol>
        <MetricForm v-model="current.item"/>

        <ElEmpty v-if="!current.item" :rows="5" class="mt-20px mb-20px">
          <ElButton type="primary" @click="addNew()">
            {{ t('entities.addNewMetric') }}
          </ElButton>
        </ElEmpty>

        <ElRow v-if="current.item">
          <ElCol>
            <div style="padding-bottom: 20px">
              <div style="text-align: right;" class="mt-20px">
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
                      {{ t('metrics.removeMetric') }}
                    </ElButton>
                  </template>
                </ElPopconfirm>
              </div>
            </div>
          </ElCol>
        </ElRow>
      </ElCol>
    </ElCol>



  </ElRow>
</template>

<style lang="less" scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

</style>
