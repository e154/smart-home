<script setup lang="ts">
import {useDesign} from '@/hooks/web/useDesign'
import {computed, nextTick, onBeforeUnmount, onMounted, ref, unref, watch} from 'vue'
import Cropper from 'cropperjs'
import 'cropperjs/dist/cropper.min.css'
import {ElDescriptions, ElDescriptionsItem, ElMessage, ElTooltip, UploadFile} from 'element-plus'
import {useDebounceFn} from '@vueuse/core'
import {BaseButton} from '@/components/Button'

const {getPrefixCls} = useDesign()

const prefixCls = getPrefixCls('image-cropping')

const props = defineProps({
  imageUrl: {
    type: String,
    default: '',
    required: true
  },
  cropBoxWidth: {
    type: Number,
    default: 200
  },
  cropBoxHeight: {
    type: Number,
    default: 200
  },
  boxWidth: {
    type: [Number, String],
    default: 425
  },
  boxHeight: {
    type: [Number, String],
    default: 320
  },
  showResult: {
    type: Boolean,
    default: true
  },
  showActions: {
    type: Boolean,
    default: true
  },
  data: {
    type: Object,
    default: () => null
  }
})

const getBase64 = useDebounceFn(() => {
  imgBase64.value = unref(cropperRef)?.getCroppedCanvas()?.toDataURL() ?? ''
}, 80)

const resetCropBox = async () => {
  const containerData = unref(cropperRef)?.getContainerData()
  unref(cropperRef)?.setCropBoxData({
    width: props.cropBoxWidth,
    // height: props.cropBoxHeight,
    left: (containerData?.width || 0) / 2 - 100,
    top: (containerData?.height || 0) / 2 - 100
  })
  imgBase64.value = unref(cropperRef)?.getCroppedCanvas()?.toDataURL() ?? ''

  if (props.data) {
    if (props.data?.canvasData) cropperRef.value?.setCanvasData(props.data.canvasData)
    cropperRef.value?.setData(props.data)
    await nextTick()
    // resetCropBox()
  }
}

const getBoxStyle = computed(() => {
  return {
    width: `${props.boxWidth}px`,
    height: `${props.boxHeight}px`
  }
})

const getCropBoxStyle = computed(() => {
  return {
    // width: `${props.cropBoxWidth}px`,
    width: `auto`,
    height: `auto`,
    maxHeight: `250px`,
    maxWidth: `250px`
  }
})

// 获取对应的缩小倍数的宽高
const getScaleSize = (scale: number) => {
  return {
    width: props.cropBoxWidth * scale + 'px',
    height: props.cropBoxHeight * scale + 'px'
  }
}

const imgBase64 = ref('')
const imgRef = ref<HTMLImageElement>()
const cropperRef = ref<Cropper>()

const detail = ref();
const intiCropper = async () => {
  console.log('intiCropper')
  if (!unref(imgRef)) return
  const imgEl = unref(imgRef)!
  cropperRef.value = new Cropper(imgEl, {
    responsive: true,
    scalable: true,
    zoomable: true,
    // aspectRatio: 1,
    initialAspectRatio: 1,
    viewMode: 1,
    dragMode: 'move',
    // cropBoxResizable: false,
    // cropBoxMovable: false,
    toggleDragModeOnDblclick: false,
    checkCrossOrigin: false,
    // zoomable: false,
    ready() {
      resetCropBox()
    },
    cropmove() {
      getBase64()
    },
    cropstart() {

    },
    zoom() {
      getBase64()
    },
    crop(e) {
      const data = e.detail;
      detail.value = {
        x: Math.round(data.x),
        y: Math.round(data.y),
        height: Math.round(data.height),
        width: Math.round(data.width),
      }
      getBase64()
    }
  })
  cropperRef.value?.getCroppedCanvas({
    imageSmoothingEnabled: false,
    imageSmoothingQuality: 'low'
  })

}

const uploadChange = (uploadFile: UploadFile) => {
  // 判断是否是图片
  if (uploadFile?.raw?.type.indexOf('image') === -1) {
    ElMessage.error('请上传图片格式的文件')
    return
  }
  if (!uploadFile.raw) return
  // 获取图片的访问地址
  const url = URL.createObjectURL(uploadFile.raw)
  unref(cropperRef)?.replace(url)
}

const reset = () => {
  unref(cropperRef)?.reset()
}

const rotate = (deg: number) => {
  unref(cropperRef)?.rotate(deg)
}

const scaleX = ref(1)
const scaleY = ref(1)
const scale = (type: 'scaleX' | 'scaleY') => {
  if (type === 'scaleX') {
    scaleX.value = scaleX.value === 1 ? -1 : 1
    unref(cropperRef)?.[type](unref(scaleX))
  } else {
    scaleY.value = scaleY.value === 1 ? -1 : 1
    unref(cropperRef)?.[type](unref(scaleY))
  }
}

const zoom = (num: number) => {
  unref(cropperRef)?.zoom(num)
}

onMounted(() => {
  intiCropper()
})

