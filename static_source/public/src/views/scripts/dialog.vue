<template>
  <el-dialog
    v-if="script"
    :title="script.name"
    :visible.sync="vis"
    width="80%"
    append-to-body
    destroy-on-close
  >
    <script-edit-modal :id="script.id"/>
  </el-dialog>
</template>

<script lang="ts">

import { Component, Prop, Vue } from 'vue-property-decorator'
import { ApiAttribute, ApiScript } from '@/api/stub'
import ScriptEditModal from '@/views/scripts/edit-modal.vue'

@Component({
  name: 'ScriptDialog',
  components: {
    ScriptEditModal
  }
})
export default class extends Vue {
  @Prop({ required: true }) private visible = false;
  @Prop({ required: true }) private script?: ApiScript;

  get vis(): boolean {
    return this.visible
  }

  set vis(value: boolean) {
    this.$emit('on-close', false)
  }
}
</script>
