<template>
  <div v-loading="listLoading">
    <el-row :gutter="20">
      <el-col :span="24"
              :xs="24"
              class="json-editor">
        <textarea ref="textarea"/>
      </el-col>
    </el-row>

    <el-row>
      <el-col :span="24" align="right" class="buttons">
        <el-button type="primary" @click.prevent.stop="save">{{ $t('main.save') }}</el-button>
        <el-button @click.prevent.stop="reload">{{ $t('main.reload') }}</el-button>
      </el-col>
    </el-row>

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
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'
import api from '@/api/api'
import { ApiScript } from '@/api/stub'

// HACK: have to use script-loader to load jsonlint
/* eslint-disable import/no-webpack-loader-syntax */
require('script-loader!jsonlint')

@Component({
  name: 'ScriptEditModal'
})
export default class extends Vue {
  @Prop({ required: true }) private id!: number;

  private jsonEditor?: Editor;
  private value?: string = '';
  private listLoading = true;

  private currentScript: ApiScript = {
    name: '',
    description: '',
    source: '',
    lang: 'coffeescript'
  };

  private async fetch() {
    this.listLoading = true
    const { data } = await api.v1.scriptServiceGetScriptById(this.id)
    this.currentScript = data
    this.setValue(data.source)
    this.listLoading = false
  }

  created() {
    this.fetch()
  }

  private async save() {
    const script: ApiScript = {
      name: this.currentScript.name,
      lang: this.currentScript.lang,
      source: this.getValue(),
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
  }

  private reload() {
    this.fetch()
  }

  @Watch('id')
  private onIdChange(id: number) {
    this.fetch()
  }

  @Watch('value')
  private onValueChange(value: string) {
    if (this.jsonEditor) {
      const editorValue = this.jsonEditor.getValue()
      if (value !== editorValue) {
        this.jsonEditor.setValue(JSON.stringify(this.value, null, 2))
      }
    }
  }

  mounted() {
    this.jsonEditor = CodeMirror.fromTextArea(this.$refs.textarea as HTMLTextAreaElement, {
      lineNumbers: true,
      mode: 'application/vnd.coffeescript',
      gutters: ['CodeMirror-lint-markers'],
      theme: 'mdn-like',
      lint: true
    })

    this.jsonEditor.setValue(JSON.stringify(this.value, null, 2))
    this.jsonEditor.on('change', editor => {
      this.$emit('changed', editor.getValue())
      this.$emit('input', editor.getValue())
    })
  }

  public setValue(value: string) {
    if (this.jsonEditor) {
      this.jsonEditor.setValue(value)
    }
  }

  public getValue() {
    if (this.jsonEditor) {
      return this.jsonEditor.getValue()
    }
    return ''
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
