<script setup lang="ts">
import {computed, PropType, ref, watch} from 'vue'
import {Trigger} from "@/views/Automation/components/types";
import {ApiArea, ApiPlugin, ApiScript} from "@/api/stub";
import api from "@/api/api";
import {
  ElCol,
  ElDivider,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElOption,
  ElRow,
  ElSelect,
  ElSwitch
} from "element-plus";
import AreaSearch from "@/components/AreaSearch/src/AreaSearch.vue";
import ScriptSearch from "@/components/ScriptSearch/src/ScriptSearch.vue";
import EntitiesSearch from "@/components/EntitiesSearch/src/EntitiesSearch.vue";
import ScriptFormHelper from "@/components/ScriptFormHelper/src/ScriptFormHelper.vue";
import {GetApiAttributeValue} from "@/components/Attributes";
import {useCache} from "@/hooks/web/useCache";
import {useAppStore} from "@/store/modules/app";

const props = defineProps({
  trigger: {
    type: Object as PropType<Nullable<Trigger>>,
    default: () => null
  }
})

const appStore = useAppStore()
const className = computed(() => appStore.getIsDark ? 'dark' : 'light')
const currentTrigger = computed(() => props.trigger as Trigger)
const plugin = ref<ApiPlugin>(null)
const pluginList = ref<ApiPlugin[]>([])

const rules = ref({
  name: [{required: true}],
  pluginName: [{required: true}],
});

const validator = (val: unknown): boolean => {
  for (const attr of currentTrigger.value.attributes) {
    if (attr?.name === val?.field) {
      const currentValue = GetApiAttributeValue(attr)
      return !!currentValue
    }
  }
  return false
}

const getPlugin = async (name: string) => {
  const res = await api.v1.pluginServiceGetPlugin(name)
    .catch(() => {
    })
    .finally(() => {

    })
  if (res) {
    plugin.value = res.data as ApiPlugin;

    // update required
    let _rules = {
      name: [{required: true}],
      pluginName: [{required: true}],
    };
    if (res.data.options?.triggerParams?.required) {
      for (const item of res.data.options.triggerParams.required) {
        _rules[item] = [{
          required: true,
          validator: validator,
          // type: item?.type?.toLowerCase()
        }]
      }
    }
    rules.value = _rules

    // attributes
    let attributes = []
    if (res.data.options?.triggerParams?.attributes) {
      for (const key in res.data.options.triggerParams.attributes) {
        let attr = res.data.options.triggerParams.attributes[key]

        if (attr.type == 'notice') {
          getReadme(plugin.value.name, attr.name)
        }

        // restore value from server
        if (currentTrigger?.value?.attributes) {
          for (const key2 in currentTrigger.value.attributes) {
            if (key2 == key) {
              attr = currentTrigger.value.attributes[key2]
            }
          }
        } // \restore value from server

        attributes.push(attr);
      }
    }
    currentTrigger.value.attributes = attributes;
  }
}

const getPluginList = async () => {

  let params = {
    sort: '-name',
    page: 1,
    limit: 99,
    triggers: true,
  }
  const res = await api.v1.pluginServiceGetPluginList(params)
    .catch(() => {
    })
    .finally(() => {

    })
  if (res) {
    const {items} = res.data;
    pluginList.value = items as ApiPlugin[];
  }
}

const changedArea = async (area: ApiArea) => {
  // console.log(area)
  currentTrigger.value.areaId = area?.id
  currentTrigger.value.area = area
}

const changedScript = async (script: ApiScript) => {
  // console.log(script)
  currentTrigger.value.scriptId = script?.id
  currentTrigger.value.script = script
}

const changedPlugin = async (name: string) => {
  // console.log(name)
  if (!name) {
    plugin.value = undefined
    return
  }
  getPlugin(name)
}

getPluginList()

watch(
  () => props.trigger,
  (val) => {
    if (!val) return
    if (!val.pluginName) {
      plugin.value = undefined
      return;
    }
    if (val.pluginName == plugin.value?.name) {
      return;
    }
    getPlugin(val.pluginName)
  },
  {
    deep: true,
    immediate: true
  }
)

const readme = ref('');
const notes = ref({})
const {wsCache} = useCache()
const getReadme = async (pluginName: string, note: string) => {
  const lang = wsCache.get('lang') || 'en';
  const res = await api.v1.pluginServiceGetPluginReadme(pluginName, {note: note, lang: lang})
    .catch(() => {
    })
    .finally(() => {

    })
  if (res && res.data) {
    notes.value[note] = res.data
  }
}

const formRef = ref(null)
defineExpose({
  form: formRef
})

</script>

