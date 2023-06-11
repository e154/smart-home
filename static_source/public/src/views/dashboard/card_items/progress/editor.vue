<template>
  <div>
    <common-editor :item="item" :board="board"></common-editor>

    <el-divider content-position="left">Progress options</el-divider>

    <el-form-item :label="$t('dashboard.editor.type')" prop="type">
      <el-select
        v-model="item.payload.progress.type"
        placeholder="please select type"
        style="width: 100%"
      >
        <el-option label="linear" value=""></el-option>
        <el-option label="circle" value="circle"></el-option>
        <el-option label="dashboard" value="dashboard"></el-option>
      </el-select>
    </el-form-item>

    <el-form-item :label="$t('dashboard.editor.textInside')" prop="textInside">
      <el-switch
        v-model="item.payload.progress.textInside"></el-switch>
    </el-form-item>

    <el-form-item :label="$t('dashboard.editor.strokeWidth')" prop="strokeWidth">
      <el-input-number size="small"
                       v-model="item.payload.progress.strokeWidth"></el-input-number>
    </el-form-item>

    <el-form-item :label="$t('dashboard.editor.width')" prop="width">
      <el-input-number size="small"
                       v-model="item.payload.progress.width"></el-input-number>
    </el-form-item>

    <el-form-item :label="$t('dashboard.editor.value')" prop="value">
      <el-input size="small"
                v-model="item.payload.progress.value"></el-input>
    </el-form-item>

    <el-row style="padding-bottom: 20px">
      <el-col>
        <event-viewer :item="item"></event-viewer>
      </el-col>
    </el-row>

  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import { CardItem, Core } from '@/views/dashboard/core'
import CommonEditor from '@/views/dashboard/card_items/common/editor.vue'
import EventViewer from '@/views/dashboard/card_items/common/event_viewer.vue'

@Component({
  name: 'ProgressEditor',
  components: { EventViewer, CommonEditor }
})
export default class extends Vue {
  @Prop() private item!: CardItem;
  @Prop() private board!: Core;

  private created() {
  }

  private mounted() {
  }
}
</script>

<style scoped>

</style>
