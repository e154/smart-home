<script setup lang="ts">
import {ElButton, ElRow, ElCol, ElBadge, ElImage, ElUpload, UploadProps, ElDialog, ElMessage} from 'element-plus'
import { useI18n } from '@/hooks/web/useI18n'
import {PropType, reactive, ref, unref} from "vue";
import {ApiCondition, ApiImage, GetImageFilterListResultfilter} from "@/api/stub";
import api from "@/api/api";
import {useEmitt} from "@/hooks/web/useEmitt";
import {createImageViewer} from "@/components/ImageViewer";
import {propTypes} from "@/utils/propTypes";
import {useCache} from "@/hooks/web/useCache";
const {wsCache} = useCache()

const { t } = useI18n()

interface ViewerObject {
  loading: boolean
  filterList?:  GetImageFilterListResultfilter[]
  currentFilter?: GetImageFilterListResultfilter
  imageList?: ApiImage[]
  selected?: ApiImage
}

const props = defineProps({
  id: propTypes.string.def(''),
})

const viewerObject = reactive<ViewerObject>(
    {
      loading: false,
    }
)

const fetch =  async() => {
  await getFilterList()
  if (viewerObject.filterList && viewerObject.filterList.length > 0) {
    getList(viewerObject.filterList[viewerObject.filterList.length - 1])
  }
}

const getList =  async(filter?: GetImageFilterListResultfilter) => {
  if (!filter) {
    return
  }
  viewerObject.currentFilter = filter;
  viewerObject.loading = true
  const res = await api.v1.imageServiceGetImageListByDate({ filter: filter.date })
      .catch(() => {
      })
      .finally(() => {
        viewerObject.loading = false
      })

  let {items} = unref(res.data);
  for (const key in items) {
    items[key].url = getUrl(items[key].url)
  }
  viewerObject.imageList = items
}

const getFilterList =  async() => {
  viewerObject.loading = true

  const res = await api.v1.imageServiceGetImageFilterList()
      .catch(() => {
      })
      .finally(() => {
        viewerObject.loading = false
      })
  if (res) {
    const {items} = res.data;
    viewerObject.filterList = items;
  }
}

const getUrl = (url: string): string => {
  return import.meta.env.VITE_API_BASEPATH as string + url
}

const getActiveFilter = (item: GetImageFilterListResultfilter): boolean => {
  if (viewerObject.currentFilter) {
    return viewerObject.currentFilter.date === item.date
  }
  return false
}

const { emitter } = useEmitt()
const select = (image: ApiImage) => {
  if (image) {
    if (viewerObject.selected && viewerObject.selected.id === image.id) {
      viewerObject.selected = undefined
    }
    viewerObject.selected = image
  } else {
    viewerObject.selected = undefined
  }
  const output = Object.assign({}, unref(viewerObject)?.selected) as ApiImage
  output.url = output.url.replace(import.meta.env.VITE_API_BASEPATH,'');
  emitter.emit('imageSelected', {id: props.id, image: output})
  //todo: fix
  ElMessage({
    message: t('message.selectedImage') + ` ${output.id}`,
    type: 'success',
    duration: 2000
  })
}

const removeFromServer = async (image: ApiImage) => {
  await api.v1.imageServiceDeleteImageById(image.id)
  // await getList(viewerObject.currentFilter)
  await getFilterList()
}

const onSuccess: UploadProps['onSuccess'] = (image: ApiImage, uploadFile) => {
  getFilterList()
  getList(viewerObject.currentFilter)
  ElMessage({
    message: t('message.uploadSuccessfully'),
    type: 'success',
    duration: 2000
  })
}

const getUploadURL = () => {
  const uri = import.meta.env.VITE_API_BASEPATH as string || window.location.origin;
  const accessToken = wsCache.get("accessToken")
  return uri + '/v1/image/upload?access_token=' + accessToken;
}

const handleRemove: UploadProps['onRemove'] = ( image: ApiImage, uploadFiles) => {
  removeFromServer(image)
}

const handlePictureCardPreview: UploadProps['onPreview'] = (image) => {
  createImageViewer({
    urlList: [
      image.url!
    ]
  })

  select(image)

}

fetch()

</script>

<template>
  <ElRow class="file-manager-body">
    <ElCol :span="6" :xs="24" class="mb-20px">
      <ul class="list-unstyled filters">
        <li v-for="item in viewerObject.filterList?.slice().reverse()" :key="item.date">
          <ElBadge :value="item.count" class="item" type="info">
            <ElButton :link="true" @click.prevent.stop="getList(item)">
              {{ item.date }}
            </ElButton>
          </ElBadge>
        </li>
      </ul>
    </ElCol>
    <ElCol :span="18" :xs="24">

      <ElUpload
          v-model:file-list="viewerObject.imageList"
          list-type="picture-card"
          :multiple="true"
          ref="upload"
          :action="getUploadURL()"
          :on-success="onSuccess"
          :on-preview="handlePictureCardPreview"
          :on-remove="handleRemove"
          :auto-upload="true">
        <Icon icon="ic:baseline-plus" />
      </ElUpload>

    </ElCol>
  </ElRow>

</template>

<style lang="less">

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
