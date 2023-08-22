<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {PropType, ref, unref, watch} from 'vue'
import {ElCol, ElRow, ElStatistic} from 'element-plus'
import {ApiStatistic, ApiStatistics} from "@/api/stub";
import {propTypes} from "@/utils/propTypes";

const statistic = ref<Statistic>({items: []})
const rowGutter = ref(24)
const colSpan = ref(8)
const apiStatistics = ref<Nullable<ApiStatistics>>(null)

const {t} = useI18n()

interface Statistic {
  items: ApiStatistic[][];
}

const colorList = [
    'linear-gradient(rgb(61, 73, 46) 0%, rgb(38, 56, 39) 100%)',
    'linear-gradient(rgb(40, 73, 145) 0%, rgb(18, 43, 98) 100%)',
    'linear-gradient(rgb(49, 37, 101) 0%, rgb(32, 25, 54) 100%)',
    'linear-gradient(#457b9d 0%, #30556d 100%)',
    'linear-gradient(#dda15e 0%, #78562f 100%)',
    'linear-gradient(#40916c 0%, #255640 100%)'
]

const props = defineProps({
  modelValue: {
    type: Object as PropType<Nullable<ApiStatistics>>,
    default: () => null
  },
  cols: propTypes.number.def(3),
})

watch(
    () => props.cols,
    (val?: number) => {
      colSpan.value = rowGutter.value / val
    }
)

watch(
    () => props.modelValue,
    (val?: ApiStatistics) => {
      if (val === unref(apiStatistics)) return;
      apiStatistics.value = val || null;
      const items = val?.items || [];
      let row = 0;
      statistic.value.items = []
      for (const index in val?.items) {
        row = Math.floor(index/props.cols)
        if (!statistic.value.items[row]) {
          statistic.value.items[row] = []
        }
        statistic.value.items[row].push(items[index])
      }
    },
)

const getStyle = (index, index2) => {
  return {'background': colorList[index2]}
}

</script>

<template>
  <div class="ml-20px mr-20px" v-if="statistic">
    <ElRow :gutter="rowGutter" v-for="(cols, $index) in statistic.items" :key="$index" >
      <ElCol :span="colSpan" v-for="(col, $index2) in cols" :key="$index2" class="mt-20px">
        <div class="statistic-card" :style="getStyle($index, $index2)">
          <ElStatistic :value="col.value">
            <template #title>
              <div style="display: inline-flex; align-items: center">
                {{ $t(col.name) }}
              </div>
            </template>
          </ElStatistic>
        </div>
      </ElCol>
    </ElRow>
  </div>
</template>

<style lang="less">

.el-statistic {
  --el-statistic-content-font-size: 28px;
}

.statistic-card {
  height: 100%;
  padding: 20px;
  border-radius: 4px;
  background-color: var(--el-bg-color-overlay);
}

.statistic-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  font-size: 12px;
  color: var(--el-text-color-regular);
  margin-top: 16px;
}

.statistic-footer .footer-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.statistic-footer .footer-item span:last-child {
  display: inline-flex;
  align-items: center;
  margin-left: 4px;
}

.green {
  color: var(--el-color-success);
}

.red {
  color: var(--el-color-error);
}

.el-statistic__head {
  color: #bbb;
}
.el-statistic__content {
  color: #fff;
}
</style>
