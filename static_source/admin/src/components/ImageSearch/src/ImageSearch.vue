<script setup lang="ts">

import {onMounted, PropType, ref, unref, watch} from "vue";
import {ApiImage} from "@/api/stub";
import {ElButton, ElCol, ElIcon, ElImage, ElRow} from 'element-plus'
import Browser from "@/views/Images/components/Browser.vue";
import {useI18n} from "@/hooks/web/useI18n";
import {UUID} from "uuid-generator-ts";
import {GetFullImageUrl} from "@/utils/serverId";
import {DraggableContainer} from "@/components/DraggableContainer";
import {CloseBold} from "@element-plus/icons-vue";
import {useEventBus} from "@/hooks/event/useEventBus";

const emit = defineEmits(['change', 'update:modelValue'])
const {t} = useI18n()

const props = defineProps({
  modelValue: {
    type: Object as PropType<Nullable<ApiImage>>,
    default: () => undefined
  }
})

const dialogVisible = ref(false)
const currentImage = ref<Nullable<ApiImage>>(null)

const currentID = ref('')
onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()
})

watch(
  () => props.modelValue,
  (val?: ApiImage) => {
    if (val === unref(currentImage)) return;
    if (!val?.url) {
      currentImage.value = null
      return;
    }
    ;
    currentImage.value = unref(val) || null
  },
  {
    immediate: true
  }
)

// 监听
watch(
  () => currentImage.value,
  (val?: ApiImage) => {
    emit('update:modelValue', unref(val) || null)
    emit('change', unref(val) || null)
  }
)

const remove = () => {
  currentImage.value = null
}

useEventBus({
  name: 'keydown',
  callback: (event) => {
    const {key} = event
    if (key === 'Escape') {
      dialogVisible.value = false
    }
  }
})

const imageSelected = (event: { id: number, image: ApiImage }) => {
  const {id, image} = event
  if (id && id != currentID.value) return;
  currentImage.value = image || null
}

const showBrowser = () => {
  dialogVisible.value = true
}

</script>

<template>
  <ElRow class="w-[100%]">
    <ElCol class="w-[100%]">
      <div v-if="currentImage" class="image-preview">
        <ElImage class="w-[100%]" :src="GetFullImageUrl(currentImage)"
                 :preview-src-list="[GetFullImageUrl(currentImage)]"/>
        <a href="#" class="cross delete-btn" @click.prevent.stop="remove()"></a>
      </div>
      <div v-else>
        <ElButton style="width: 100%" @click.prevent.stop="showBrowser()">
          <Icon icon="ph:upload"/>
        </ElButton>
      </div>

      <DraggableContainer
        v-if="dialogVisible"
        class-name="image-browser-modal"
        name="modal-image-browser"
        :initial-width="1024"
        :initial-height="600"
        :modal="true"
      >
        <template #header>
          <div class="w-[100%]">
            <div style="float: left">{{ $t('images.imageBrowser') }}</div>
            <div style="float: right; text-align: right">
              <a href="#" @click.prevent.stop='dialogVisible= false'>
                <ElIcon class="mr-5px">
                  <CloseBold/>
                </ElIcon>
              </a>
            </div>
          </div>
        </template>
        <template #default>
          <Browser :id="currentID" @imageSelected="imageSelected" :select-mode="true"/>
        </template>
      </DraggableContainer>
    </ElCol>
  </ElRow>

</template>

<style lang="less">

.image-preview {

  max-height: 60px;
  min-height: 60px;
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

.image-browser-modal {
  .draggable-container-content {
    background-color: var(--el-bg-color);
  }
}

</style>
