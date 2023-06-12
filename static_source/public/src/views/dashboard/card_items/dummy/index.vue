<template>
  <transition name="fade">
    [Unknown item type "{{item.type}}"]
  </transition>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import {CardItem, requestCurrentState} from '@/views/dashboard/core';

@Component({
  name: 'Dummy',
  components: {
  }
})
export default class extends Vue {
  @Prop() private item?: CardItem;

  private created() {
    requestCurrentState(this.item?.entityId);
  }
  private mounted() {}
}
</script>

<style scoped>

</style>
