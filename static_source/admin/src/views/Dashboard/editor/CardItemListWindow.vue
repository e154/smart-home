<script setup lang="ts">
import {computed, onMounted, onUnmounted, PropType, ref, watch} from 'vue'
import {CloseBold} from "@element-plus/icons-vue";
import {useI18n} from '@/hooks/web/useI18n'
import {Card, CardItem, Core, eventBus} from "@/views/Dashboard/core";
import {DraggableContainer} from "@/components/DraggableContainer";
import {ElButton, ElButtonGroup, ElIcon, ElMenu, ElMenuItem, ElTag,} from 'element-plus'

const {t} = useI18n()

const cardItem = ref<CardItem>(null)

const props = defineProps({
  core: {
    type: Object as PropType<Core>,
  },
  card: {
    type: Object as PropType<Nullable<Card>>,
    default: () => null
  },
})

const activeCard = computed({
  get(): Card {
    return props.card as Card
  },
  set(val: Card) {
  }
})

const currentCore = computed({
  get(): Core {
    return props.core as Core
  },
  set(val: Core) {
  }
})

watch(
  () => props.card,
  (val?: Card) => {
    if (!val) return
    activeCard.value = val
    if (val?.selectedItem > -1) {
      cardItem.value = val?.items[val?.selectedItem] || null
    } else {
      cardItem.value = null
    }

  },
  {
    deep: true,
    immediate: true
  }
)

const menuCardItemClick = (index: number) => {
  if (currentCore.value.activeTabIdx < 0 || currentCore.value.activeCard == undefined) {
    return;
  }

  activeCard.value.selectedItem = index;

  eventBus.emit('selectedCardItem', index)
}

const sortCardItemUp = (item: CardItem, index: number) => {
  activeCard.value.sortCardItemUp(item, index)
  currentCore.value.updateCard();
}

const sortCardItemDown = (item: CardItem, index: number) => {
  activeCard.value.sortCardItemDown(item, index)
  currentCore.value.updateCard();
}

const showMenuWindow = ref(false)

const eventHandler = () => {
  showMenuWindow.value = !showMenuWindow.value
}

onMounted(() => {
  eventBus.subscribe('toggleCardItemsMenu', eventHandler)
})

onUnmounted(() => {
  eventBus.unsubscribe('toggleCardItemsMenu', eventHandler)
})

</script>

<template>
  <DraggableContainer :name="'editor-card-items'" :initial-width="280" :min-width="280" v-show="showMenuWindow">
    <template #header>
      <div class="w-[100%]">
        <div style="float: left">Card Items</div>
        <div style="float: right; text-align: right">
          <a href="#" @click.prevent.stop='showMenuWindow= false'>
            <ElIcon class="mr-5px">
              <CloseBold/>
            </ElIcon>
          </a>
        </div>
      </div>
    </template>
    <template #default>
      <ElMenu
        v-if="activeCard && activeCard.id"
        ref="tabMenu"
        :default-active="activeCard.selectedItem + ''"
        v-model="activeCard.selectedItem"
        class="el-menu-vertical-demo box-card">
        <ElMenuItem
          :index="index + ''"
          :key="index"
          v-for="(item, index) in activeCard.items"
          @click="menuCardItemClick(index)">
          <div class="w-[100%] menu-item">
                <span>
                  {{ item.title }}
                <ElTag type="info" size="small">
                  {{ item.type }}
                </ElTag>
                </span>
            <ElButtonGroup class="buttons">
              <ElButton @click.prevent.stop="sortCardItemUp(item, index)" text size="small">
                <Icon icon="teenyicons:up-solid"/>
              </ElButton>
              <ElButton @click.prevent.stop="sortCardItemDown(item, index)" text size="small">
                <Icon icon="teenyicons:down-solid"/>
              </ElButton>
            </ElButtonGroup>
          </div>
        </ElMenuItem>
      </ElMenu>

    </template>
  </DraggableContainer>

</template>

<style lang="less">

</style>
