import type {Component} from 'vue'
import {
    ElAutocomplete,
    ElCascader,
    ElCheckboxGroup,
    ElColorPicker,
    ElDatePicker,
    ElDivider,
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
    ElTransfer
} from 'element-plus'
import {InputPassword} from '@/components/InputPassword'
import {Editor} from '@/components/Editor'
import {ComponentName} from '@/types/components'
import {EntitySearch} from '@/components/EntitySearch'
import {EntitiesSearch} from '@/components/EntitiesSearch'
import {PluginSearch} from '@/components/PluginSearch'
import {TriggerSearch} from '@/components/TriggerSearch'
import {TriggersSearch} from '@/components/TriggersSearch'
import {ConditionSearch} from '@/components/ConditionSearch'
import {ConditionsSearch} from '@/components/ConditionsSearch'
import {ActionSearch} from '@/components/ActionSearch'
import {ActionsSearch} from '@/components/ActionsSearch'
import {ScriptsSearch} from '@/components/ScriptsSearch'
import {ScriptFormHelper} from '@/components/ScriptFormHelper'
import {TagsSearch} from '@/components/TagsSearch'
import {ScriptSearch} from '@/components/ScriptSearch'
import {RoleSearch} from '@/components/RoleSearch'
import {AreaSearch} from "@/components/AreaSearch";
import {ImageSearch} from "@/components/ImageSearch";

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

export {componentMap}
