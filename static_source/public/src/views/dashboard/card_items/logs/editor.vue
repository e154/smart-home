<template>
  <div>
    <common-editor :item="item" :board="board"></common-editor>

    <el-divider content-position="left">Logs options</el-divider>

    <el-form label-position="top"
             :model="item.payload.logs"
             style="width: 100%"
             :rules="logsRules"
             ref="cardForm">

      <el-form-item :label="$t('dashboard.editor.limit')" prop="limit">
        <el-input-number size="small"
                         v-model="item.payload.logs.limit"></el-input-number>
      </el-form-item>

    </el-form>

  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import { CardItem, Core } from '@/views/dashboard/core'
import CommonEditor from '@/views/dashboard/card_items/common/editor.vue'

@Component({
  name: 'LogsEditor',
  components: { CommonEditor }
})
export default class extends Vue {
  @Prop() private item!: CardItem;
  @Prop() private board!: Core;

  private logsRules = {
    limit: [
      { required: true, trigger: 'blur' },
      { min: 1, max: 255, trigger: 'blur' }
    ]
  };

  private created() {
  }

  private mounted() {
  }
}
</script>

<style scoped>

</style>
