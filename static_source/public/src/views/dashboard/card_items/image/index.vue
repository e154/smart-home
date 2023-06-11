<template>
  <transition name="fade">
    <el-image
      v-show="!item.hidden"
      :src="getUrl()">
      <div slot="error" class="image-slot">
        <i class="el-icon-picture-outline"></i>
      </div>
    </el-image>
  </transition>
</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator';
import {CardItem, requestCurrentState} from '@/views/dashboard/core';
import {basePath} from '@/utils';
import {RenderText} from '@/views/dashboard/render';

@Component({
  name: 'IImage',
  components: {}
})
export default class extends Vue {
  @Prop() private item?: CardItem;

  private created() {
    requestCurrentState(this.item?.entityId);
  }

  private mounted() {
  }

  private getUrl(): string {
    if (!this.item?.payload?.image) {
      return '';
    }
    if (this.item?.payload.image.attrField) {
      const imageUrl = RenderText([this.item?.payload.image.attrField], '[' + this.item?.payload.image.attrField + ']', this.item?.lastEvent);
      return basePath + imageUrl;
    }
    return basePath + this.item?.payload.image?.image?.url || '';
  }
}
</script>

<style>
.el-image__error, .el-image__placeholder, .el-image__inner {
  height: auto;
}

.el-image.item-element {
  overflow: visible;
}
</style>
