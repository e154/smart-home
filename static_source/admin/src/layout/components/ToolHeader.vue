<script lang="tsx">
import {defineComponent, computed, ref} from 'vue'
import { Collapse } from '@/components/Collapse'
import { LocaleDropdown } from '@/components/LocaleDropdown'
import { SizeDropdown } from '@/components/SizeDropdown'
import { UserInfo } from '@/components/UserInfo'
import { Screenfull } from '@/components/Screenfull'
import { GateToggle } from '@/components/Gate'
import { InstallPWA } from '@/components/InstallPWA'
import { NetworkStatus } from '@/components/NetworkStatus'
import { TerminalToggle } from '@/components/Terminal'
import { Breadcrumb } from '@/components/Breadcrumb'
import { useAppStore } from '@/store/modules/app'
import { useDesign } from '@/hooks/web/useDesign'

const { getPrefixCls, variables } = useDesign()
const prefixCls = getPrefixCls('tool-header')
const appStore = useAppStore()

const breadcrumb = computed(() => appStore.getBreadcrumb)
const hamburger = computed(() => appStore.getHamburger)
const screenfull = computed(() => appStore.getScreenfull)
const serverId = computed(() => appStore.getServerId)
const size = computed(() => appStore.getSize)
const layout = computed(() => appStore.getLayout)
const locale = computed(() => appStore.getLocale)
// const standalone = computed(() => appStore.getStandalone)

// const showInstallMenu = ref();
// window.addEventListener("beforeinstallprompt", (e) => {
//   e.preventDefault();
//   showInstallMenu.value = true;
// });
//
// const appInstalled = ref();
// window.addEventListener("appinstalled", () => {
//   appInstalled.value = true
// });

export default defineComponent({
  name: 'ToolHeader',
  setup() {
    return () => (
      <div
        id={`${variables.namespace}-tool-header`}
        class={[
          prefixCls,
          'h-[var(--top-tool-height)] relative px-[var(--top-tool-p-x)] flex items-center justify-between',
          'dark:bg-[var(--el-bg-color)]'
        ]}
      >
        {layout.value !== 'top' ? (
          <div class="h-full flex items-center">
            {hamburger.value && layout.value !== 'cutMenu' ? (
              <Collapse class="hover-trigger" color="var(--top-header-text-color)"></Collapse>
            ) : undefined}
            {breadcrumb.value ? <Breadcrumb class="<md:hidden"></Breadcrumb> : undefined}
          </div>
        ) : undefined}
        <div class="h-full flex items-center">
            <TerminalToggle class="hover-trigger" color="var(--top-header-text-color)"></TerminalToggle>
          {serverId.value ? (
            <GateToggle class="hover-trigger" color="var(--top-header-text-color)"></GateToggle>
          ) : undefined}
          {screenfull.value ? (
            <Screenfull class="hover-trigger" color="var(--top-header-text-color)"></Screenfull>
          ) : undefined}
          <NetworkStatus class="hover-trigger" color="var(--top-header-text-color)"></NetworkStatus>
          {size.value ? (
            <SizeDropdown class="hover-trigger" color="var(--top-header-text-color)"></SizeDropdown>
          ) : undefined}
          {locale.value ? (
            <LocaleDropdown
              class="hover-trigger"
              color="var(--top-header-text-color)"
            ></LocaleDropdown>
          ) : undefined}
          <UserInfo class="hover-trigger"></UserInfo>
        </div>
      </div>
    )
  }
})
</script>

<!--  {!standalone.value && showInstallMenu.value && !appInstalled.value ? (
            <InstallPWA class="hover-trigger" color="var(--top-header-text-color)"></InstallPWA>
          ) : undefined}-->

<style lang="less" scoped>
@prefix-cls: ~'@{namespace}-tool-header';

.@{prefix-cls} {
  transition: left var(--transition-time-02);
}
</style>
