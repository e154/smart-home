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
        <el-form-item :label="$t('entities.table.description')" prop="description">
          <el-input v-model="currentItem.description"></el-input>
        </el-form-item>
        <el-form-item :label="$t('entities.table.icon')" prop="icon">
          <el-input v-model="currentItem.icon"></el-input>
        </el-form-item>
        <el-form-item :label="$t('entities.table.image')" prop="image">
          <image-preview :image="currentItem.image" @on-select="onSelectImage"/>
        </el-form-item>
        <el-form-item :label="$t('entities.table.script')" prop="script">
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

        <el-form-item>
          <el-button v-if="mode == 'NEW'" type="primary" @click="submitForm()">{{$t('entities.addAction') }}</el-button>
          <el-button v-if="mode == 'EDIT'" type="primary" @click="submitForm()">{{ $t('main.update') }}</el-button>
          <el-button @click="resetForm()">{{ $t('main.cancel') }}</el-button>
          <el-button v-if="mode == 'EDIT'" type="danger" @click="removeItem()">{{ $t('main.remove') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div v-if="mode == 'VIEW'">
      <el-main>
        <el-row>
          <el-col>
            <el-button
              @click='add()'>
              <i class="el-icon-plus"/> {{ $t('entities.addAction') }}
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
                :label="$t('entities.table.name')"
                prop="name"
                align="left"
                width="185px"
              >
                <template slot-scope="{row}">
                  <div>{{ row.name }}</div>
                </template>
              </el-table-column>

              <el-table-column
                :label="$t('entities.table.image')"
                prop="image"
                align="center"
                width="60px"
              >
                <template slot-scope="{row}">
                  <i v-if="row.image" :class="'el-icon-check'"/>
                  <i v-if="!row.image" :class="'el-icon-minus'"/>
                </template>
              </el-table-column>

              <el-table-column
                :label="$t('entities.table.icon')"
                prop="icon"
                align="center"
                width="50px"
              >
                <template slot-scope="{row}">
                  <i v-if="row.icon" :class="'el-icon-check'"/>
                  <i v-if="!row.icon" :class="'el-icon-minus'"/>
                </template>
              </el-table-column>

              <el-table-column
                :label="$t('entities.table.script')"
                prop="script"
                align="center"
                width="70px"
              >
                <template slot-scope="{row}">
                  <i v-if="row.script" :class="'el-icon-check'"/>
                  <i v-if="!row.script" :class="'el-icon-minus'"/>
                </template>
              </el-table-column>

              <el-table-column
                :label="$t('entities.table.description')"
                prop="description"
                align="left"
                width="auto"
              >
                <template slot-scope="{row}">
                  <div>{{ row.description }}</div>
                </template>
              </el-table-column>

              <el-table-column
                :label="$t('entities.table.operations')"
                prop="description"
                align="left"
                width="180px"
              >
                <template slot-scope="{row, $index}">
                  <el-button @click="callAction(row)" type="text" size="small">{{ $t('entities.callAction') }}</el-button>
                  <el-button type="text" size="small" @click='editAction(row, $index)'>{{ $t('main.edit') }}</el-button>
                </template>
              </el-table-column>

            </el-table>

          </el-col>
        </el-row>
      </el-main>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import { ApiEntityAction, ApiGetPluginOptionsResultEntityAction, ApiImage, ApiScript } from '@/api/stub'
import ScriptSearch from '@/views/scripts/components/script_search.vue'
import ScriptEditModal from '@/views/scripts/edit-modal.vue'
import ScriptDialog from '@/views/scripts/dialog.vue'
import ImagePreview from '@/views/images/preview.vue'
import { Form } from 'element-ui'

export enum Mode {
  VIEW = 'VIEW',
  EDIT = 'EDIT',
  NEW = 'NEW'
}

@Component({
  name: 'Actions',
  components: {
    ScriptSearch,
    ScriptEditModal,
    ScriptDialog,
    ImagePreview
  }
})
export default class extends Vue {
  @Prop() private value?: ApiEntityAction[];
  @Prop() private settings?: Map<string, ApiGetPluginOptionsResultEntityAction>;
  @Prop({ default: false }) private customActions?: boolean;

  private mode: Mode = Mode.VIEW;
  private currentItem: ApiEntityAction = {} as ApiEntityAction;
  private currentItemIndex?: number;
  private dialogScriptVisible = false;

  private rules = {
    name: [
      { required: true, trigger: 'blur' },
      { min: 4, max: 255, trigger: 'blur' }
    ],
    description: [
      { required: false, trigger: 'blur' },
      { max: 255, trigger: 'blur' }
    ]
  };

  private changedScript(value: ApiScript, event: any) {
    if (value) {
      this.$set(this.currentItem, 'script', value)
    } else {
      this.$set(this.currentItem, 'script', undefined)
    }
  }

  private editAction(action: ApiEntityAction, index: number) {
    this.currentItem = Object.assign({}, action)
    this.currentItemIndex = index
    this.mode = Mode.EDIT
  }

  private add() {
    this.currentItem = {} as ApiEntityAction
    this.mode = Mode.NEW
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

  private resetForm() {
    this.currentItem = {} as ApiEntityAction
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
    this.currentItem = {} as ApiEntityAction
  }

  private callAction(action: ApiEntityAction) {
    this.$emit('call-action', action.name)
  }

  private onSelectImage(value: ApiImage, event?: any) {
    this.$set(this.currentItem, 'image', value)
  }
}
</script>

<style>
.el-main {
  padding: 20px 0 0 0;
}
</style>
