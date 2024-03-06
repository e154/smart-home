<script setup lang="ts">
import {computed, onMounted, onUnmounted, reactive, ref, unref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {ElButton, ElEmpty, ElMessage, ElTabPane, ElTabs} from 'element-plus'
import {useRoute} from 'vue-router'
import api from "@/api/api";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";
import {Card, Core, eventBus, EventContextMenu, Tab} from "@/views/Dashboard/core";
import ViewTab from "@/views/Dashboard/editor/ViewTab.vue";
import {DraggableContainer} from "@/components/DraggableContainer";
import TabSettings from "@/views/Dashboard/editor/TabSettings.vue";
import TabEditor from "@/views/Dashboard/editor/TabEditor.vue";
import TabCardItem from "@/views/Dashboard/editor/TabCardItem.vue";
import TabCard from "@/views/Dashboard/editor/TabCard.vue";
import {EventStateChange} from "@/api/types";
import {useAppStore} from "@/store/modules/app";
import {GetFullImageUrl} from "@/utils/serverId";
import {SecondMenu} from "@/views/Dashboard/core/src/secondMenu";
import {JsonEditor} from "@/components/JsonEditor";
import {Dialog} from "@/components/Dialog";
import {ApiDashboardTab} from "@/api/stub";
import CardListWindow from "@/views/Dashboard/editor/CardListWindow.vue";

const route = useRoute();
const {t} = useI18n()
const appStore = useAppStore()

// ---------------------------------
// common
// ---------------------------------

const loading = ref(true)
const dashboardId = computed(() => parseInt(route.params.id as string) as number);
const core = reactive<Core>(new Core());
const currentID = ref('')

// context menu
const contextMenu = reactive<SecondMenu>(new SecondMenu(unref(core)));

const onStateChanged = (event: EventStateChange) => {
  eventBus.emit('stateChanged', event);
  core.onStateChanged(event);
}

const eventHandler = (event: string, args: any[]) => {
  switch (event) {
    case 'showTabImportDialog':
      importDialogVisible.value = true
      break;
    case 'fetchDashboard':
      fetchDashboard()
      break;
  }
}

const eventBusHandler =  (event: string, args: any[]) => {
  core.eventBusHandler(event, args)
}

onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  fetchDashboard()

  stream.subscribe('state_changed', currentID.value, onStateChanged);
  eventBus.subscribe(['showTabImportDialog', 'fetchDashboard'], eventHandler)
  eventBus.subscribe(undefined, eventBusHandler)
  eventBus.subscribe('eventContextMenu', contextMenu.eventHandler)
})

onUnmounted(() => {
  core.shutdown()

  eventBus.unsubscribe(['showTabImportDialog', 'fetchDashboard'], eventHandler)
  stream.unsubscribe('state_changed', currentID.value);
  eventBus.unsubscribe(undefined, eventBusHandler)
  eventBus.unsubscribe('eventContextMenu', contextMenu.eventHandler)
})

// ---------------------------------
// dashboard
// ---------------------------------

const fetchDashboard = async () => {
  loading.value = true;
  const res = await api.v1.dashboardServiceGetDashboardById(dashboardId.value)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false;
      })
  core.currentBoard(res.data);
}

// ---------------------------------
// tabs
// ---------------------------------

const updateCurrentTab = (tab: any, ev: any) => {
  const {index} = tab;
  if (core.activeTabIdx === index) return;
  core.activeTabIdx = index;
  core.updateCurrentTab();
}

const activeTabIdx = computed({
  get(): string {
    return core.activeTabIdx + ''
  },
  set(value: string) {
    core.activeTabIdx = parseInt(value)
  }
})

const activeTab = computed<Tab>(() => core.getActiveTab as Tab)
const activeCard = computed<Card>(() => core.getActiveTab.cards[core.activeCard] as Card)

const getTabStyle = () => {
  const style = {
    margin: 0
  }
  if (core.getActiveTab?.background) {
    style['background-color'] = core.getActiveTab?.background
  } else {
    if (core.getActiveTab?.backgroundAdaptive) {
      style['background-color'] = appStore.isDark ? '#333335' : '#FFF'
    }
  }

  if (activeTab.value?.backgroundImage) {
    style['background-image'] = `url(${GetFullImageUrl(activeTab.value.backgroundImage)})`
    style['background-repeat'] = 'repeat';
    style['background-position'] = 'center';
    // style['background-size'] = 'cover';
  }
  return style
}

const tagsView = computed(() => tagsView.value ? 37 : 0)

const createTab = async () => {
  await core.createTab();

  ElMessage({
    title: t('Success'),
    message: t('message.createdSuccessfully'),
    type: 'success',
    duration: 2000
  });
}

const addCard = () => {
  core.createCard();
}

