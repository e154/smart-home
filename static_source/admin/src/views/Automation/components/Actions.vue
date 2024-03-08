<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {PropType, reactive, ref, unref, watch} from 'vue'
import {TableColumn} from '@/types/table'
import {ElButton, ElTag, ElPopconfirm} from 'element-plus'
import {ApiAction, ApiAttribute, ApiTrigger, ApiTypes} from "@/api/stub";
import {useEmitt} from "@/hooks/web/useEmitt";
import ActionForm from "@/views/Automation/components/ActionForm.vue";
import {useRouter} from "vue-router";
import {Trigger} from "@/views/Automation/components/types";

const {currentRoute, addRoute, push} = useRouter()
const props = defineProps({
  actions: {
    type: Array as PropType<ApiAction[]>,
    default: () => []
  }
})

const writeRef = ref<ComponentRef<typeof ActionForm>>()
let currentAction = reactive<Nullable<ApiAction>>(null)
let currentActionIndex = ref(0)

enum Mode {
  VIEW = 'VIEW',
  EDIT = 'EDIT',
  NEW = 'NEW'
}

const mode = ref<Mode>(Mode.VIEW)
const {t} = useI18n()

interface TableObject {
  tableList: ApiAction[]
}

const tableObject = reactive<TableObject>(
    {
      tableList: [],
    }
);

watch(
    () => props.actions,
    (val: ApiAction[]) => {
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
    label: t('actions.name'),
  },
  {
    field: 'script',
    label: t('actions.script'),
  },
  {
    field: 'entity',
    label: t('actions.entity'),
  },
  {
    field: 'action',
    label: t('actions.action'),
  },
  {
    field: 'operations',
    label: t('actions.operations'),
    width: "200px",
    type: 'time',

  },
]

const addNew = () => {
  currentAction = null
  mode.value = Mode.NEW
}

const {emitter} = useEmitt()
const call = (action: ApiAction) => {
  emitter.emit('callAction', action.name)
}

const updateActions = () => {
  const actions = unref(tableObject.tableList)
  emitter.emit('updateActions', actions)
}

const edit = (action: ApiAction, $index) => {
  currentActionIndex.value = $index
  currentAction = unref(action)
  mode.value = Mode.EDIT
}

const save = async () => {
  const write = unref(writeRef)
  const validate = await write?.elFormRef?.validate()?.catch(() => {
  })
  const data = (await write?.getFormData()) as ApiAction
  if (validate) {
    const action = {
      id: data?.id,
      name: data.name,
      entity: data.entity,
      entityId: data.entity?.id || null,
      script: data.script,
      scriptId: data.script?.id || null,
      entityActionName: data.entityActionName,
    } as ApiAction
    if (mode.value === Mode.NEW) {
      tableObject.tableList.push(action)
    } else {
      tableObject.tableList[currentActionIndex.value] = action
    }
    mode.value = Mode.VIEW
    updateActions()
  }
}

const resetForm = () => {
  mode.value = Mode.VIEW
  currentAction = null
}

const removeItem = () => {
  tableObject.tableList.splice(currentActionIndex.value, 1);
  mode.value = Mode.VIEW
}

const goToEntity = (action: ApiAction) => {
  if (!action.entity?.id) {
    return
  }
  push(`/entities/edit/${action.entity?.id}`)
}

const goToScript = (action: ApiAction) => {
  if (!action.script?.id) {
    return
  }
  push(`/scripts/edit/${action.script?.id}`)
}
</script>

<template>
  <ElButton class="flex mb-20px items-left" @click="addNew()" plain>
    <Icon icon="ep:plus" class="mr-5px"/>
    {{ t('actions.addNew') }}
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

    <template #action="{row}">
      <span v-if="row.entityActionName">{{ row.entityActionName }}</span>
      <span v-else>-</span>
    </template>

    <template #operations="{ row }">
      <ElButton :link="true" @click.prevent.stop="call(row)">
        {{ $t('main.call') }}
      </ElButton>

      <ElButton :link="true" @click.prevent.stop="edit(row)">
        {{ $t('main.edit') }}
      </ElButton>

    </template>

  </Table>

  <ActionForm ref="writeRef" :action="currentAction" v-if="mode!=='VIEW'"/>

  <div style="text-align: right" v-if="mode!=='VIEW'">

    <ElButton v-if="mode === 'NEW'" type="primary" @click="save()">
      {{ t('actions.addAction') }}
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
