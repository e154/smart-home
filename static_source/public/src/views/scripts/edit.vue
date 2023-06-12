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
                       ref="currentScript"
                       :model="currentScript"
                       :rules="rules"
                       style="width: 100%">
                <el-form-item :label="$t('scripts.table.name')" prop="name">
                  <el-input
                    size="medium"
                    placeholder="Name"
                    label="Name"
                    v-model="currentScript.name">
                  </el-input>
                </el-form-item>

                <el-form-item :label="$t('scripts.table.lang')" prop="lang">
                  <el-select v-model="currentScript.lang"
                             placeholder="Language"
                             label="Language"
                             style="width: 100%"
                  >
                    <el-option
                      v-for="item in options"
                      :key="item.value"
                      :label="item.label"
                      :value="item.value">
                    </el-option>
                  </el-select>
                </el-form-item>
                <el-form-item :label="$t('scripts.table.description')" prop="description">
                  <el-input
                    type="textarea"
                    size="medium"
                    placeholder="Description"
                    label="Description"
                    v-model="currentScript.description">
                  </el-input>
                </el-form-item>
              </el-form>
            </el-tab-pane>
            <el-tab-pane
              label="Source"
              name="source"
            >
              <script-editor
                :value="currentScript.source"
                @changed="changed"
              />
            </el-tab-pane>
          </el-tabs>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="24" align="right" class="buttons">
          <el-button type="success" @click.prevent.stop="exec">{{ $t('main.exec') }}</el-button>
          <el-button type="primary" @click.prevent.stop="save">{{ $t('main.save') }}</el-button>
          <el-button @click.prevent.stop="copy">{{ $t('main.copy') }}</el-button>
          <el-button @click.prevent.stop="fetch">{{ $t('main.load_from_server') }}</el-button>
          <el-button @click.prevent.stop="cancel">{{ $t('main.cancel') }}</el-button>
          <el-popconfirm
            :confirm-button-text="$t('main.ok')"
            :cancel-button-text="$t('main.no')"
            icon="el-icon-info"
            icon-color="red"
            style="margin-left: 10px;"
            :title="$t('main.are_you_sure_to_do_want_this?')"
            v-on:confirm="remove"
          >
            <el-button type="danger" icon="el-icon-delete" slot="reference">{{ $t('main.remove') }}</el-button>
          </el-popconfirm>
        </el-col>
      </el-row>
    </card-wrapper>
  </div>
</template>

<script lang="ts">
import 'codemirror/addon/lint/lint.css'
import 'codemirror/lib/codemirror.css'
import 'codemirror/theme/mdn-like.css'
import 'codemirror/mode/coffeescript/coffeescript'
import 'codemirror/addon/lint/lint'
import 'codemirror/addon/lint/json-lint'
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'
import api from '@/api/api'
import { ApiScript } from '@/api/stub'
import router from '@/router'
import { Form } from 'element-ui'
import CardWrapper from '@/components/card-wrapper/index.vue'
import ScriptEditor from '@/components/script_editor/index.vue'

// HACK: have to use script-loader to load jsonlint
/* eslint-disable import/no-webpack-loader-syntax */
require('script-loader!jsonlint')

class elementOption {
  public value = '';
  public label = '';
}

@Component({
  name: 'ScriptEdit',
  components: { ScriptEditor, CardWrapper }
})
export default class extends Vue {
  @Prop({ required: true }) private id!: number;

  private listLoading = true;
  private options: elementOption[] = [
    { value: 'coffeescript', label: 'coffeescript' },
    { value: 'javascript', label: 'javascript' },
    { value: 'typescript', label: 'typescript' }
  ];

  private internal = {
    activeTab: 'source'
  };

  private currentScript: ApiScript = {
    name: '',
    description: '',
    source: '',
    lang: 'coffeescript'
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
    lang: [
      { required: true, trigger: 'blur' },
      { max: 255, trigger: 'blur' }
    ]
  };

  private async fetch() {
    this.listLoading = true
    const { data } = await api.v1.scriptServiceGetScriptById(this.id)
    this.currentScript = data
    this.listLoading = false
  }

  created() {
    this.fetch()
  }

  private async save() {
    (this.$refs.currentScript as Form).validate(async valid => {
      if (!valid) {
        return
      }

      const script: ApiScript = {
        name: this.currentScript.name,
        lang: this.currentScript.lang,
        source: this.currentScript.source,
        description: this.currentScript.description
      }
      const { data } = await api.v1.scriptServiceUpdateScriptById(this.id, script)
      if (data) {
        this.$notify({
          title: 'Success',
          message: 'script updated successfully',
          type: 'success',
          duration: 2000
        })
      }
    })
  }

  private async copy() {
    const { data } = await api.v1.scriptServiceCopyScriptById(this.id)
    router.push({ path: `/scripts/edit/${data.id}` })
  }

  private async remove() {
    await api.v1.scriptServiceDeleteScriptById(this.id)
    this.$notify({
      title: 'Success',
      message: 'Delete Successfully',
      type: 'success',
      duration: 2000
    })
    router.push({ path: '/scripts' })
  }

  private changed(value: string, event?: any) {
    this.currentScript.source = value
  }

  private async exec() {
    await api.v1.scriptServiceExecSrcScriptById({
      name: this.currentScript.name,
      source: this.currentScript.source,
      lang: this.currentScript.lang
    })
    this.$notify({
      title: 'Success',
      message: 'Call Successfully',
      type: 'success',
      duration: 2000
    })
  }

  private cancel() {
    router.go(-1)
  }
}
</script>

<style lang="scss">

.buttons {
  margin: 20px 0;
}

.CodeMirror {
  height: auto;
  min-height: 300px;
  font-family: inherit;
}

.CodeMirror-scroll {
  min-height: 300px;
}

.cm span.cm-string {
  color: #F08047;
}
</style>

<style lang="scss" scoped>
.json-editor {
  height: 100%;
  position: relative;
}
</style>
