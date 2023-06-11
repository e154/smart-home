<template>
  <transition name="fade">
    <el-button
      style="width: 100%; height: 100%"
      v-if="item.enabled" v-show="!item.hidden"
      :size="item.payload.button.size"
      :type="item.payload.button.type"
      :icon="item.payload.button.icon"
      :round="item.payload.button.round"
      @click.prevent.stop="onClick"
      :disabled="disabled"
    >{{ item.payload.button.text }}
    </el-button>
  </transition>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import {CardItem, requestCurrentState} from '@/views/dashboard/core';
import api from '@/api/api'

@Component({
  name: 'IButton',
  components: {}
})
export default class extends Vue {
  @Prop() private item?: CardItem;
  @Prop({ default: false }) private disabled!: boolean;

  private created() {
    requestCurrentState(this.item?.entityId);
  }

  private mounted() {
  }

  private async callAction() {
    await api.v1.interactServiceEntityCallAction({
      id: this.item?.entityId,
      name: this.item?.payload.button?.action || ''
    })
    this.$notify({
      title: 'Success',
      message: 'Call Successfully',
      type: 'success',
      duration: 2000
    })
  }

  private onClick() {
    this.callAction()
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