watch(
    () => props.imageUrl,
    async (url) => {
      if (url) {
        unref(cropperRef)?.replace(url)
        await nextTick()
        resetCropBox()
      }
    }
)

watch(
    () => props.data,
    async (val) => {
      if (props.data?.canvasData) cropperRef.value?.setCanvasData(props.data.canvasData)
      cropperRef.value?.setData(props.data)
      await nextTick()
    },
    {
      deep: true,
      immediate: true
    }
)

onBeforeUnmount(() => {
  unref(cropperRef)?.destroy()
})

defineExpose({
  cropperExpose: cropperRef
})

</script>

<template>
  <div
      :class="{
      [prefixCls]: true,
      'flex items-top': showResult
    }"
  >
    <div>
      <div :style="getBoxStyle" class="flex justify-center items-center">
        <img
            v-show="imageUrl"
            ref="imgRef"
            :src="imageUrl"
            class="block max-w-full"
            crossorigin="anonymous"
            alt=""
            srcset=""
        />
      </div>
      <div v-if="showActions" class="mt-10px flex items-center">
        <!--        <div class="flex items-center">-->
        <!--          <ElTooltip content="选择文件" placement="bottom">-->
        <!--            <ElUpload-->
        <!--                action="''"-->
        <!--                accept="image/*"-->
        <!--                :auto-upload="false"-->
        <!--                :show-file-list="false"-->
        <!--                :on-change="uploadChange"-->
        <!--            >-->
        <!--              <BaseButton size="small" type="primary" class="mt-2px"-->
        <!--              ><Icon icon="ep:upload-filled"-->
        <!--              /></BaseButton>-->
        <!--            </ElUpload>-->
        <!--          </ElTooltip>-->
        <!--        </div>-->
        <div class="flex items-center justify-end flex-1">
          <ElTooltip :content="$t('dashboard.editor.refresh')" placement="bottom">
            <BaseButton size="small" type="primary" @click="reset"
            >
              <Icon icon="ep:refresh"
              />
            </BaseButton>
          </ElTooltip>
          <ElTooltip :content="$t('dashboard.editor.rotateLeft')" placement="bottom">
            <BaseButton size="small" type="primary" @click="rotate(-45)"
            >
              <Icon icon="ant-design:rotate-left-outlined"
              />
            </BaseButton>
          </ElTooltip>
          <ElTooltip :content="$t('dashboard.editor.rotateRight')" placement="bottom">
            <BaseButton size="small" type="primary" @click="rotate(45)"
            >
              <Icon icon="ant-design:rotate-right-outlined"
              />
            </BaseButton>
          </ElTooltip>
          <ElTooltip :content="$t('dashboard.editor.swapX')" placement="bottom">
            <BaseButton size="small" type="primary" @click="scale('scaleX')"
            >
              <Icon icon="vaadin:arrows-long-h"
              />
            </BaseButton>
          </ElTooltip>
          <ElTooltip :content="$t('dashboard.editor.swapY')" placement="bottom">
            <BaseButton size="small" type="primary" @click="scale('scaleY')"
            >
              <Icon icon="vaadin:arrows-long-v"
              />
            </BaseButton>
          </ElTooltip>
          <ElTooltip :content="$t('dashboard.editor.zoomIn')" placement="bottom">
            <BaseButton size="small" type="primary" @click="zoom(0.1)"
            >
              <Icon icon="ant-design:zoom-in-outlined"
              />
            </BaseButton>
          </ElTooltip>
          <ElTooltip :content="$t('dashboard.editor.zoomOut')" placement="bottom">
            <BaseButton size="small" type="primary" @click="zoom(-0.1)"
            >
              <Icon icon="ant-design:zoom-out-outlined"
              />
            </BaseButton>
          </ElTooltip>
        </div>
      </div>
    </div>
    <!-- preview -->
    <div v-if="imgBase64 && showResult" class="ml-20px mr-20px w-[100%]">
      <div>
        <!-- info -->
        <ElDescriptions v-if="detail"
                        class="ml-10px mr-10px mb-20px w-[100%]"
                        :title="$t('dashboard.editor.imageCropInfo')"
                        direction="vertical"
                        :column="2"
                        border
        >
          <ElDescriptionsItem :label="$t('dashboard.editor.imageCropX')">{{ detail.x }}</ElDescriptionsItem>
          <ElDescriptionsItem :label="$t('dashboard.editor.imageCropY')">{{ detail.y }}</ElDescriptionsItem>
          <ElDescriptionsItem :label="$t('dashboard.editor.imageCropHeight')">{{ detail.height }}</ElDescriptionsItem>
          <ElDescriptionsItem :label="$t('dashboard.editor.imageCropWidth')">{{ detail.width }}</ElDescriptionsItem>
        </ElDescriptions>
        <!-- /info -->
      </div>
      <div class="flex justify-center items-center">
        <img :src="imgBase64" :style="getCropBoxStyle"/>
      </div>
    </div>
    <!-- \preview -->
  </div>
</template>
