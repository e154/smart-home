<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {PropType, reactive, ref, unref, watch} from 'vue'
import {TableColumn} from '@/types/table'
import {ElButton, ElTag, ElPopconfirm} from 'element-plus'
import {ApiAttribute, ApiTrigger, ApiTypes} from "@/api/stub";
import {useEmitt} from "@/hooks/web/useEmitt";
import TriggerForm from "./TriggerForm.vue";
import {Trigger} from "@/views/Automation/components/types";
import {useRouter} from "vue-router";

const {currentRoute, addRoute, push} = useRouter()
const props = defineProps({
  triggers: {
    type: Array as PropType<ApiTrigger[]>,
    default: () => []
  }
})

const writeRef = ref<ComponentRef<typeof TriggerForm>>()
let currentTrigger = reactive<Nullable<ApiTrigger>>(null)
let currentTriggerIndex = ref(0)

enum Mode {
  VIEW = 'VIEW',
  EDIT = 'EDIT',
  NEW = 'NEW'
}

const mode = ref<Mode>(Mode.VIEW)
const {t} = useI18n()

interface TableObject {
  tableList: ApiTrigger[]
}

const tableObject = reactive<TableObject>(
    {
      tableList: [],
    }
);

watch(
    () => props.triggers,
    (val: ApiTrigger[]) => {
      if (val == unref(tableObject.tableList)) return;
      tableObject.tableList = val
      if (val.length) {
        currentTrigger = val[0]
      }
    },
    {
      immediate: true
    }
)

const columns: TableColumn[] = [
  {
    field: 'name',
    label: t('triggers.name'),
  },
  {
    field: 'script',
    label: t('triggers.script'),
  },
  {
    field: 'entity',
    label: t('triggers.entity'),
  },
  {
    field: 'plugin',
    label: t('triggers.plugin'),
  },
  {
    field: 'operations',
    label: t('triggers.operations'),
    width: "200px",
    type: 'time',

  },
]

const addNew = () => {
  currentTrigger = null
  mode.value = Mode.NEW
}

const {emitter} = useEmitt()
const call = (trigger: ApiTrigger) => {
  emitter.emit('callTrigger', trigger.name)
}

const updateTriggers = () => {
  const triggers = unref(tableObject.tableList)
  emitter.emit('updateTriggers', triggers)
}

const edit = (trigger: ApiTrigger, $index) => {
  currentTriggerIndex.value = $index
  currentTrigger = unref(trigger)
  mode.value = Mode.EDIT
}

const save = async () => {
  const write = unref(writeRef)
  const validate = await write?.elFormRef?.validate()?.catch(() => {
  })
  const data = (await write?.getFormData()) as Trigger
  if (validate) {
    let attributes: { [key: string]: ApiAttribute } = {};
    if (data.timePluginOptions) {
      attributes.cron = {
        name: 'cron',
        type: ApiTypes.STRING,
        string: data.timePluginOptions
      };
    }
    if (data.alexaPluginOptions) {
      attributes.skillId = {
        name: 'skillId',
        type: ApiTypes.INT,
        int: data.alexaPluginOptions
      };
    }
    const trigger = {
      id: data?.id,
      name: data.name,
      pluginName: data.pluginName,
      entity: data.entity,
      entityId: data.entity?.id || null,
      script: data.script,
      scriptId: data.script?.id || null,
      attributes: attributes,
    } as ApiTrigger
    if (mode.value === Mode.NEW) {
      tableObject.tableList.push(trigger)
    } else {
      tableObject.tableList[currentTriggerIndex.value] = trigger
    }
    mode.value = Mode.VIEW
    updateTriggers()
  }
}

const resetForm = () => {
  mode.value = Mode.VIEW
  currentTrigger = null
}

const removeItem = () => {
  tableObject.tableList.splice(currentTriggerIndex.value, 1);
  mode.value = Mode.VIEW
}

const goToEntity = (trigger: ApiTrigger) => {
  if (!trigger.entity?.id) {
    return
  }
  push(`/entities/edit/${trigger.entity?.id}`)
}

const goToScript = (trigger: ApiTrigger) => {
  if (!trigger.script?.id) {
    return
  }
  push(`/scripts/edit/${trigger.script?.id}`)
}

</script>

<template>
  <ElButton class="flex mb-20px items-left"  @click="addNew()" plain v-if="mode==='VIEW'">
    <Icon icon="ep:plus" class="mr-5px"/>
    {{ t('triggers.addNew') }}
  </ElButton>

  <Table
      v-if="mode==='VIEW'"
      :selection="false"
      :columns="columns"
      :data="tableObject.tableList"
      style="width: 100%"
  >

    <template #script="{row}">
      <ElButton v-if="row.script" :link="true" @click.prevent.stop="goToScript(row)">
        {{ row.script.name }}
      </ElButton>
      <span v-else>-</span>
    </template>

    <template #entity="{row}">
      <ElButton v-if="row.entity" :link="true" @click.prevent.stop="goToEntity(row)">
        {{ row.entity.id }}
      </ElButton>
      <span v-else>-</span>
    </template>

    <template #plugin="{row}">
      <ElTag type="info">
        {{ row.pluginName }}
      </ElTag>
    </template>

    <template #operations="{ row, $index }">

      <ElButton :link="true" @click.prevent.stop="call(row)">
        {{ $t('main.call') }}
      </ElButton>

      <ElButton :link="true" @click.prevent.stop="edit(row, $index)">
        {{ $t('main.edit') }}
      </ElButton>

    </template>

  </Table>

  <TriggerForm ref="writeRef" :trigger="currentTrigger" v-if="mode!=='VIEW'"/>

  <div style="text-align: right" v-if="mode!=='VIEW'">

    <ElButton v-if="mode === 'NEW'" type="primary" @click="save()">
      {{ t('triggers.addTrigger') }}
    </ElButton>

    <ElButton v-if="mode === 'EDIT'" type="primary" @click="save()">
      {{ t('main.update') }}
    </ElButton>

    <ElButton type="default" @click="resetForm()">
      {{ t('main.cancel') }}
    </ElButton>

    <ElPopconfirm
        v-if="mode === 'EDIT'"
        :confirm-button-text="$t('main.ok')"
        :cancel-button-text="$t('main.no')"
        width="250"
        style="margin-left: 10px;"
        :title="$t('main.are_you_sure_to_do_want_this?')"
        @confirm="removeItem"
    >
      <template #reference>
        <ElButton class="mr-10px" type="danger" plain>
          <Icon icon="ep:delete" class="mr-5px"/>
          {{ t('main.remove') }}
        </ElButton>
      </template>
    </ElPopconfirm>

  </div>
</template>

<style lang="less">

</style>
