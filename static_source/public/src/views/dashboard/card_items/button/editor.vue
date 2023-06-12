<template>
  <div>
    <common-editor :item="item" :board="board"></common-editor>

    <el-divider content-position="left">{{$t('dashboard.editor.buttonOptions')}}</el-divider>

    <el-row :gutter="20">
      <el-col
        :span="8"
        :xs="8"
      >
        <el-form-item :label="$t('dashboard.editor.icon')" prop="icon">
          <el-input size="small"
                    v-model="item.payload.button.icon"></el-input>
        </el-form-item>

        <el-form-item :label="$t('dashboard.editor.text')" prop="text">
          <el-input size="small"
                    v-model="item.payload.button.text"></el-input>
        </el-form-item>

        <el-form-item :label="$t('dashboard.editor.action')" prop="action" :aria-disabled="!item.entity">

          <el-select
            v-model="item.payload.button.action"
            clearable
            :placeholder="$t('dashboard.editor.selectAction')"
            style="width: 100%"
          >
            <el-option
              v-for="item in item.entityActions"
              :key="item.value"
              :label="item.label"
              :value="item.value">
            </el-option>
          </el-select>

        </el-form-item>

      </el-col>
      <el-col
        :span="8"
        :xs="8"
      >

        <el-form-item :label="$t('dashboard.editor.type')" prop="type">
          <el-select
            v-model="item.payload.button.type"
            placeholder="please select type"
            style="width: 100%"
          >
            <el-option label="primary" value="primary"></el-option>
            <el-option label="success" value="success"></el-option>
            <el-option label="info" value="info"></el-option>
            <el-option label="warning" value="warning"></el-option>
            <el-option label="danger" value="danger"></el-option>
          </el-select>
        </el-form-item>

        <el-form-item :label="$t('dashboard.editor.size')" prop="size">
          <el-select
            v-model="item.payload.button.size"
            placeholder="please select type"
            style="width: 100%"
          >
            <el-option label="mini" value="mini"></el-option>
            <el-option label="small" value="small"></el-option>
            <el-option label="medium" value="medium"></el-option>
            <el-option label="default" value="default"></el-option>
          </el-select>
        </el-form-item>

      </el-col>
      <el-col
        :span="8"
        :xs="8"
      >

        <el-form-item :label="$t('dashboard.editor.round')" prop="round">
          <el-switch
            v-model="item.payload.button.round"></el-switch>
        </el-form-item>
      </el-col>
    </el-row>

  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import { CardItem, Core } from '@/views/dashboard/core'
import EntitySearch from '@/views/entities/components/entity_search.vue'
import CommonEditor from '@/views/dashboard/card_items/common/editor.vue'
import { UUID } from 'uuid-generator-ts'
import { EventStateChange } from '@/api/stream_types'

@Component({
  name: 'IButtonEditor',
  components: { CommonEditor, EntitySearch }
})
export default class extends Vue {
  @Prop() private item!: CardItem;
  @Prop() private board!: Core;

  private currentID = '';

  private created() {
    const uuid = new UUID()
    this.currentID = uuid.getDashFreeUUID()

    this.board.bus.$on('state_changed', this.onStateChanged)
  }

  private mounted() {
  }

  private destroyed() {
    this.board.bus.$off('state_changed', this.onStateChanged)
  }

  private onStateChanged(m: EventStateChange) {
    if (!this.item.entityId || m.entity_id != this.item.entityId) {

    }
    // console.log(m)
  }
}
</script>

<style scoped>

.clearfix:before,
.clearfix:after {
  display: table;
  content: "";
}

.clearfix:after {
  clear: both
}

</style>
