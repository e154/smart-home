<script setup lang="ts">
import {computed, onMounted, onUnmounted, PropType, ref} from 'vue'
import {ElButton, ElButtonGroup, ElIcon, ElMenu, ElMenuItem,} from 'element-plus'
import {CloseBold} from '@element-plus/icons-vue'
import {Card, Core, Tab, eventBus} from "@/views/Dashboard/core";
import {DraggableContainer} from "@/components/DraggableContainer";

const props = defineProps({
  core: {
    type: Object as PropType<Nullable<Core>>,
  },
})

const currentCore = computed(() => props.core as Core)

const activeTab = computed({
  get(): Tab {
    return currentCore.value.getActiveTab as Tab
  },
  set(val: Tab) {
  }
})

// ---------------------------------
// common
// ---------------------------------

const onSelectedCard = (id: number) => {
  currentCore.value.onSelectedCard(id);
  eventBus.emit('unselectedCardItem')
}

const menuCardsClick = (card) => {
  onSelectedCard(card.id)
}

const sortCardUp = (card: Card, index: number) => {
  activeTab.value.sortCardUp(card, index)
  eventBus.emit('updateGrid', activeTab.value.id)
}

const sortCardDown = (card: Card, index: number) => {
  activeTab.value.sortCardDown(card, index)
  eventBus.emit('updateGrid', activeTab.value.id)
}

const showMenuWindow = ref(false)
const eventHandler = () => {
  showMenuWindow.value = !showMenuWindow.value
}
onMounted(() => {
  eventBus.subscribe('toggleCardsMenu', eventHandler)
})

onUnmounted(() => {
  eventBus.unsubscribe('toggleCardsMenu', eventHandler)
})

</script>

<template>

  <DraggableContainer :name="'editor-cards'" :initial-width="280" :min-width="280" v-show="showMenuWindow">
    <template #header>
      <div class="w-[100%]">
        <div style="float: left">Cards</div>
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

      <ElMenu v-if="currentCore.activeTabIdx > -1 && activeTab.cards.length"
              :default-active="currentCore.activeCard + ''" v-model="currentCore.activeCard"
              class="el-menu-vertical-demo">
        <ElMenuItem :index="index + ''" :key="index" v-for="(card, index) in activeTab.cards"
                    @click="menuCardsClick(card)">
          <div class="w-[100%] menu-item">
            <span>{{ card.title }}</span>
            <ElButtonGroup class="buttons">
              <ElButton @click.prevent.stop="sortCardUp(card, index)" text size="small">
                <Icon icon="teenyicons:up-solid"/>
              </ElButton>
              <ElButton @click.prevent.stop="sortCardDown(card, index)" text size="small">
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
