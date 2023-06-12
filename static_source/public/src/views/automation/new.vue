<template>
  <div class="app-container">
    <card-wrapper>
      <el-row :gutter="20">
        <el-col
          :span="24"
          :xs="24"
        >

          <el-tabs v-model="internal.activeTab">

            <el-tab-pane
              label="Main"
              name="main"
            >
              <el-form label-position="top"
                       ref="currentTask"
                       :model="currentTask"
                       :rules="rules"
                       style="width: 100%">
                <el-form-item :label="$t('automation.table.name')" prop="name">
                  <el-input v-model.trim="currentTask.name"/>
                </el-form-item>
                <el-form-item :label="$t('automation.table.description')" prop="description">
                  <el-input v-model.trim="currentTask.description"/>
                </el-form-item>
                <el-form-item :label="$t('automation.table.enabled')" prop="autoLoad">
                  <el-switch v-model="currentTask.enabled"></el-switch>
                </el-form-item>
                <el-form-item :label="$t('automation.table.condition')" prop="icon">
                  <el-select v-model="currentTask.condition" placeholder="please select type">
                    <el-option label="OR" value="or"></el-option>
                    <el-option label="AND" value="and"></el-option>
                  </el-select>
                </el-form-item>

                <el-form-item :label="$t('automation.table.area')" prop="area">
                  <area-search
                    :multiple=false
                    v-model="currentTask.area"
                    @update-value="changedArea"
                  />
                </el-form-item>

              </el-form>
            </el-tab-pane>

            <el-tab-pane
              label="Triggers"
              name="triggers"
            >
              <triggers
                v-model="currentTask.triggers"
                @update-value="changedTriggers"
              />
            </el-tab-pane>

            <el-tab-pane
              label="Conditions"
              name="conditions"
            >
              <conditions
                v-model="currentTask.conditions"
                @update-value="changedConditions"
              />
            </el-tab-pane>

            <el-tab-pane
              label="Actions"
              name="actions"
            >
              <actions
                v-model="currentTask.actions"
                @update-value="changedActions"
              />
            </el-tab-pane>

          </el-tabs>

        </el-col>
      </el-row>
      <el-row style="margin-top: 20px">
        <el-col :span="24" align="right">
          <el-button type="primary" @click.prevent.stop="save">{{ $t('main.save') }}</el-button>
          <el-button @click.prevent.stop="cancel">{{ $t('main.cancel') }}</el-button>
        </el-col>
      </el-row>
    </card-wrapper>
  </div>
</template>

<script lang="ts">

import { Component, Vue } from 'vue-property-decorator'
import api from '@/api/api'
import { ApiAction, ApiArea, ApiCondition, ApiNewTaskRequest, ApiTrigger } from '@/api/stub'
import router from '@/router'
import ScriptSearch from '@/views/scripts/components/script_search.vue'
import AreaSearch from '@/views/areas/components/areas_search.vue'
import { Form } from 'element-ui'
import Triggers from '@/views/automation/components/triggers.vue'
import Conditions from '@/views/automation/components/conditions.vue'
import Actions from '@/views/automation/components/actions.vue'
import CardWrapper from '@/components/card-wrapper/index.vue'

@Component({
  name: 'Editor',
  components: {
    CardWrapper,
    Triggers,
    ScriptSearch,
    AreaSearch,
    Conditions,
    Actions
  }
})
export default class extends Vue {
  private internal = {
    activeTab: 'main',
    pluginOptions: undefined
  };

  // entity params
  private currentTask: ApiNewTaskRequest = {
    name: '',
    enabled: true,
    condition: 'and',
    triggers: [],
    conditions: [],
    actions: []
  };

  private rules = {
    name: [
      { required: true, trigger: 'blur' },
      { min: 4, max: 255, trigger: 'blur' }
    ],
    description: [
      { required: false, trigger: 'blur' },
      { max: 255, trigger: 'blur' }
    ],
    plugin: [
      { required: false, trigger: 'blur' },
      { max: 255, trigger: 'blur' }
    ]
  };

  private changedArea(values: ApiArea, event?: any) {
    if (values) {
      this.$set(this.currentTask, 'area', values)
    } else {
      this.$set(this.currentTask, 'area', undefined)
    }
  }

  private changedTriggers(values: ApiTrigger[], event?: any) {
    if (values) {
      this.$set(this.currentTask, 'triggers', values)
    } else {
      this.$set(this.currentTask, 'triggers', undefined)
    }
  }

  private changedConditions(values: ApiCondition[], event?: any) {
    if (values) {
      this.$set(this.currentTask, 'conditions', values)
    } else {
      this.$set(this.currentTask, 'conditions', undefined)
    }
  }

  private changedActions(values: ApiAction[], event?: any) {
    if (values) {
      this.$set(this.currentTask, 'action', values)
    } else {
      this.$set(this.currentTask, 'action', undefined)
    }
  }

  private async save() {
    (this.$refs.currentTask as Form).validate(async valid => {
      if (!valid) {
        return
      }
      const task = {
        name: this.currentTask.name,
        description: this.currentTask.description,
        enabled: this.currentTask.enabled,
        condition: this.currentTask.condition,
        triggers: this.currentTask.triggers,
        conditions: this.currentTask.conditions,
        actions: this.currentTask.actions,
        area: this.currentTask.area
      }
      const { data } = await api.v1.automationServiceAddTask(task)
      if (data) {
        this.$notify({
          title: 'Success',
          message: 'task created successfully',
          type: 'success',
          duration: 2000
        })
        router.push({ path: `/automation/edit/${data.id}` })
      }
    })
  }

  private cancel() {
    router.go(-1)
  }
}
</script>
