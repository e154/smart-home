<template>
  <el-dialog
    class="export-tool"
    :title="title"
    :visible.sync="vis"
    width="80%"
    append-to-body
    destroy-on-close
    :close-on-press-escape="true"
  >
      <json-editor :value="value" @changed="changed" />

    <span slot="footer" class="dialog-footer">
      <el-button v-if="importDialog" type="primary" @click.prevent.stop="_import">{{$t('main.import')}}</el-button>
      <el-button @click.prevent.stop="vis = false">{{$t('main.cancel')}}</el-button>
    </span>

  </el-dialog>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'
import JsonEditor from '@/components/JsonEditor/index.vue'

@Component({
  name: 'ExportTool',
  components: {
    JsonEditor
  }
})
export default class extends Vue {
  @Prop({ default: () => false }) private importDialog!: boolean;
  @Prop() private visible!: boolean;
  @Prop() private value?: string;
  @Prop({ default: () => '' }) private title!: string;

  get vis(): boolean {
    return this.visible
  }

  set vis(value: boolean) {
    this.$emit('on-close', false)
  }

  private currentValue = '';
  private changed(value: string) {
    this.currentValue = value
  }

  private _import() {
    this.$emit('on-import', this.currentValue)
  }
}
</script>

<style lang="scss">
.export-tool {
  .el-dialog__title {
    color: #a9a9a9;
  }
  .el-dialog {
    position: relative;
    margin: 0 auto 50px;
    background: #3a3a3a;
    border-radius: 2px;
    box-shadow: 0 1px 3px rgb(0 0 0 / 30%);
    box-sizing: border-box;
    width: 50%;
  }
  .el-dialog__body {
    padding: 10px 0;
  }
}

</style>
