import type { Component } from 'vue'
import {
  ElCascader,
  ElCheckboxGroup,
  ElColorPicker,
  ElDatePicker,
  ElInput,
  ElInputNumber,
  ElRadioGroup,
  ElRate,
  ElSelect,
  ElSelectV2,
  ElSlider,
  ElSwitch,
  ElTimePicker,
  ElTimeSelect,
  ElTransfer,
  ElAutocomplete,
  ElDivider
} from 'element-plus'
import { InputPassword } from '@/components/InputPassword'
import { Editor } from '@/components/Editor'
import { ComponentName } from '@/types/components'
import ImageSearch from '@/views/Images/components/ImageSearch.vue'
import AreaSearch from '@/views/Areas/components/AreaSearch.vue'
import RoleSearch from '@/views/Users/components/RoleSearch.vue'
import ScriptSearch from '@/views/Scripts/components/ScriptSearch.vue'
import TagsSearch from '@/views/Tags/components/TagsSearch.vue'
import ScriptFormHelper from '@/views/Scripts/components/ScriptFormHelper.vue'
import ScriptsSearch from '@/views/Scripts/components/ScriptsSearch.vue'
import EntitySearch from '@/views/Entities/components/EntitySearch.vue'
import EntitiesSearch from '@/views/Entities/components/EntitiesSearch.vue'
import PluginSearch from '@/views/Plugins/components/PluginSearch.vue'
import TriggerSearch from '@/views/Automation/components/TriggerSearch.vue'
import TriggersSearch from '@/views/Automation/components/TriggersSearch.vue'
import ConditionSearch from '@/views/Automation/components/ConditionSearch.vue'
import ConditionsSearch from '@/views/Automation/components/ConditionsSearch.vue'
import ActionSearch from '@/views/Automation/components/ActionSearch.vue'
import ActionsSearch from '@/views/Automation/components/ActionsSearch.vue'

const componentMap: Recordable<Component, ComponentName> = {
  Radio: ElRadioGroup,
  Checkbox: ElCheckboxGroup,
  CheckboxButton: ElCheckboxGroup,
  Input: ElInput,
  Autocomplete: ElAutocomplete,
  InputNumber: ElInputNumber,
  Select: ElSelect,
  Cascader: ElCascader,
  Switch: ElSwitch,
  Slider: ElSlider,
  TimePicker: ElTimePicker,
  DatePicker: ElDatePicker,
  Rate: ElRate,
  ColorPicker: ElColorPicker,
  Transfer: ElTransfer,
  Divider: ElDivider,
  TimeSelect: ElTimeSelect,
  SelectV2: ElSelectV2,
  RadioButton: ElRadioGroup,
  InputPassword: InputPassword,
  Editor: Editor,
  Image: ImageSearch,
  Area: AreaSearch,
  Role: RoleSearch,
  Script: ScriptSearch,
  Scripts: ScriptsSearch,
  Tags: TagsSearch,
  ScriptHelper: ScriptFormHelper,
  Entity: EntitySearch,
  Entities: EntitiesSearch,
  Plugin: PluginSearch,
  Trigger: TriggerSearch,
  Triggers: TriggersSearch,
  Condition: ConditionSearch,
  Conditions: ConditionsSearch,
  Action: ActionSearch,
  Actions: ActionsSearch,
}

export { componentMap }