<template>

  <ElForm
    v-if="currentTrigger"
    label-position="top"
    :model="currentTrigger"
    style="width: 100%"
    :rules="rules"
    ref="formRef"
  >

    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('automation.triggers.name')" prop="name">
          <ElInput placeholder="Please input name" v-model="currentTrigger.name" clearable/>
        </ElFormItem>
      </ElCol>
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('automation.description')" prop="description">
          <ElInput placeholder="Please input description" v-model="currentTrigger.description" clearable/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('automation.enabled')" prop="enabled">
          <ElSwitch v-model="currentTrigger.enabled"/>
        </ElFormItem>
      </ElCol>
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('automation.area')" prop="area">
          <AreaSearch v-model="currentTrigger.area" @change="changedArea($event)"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('automation.triggers.pluginName')" prop="pluginName">
          <ElSelect v-model="currentTrigger.pluginName" placeholder="Please select plugin" clearable class="w-[100%]"
                    @change="changedPlugin($event)">
            <ElOption
              v-for="(prop, index) in pluginList"
              :key="index"
              :label="prop.name"
              :value="prop.name"/>
          </ElSelect>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow :gutter="24" v-if="plugin">
      <ElCol :span="12" :xs="12">
        <ElDivider content-position="left">{{ $t('automation.triggers.pluginOptions') }}</ElDivider>
      </ElCol>
    </ElRow>

    <ElRow :gutter="24" v-if="plugin?.options?.triggerParams?.script">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('automation.triggers.script')" prop="script">
          <ScriptSearch v-model="currentTrigger.script" @change="changedScript($event)"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow :gutter="24" v-if="plugin?.options?.triggerParams?.script && currentTrigger?.script">
      <ElCol :span="12" :xs="12">
        <ElFormItem label="" prop="help">
          <ScriptFormHelper v-model="currentTrigger.script"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow :gutter="24" v-if="plugin?.options?.triggerParams?.entities">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.entityStorage.entities')" prop="entityIds">
          <EntitiesSearch v-model="currentTrigger.entityIds"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <div v-if="currentTrigger?.attributes">
      <ElRow :gutter="24" v-for="(row, index) in currentTrigger.attributes"
             :name="index"
             :key="index">
        <ElCol :span="12" :xs="12">

          <ElFormItem :label="$t(`automation.triggers.${currentTrigger.pluginName}.${row.name}`)" :prop="row.name"
                      v-if="row.type !== 'notice'">

            <div v-if="row.type === 'STRING'" class="w-[100%]">
              <ElInput type="string" v-model="row.string" clearable/>
            </div>
            <div v-if="row.type === 'IMAGE'" class="w-[100%]">
              <ElInput type="string" v-model="row.imageUrl" clearable/>
            </div>
            <div v-if="row.type === 'ICON'" class="w-[100%]">
              <ElInput type="string" v-model="row.icon" clearable/>
            </div>
            <div v-if="row.type === 'INT'" class="w-[100%]">
              <ElInputNumber v-model="row.int" class="w-[100%]" clearable/>
            </div>
            <div v-if="row.type === 'FLOAT'">
              <ElInputNumber v-model="row.float" clearable/>
            </div>

            <div v-if="row.type === 'POINT'" class="w-[100%]">
              <ElInput type="string" v-model="row.point" clearable/>
            </div>
            <div v-if="row.type === 'ENCRYPTED'">
              <ElInput type="password" v-model="row.encrypted" show-password clearable/>
            </div>
            <ElSelect
              v-model="row.bool"
              placeholder="please select value"
              v-if="row.type === 'BOOL'"
              class="w-[100%]" clearable>
              <ElOption label="TRUE" :value="true"/>
              <ElOption label="FALSE" :value="false"/>
            </ElSelect>

            <div v-if="row.type === 'TIME'" class="w-[100%]">
              <ElInput type="string" v-model="row.time" clearable/>
            </div>

          </ElFormItem>

        </ElCol>
      </ElRow>

      <div v-for="(row, index) in currentTrigger.attributes"
           :name="index"
           :key="index">
        <div class="mt-10px mb-25px markdown" :class="className" v-if="row.type === 'notice'"
             v-html="notes[row.name]"></div>

      </div>
    </div>


  </ElForm>

</template>

<style lang="less">
.markdown {
  font-size: 0.8rem;

  &.light {
    padding: 8px 16px;
    background-color: #409eff1a;
    border-radius: 4px;
    border-left: 5px solid var(--el-color-primary);
  }

  &.dark {
    padding: 8px 16px;
    background-color: #409eff1a;
    border-radius: 4px;
    border-left: 5px solid var(--el-color-primary);
  }

  p {
    margin: 10px 0 !important;
  }

  table {
    thead th {
      border: 1px solid ;
      padding: 0 10px;
    }

    tbody td {
      border: 1px solid ;
      padding: 0 10px;
    }
  }
}
</style>
