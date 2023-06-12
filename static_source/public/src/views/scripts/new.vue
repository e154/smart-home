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
                             size="medium"
                             placeholder="Language"
                             label="Language"
                             style="width: 100%">
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
              :lazy="true"
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
        <el-col style="text-align: right" class="buttons">
          <el-button type="primary" @click.prevent.stop="save">{{ $t('main.create') }}</el-button>
          <el-button @click.prevent.stop="cancel">{{ $t('main.cancel') }}</el-button>
        </el-col>
      </el-row>

    </card-wrapper>

  </div>
</template>

<script lang="ts">
import CodeMirror, { Editor } from 'codemirror'
import 'codemirror/addon/lint/lint.css'
import 'codemirror/lib/codemirror.css'
import 'codemirror/theme/mdn-like.css'
import 'codemirror/mode/coffeescript/coffeescript'
import 'codemirror/addon/lint/lint'
import 'codemirror/addon/lint/json-lint'
import { Component, Vue, Watch } from 'vue-property-decorator'
import api from '@/api/api'
import { ApiScript } from '@/api/stub'
import router from '@/router'
import { Form } from 'element-ui'
import CardWrapper from '@/components/card-wrapper/index.vue'
import ScriptEditor from '@/components/script_editor/index.vue'
import stream from '@/api/stream'

// HACK: have to use script-loader to load jsonlint
/* eslint-disable import/no-webpack-loader-syntax */
require('script-loader!jsonlint')

class elementOption {
  public value = '';
  public label = '';
}

@Component({
  name: 'ScriptNew',
  components: { ScriptEditor, CardWrapper }
})
export default class extends Vue {
  private options: elementOption[] = [
    { value: 'coffeescript', label: 'coffeescript' },
    { value: 'javascript', label: 'javascript' },
    { value: 'typescript', label: 'typescript' }
  ];

  private internal = {
    activeTab: 'main'
  };

  private currentScript: ApiScript = {
    name: '',
    description: '',
    source: 'main =->\n\n\n\n\n\n\n\n\n\n\n\n',
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

  created() {

  }

  private changed(value: string, event?: any) {
    this.currentScript.source = value
  }

  private async save() {
    (this.$refs.currentScript as Form).validate(async valid => {
      if (!valid) {
        return
      }

      const { data } = await api.v1.scriptServiceAddScript(this.currentScript)
      if (data) {
        this.$notify({
          title: 'Success',
          message: 'script created successfully',
          type: 'success',
          duration: 2000
        })
        router.push({ path: `/scripts/edit/${data.id}` })
      }
    })
  }

  private cancel() {
    router.go(-1)
  }
}
</script>

<style lang="scss">
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

.buttons {
  margin: 20px 0;
}
</style>

<style lang="scss" scoped>
.json-editor {
  height: 100%;
  position: relative;
}
</style>
