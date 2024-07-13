import type {Component} from 'vue'
import {
    ElAutocomplete,
    ElCascader,
    ElCheckboxGroup,
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
import {SystemEventsHelper} from '@/components/SystemEventsHelper'
import {CronFormHelper} from '@/components/CronFormHelper'
import {TagsSearch} from '@/components/TagsSearch'
import {ScriptSearch} from '@/components/ScriptSearch'
import {RoleSearch} from '@/components/RoleSearch'
import {AreaSearch} from "@/components/AreaSearch";
import {ImageSearch} from "@/components/ImageSearch";
import {ColorPicker} from "@/components/ColorPicker";
import {VariableSearch} from "@/components/VariableSearch";
import {VariablesSearch} from "@/components/VariablesSearch";

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
    ColorPicker: ColorPicker,
    Transfer: ElTransfer,
    Divider: ElDivider,
    TimeSelect: ElTimeSelect,
    SelectV2: ElSelectV2,
    RadioButton: ElRadioGroup,
    InputPassword: InputPassword,
    Image: ImageSearch,
    Area: AreaSearch,
    Role: RoleSearch,
    Script: ScriptSearch,
    Scripts: ScriptsSearch,
    Tags: TagsSearch,
    ScriptHelper: ScriptFormHelper,
    SystemEventsHelper: SystemEventsHelper,
    CronHelper: CronFormHelper,
    Entity: EntitySearch,
    Entities: EntitiesSearch,
    Plugin: PluginSearch,
    Trigger: TriggerSearch,
    Triggers: TriggersSearch,
    Condition: ConditionSearch,
    Conditions: ConditionsSearch,
    Action: ActionSearch,
    Actions: ActionsSearch,
    Variable: VariableSearch,
    Variables: VariablesSearch,
}

export {componentMap}
