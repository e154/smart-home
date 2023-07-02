<template>
  <div>

    <div v-if="mode != 'VIEW'">
      <el-form label-position="top" label-width="100px"
               ref="currentItem"
               :model="currentItem"
               :rules="rules"
               style="width: 100%">
        <el-form-item :label="$t('entities.table.name')" prop="name">
          <el-input v-model="currentItem.name"></el-input>
        </el-form-item>

        <el-form-item :label="$t('automation.table.script')" prop="script">

          <span slot="label" v-if="currentItem.script && currentItem.script.id">
            {{ $t('entities.table.script') }}
            <script-dialog
              :visible.sync="dialogScriptVisible"
              :script="currentItem.script"
              @on-close="dialogScriptVisible=false"
            />
            <el-button
              type="text"
              @click="dialogScriptVisible=true">
             {{ $t('scripts.view') }}   <svg-icon name="link"/>
            </el-button>
          </span>

          <script-search
            :multiple="false"
            v-model="currentItem.script"
            @update-value="changedScript"
          />
        </el-form-item>

        <el-form-item :label="$t('dashboard.editor.entity')" prop="entity">
          <entity-search
            v-model="currentItem.entity"
            @update-value="changedEntity"
          />
        </el-form-item>

        <el-select
          v-if="currentItem.entity"
          v-model="currentItem.entityActionName"
          clearable
          :placeholder="$t('dashboard.editor.selectAction')"
          style="width: 100%"
        >
          <el-option
            v-for="item in buffer.entityActions"
            :key="item.value"
            :label="item.label"
            :value="item.value">
          </el-option>
        </el-select>


        <el-form-item style="margin-top: 20px">
          <el-button v-if="mode == 'NEW'" type="primary" @click="submitForm()">{{
              $t('entities.addAction')
            }}
          </el-button>
          <el-button v-if="mode == 'EDIT'" type="primary" @click="submitForm()">{{ $t('main.update') }}</el-button>
          <el-button @click="resetForm()">{{ $t('main.cancel') }}</el-button>
          <el-button v-if="mode == 'EDIT'" type="danger" plain @click="removeItem()">{{ $t('main.remove') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div v-if="mode == 'VIEW'">
      <el-row>
        <el-col>
          <el-button
            @click='add()'>
            <i class="el-icon-plus"/> {{ $t('automation.addAction') }}
          </el-button>
        </el-col>
      </el-row>

      <el-row>
        <el-col>
          <el-table
            key="key"
            :data="value"
            style="width: 100%;"
          >
            <el-table-column
              :label="$t('automation.table.name')"
              prop="name"
              align="left"
              width="200px"
            >
              <template slot-scope="{row}">
                <div>{{ row.name }}</div>
              </template>
            </el-table-column>

            <el-table-column
              :label="$t('automation.table.script')"
              prop="script"
              align="left"
              width="auto"
            >
              <template slot-scope="{row}">
                <span v-if="row.script && row.script.name">{{ row.script.name }}</span>
                <span v-else>-</span>
              </template>
            </el-table-column>

            <el-table-column
              :label="$t('automation.table.entity')"
              prop="entity"
              align="left"
              width="auto"
            >
              <template slot-scope="{row}">
                <span v-if="row.entityId">{{ row.entityId }}</span>
                <span v-else>-</span>
              </template>
            </el-table-column>

            <el-table-column
              :label="$t('automation.table.action')"
              prop="action"
              align="left"
              width="auto"
            >
              <template slot-scope="{row}">
                <span v-if="row.entityActionName">{{ row.entityActionName }}</span>
                <span v-else>-</span>
              </template>
            </el-table-column>

            <el-table-column
              :label="$t('entities.table.operations')"
              prop="description"
              align="left"
              width="200px"
            >
              <template slot-scope="{row, $index}">
                <el-button type="text" size="small" @click='callAction(row, $index)'>{{ $t('main.call') }}</el-button>
                <el-button type="text" size="small" @click='editAction(row, $index)'>{{ $t('main.edit') }}</el-button>
              </template>
            </el-table-column>

          </el-table>

        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator'
import {ApiAction, ApiEntity, ApiEntityAction, ApiScript} from '@/api/stub'
import {Form} from 'element-ui'
import ScriptSearch from '@/views/scripts/components/script_search.vue'
import ScriptEditModal from '@/views/scripts/edit-modal.vue'
import ScriptDialog from '@/views/scripts/dialog.vue'
import EntitySearch from "@/views/entities/components/entity_search.vue";
import api from "@/api/api";

export interface Action {
  value: string;
  label: string;
}

export enum Mode {
  VIEW = 'VIEW',
  EDIT = 'EDIT',
  NEW = 'NEW'
}

@Component({
  name: 'Actions',
  components: {
    EntitySearch,
    ScriptSearch,
    ScriptEditModal,
    ScriptDialog
  }
})
export default class extends Vue {
  @Prop({required: false, default: () => []}) private value?: ApiAction[];

  private mode: Mode = Mode.VIEW;
  private currentItem: ApiAction = {
    entityActionName: undefined
  };
  private currentItemIndex?: number;
  private dialogScriptVisible = false;
  private buffer: {
    entity?: ApiEntity
    entityActions: Action[];
  } = {
    entity: undefined,
    entityActions: []
  }

  private rules = {
    name: [
      {required: true, trigger: 'blur'},
      {min: 4, max: 255, trigger: 'blur'}
    ],
    script: [
      {required: false, trigger: 'blur'}
    ],
    entity: [
      {required: false, trigger: 'blur'}
    ]
  };

  private editAction(action: ApiAction, index: number) {
    this.currentItem = Object.assign({}, action)
    this.currentItemIndex = index
    this.mode = Mode.EDIT
    if (this.currentItem.entityId) {
      this.fetchEntity(this.currentItem.entityId)
    }
  }

  private callAction(action: ApiAction, index: number) {
    this.$emit('call-action', action.name)
  }

  private add() {
    this.currentItem = {}
    this.mode = Mode.NEW
  }

  private resetForm() {
    this.currentItem = {}
    this.mode = Mode.VIEW
    this.currentItemIndex = undefined
  }

  private removeItem() {
    if (this.value) {
      for (const index in this.value) {
        if (this.currentItem && this.value[index].name == this.currentItem.name) {
          this.value.splice(+index, 1)
        }
      }
    }
    this.mode = Mode.VIEW
    this.currentItem = {}
  }

  private changedScript(value: ApiScript, event?: any) {
    if (value) {
      this.$set(this.currentItem, 'script', value)
      this.$set(this.currentItem, 'scriptId', value.id)
    } else {
      this.$set(this.currentItem, 'script', undefined)
      this.$set(this.currentItem, 'scriptId', undefined)
    }
  }

  private changedEntity(value: ApiEntity, event?: any) {
    if (value && value.id) {
      this.$set(this.currentItem, 'entity', value)
      this.$set(this.currentItem, 'entityId', value.id)
      this.fetchEntity(value.id)
    } else {
      this.$set(this.currentItem, 'entity', undefined)
      this.$set(this.currentItem, 'entityId', undefined)
      this.$set(this.currentItem, 'entityActionName', undefined)
    };
  }

  private submitForm() {
    (this.$refs.currentItem as Form).validate(valid => {
      if (!valid) {
        return
      }

      if (this.mode === Mode.NEW) {
        if (this.value) {
          this.value.push(Object.assign({}, this.currentItem))
        }
      } else if (this.mode === Mode.EDIT) {
        if (this.value) {
          if (this.currentItemIndex != undefined) {
            this.value[this.currentItemIndex] = Object.assign({}, this.currentItem)
          }
        }
      }
      this.resetForm()
    })
  }

  private async fetchEntity(id: string) {
    const {data} = await api.v1.entityServiceGetEntity(id);
    this.currentItem.entity = data;
    this.currentItem.entityId = data.id;

    this.buffer.entityActions = [];
    if (data.actions) {
      for (const item of data.actions) {
        this.buffer.entityActions.push({label: item.description || item.name, value: item.name || 'no name'});
      }
    }

  }

}
</script>

<style>
.el-main {
  padding: 20px 0 0 0;
}
</style>
