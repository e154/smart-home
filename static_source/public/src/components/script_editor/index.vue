<template>
  <div class="script-editor">
    <textarea ref="textarea" />
  </div>
</template>

<script lang="ts">
import CodeMirror, { Editor } from 'codemirror'
import 'codemirror/addon/lint/lint.css'
import 'codemirror/lib/codemirror.css'
import 'codemirror/theme/mdn-like.css'
import 'codemirror/mode/coffeescript/coffeescript'
import 'codemirror/addon/lint/lint'
import 'codemirror/addon/lint/coffeescript-lint'
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'

// HACK: have to use script-loader to load jsonlint
/* eslint-disable import/no-webpack-loader-syntax */
require('script-loader!jsonlint')

@Component({
  name: 'ScriptEditor'
})
export default class extends Vue {
  @Prop({ required: true }) private value!: string

  private jsonEditor?: Editor

  @Watch('value')
  private onValueChange(value: string) {
    if (this.jsonEditor) {
      const editorValue = this.jsonEditor.getValue()
      if (value !== editorValue) {
        this.jsonEditor.setValue(this.value)
      }
    }
  }

  mounted() {
    this.jsonEditor = CodeMirror.fromTextArea(this.$refs.textarea as HTMLTextAreaElement, {
      lineNumbers: true,
      mode: 'application/vnd.coffeescript',
      gutters: ['CodeMirror-lint-markers'],
      theme: 'mdn-like',
      lint: false
    })

    this.jsonEditor.setValue(this.value)
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
.script-editor {
  height: 100%;
  position: relative;
}
</style>