const toggleMenu = (menu: string): void => {
  switch (menu) {
    case 'tabs':
      eventBus.emit('toggleTabsMenu');
      break
    case 'cards':
      eventBus.emit('toggleCardsMenu');
      break
    case 'cardItems':
      eventBus.emit('toggleCardItemsMenu');
      break
  }
}

const onContextMenu = (e: MouseEvent, owner: 'editor' | 'tab', tabId?: number) => {
  e.preventDefault();
  e.stopPropagation();
  eventBus.emit('eventContextMenu', {
    event: e,
    owner: owner,
    tabId: tabId,
  } as EventContextMenu)
}

// ---------------------------------
// import/export
// ---------------------------------

const importedTab = ref(null)
const importDialogVisible = ref(false)

const importHandler = (val: any) => {
  if (importedTab.value == val) {
    return
  }
  importedTab.value = val
}

const importTab = async () => {
  let card: ApiDashboardTab
  try {
    if (importedTab.value?.json) {
      card = importedTab.value.json as ApiDashboardTab;
    } else if (importedTab.value.text) {
      card = JSON.parse(importedTab.value.text) as ApiDashboardTab;
    }
  } catch {
    ElMessage({
      title: t('Error'),
      message: t('message.corruptedJsonFormat'),
      type: 'error',
      duration: 2000
    });
    return
  }
  const res = await core.importTab(card);
  if (res) {
    ElMessage({
      title: t('Success'),
      message: t('message.importedSuccessful'),
      type: 'success',
      duration: 2000
    })
  }
  importDialogVisible.value = false
}

defineOptions({
  inheritAttrs: false
})
</script>

<template>

  <div class="dashboard-container"
       v-if="!loading"
       :style="getTabStyle()"
       @contextmenu="onContextMenu($event, 'editor', undefined)">

    <ElTabs
        v-model="activeTabIdx"
        @tab-click="updateCurrentTab"
        class="ml-20px"
        :lazy="true">
      <ElTabPane
          v-for="(tab, index) in core.tabs"
          :label="tab.name"
          :key="index"
          :disabled="!tab.enabled"
          :class="[{'gap': tab.gap}]"
          :lazy="true"
          @contextmenu="onContextMenu($event, 'tab', tab.id)"
      >
        <ViewTab :tab="tab" :key="index" :core="core"/>
      </ElTabPane>
    </ElTabs>

    <!-- main menu -->
    <DraggableContainer :name="'editor-main'">
      <template #header>
        <div class="w-[100%]">
          <div style="float: left">Main menu</div>
          <div style="float: right; text-align: right">
            <a href="#" @click.prevent.stop='toggleMenu("tabs")'>
              <Icon icon="vaadin:tabs" class="mr-5px" @click.prevent.stop='toggleMenu("tabs")'/>
            </a>
            <a href="#" class="mr-5px" @click.prevent.stop='toggleMenu("cards")'>
              <Icon icon="material-symbols:cards-outline"/>
            </a>
            <a href="#" @click.prevent.stop='toggleMenu("cardItems")'>
              <Icon icon="icon-park-solid:add-item"/>
            </a>
          </div>
        </div>
      </template>
      <template #default>
        <ElTabs v-model="core.mainTab">
          <!-- main -->
          <ElTabPane :label="$t('dashboard.mainTab')" name="main">
            <template #label>
              <Icon icon="wpf:maintenance"/>
            </template>
            <TabSettings v-if="core.current" :core="core"/>
          </ElTabPane>
          <!-- /main -->

          <!-- tabs -->
          <ElTabPane :label="$t('dashboard.tabsTab')" name="tabs">
            <template #label>
              <Icon icon="vaadin:tabs"/>
            </template>
            <TabEditor v-if="core.current && activeTab" :tab="activeTab" :core="core"/>
            <ElEmpty v-if="!core.tabs.length" :rows="5">
              <ElButton type="primary" @click="createTab()">
                {{ t('dashboard.addNewTab') }}
              </ElButton>
            </ElEmpty>
          </ElTabPane>
          <!-- /tabs -->

          <!-- cards -->
          <ElTabPane :label="$t('dashboard.cardsTab')" name="cards">
            <template #label>
              <Icon icon="material-symbols:cards-outline"/>
            </template>
            <TabCard v-if="core.current && activeTab" :tab="activeTab" :core="core"/>
            <ElEmpty v-if="!core.tabs.length" :rows="5">
              <ElButton type="primary" @click="createTab()">
                {{ t('dashboard.addNewTab') }}
              </ElButton>
            </ElEmpty>
          </ElTabPane>
          <!-- /cards -->

          <!-- cardItems -->
          <ElTabPane :label="$t('dashboard.cardItemsTab')" name="cardItems">
            <template #label>
              <Icon icon="icon-park-solid:add-item"/>
            </template>
            <TabCardItem v-if="core.current && activeTab && activeCard" :card="activeCard" :core="core"/>
            <ElEmpty v-if="!core.tabs.length" :rows="5">
              <ElButton type="primary" @click="createTab()">
                {{ t('dashboard.addNewTab') }}
              </ElButton>
            </ElEmpty>
            <ElEmpty v-if="core.tabs.length && !(core.activeCard >= 0)" :rows="5">
              <ElButton type="primary" @click="addCard()">
                {{ t('dashboard.addNewCard') }}
              </ElButton>
            </ElEmpty>
          </ElTabPane>
          <!-- /cardItems -->

        </ElTabs>
      </template>
    </DraggableContainer>
    <!-- /main menu -->

    <!-- card list window -->
    <CardListWindow v-if="core.current && activeTab" :core="core"/>
    <!-- /card list window -->

    <!-- import dialog -->
    <Dialog v-model="importDialogVisible" :title="t('main.dialogImportTitle')" :maxHeight="400" width="80%"
            custom-class>
      <JsonEditor @change="importHandler"/>
      <template #footer>
        <ElButton type="primary" @click="importTab()" plain>{{ t('main.import') }}</ElButton>
        <ElButton @click="importDialogVisible = false">{{ t('main.closeDialog') }}</ElButton>
      </template>
    </Dialog>
    <!-- /import dialog -->

  </div>

