<script setup lang="ts">

import {onMounted, PropType, ref, unref, watch} from "vue";
import {ApiImage} from "@/api/stub";
import {ElButton, ElCol, ElDialog, ElImage, ElRow} from 'element-plus'
import Browser from "@/views/Images/components/Browser.vue";
import {useEmitt} from "@/hooks/web/useEmitt";
import {useI18n} from "@/hooks/web/useI18n";
import {UUID} from "uuid-generator-ts";
import {GetFullImageUrl} from "@/utils/serverId";

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
      if (val === unref(currentImage)) return
      currentImage.value = val || null
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
    }
)

const remove = () => {
  currentImage.value = null
}

useEmitt({
  name: 'imageSelected',
  callback: (val) => {
    const {id, image} = val;
    if (id && id != currentID.value) return;
    currentImage.value = image || null
  }
})

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

      <ElDialog v-model="dialogVisible" append-to-body :title="$t('images.imageBrowser')" :maxHeight="400" width="80%">
        <Browser :id="currentID"/>
      </ElDialog>
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


</style>
