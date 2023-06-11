<template>
  <vuuri
    v-model="tab.cards"
    item-key="id"
    :get-item-width="getItemWidth"
    :get-item-height="getItemHeight"
    @updated="onUpdated"
    @move="move"
    :drag-enabled="false"
    ref="grid"
  >
    <template #item="{item}">
      <dashboard-card :card="item" :bus="bus"/>
    </template>
  </vuuri>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import DashboardCard from './card.vue'
import { ApiDashboardCard } from '@/api/stub'
// @ts-ignore
import vuuri from 'vuuri'
import { Card, Core, Tab } from '@/views/dashboard/core'

// register globally
Vue.component('vuuri', vuuri)

@Component({
  name: 'DashboardTabMuu',
  components: {
    DashboardCard,
    vuuri
  }
})
export default class extends Vue {
  @Prop() private tab!: Tab;
  @Prop() private bus!: Vue;

  created() {
    // this.tab = this.dashboardEditor.tabs[this.tabIndex];
    this.bus.$on('update_tab', this.update)
  }

  getItemWidth(card: Card) {
    // console.log('getItemWidth', this.tab.columnWidth)
    if (card.width > 0) {
      return `${card.width}px`
    }
    return `${this.tab.columnWidth}px`
  }

  getItemHeight(item: ApiDashboardCard) {
    // console.log('getItemHeight', item.height)
    return `${item.height}px`
  }

  update(tabId: number) {
    if (this.tab.id != tabId || !this.tab.cards || this.tab.cards.length == 0) {
      return
    }
    // todo
    // TypeError: Cannot read properties of undefined (reading 'update')
    let grid = this.$refs.grid as vuuri
    grid.update()
  }

  onUpdated() {

  }

  move(event: any) {
    console.log('move', event)
    console.log(this.tab.cards)
  }
}
</script>

<style lang="scss">
/*.gap {*/
/*.muuri-item {*/
/*  padding: 5px;*/
/*.muuri-item-content {*/
/*  !*border: 1px #e9e9e9 solid;*!*/
/*  -webkit-box-shadow: 1px 1px 9px 0px #777;*/
/*  -moz-box-shadow: 1px 1px 9px 0px #777;*/
/*  box-shadow: 1px 1px 9px 0px #777;*/
/*}*/
/*}*/
/*}*/
.gap {
.muuri-item {
  padding: 5px;
.muuri-item-content {
  /*border: 1px #e9e9e9 solid;*/
  border: 1px solid #e9edf3;
}
}
}
</style>