</template>

<style lang="less">


/* Track */
::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.dashboard-container {
  position: relative;
  min-height: calc(100vh - 87px);
}

p {
  display: block;
  margin-block-start: 1em;
  margin-block-end: 1em;
  margin-inline-start: 0;
  margin-inline-end: 0;
}

h1 {
  display: block;
  font-size: 2em;
  margin-block-start: 0.67em;
  margin-block-end: 0.67em;
  margin-inline-start: 0;
  margin-inline-end: 0;
  font-weight: 700;
}

h2 {
  display: block;
  font-size: 1.5em;
  margin-block-start: 0.67em;
  margin-block-end: 0.67em;
  margin-inline-start: 0;
  margin-inline-end: 0;
  font-weight: 700;
}

html {
  line-height: 1.15;
}

.el-tabs {
  height: inherit;
  height: -webkit-fill-available;
}

html.dark {
  .draggable-container {
    &.container-editor-main {


      .el-card {
        .el-divider__text {
          background-color: var(--el-bg-color-overlay);
        }
      }
    }

    &.container-editor-cards,
    &.container-editor-tabs,
    &.container-editor-card-items,
    &.container-editor-main,
    &.container-frame-editor {
      .draggable-container-content,
      .el-divider__text {
        background-color: hsl(230, 7%, 17%);
      }
    }
  }

}

// custom style
.draggable-container.container-editor-main {
  .el-main {
    padding: 2px !important;
  }

  .el-card__header {
    padding: 18px 20px !important;
  }

  .el-card {
    --el-card-padding: 2px 5px;
  }

  .el-form-item--small {
    margin-bottom: 5px;
  }

  .el-divider--horizontal {
    margin: 11px 0;
  }

  .el-col.el-col-12 {
    padding-right: 6px;
    padding-left: 6px;
  }

  .el-menu-item {
    padding: 0 2px;
    line-height: 14px !important;
    height: 14px !important;
  }

  .el-menu--vertical:not(.el-menu--collapse):not(.el-menu--popup-container) .el-menu-item, .el-menu--vertical:not(.el-menu--collapse):not(.el-menu--popup-container) .el-menu-item-group__title, .el-menu--vertical:not(.el-menu--collapse):not(.el-menu--popup-container) .el-sub-menu__title {
    padding-left: 2px;
  }

  .el-col.el-col-24.is-guttered {
    padding: 0 !important;
  }

  .el-button {
    margin-bottom: 10px !important;
  }

  .el-collapse-item__content {
    padding-bottom: 10px !important;
  }
}

.container-editor-main {
  .draggable-container-content {
    padding-top: 0;

  }
}

.draggable-container {
  &.container-editor-cards,
  &.container-editor-tabs,
  &.container-editor-card-items,
  &.container-editor-main,
  &.container-frame-editor {
    .draggable-container-header {
      font-size: 12px;
    }
  }
}

.draggable-container {
  &.container-editor-cards,
  &.container-editor-tabs,
  &.container-editor-card-items,
  &.container-frame-editor {
    .el-menu-item {
      padding-left: 5px !important;
      padding-right: 5px !important;
      line-height: 30px;
      height: 30px;
      font-size: 12px;
    }

    .el-menu-item * {
      vertical-align: baseline;
    }
  }
}

// menu
.menu-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.buttons {
  display: none;
  position: absolute;
  right: 0;
  background: var(--el-bg-color);
}

.el-menu-item:hover .buttons {
  display: block;
  color: red;
}
</style>
