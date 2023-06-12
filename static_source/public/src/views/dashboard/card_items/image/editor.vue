<template>
  <div>
    <common-editor :item="item" :board="board"></common-editor>

    <el-divider content-position="left">{{ $t('dashboard.editor.imageOptions') }}</el-divider>

    <el-form-item :label="$t('dashboard.editor.image')" prop="image">
      <image-preview
        :image="item.payload.image.image"
        @on-select="onSelectImage(index, ...arguments)"/>
    </el-form-item>


    <el-form-item :label="$t('dashboard.editor.attrField')" prop="text">
      <el-input size="small"
                v-model="item.payload.image.attrField"></el-input>
    </el-form-item>

    <el-row style="padding-bottom: 20px">
      <el-col>
        <event-viewer :item="item"></event-viewer>
      </el-col>
    </el-row>

  </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator';
import {CardItem, Core} from '@/views/dashboard/core';
import {ApiImage} from '@/api/stub';
import ImagePreview from '@/views/images/preview.vue';
import CommonEditor from '@/views/dashboard/card_items/common/editor.vue';
import EventViewer from '@/views/dashboard/card_items/common/event_viewer.vue';

@Component({
  name: 'IImageEditor',
  components: {
    EventViewer,
    CommonEditor,
    ImagePreview
  }
})
export default class extends Vue {
  @Prop() private item!: CardItem;
  @Prop() private board!: Core;
  @Prop() private index!: number;

  private created() {
  }

  private mounted() {
  }

  private onSelectImage(index: number, image: ApiImage) {
    if (!this.item.payload?.image) {
      return;
    }
    // console.log('select image', index, image);
    this.item.payload.image.image = image || undefined;
  }
}
</script>

<style>

</style>
