<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {PropType, reactive, ref, unref, watch} from 'vue'
import {TableColumn} from '@/types/table'
import {ElButton, ElPopconfirm} from 'element-plus'
import {useEmitt} from "@/hooks/web/useEmitt";
import {useRouter} from "vue-router";
import StateForm from "@/views/Entities/components/StateForm.vue";
import {EntityState} from "@/views/Entities/components/types";
import {propTypes} from "@/utils/propTypes";

const {push} = useRouter()
const props = defineProps({
  states: {
    type: Array as PropType<EntityState[]>,
    default: () => []
  },
  pluginStates: {
    type: Array as PropType<EntityState[]>,
    default: () => []
  },
  customStates: propTypes.bool.def(false),
})

const writeRef = ref<ComponentRef<typeof StateForm>>()
let currentState = reactive<Nullable<EntityState>>(null)
let currentStateIndex = ref(0)

enum Mode {
  VIEW = 'VIEW',
  EDIT = 'EDIT',
  NEW = 'NEW'
}

const mode = ref<Mode>(Mode.VIEW)
const {t} = useI18n()

interface TableObject {
  tableList: EntityState[]
}

const tableObject = reactive<TableObject>(
    {
      tableList: [],
    }
);

watch(
    () => props.states,
    (val: EntityState[]) => {
      if (val == unref(tableObject.tableList)) return;
      tableObject.tableList = val
      if (val.length) {
        currentState = val[0]
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
    width: "200px",
  },
  {
    field: 'image',
    label: t('entities.image'),
    width: "110px",
  },
  {
    field: 'description',
    label: t('entities.description'),
  },
  {
    field: 'operations',
    label: t('entities.operations'),
    width: "200px",
    type: 'time',

  },
]

const addNew = () => {
  currentState = null
  mode.value = Mode.NEW
}

const {emitter} = useEmitt()
const setState = (state: EntityState) => {
  emitter.emit('setState', state.name)
}

const updateTriggers = () => {
  const states = unref(tableObject.tableList)
  emitter.emit('updateStates', states)
}

const edit = (state: EntityState, $index) => {
  currentStateIndex.value = $index
  currentState = unref(state)
  mode.value = Mode.EDIT
}

const save = async () => {
  const write = unref(writeRef)
  const validate = await write?.elFormRef?.validate()?.catch(() => {
  })
  const data = (await write?.getFormData()) as EntityState
  if (validate) {
    const state = {
      name: data.name,
      description: data.description,
      imageId: data.image?.id,
      image: data.image,
    } as EntityState
    if (mode.value === Mode.NEW) {
      tableObject.tableList.push(state)
    } else {
      tableObject.tableList[currentStateIndex.value] = state
    }
    mode.value = Mode.VIEW
    updateTriggers()
  }
}

const resetForm = () => {
  mode.value = Mode.VIEW
  currentState = null
}

const removeItem = () => {
  tableObject.tableList.splice(currentStateIndex.value, 1);
  mode.value = Mode.VIEW
}

const goToScript = (state: EntityState) => {
  if (!state.script?.id) {
    return
  }
  push(`/scripts/edit/${state.script?.id}`)
}

const loadFromPlugin = async () => {
  emitter.emit('updateStates', unref(props.pluginStates))
}

</script>

<template>
  <ElButton class="flex mb-20px items-left" @click="addNew()" plain v-if="mode ==='VIEW' && customStates">
    <Icon icon="ep:plus" class="mr-5px"/>
    {{ t('entities.addNewState') }}
  </ElButton>

  <ElButton class="flex mb-20px items-left" @click="loadFromPlugin()" plain
            v-if="mode ==='VIEW' && pluginStates.length">
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

    <template #operations="{ row, $index }">

      <ElButton :link="true" @click.prevent.stop="setState(row)">
        {{ $t('entities.setState') }}
      </ElButton>

      <ElButton :link="true" @click.prevent.stop="edit(row, $index)">
        {{ $t('main.edit') }}
      </ElButton>

    </template>

  </Table>

  <StateForm ref="writeRef" :state="currentState" v-if="mode!=='VIEW'"/>

  <div style="text-align: right" v-if="mode!=='VIEW'">

    <ElButton v-if="mode === 'NEW'" type="primary" @click="save()">
      {{ t('entities.addNewState') }}
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
