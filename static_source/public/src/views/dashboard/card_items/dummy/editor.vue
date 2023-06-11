<template>
  <div>
    <common-editor :item="item" :board="board"></common-editor>

    <el-divider content-position="left">Dummy options</el-divider>

    [Unknown item type "{{ item.type }}"]

  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import { CardItem, Core } from '@/views/dashboard/core'
import CommonEditor from '@/views/dashboard/card_items/common/editor.vue'

@Component({
  name: 'DummyEditor',
  components: { CommonEditor }
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
