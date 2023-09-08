<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {PropType, reactive, ref, unref, watch} from 'vue'
import {TableColumn} from '@/types/table'
import {ElButton, ElTag, ElPopconfirm} from 'element-plus'
import {ApiCondition} from "@/api/stub";
import {useRouter} from "vue-router";
import ConditionForm from "./ConditionForm.vue";
import {useEmitt} from "@/hooks/web/useEmitt";

const {currentRoute, addRoute, push} = useRouter()
const props = defineProps({
  conditions: {
    type: Array as PropType<ApiCondition[]>,
    default: () => []
  }
})

const writeRef = ref<ComponentRef<typeof ConditionForm>>()
let currentCondition = reactive<Nullable<ApiCondition>>(null)
let currentConditionIndex = ref(0)

enum Mode {
  VIEW = 'VIEW',
  EDIT = 'EDIT',
  NEW = 'NEW'
}

const mode = ref<Mode>(Mode.VIEW)
const {t} = useI18n()

interface TableObject {
  tableList: ApiCondition[]
}

const tableObject = reactive<TableObject>(
    {
      tableList: [],
    }
);

watch(
    () => props.conditions,
    (val: ApiCondition[]) => {
      if (val == unref(tableObject.tableList)) return;
      tableObject.tableList = val
    },
    {
      immediate: true
    }
)

const columns: TableColumn[] =  [
  {
    field: 'name',
    label: t('conditions.name'),
  },
  {
    field: 'script',
    label: t('conditions.script'),
  },
  {
    field: 'operations',
    label: t('conditions.operations'),
    width: "200px",
    type: 'time',

  },
]

const addNew = () => {
  currentCondition = null
  mode.value = Mode.NEW
}

const {emitter} = useEmitt()
const updateConditions = () => {
  const triggers = unref(tableObject.tableList)
  emitter.emit('updateConditions', triggers)
}

const edit = (condition: ApiCondition, $index) => {
  currentConditionIndex.value = $index
  currentCondition = unref(condition)
  mode.value = Mode.EDIT
}

const save = async () => {
  const write = unref(writeRef)
  const validate = await write?.elFormRef?.validate()?.catch(() => {
  })
  const data = (await write?.getFormData())
  if (validate) {
    const condition = {
      id: data.id,
      name: data.name,
      script: data.script,
      scriptId: data.script?.id || null,
    } as ApiCondition
    if (mode.value === Mode.NEW) {
      tableObject.tableList.push(condition)
    } else {
      tableObject.tableList[currentConditionIndex.value] = condition
    }
    mode.value = Mode.VIEW
    updateConditions()
  }
}

const resetForm = () => {
  mode.value = Mode.VIEW
  currentCondition = null
}

const removeItem = () => {
  tableObject.tableList.splice(currentConditionIndex.value, 1);
  mode.value = Mode.VIEW
}

const goToScript = (trigger: ApiCondition) => {
  if (!trigger.script?.id) {
    return
  }
  push(`/scripts/edit/${trigger.script?.id}`)
}

</script>

<template>
  <ElButton class="flex mb-20px items-left" @click="addNew()" plain v-if="mode==='VIEW'">
    <Icon icon="ep:plus" class="mr-5px"/>
    {{ t('conditions.addNew') }}
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

    <template #operations="{ row }">
      <ElButton :link="true" @click.prevent.stop="edit(row)">
        {{ $t('main.edit') }}
      </ElButton>

    </template>

  </Table>

  <ConditionForm ref="writeRef" :condition="currentCondition" v-if="mode!=='VIEW'"/>

  <div style="text-align: right" v-if="mode!=='VIEW'">

    <ElButton v-if="mode === 'NEW'" type="primary" @click="save()">
      {{ t('conditions.addCondition') }}
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
