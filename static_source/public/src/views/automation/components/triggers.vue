<template>
  <div>

    <div v-if="mode != 'VIEW'">
      <el-form label-position="top" label-width="100px"
               ref="currentItem"
               :model="currentItem"
               :rules="rules"
               style="width: 100%">
        <el-form-item :label="$t('automation.table.name')" prop="name">
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

        <el-form-item :label="$t('automation.table.entity')" prop="entity">
          <entity-search
            v-model="currentItem.entity"
            @update-value="changedEntity"
          />
        </el-form-item>

        <el-form-item :label="$t('automation.table.pluginName')" prop="pluginName">
          <el-select
            v-model="currentItem.pluginName"
            placeholder="please select type"
            style="width: 100%"
            @change="changedPlugin"
          >
            <el-option label="STATE_CHANGE" value="state_change"></el-option>
            <el-option label="TIME" value="time"></el-option>
            <el-option label="SYSTEM" value="system"></el-option>
            <el-option label="ALEXA" value="alexa"></el-option>
          </el-select>
        </el-form-item>

        <el-form-item
          v-if="currentItem.pluginName==='time'"
          :label="$t('automation.trigger.pluginOptions')"
          prop="attributes"
        >
          <el-input
            v-model="attributes.time"
            @change="changedAttrParams"
          ></el-input>
        </el-form-item>

        <el-form-item
          v-if="currentItem.pluginName==='alexa'"
          :label="$t('automation.trigger.pluginOptions')"
          prop="attributes"
        >
          <el-input
            type="number"
            v-model="attributes.alexa"
            @change="changedAttrParams"
          ></el-input>
        </el-form-item>

        <el-form-item>
          <el-button v-if="mode == 'NEW'" type="primary" @click="submitForm()">{{
              $t('entities.addAction')
            }}
          </el-button>
          <el-button v-if="mode == 'EDIT'" type="primary" @click="submitForm()">{{ $t('main.update') }}</el-button>
          <el-button @click="resetForm()">{{ $t('main.cancel') }}</el-button>
          <el-button v-if="mode == 'EDIT'" type="danger" @click="removeItem()">{{ $t('main.remove') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div v-if="mode == 'VIEW'">
      <el-row>
        <el-col>
          <el-button
            @click='add()'>
            <i class="el-icon-plus"/> {{ $t('automation.addTrigger') }}
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
                <span> {{ row.name }}</span>
              </template>
            </el-table-column>

            <el-table-column
              :label="$t('automation.table.script')"
              prop="script"
              align="left"
              width="120px"
            >
              <template slot-scope="{row}">
                <span v-if="row.script && row.script.name">{{ row.script.name }}</span>
              </template>
            </el-table-column>

            <el-table-column
              :label="$t('automation.table.entity')"
              prop="entity"
              align="left"
              width="120px"
            >
              <template slot-scope="{row}">
                {{ row.entity.id }}
              </template>
            </el-table-column>

            <el-table-column
              :label="$t('automation.table.pluginName')"
              prop="pluginName"
              align="left"
              width="120px"
            >
              <template slot-scope="{row}">
                <el-tag type="info">
                  {{ row.pluginName }}
                </el-tag>
              </template>
            </el-table-column>

            <el-table-column
              :label="$t('entities.table.operations')"
              prop="description"
              align="right"
              width="auto"
            >
              <template slot-scope="{row, $index}">
                <el-button type="text" size="small" @click='callTrigger(row, $index)'>{{ $t('main.call') }}</el-button>
                <el-button type="text" size="small" @click='editTrigger(row, $index)'>{{ $t('main.edit') }}</el-button>
              </template>
            </el-table-column>

          </el-table>

        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator';
import {ApiAttribute, ApiEntity, ApiScript, ApiTrigger, ApiTypes} from '@/api/stub';
import {Form} from 'element-ui';
import ScriptSearch from '@/views/scripts/components/script_search.vue';
import EntitySearch from '@/views/entities/components/entity_search.vue';
import AttributesEditor from '@/views/entities/components/attributes_editor.vue';
import ScriptEditModal from '@/views/scripts/edit-modal.vue';
import ScriptDialog from '@/views/scripts/dialog.vue';

export enum Mode {
  VIEW = 'VIEW',
  EDIT = 'EDIT',
  NEW = 'NEW'
}

@Component({
  name: 'Triggers',
  components: {
    ScriptSearch,
    EntitySearch,
    AttributesEditor,
    ScriptEditModal,
    ScriptDialog
  }
})
export default class extends Vue {
  @Prop({required: false, default: () => []}) private value?: ApiTrigger[];

  private mode: Mode = Mode.VIEW;
  private currentItem: ApiTrigger = {
    name: ''
  };

  private currentItemIndex?: number;
  private attributes: {
    cron: string
    alexa: number
    time: string
  } = {
    cron: '',
    alexa: 0,
    time: ''
  };

  private dialogScriptVisible = false;

  private rules = {
    name: [
      {required: true, trigger: 'blur'},
      {min: 4, max: 255, trigger: 'blur'}
    ],
    script: [
      {required: true, trigger: 'blur'}
    ],
    entity: [
      {required: true, trigger: 'blur'}
    ],
    pluginName: [
      {required: true, trigger: 'blur'}
    ]
  };

  private editTrigger(trigger: ApiTrigger, index: number) {
    this.currentItem = Object.assign({}, trigger);
    this.currentItemIndex = index;
    this.mode = Mode.EDIT;
    if (this.currentItem && this.currentItem.pluginName) {
      this.changedPlugin(this.currentItem.pluginName);
    }
  }

  private add() {
    this.currentItem = {
      name: ''
    };
    this.mode = Mode.NEW;
  }

  private resetForm() {
    this.currentItem = {};
    this.mode = Mode.VIEW;
    this.currentItemIndex = undefined;
  }

  private removeItem() {
    if (this.value) {
      for (const index in this.value) {
        if (this.currentItem && this.value[index].name == this.currentItem.name) {
          this.value.splice(+index, 1);
        }
      }
    }
    this.mode = Mode.VIEW;
    this.currentItem = {};
  }

  private changedScript(value: ApiScript, event?: any) {
    if (value) {
      this.$set(this.currentItem, 'script', {id: value.id, name: value.name});
    } else {
      this.$set(this.currentItem, 'script', undefined);
    }
  }

  private changedEntity(value: ApiEntity, event?: any) {
    if (value) {
      this.$set(this.currentItem, 'entity', {id: value.id});
    } else {
      this.$set(this.currentItem, 'entity', undefined);
    }
  }

  private changedPlugin(value: string) {
    switch (value) {
      case 'state_change':
        break;
      case 'time':
        if (this.currentItem.attributes && this.currentItem.attributes.cron) {
          this.$set(this.attributes, 'time', this.currentItem.attributes.cron.string);
        }
        break;
      case 'system':
        break;
      case 'alexa':
        if (this.currentItem.attributes && this.currentItem.attributes.skillId) {
          this.$set(this.attributes, 'alexa', this.currentItem.attributes.skillId.int);
        }
        break;
      default:
        console.log('unknown plugin name', value);
    }
  }

  private changedAttrParams(value: any) {
    if (value) {
      let attributes: { [key: string]: ApiAttribute } = {};

      switch (this.currentItem.pluginName) {
        case 'state_change':
          break;
        case 'time':
          attributes.cron = {
            name: 'cron',
            type: ApiTypes.STRING,
            string: value
          };
          break;
        case 'system':
          break;
        case 'alexa':
          attributes.skillId = {
            name: 'skillId',
            type: ApiTypes.INT,
            int: value
          };
          break;
        default:
          console.log('unknown plugin name', value);
      }

      this.$set(this.currentItem, 'attributes', attributes);
    } else {
      this.$set(this.currentItem, 'attributes', undefined);
    }
  }

  private submitForm() {
    (this.$refs.currentItem as Form).validate(valid => {
      if (!valid) {
        return;
      }

      if (this.mode === Mode.NEW) {
        if (this.value) {
          this.value.push(Object.assign({}, this.currentItem));
        }
      } else if (this.mode === Mode.EDIT) {
        if (this.value) {
          if (this.currentItemIndex != undefined) {
            this.value[this.currentItemIndex] = Object.assign({}, this.currentItem);
          }
        }
      }
      this.resetForm();
    });
  }

  private callTrigger(tr: ApiTrigger, index: number) {
    this.$emit('call-trigger', tr.name);
  }
}
</script>

<style>
.el-main {
  padding: 20px 0 0 0;
}
</style>
