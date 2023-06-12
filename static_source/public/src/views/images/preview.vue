<template>
  <div>
    <div v-if="currentUrl !== ''" class="image-preview">
      <el-image
        :src="currentUrl"
        fit="fil"
        :preview-src-list="[currentUrl]"
      >
      </el-image>
      <a href="#" class="cross delete-btn"
         @click.prevent.stop="remove()">
      </a>
    </div>
    <div v-else>
      <el-button
        style="width: 100%"
        @click="visible=true">
        <i class="el-icon-upload"/> {{$t('upload')}}
      </el-button>
    </div>
    <image-dialog
      :visible.sync="visible"
      @on-select="onSelect"
      @on-close="visible=false"
    />
  </div>
</template>

<script lang="ts">

import { Component, Prop, Vue, Watch } from 'vue-property-decorator'
import ImageBrowser from '@/views/images/browser.vue'
import { ApiImage } from '@/api/stub'
import ImageDialog from '@/views/images/dialog.vue'

@Component({
  name: 'ImagePreview',
  components: {
    ImageBrowser,
    ImageDialog
  }
})
export default class extends Vue {
  @Prop() private image?: ApiImage;

  private currentUrl = '';
  private currentImage?: ApiImage;
  private visible = false;
  private basePath: string = process.env.VUE_APP_BASE_API || window.location.origin;

  private created() {
    if (this.image) {
      this.currentImage = this.image
      this.getUrl()
    }
  }

  @Watch('image')
  private watchImage(image: ApiImage) {
    this.currentImage = image
    this.getUrl()
  }

  private getUrl() {
    if (this.currentImage) {
      this.currentUrl = this.basePath + this.currentImage.url
    } else {
      this.currentUrl = ''
    }
  }

  private onSelect(image: ApiImage, event?: any) {
    this.$emit('on-select', image)
  }

  private remove() {
    this.$emit('on-select', undefined)
  }
}
</script>

<style lang="scss" scoped>
.image-preview {
  max-height: 60px;
  overflow: hidden;
  position: relative;
  border: 1px solid #DCDFE6;
  text-align: center;

  .cross.delete-btn {
    background-color: #FFFFFF;
    position: absolute;
    top: 0;
    right: 0;
    cursor: pointer;
  }

  &:hover {
    .cross.delete-btn {
      opacity: 0.7;
      -webkit-transition: opacity 0.6s ease-in-out;
      -moz-transition: opacity 0.6s ease-in-out;
      -ms-transition: opacity 0.6s ease-in-out;
      -o-transition: opacity 0.6s ease-in-out;
      transition: opacity 0.6s ease-in-out;
    }
  }
}

</style>
