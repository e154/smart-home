<template>
  <div :class="{'has-logo': showLogo}">
    <sidebar-logo
      v-if="showLogo"
      :collapse="isCollapse"
    />
<!--    <hamburger-->
<!--      id="hamburger-container"-->
<!--      :is-active="sidebar.opened"-->
<!--      class="hamburger-container"-->
<!--      @toggle-click="toggleSideBar"-->
<!--      :style="{background: variables.menuBg, color: variables.menuText}"-->
<!--    />-->
    <el-scrollbar wrap-class="scrollbar-wrapper">
      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse"
        :background-color="variables.menuBg"
        :text-color="variables.menuText"
        :active-text-color="menuActiveTextColor"
        :unique-opened="false"
        :collapse-transition="false"
        mode="vertical"
      >
        <sidebar-item
          v-for="route in routes"
          :key="route.path"
          :item="route"
          :base-path="route.path"
          :is-collapse="isCollapse"
        />
      </el-menu>
    </el-scrollbar>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import { AppModule } from '@/store/modules/app'
import { PermissionModule } from '@/store/modules/permission'
import { SettingsModule } from '@/store/modules/settings'
import SidebarItem from './SidebarItem.vue'
import SidebarLogo from './SidebarLogo.vue'
import variables from '@/styles/_variables.scss'
import { RouteConfig } from 'vue-router'
import Hamburger from '@/components/Hamburger/index.vue'

@Component({
  name: 'SideBar',
  components: {
    SidebarItem,
    SidebarLogo,
    Hamburger
  }
})
export default class extends Vue {
  @Prop({ required: true }) private routes!: RouteConfig[];

  get sidebar() {
    return AppModule.sidebar
  }

  // get routes() {
  //   return PermissionModule.routes
  // }

  private toggleSideBar() {
    AppModule.ToggleSideBar(false)
  }

  get showLogo() {
    return SettingsModule.showSidebarLogo
  }

  get menuActiveTextColor() {
    if (SettingsModule.sidebarTextTheme) {
      return SettingsModule.theme
    } else {
      return variables.menuActiveText
    }
  }

  get variables() {
    return variables
  }

  get activeMenu() {
    const route = this.$route
    const { meta, path } = route
    // if set path, the sidebar will highlight the path you set
    if (meta.activeMenu) {
      return meta.activeMenu
    }
    return path
  }

  get isCollapse() {
    return !this.sidebar.opened
  }
}
</script>

<style lang="scss">
.sidebar-container {

  .hamburger-container {
    line-height: 46px;
    /*height: 100%;*/
    /*float: left;*/

    padding: 0 15px;
    cursor: pointer;
    transition: background .3s;
    -webkit-tap-highlight-color:transparent;

  &:hover {
     background: rgba(0, 0, 0, .025)
   }
  }

// reset element-ui css
  .horizontal-collapse-transition {
    transition: 0s width ease-in-out, 0s padding-left ease-in-out, 0s padding-right ease-in-out;
  }

  .scrollbar-wrapper {
    overflow-x: hidden !important;
  }

  .el-scrollbar__view {
    height: 100%
  }

  .el-scrollbar__bar {
    &.is-vertical {
      right: 0px;
    }

    &.is-horizontal {
      display: none;
    }
  }
}
</style>

<style lang="scss" scoped>
.el-scrollbar {
  height: 100%
}

.has-logo {
  .el-scrollbar {
    height: calc(100% - 50px);
  }
}

.el-menu {
  border: none;
  height: 100%;
  width: 100% !important;
}
</style>
