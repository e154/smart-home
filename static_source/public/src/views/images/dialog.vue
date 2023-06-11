<template>
  <el-dialog
    :visible.sync="vis"
    width="80%"
    append-to-body
    destroy-on-close
  >
    <image-browser @on-select="onSelect"/>
  </el-dialog>
</template>

<script lang="ts">

import { Component, Prop, Vue } from 'vue-property-decorator'
import ImageBrowser from '@/views/images/browser.vue'
import { ApiImage } from '@/api/stub'

@Component({
  name: 'ImageDialog',
  components: {
    ImageBrowser
  }
})
export default class extends Vue {
  @Prop() private visible = false;

  get vis(): boolean {
    return this.visible
  }

  set vis(value: boolean) {
    this.$emit('on-close', false)
  }

  private onSelect(value: ApiImage, event?: any) {
    this.$emit('on-select', value)
  }
}
</script>
