<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {PropType, reactive, ref, unref, watch} from 'vue'
import {TableColumn} from '@/types/table'
import {ElButton, ElPopconfirm} from 'element-plus'
import {useEmitt} from "@/hooks/web/useEmitt";
import {useRouter} from "vue-router";
import ActionForm from "@/views/Entities/components/ActionForm.vue";
import {EntityAction} from "@/views/Entities/components/types";
import {propTypes} from "@/utils/propTypes";

const {currentRoute, addRoute, push} = useRouter()
const props = defineProps({
  actions: {
    type: Array as PropType<EntityAction[]>,
    default: () => []
  },
  pluginActions: {
    type: Array as PropType<EntityAction[]>,
    default: () => []
  },
  customActions: propTypes.bool.def(false),
})

const writeRef = ref<ComponentRef<typeof ActionForm>>()
let currentAction = reactive<Nullable<EntityAction>>(null)
let currentActionIndex = ref(0)

enum Mode {
  VIEW = 'VIEW',
  EDIT = 'EDIT',
  NEW = 'NEW'
}

const mode = ref<Mode>(Mode.VIEW)
const {t} = useI18n()

interface TableObject {
  tableList: EntityAction[]
}

const tableObject = reactive<TableObject>(
    {
      tableList: [],
    }
);

watch(
    () => props.actions,
    (val: EntityAction[]) => {
      if (val == unref(tableObject.tableList)) return;
      tableObject.tableList = val
      if (val.length) {
        currentAction = val[0]
      }
    },
    {
      immediate: true
    }
)

const columns: TableColumn[] = [
  {
    field: 'name',
    label: t('entities.name'),
  },
  {
    field: 'image',
    label: t('entities.image'),
    width: "110px",
  },
  {
    field: 'icon',
    label: t('entities.icon'),
    width: "80px",
  },
  {
    field: 'script',
    label: t('entities.script'),
    width: "170px",
  },
  {
    field: 'operations',
    label: t('entities.operations'),
    width: "200px",
    type: 'time',

  },
]

const addNew = () => {
  currentAction = null
  mode.value = Mode.NEW
}

const {emitter} = useEmitt()
const call = (action: EntityAction) => {
  emitter.emit('callAction', action.name)
}

const edit = (actions: EntityAction, $index) => {
  currentActionIndex.value = $index
  currentAction = unref(actions)
  mode.value = Mode.EDIT
}

const save = async () => {
  const write = unref(writeRef)
  const validate = await write?.elFormRef?.validate()?.catch(() => {
  })
  const data = (await write?.getFormData()) as EntityAction
  if (validate) {
    const action = {
      name: data.name,
      description: data.description,
      icon: data.icon,
      imageId: data.image?.id,
      image: data.image,
      scriptId: data.script?.id,
      script: data.script,
    } as EntityAction
    if (mode.value === Mode.NEW) {
      tableObject.tableList.push(action)
    } else {
      tableObject.tableList[currentActionIndex.value] = action
    }
    mode.value = Mode.VIEW
    emitter.emit('updateActions', unref(tableObject.tableList))
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

const goToScript = (action: EntityAction) => {
  if (!action.script?.id) {
    return
  }
  push(`/scripts/edit/${action.script?.id}`)
}

const loadFromPlugin = async () => {
  emitter.emit('updateActions', unref(props.pluginActions))
}

</script>

<template>
  <ElButton class="flex mb-20px items-left" @click="addNew()" plain v-if="mode ==='VIEW' && customActions">
    <Icon icon="ep:plus" class="mr-5px"/>
    {{ t('entities.addNewAction') }}
  </ElButton>

  <ElButton class="flex mb-20px items-left" @click="loadFromPlugin()" plain
            v-if="mode ==='VIEW' && pluginActions.length">
    {{ t('entities.loadFromPlugin') }}
  </ElButton>

  <Table
      v-if="mode==='VIEW'"
      :selection="false"
      :columns="columns"
      :data="tableObject.tableList"
      style="width: 100%"
  >

    <template #image="{row}">
      <span v-if="row.image"><Icon icon="ic:sharp-check" class="mr-5px"/></span>
      <span v-else>-</span>
    </template>

    <template #icon="{row}">
      <span v-if="row.icon"><Icon icon="ic:sharp-check" class="mr-5px"/></span>
      <span v-else>-</span>
    </template>

    <template #script="{row}">
      <ElButton v-if="row.script" :link="true" @click.prevent.stop="goToScript(row)">
        {{ row.script.name }}
      </ElButton>
      <span v-else>-</span>
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

  <ActionForm ref="writeRef" :actions="currentAction" v-if="mode!=='VIEW'"/>

  <div style="text-align: right" v-if="mode!=='VIEW'">

    <ElButton v-if="mode === 'NEW'" type="primary" @click="save()">
      {{ t('entities.addNewAction') }}
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
