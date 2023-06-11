<template>
  <div class="file-manager-body">
    <el-row :gutter="20"
            v-if="mode==='view'">
      <el-col
        v-loading.fullscreen.lock="filterListLoading"
        :span="6"
        :xs="24"
        element-loading-text="Loading ..."
      >
        <ul class="list-unstyled filters">
          <li
            v-for="item of filterList.slice().reverse()"
            v-bind:class="{selected: getActiveFilter(item)}"
          >
            <el-button
              type="text"
              @click.prevent.stop="getList(item)"
            >
              {{ item.date }}
              <span class="badge">{{ item.count }}</span>
            </el-button>
          </li>
        </ul>
      </el-col>
      <el-col
        v-loading.fullscreen.lock="listLoading"
        :span="18"
        :xs="24"
        element-loading-text="Loading ..."
      >
        <ul class="list-unstyled file-manager-items">
          <a @click.prevent.stop="mode='upload'">
            <li class="drop-box file-manager-item">
              <i class="el-icon-plus"/>
              <div class="title">upload</div>
            </li>
          </a>
          <li class="file-manager-item"
              v-for="image of list"
              @click.prevent.stop="select(image)"
          >
            <el-image
              :src="getUrl(image)"
              fit="fil">
            </el-image>
            <div class="file-manager-item-title">{{ image.name }}</div>
            <div class="cross close-button" @click.prevent.stop="removeFromServer(image)"></div>
            <div class="is_selected" v-if="selected && selected.id === image.id">
              <i class="el-icon-check"/>
            </div>
          </li>
        </ul>
      </el-col>
    </el-row>
    <el-row :gutter="20"
            v-if="mode==='upload'">
      <el-col
        :span="6"
        :xs="24"
      >
        <ul class="list-unstyled filters">
          <li>
            <el-button
              type="text"
              @click.prevent.stop="mode='view'"
            >
              archive
            </el-button>
          </li>
        </ul>
      </el-col>
      <el-col
        :span="18"
        :xs="24"
      >
        <el-upload
          list-type="picture-card"
          :multiple="true"
          ref="upload"
          :on-success="onSuccess"
          :action="getUploadURL()"
          :auto-upload="true">
          <i slot="default" class="el-icon-plus"></i>
          <div slot="file" slot-scope="{file}">
            <img
              class="el-upload-list__item-thumbnail"
              :src="file.url" alt=""
            >
            <span class="el-upload-list__item-actions">
               <span
                 class="el-upload-list__item-preview"
                 @click="handlePictureCardPreview(file)"
               >
                        <i class="el-icon-zoom-in"></i>
                </span>
            </span>
          </div>
        </el-upload>
        <el-dialog :visible.sync="dialogVisible">
          <img width="100%" :src="dialogImageUrl" alt="">
        </el-dialog>
      </el-col>
    </el-row>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import api from '@/api/api'
import { ApiImage, GetImageFilterListResultfilter } from '@/api/stub'

@Component({
  name: 'ImageBrowser'
})
export default class extends Vue {
  private mode = 'view';
  private list: ApiImage[] = [];
  private filterList: GetImageFilterListResultfilter[] = [];
  private filterListLoading = true;
  private listLoading = true;
  private basePath: string = process.env.VUE_APP_BASE_API || window.location.origin;
  private selected?: ApiImage = {};
  private currentFilter?: GetImageFilterListResultfilter;

  private dialogImageUrl = '';
  private dialogVisible = false;

  async created() {
    this.fetch()
  }

  private async fetch() {
    await this.getFilterList()
    if (this.filterList && this.filterList.length > 0) {
      this.getList(this.filterList[this.filterList.length - 1])
    }
  }

  private async getList(filter?: GetImageFilterListResultfilter) {
    if (filter) {
      this.currentFilter = filter
      this.listLoading = true
      const { data } = await api.v1.imageServiceGetImageListByDate({ filter: filter.date })
      this.list = data.items
      this.listLoading = false
    }
  }

  private async getFilterList() {
    this.filterListLoading = true
    const { data } = await api.v1.imageServiceGetImageFilterList()
    this.filterList = data.items
    this.filterListLoading = false
  }

  private getUrl(image: ApiImage): string {
    return this.basePath + image.url
  }

  private select(image: ApiImage) {
    if (image) {
      if (this.selected && this.selected.id === image.id) {
        this.selected = undefined
        return
      }
      this.selected = image
    } else {
      this.selected = undefined
    }
    this.$emit('on-select', this.selected)
  }

  private async removeFromServer(image: ApiImage) {
    await api.v1.imageServiceDeleteImageById(image.id || 0)
    this.getList(this.currentFilter)
    this.getFilterList()
  }

  private handlePictureCardPreview(file: any) {
    this.dialogImageUrl = file.url
    this.dialogVisible = true
  }

  private getUploadURL(): string {
    return (process.env.VUE_APP_BASE_API || '') + '/v1/image/upload'
  }

  private onSuccess() {
    this.getFilterList()
    this.getList(this.currentFilter)
    this.$notify({
      title: 'Success',
      message: 'Upload successfully',
      type: 'success',
      duration: 2000
    })
  }

  private getActiveFilter(item: GetImageFilterListResultfilter): boolean {
    if (this.currentFilter) {
      return this.currentFilter.date === item.date
    }
    return false
  }
}
</script>

<style lang="scss" scoped>

.list-unstyled {
  list-style: none;
}

.file-manager-body {
  padding: 20px;
  position: relative;

.el-image {
  display: block;
}

.drop-box {
  background: #F8F8F8;
  border: 1px dashed #DDD;
  text-align: center;
  padding-top: 25px;
  cursor: pointer;

.title {
  font-size: 10px;
}

}
ul.filters {

li {
  position: relative;

.badge {
  font-size: 10px;
  position: absolute;
  top: 0;
}

&.selected {
  font-weight: 600;
}

}
}

ul.file-manager-items {
  margin: 0;
  padding: 0;
  list-style: none;

.fa-cloud-upload {
  font-size: 3rem;
}

li.file-manager-item {
  width: 150px;
  height: 100px;
  float: left;
  margin: 5px;
  overflow: hidden;
  position: relative;

img {
  background-repeat: no-repeat;
  background-size: 100% auto;
  width: 100%;
}

.file-manager-item-title {
  position: absolute;
  bottom: 0;
  background: black;
  opacity: 0.5;
  color: #ffffff;
  font-size: 10px;
  left: 0;
  right: 0;
  padding: 2px 8px;
}

.cross.close-button {
  opacity: 0;
  background-color: #FFFFFF;
  position: absolute;
  top: 0;
  right: 0;
  cursor: pointer;

}

.is_selected {
  position: absolute;
  top: 0;
  left: 0;
  bottom: 0;
  right: 0;
  background: rgba(0, 0, 0, 0.5);
  text-align: center;
  font-size: 35px;
  color: #FFF;
  padding-top: 31px;
}

&:hover {

.cross.close-button {
  opacity: 0.7;
  -webkit-transition: opacity 0.6s ease-in-out;
  -moz-transition: opacity 0.6s ease-in-out;
  -ms-transition: opacity 0.6s ease-in-out;
  -o-transition: opacity 0.6s ease-in-out;
  transition: opacity 0.6s ease-in-out;
}

}
}
}
}

</style>
