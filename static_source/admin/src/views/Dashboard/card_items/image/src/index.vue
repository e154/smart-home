<script setup lang="ts">
import {onMounted, onUnmounted, PropType, ref} from "vue";
import {CardItem} from "@/views/Dashboard/core";
import {RenderVar} from "@/views/Dashboard/render";
import {ElIcon, ElImage} from "element-plus";
import {Picture as IconPicture} from '@element-plus/icons-vue'
import {useCache} from "@/hooks/web/useCache";
import {prepareUrl} from "@/utils/serverId";

const {wsCache} = useCache()

// ---------------------------------
// common
// ---------------------------------
const props = defineProps({
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const el = ref(null)
onMounted(() => {
  // store dom element moveable
  props.item.setTarget(el.value)
})

onUnmounted(() => {

})

// ---------------------------------
// component methods
// ---------------------------------

const getUrl = (): string => {
  if (!props.item?.payload?.image) {
    return '';
  }
  if (props.item?.payload.image.attrField) {
    const imageUrl = RenderVar(props.item?.payload.image.attrField, props.item?.lastEvent);
    return prepareUrl(import.meta.env.VITE_API_BASEPATH as string + imageUrl);
  }
  return prepareUrl(import.meta.env.VITE_API_BASEPATH as string + props.item?.payload.image?.image?.url || '');
}

const getTileStyle = () => {
  const uri = getUrl();
  return {
    "background": `url(${uri})`,
  }
}
</script>

<template>
  <div ref="el" class="w-[100%] h-[100%] overflow-hidden">
    <ElImage v-if="!item.payload.image.background" v-show="!item.hidden" :src="getUrl()">
      <template #error>
        <div class="image-slot">
          <ElIcon>
            <icon-picture/>
          </ElIcon>
        </div>
      </template>
    </ElImage>
    <div v-else :style="getTileStyle()" class="w-[100%] h-[100%]"></div>
  </div>

</template>

<style lang="less">
.el-image__error, .el-image__placeholder, .el-image__inner {
  height: auto;
}

.el-image.item-element {
  overflow: visible;
}
</style>
