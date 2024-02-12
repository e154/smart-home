<script setup lang="ts">
import {onMounted, onUnmounted, PropType, ref} from "vue";
import {CardItem} from "@/views/Dashboard/core/core";
import {RenderVar} from "@/views/Dashboard/core/render";
import {ElIcon, ElImage} from "element-plus";
import {Picture as IconPicture} from '@element-plus/icons-vue'
import {useCache} from "@/hooks/web/useCache";
import {GetFullUrl} from "@/utils/serverId";

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

const getTileUrl = (): string => {
  if (!props.item?.payload?.image) {
    return '';
  }
  if (props.item?.payload.image.attrField) {
    const imageUrl = RenderVar(props.item?.payload.image.attrField, props.item?.lastEvent);
    return GetFullUrl(imageUrl);
  }
  return GetFullUrl(props.item?.payload.image?.image?.url);
}

const getTileStyle = () => {
  const uri = getTileUrl();
  return {
    "background": `url(${uri})`,
  }
}
</script>

<template>
  <div ref="el" class="w-[100%] h-[100%] overflow-hidden">
    <ElImage v-if="!item.payload.image.background" v-show="!item.hidden" :src="getTileUrl()">
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
