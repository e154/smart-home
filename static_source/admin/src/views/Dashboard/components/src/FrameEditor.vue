<script setup lang="ts">

import {computed, PropType, ref, unref} from "vue";
import {Card, Core} from "@/views/Dashboard/core/core";
import {propTypes} from "@/utils/propTypes";
import {ImageSearch} from "@/components/ImageSearch";
import {ElButton, ElCol, ElDivider, ElForm, ElFormItem, ElIcon, ElRow} from "element-plus";
import {ApiImage} from "@/api/stub";
import {CloseBold} from "@element-plus/icons-vue";
import {DraggableContainer} from "@/components/DraggableContainer";
import {GetFullImageUrl} from "@/utils/serverId";
import {useI18n} from "@/hooks/web/useI18n";
import {FluentSquareHint16Filled, FluentSquareHint20Regular, ImageCropping} from "@/views/Dashboard/components";

const {t} = useI18n()

// ---------------------------------
// common
// ---------------------------------

const props = defineProps({
  core: {
    type: Object as PropType<Core>,
  },
  card: {
    type: Object as PropType<Nullable<Card>>,
    default: () => null
  },
  hover: propTypes.bool.def(false),
})

const currentCard = computed({
  get(): Card {
    return props.card as Card
  },
  set(val: Card) {
  }
})

// ---------------------------------
// component methods
// ---------------------------------

const frameMap: string[] = ['top-left-corner', 'top', 'top-right-corner', 'left', 'content', 'right', 'bottom-left-corner', 'bottom', 'bottom-right-corner']

const templateFrameSelectImage = (image: ApiImage) => {
  if (!currentCard.value?.templateFrame) {
    currentCard.value.templateFrame = {}
  }
  currentCard.value.templateFrame.image = image || undefined;
}

const currentItemName = ref(null)
const showEditorWindow = ref(false)
const showFrameEditor = (item: string) => {
  if (!currentCard.value.templateFrame) {
    currentCard.value.templateFrame = {
      items: {}
    }
  }
  if (!currentCard.value.templateFrame.items) {
    currentCard.value.templateFrame.items = {}
  }
  if (!currentCard.value.templateFrame.items[item]) {
    currentCard.value.templateFrame.items[item] = {}
  }

  if (currentItemName.value == item) {
    showEditorWindow.value = !showEditorWindow.value;
  } else {
    showEditorWindow.value = true;
    currentItemName.value = item
  }
}

const cropperExpose = ref<InstanceType<typeof ImageCropping>>()
const updateItem = () => {
  if (!currentItemName.value) {
    return
  }
  const data: Cropper.SetDataOptions = unref(cropperExpose)?.cropperExpose?.getData()
  const canvasData: Cropper.CanvasData = unref(cropperExpose)?.cropperExpose?.getCanvasData()
  currentCard.value.templateFrame.items[currentItemName.value] = {
    x: Math.round(data.x),
    y: Math.round(data.y),
    height: Math.round(data.height),
    width: Math.round(data.width),
    base64: !currentItemName.value.includes('-') ? unref(cropperExpose)?.cropperExpose?.getCroppedCanvas()?.toDataURL() ?? '' : undefined,
    canvasData: canvasData,
  }
}

const clear = () => {
  if (!currentItemName.value) {
    return
  }
  currentCard.value.templateFrame.items[currentItemName.value] = undefined
}

const currentItemData = computed(() => currentItemName.value ? currentCard.value.templateFrame?.items[currentItemName.value] : {})
const currentImage = computed(() => GetFullImageUrl(currentCard.value.templateFrame?.image))


</script>

<template>
  <ElRow class="mb-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.editor.CardTemplate') }}</ElDivider>
    </ElCol>
  </ElRow>

  <ElForm
    label-position="top"
    style="width: 100%"
  >
    <ElFormItem v-if="currentCard.template" :label="$t('dashboard.editor.image')" prop="image">
      <ImageSearch v-model="currentCard.templateFrame.image" @change="templateFrameSelectImage"/>
    </ElFormItem>

  </ElForm>

  <div class="frame-map" v-if="currentCard.templateFrame.image">
    <div class="frame-item" v-for="item in frameMap" :key="item" @click.prevent.stop="showFrameEditor(item)">
      <FluentSquareHint16Filled v-if="currentItemName == item"/>
      <FluentSquareHint20Regular v-else/>
    </div>
  </div>

  <DraggableContainer :name="'frame-editor'"
                      :initial-width="280"
                      :min-width="280"
                      :max-height="600"
                      :max-width="800"
                      :can-transparent="false"
                      v-if="currentCard.templateFrame.image"
                      v-show="showEditorWindow">
    <template #header>
      <div class="w-[100%]">
        <div style="float: left">Frame editor</div>
        <div style="float: right; text-align: right">
          <a href="#" @click.prevent.stop='showEditorWindow= false'>
            <ElIcon class="mr-5px">
              <CloseBold/>
            </ElIcon>
          </a>
        </div>
      </div>
    </template>
    <template #default>
      <ElRow class="mt-20px ml-10px">
        <ElCol>
          <ImageCropping
            v-if="currentItemName"
            ref="cropperExpose"
            :crop-box-width="300"
            :crop-box-height="300"
            :box-width="450"
            :box-height="450"
            :show-actions="true"
            :data="currentItemData"
            :canvas-data="currentItemData"
            :image-url="currentImage"
          />
        </ElCol>
      </ElRow>

      <ElRow class="mb-10px mt-10px mr-10px">
        <ElCol>
          <div class="text-right">
            <ElButton type="primary" @click.prevent.stop="updateItem" plain>{{ $t('main.update') }}</ElButton>
                        <ElButton @click.prevent.stop="clear" plain>{{ t('dashboard.editor.clear') }}</ElButton>
          </div>
        </ElCol>
      </ElRow>

    </template>
    <template #footer>

    </template>
  </DraggableContainer>

</template>

<style scoped lang="less">
.frame-map {
  width: 150px;
  height: 150px;
  position: relative;
  margin: 20px;
}

.frame-item {
  cursor: pointer;
  float: left;
  width: 45px;
  height: 45px;

  svg {
    width: 40px;
    height: 40px;
  }

  &:hover {
    svg {
      transform: scale(1.1);
    }
  }
}

</style>
