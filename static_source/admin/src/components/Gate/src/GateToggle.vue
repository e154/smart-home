<script setup lang="ts">
import {propTypes} from "@/utils/propTypes";
import {ElPopover} from 'element-plus'
import {computed} from "vue";
import {useAppStore} from "@/store/modules/app";

defineProps({
  color: propTypes.string.def('')
})

const appStore = useAppStore()

const serverId = computed(() => appStore.getServerId)

const toggle = () => {

}

</script>

<template>
  <div @click="toggle">
    <ElPopover
        placement="bottom"
        title="Server Id"
        :width="200"
        trigger="hover"
        :content="serverId">
      <template #reference>
        <Icon
            :size="18"
            icon="uil:server-network"
            :color="color"
        />
      </template>

    </ElPopover>
  </div>
</template>
