<template>
  <el-collapse>
    <el-collapse-item :title="$t('dashboard.editor.eventstateJSONobject')">

      <el-button type="default" @click.prevent.stop="requestCurrentState()" style="margin-bottom: 20px"><i
        class="el-icon-refresh"/> {{ $t('dashboard.editor.getEvent') }}
      </el-button>

      <json-editor
        ref="jsoneditor"
        :value="lastEvent"
      />
    </el-collapse-item>
  </el-collapse>
</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator'
import {requestCurrentState} from '@/views/dashboard/core'
import CommonEditor from '@/views/dashboard/card_items/common/editor.vue'
import JsonEditor from '@/components/JsonEditor/index.vue'
import {EventStateChange} from "@/api/stream_types";

export interface Item {
  entityId: string;
  lastEvent: EventStateChange;
}

@Component({
  name: 'EventViewer',
  components: {JsonEditor, CommonEditor}
})
export default class extends Vue {
  @Prop() private item!: Item;

  get lastEvent() {
    return this.item.lastEvent || {}
  }

  private requestCurrentState() {
    if (this.item.entityId) {
      requestCurrentState(this.item.entityId)
    }
  }
}
</script>

<style scoped>

</style>
