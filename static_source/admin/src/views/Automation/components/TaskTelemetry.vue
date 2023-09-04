<script setup lang="ts">
import {computed, defineEmits, onMounted, onUnmounted, PropType, watch, ref, nextTick} from 'vue'
import {Form} from '@/components/Form'
import {useI18n} from '@/hooks/web/useI18n'
import {ElSkeleton} from 'element-plus'
import {useForm} from '@/hooks/web/useForm'
import {useCache} from '@/hooks/web/useCache'
import {useAppStore} from '@/store/modules/app'
import {usePermissionStore} from '@/store/modules/permission'
import {useRoute, useRouter} from 'vue-router'
import {useValidator} from '@/hooks/web/useValidator'
import {FormSchema} from '@/types/form'
import {ApiPlugin, ApiTask, ApiTelemetryItem} from "@/api/stub";
import api from "@/api/api";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";
import {EventTaskCompleted, EventTriggerCompleted} from "@/api/stream_types";
import {debounce} from "lodash-es";
import {EChartsOption} from "echarts";
import { Echart } from '@/components/Echart'
import * as echarts from 'echarts/core'
import {formatTime, parseTime} from "@/utils";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const appStore = useAppStore()
const permissionStore = usePermissionStore()
const {currentRoute, addRoute, push} = useRouter()
const route = useRoute();
const {wsCache} = useCache()
const {t} = useI18n()

const telemetry = ref<ApiTelemetryItem[]>([])
const props = defineProps({
  task: {
    type: Object as PropType<Nullable<ApiTask>>,
    default: () => null
  }
})

watch(
    () => props.task,
    (val) => {
      if (!val) return
      telemetry.value = val?.telemetry || undefined;
    },
    {
      immediate: true
    }
)

const currentID = ref('')
onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()
  stream.subscribe('event_task_completed', currentID.value, onEventTaskActivated);
})

onUnmounted(() => {
  stream.unsubscribe('event_task_completed', currentID.value);
})

const fetch = debounce( async () => {
  const res = await api.v1.automationServiceGetTask(props.task?.id)
      .catch(() => {})
      .finally(() => {})
  if (res) {
    const task = res.data as ApiTask;
    telemetry.value = task?.telemetry;
    getOptions()
  }
}, 100)

const onEventTaskActivated = (event: EventTaskCompleted) => {
  if (event.id != props.task?.id) {
    return;
  }
  fetch()
}

const chartOptions = ref<EChartsOption>({})

const reloadKey = ref(0)
const reload = () => {
  reloadKey.value += 1
}

const renderItem = (params, api) => {
  let categoryIndex = api.value(0);
  let start = api.coord([api.value(1), categoryIndex]);
  let end = api.coord([api.value(2), categoryIndex]);
  let height = api.size([0, 1])[1] * 0.6;
  let rectShape = echarts.graphic.clipRectByRect(
      {
        x: start[0],
        y: start[1] - height / 2,
        width: end[0] - start[0],
        height: height
      },
      {
        x: params.coordSys.x,
        y: params.coordSys.y,
        width: params.coordSys.width,
        height: params.coordSys.height
      }
  );
  return (
      rectShape && {
        type: 'rect',
        transition: ['shape'],
        shape: rectShape,
        style: api.style()
      }
  );
}

var colors = ['#7b9ce1','#bd6d6c','#75d874','#e0bc78','#dc77dc','#72b362','#7b9ce1','#bd6d6c','#75d874','#e0bc78','#dc77dc','#72b362'];

const getID = (attrs: Record<string, string>) => {
  if (attrs['id']) {
    return `(id: ${attrs['id']})` || ''
  } else {
    return ''
  }
}
const getOptions = () => {

  let startTime = 0;
  let categories = [];
  let data = [];

  for (const item of telemetry.value) {
    const label = `level${item.level}`
    if (categories.indexOf(label) === -1) {
      categories.push(label)
    }
    const start =  item?.start;
    const end =  item?.end;
    if (startTime === 0) {
      startTime = start || 0;
    }
    const timeEstimate = item?.timeEstimate;
    data.push({
      name: item.name + getID(item.attributes),
      value: [item.level - 1, start, end, timeEstimate],
      itemStyle: {
        normal: {
          color: colors[item.num],
        }
      }
    })
  }

  let options: EChartsOption = {
    tooltip: {
      formatter: function (params) {
        return params.marker + params.name + ': ' + params.value[3] + ' ms';
      }
    },
    grid: {
      height: 300
    },
    xAxis: {
      min: startTime,
      scale: true,
      axisLabel: {
        formatter: function (val) {
          return Math.max(0, val - startTime) + ' ms';
        }
      }
    },
    yAxis: {
      data: categories
    },
    series: [
      {
        type: 'custom',
        renderItem: renderItem,
        itemStyle: {
          opacity: 0.8
        },
        encode: {
          x: [1, 2],
          y: 0
        },
        data: data
      }
    ]
  }

  nextTick(() => {
    chartOptions.value = options;
  })

}

getOptions()

</script>

<template>
  <div class="h-[100%] w-[100%]" style="height: 400px" v-if="telemetry">
    <Echart :options="chartOptions" :key="reloadKey"/>
  </div>
  <ElSkeleton v-else :rows="5" />
</template>
