<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {PropType, ref, unref, watch} from 'vue'
import {ElCol, ElIcon, ElRow, ElStatistic, ElTag, ElTooltip} from 'element-plus'
import {ApiStatistics, ApiStatistic} from "@/api/stub";
import {propTypes} from "@/utils/propTypes";

const statistic = ref<Statistic>({items: []})
const rowGutter = ref(24)
const colSpan = ref(8)
const apiStatistics = ref<Nullable<ApiStatistics>>(null)

const {t} = useI18n()

interface Statistic {
  items: ApiStatistic[][];
}

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
    {
      immediate: true
    }
)

</script>

<template>
  <div class="ml-20px mr-20px" v-if="statistic">
    <ElRow :gutter="rowGutter" v-for="(cols, $index) in statistic.items" :key="$index" >
      <ElCol :span="colSpan" v-for="(col, $index2) in cols" :key="$index2" class="mt-20px">
        <div class="statistic-card">
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

.el-table__row {
  cursor: pointer;
}

:global(h2#card-usage ~ .example .example-showcase) {
  background-color: var(--el-fill-color) !important;
}

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
</style>
