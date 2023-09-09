<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {computed, PropType, reactive, ref, unref, watch} from 'vue'
import {ElButton, ElInput, ElSelect, ElOption} from 'element-plus'
import {Attribute, EntityAction, EntityAttribute} from "@/views/Entities/components/types";
import {TableColumn} from "@/types/table";
import {propTypes} from "@/utils/propTypes";

const {t} = useI18n()

const props = defineProps({
  modelValue: {
    type: Array as PropType<Attribute[]>,
    default: () => []
  },
  pluginAttrs: {
    type: Array as PropType<Attribute[]>,
    default: () => []
  },
  customAttrs: propTypes.bool.def(false),
})

const attributes = computed(() => props.modelValue)

const columns: TableColumn[] = [
  {
    field: 'name',
    label: t('entities.name'),
  },
  {
    field: 'type',
    label: t('entities.type'),
  },
  {
    field: 'value',
    label: t('entities.value'),
  },
  {
    field: 'operations',
    label: t('entities.operations'),
    width: "150px",
    type: 'time',

  },
]

const addNew = () => {
  attributes.value.push(new EntityAttribute('new_value'))
}

const remove = (attr: Attribute, index: number) => {
  attributes.value.splice(index, 1)
}

const emit = defineEmits(['change', 'update:modelValue'])
const loadFromPlugin = async () => {
  emit('change', unref(props.pluginAttrs))
}

</script>

<template>
  <ElButton class="flex mb-20px items-left"  @click="addNew()" plain v-if="customAttrs">
    <Icon icon="ep:plus" class="mr-5px"/>
    {{ t('entities.addNewAttr') }}
  </ElButton>

  <ElButton class="flex mb-20px items-left"  @click="loadFromPlugin()" plain v-if="pluginAttrs.length">
    {{ t('entities.loadFromPlugin') }}
  </ElButton>

  <Table
      :selection="false"
      :columns="columns"
      :data="attributes"
      style="width: 100%"
  >

    <template #name="{row}">
      <ElInput v-model="row.name" placeholder="Please input" />
    </template>

    <template #type="{row}">
      <ElSelect v-model="row.type" placeholder="please select type" class="w-[100%]">
        <ElOption label="String" value="STRING"/>
        <ElOption label="Int" value="INT"/>
        <ElOption label="Float" value="FLOAT"/>
        <ElOption label="Bool" value="BOOL"/>
        <ElOption label="Array" value="ARRAY"/>
        <ElOption label="Time" value="TIME"/>
        <ElOption label="Map" value="MAP"/>
        <ElOption label="Image" value="IMAGE"/>
        <ElOption label="Point" value="POINT"/>
      </ElSelect>
    </template>

    <template #value="{row}">
      <div v-if="row.type === 'STRING'">
        <ElInput type="string" v-model="row.string"/>
      </div>
      <div v-if="row.type === 'IMAGE'">
        <ElInput type="string" v-model="row.imageUrl"/>
      </div>
      <div v-if="row.type === 'INT'">
        <ElInput type="number" v-model="row.int"/>
      </div>
      <div v-if="row.type === 'FLOAT'">
        <ElInput type="number" v-model="row.float"/>
      </div>
      <div v-if="row.type === 'ARRAY'">
        <ElInput type="string" v-model="row.array"/>
      </div>
      <div v-if="row.type === 'POINT'">
        <ElInput type="string" v-model="row.point"/>
      </div>
      <ElSelect
          v-model="row.bool"
          placeholder="please select value"
          v-if="row.type === 'BOOL'"
          class="w-[100%]">
        <ElOption label="TRUE" :value="true"/>
        <ElOption label="FALSE" :value="false"/>
      </ElSelect>

      <div v-if="row.type === 'TIME'">
        <ElInput type="string" v-model="row.time"/>
      </div>

      <div v-if="row.type === 'MAP'">
        <ElInput type="string" v-model="row.map"/>
      </div>
    </template>

    <template #operations="{ row, $index }">

      <ElButton :link="true" @click.prevent.stop="remove(row, $index)">
        {{ $t('main.remove') }}
      </ElButton>

    </template>

  </Table>



</template>

<style lang="less">

</style>
