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
        <el-form-item :label="$t('entities.table.image')" prop="image">
          <image-preview :image="currentItem.image" @on-select="onSelectImage"/>
        </el-form-item>
        <el-form-item>
          <el-button v-if="mode == 'NEW'" type="primary" @click="submitForm()">{{ $t('entities.addState') }}</el-button>
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
              @click='add()'
              v-if="customStates">
              <i class="el-icon-plus"/> {{ $t('entities.addState') }}
            </el-button>

            <el-button
              align="left"
              v-if="settings && Object.keys(settings).length"
              @click='loadFromPlugin()'>
              {{ $t('entities.loadFromPlugin') }}
            </el-button>
          </el-col>
        </el-row>
        <el-row>
          <el-col>

            <el-table
              key="key"
              :data="states"
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
                  <el-button @click="setState(row)" type="text" size="small">{{ $t('entities.setState') }}</el-button>
                  <el-button type="text" size="small" @click='editState(row, $index)'>{{ $t('main.edit') }}</el-button>
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
import { ApiEntityState, ApiGetPluginOptionsResultEntityState, ApiImage } from '@/api/stub'
import { Form } from 'element-ui'
import ImagePreview from '@/views/images/preview.vue'

export enum Mode {
  VIEW = 'VIEW',
  EDIT = 'EDIT',
  NEW = 'NEW'
}

@Component({
  name: 'States',
  components: {
    ImagePreview
  }
})
export default class extends Vue {
  @Prop() private value?: ApiEntityState[];
  @Prop() private settings?: { [key: string]: ApiGetPluginOptionsResultEntityState };
  @Prop({ default: false }) private customStates?: boolean;

  private mode: Mode = Mode.VIEW;
  private currentItem: ApiEntityState = {} as ApiEntityState;
  private currentItemIndex?: number;

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

  get states(): ApiEntityState[] {
    const states: ApiEntityState[] = []
    if (this.value) {
      for (const key in this.value) {
        states.push(this.value[key])
      }
    }
    return states
  }

  set states(value: ApiEntityState[]) {

  }

  private setState(state: ApiEntityState) {
    this.$emit('set-state', state.name)
  }

  private editState(state: ApiEntityState, index: number) {
    this.currentItem = Object.assign({}, state)
    this.currentItemIndex = index
    this.mode = Mode.EDIT
  }

  private add() {
    this.currentItem = {} as ApiEntityState
    this.mode = Mode.NEW
  }

  private submitForm() {
    (this.$refs.currentItem as Form).validate(valid => {
      if (!valid) {
        return
      }
      if (this.mode === Mode.NEW) {
        if (this.states) {
          this.states.push(this.currentItem)
        }
      } else if (this.mode === Mode.EDIT) {
        if (this.states && this.currentItemIndex != undefined) {
          this.states[this.currentItemIndex] = this.currentItem
        }
      }
      const states: ApiEntityState[] = []
      for (const index in this.states) {
        states.push(this.states[index])
      }
      this.$emit('update-value', states)
      this.resetForm()
      setTimeout(() => {
        this.$forceUpdate()
      }, 0.5 * 1000)
    })
  }

  private resetForm() {
    this.currentItem = {} as ApiEntityState
    this.mode = Mode.VIEW
    this.currentItemIndex = undefined
  }

  private removeItem() {
    if (this.states) {
      for (const index in this.states) {
        if (this.currentItem && this.states[index].name == this.currentItem.name) {
          this.states.splice(+index, 1)
        }
      }
    }
    this.mode = Mode.VIEW
    this.currentItem = {} as ApiEntityState
  }

  private loadFromPlugin() {
    const value: ApiEntityState[] = []
    if (this.settings) {
      for (const k in this.settings) {
        const item = this.settings[k]
        value.push({
          name: item.name || '',
          description: item.description,
          icon: item.icon
        })
      }
    }
    this.$emit('update-value', value)
    setTimeout(() => {
      this.$forceUpdate()
    }, 0.5 * 1000)
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
