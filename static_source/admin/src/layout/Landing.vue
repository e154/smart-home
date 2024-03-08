<script lang="tsx">
import {computed, defineComponent} from 'vue'
import {useAppStore} from '@/store/modules/app'
import {Backtop} from '@/components/Backtop'
import {useRenderLayout} from './components/useRenderLayout'
import {useDesign} from '@/hooks/web/useDesign'

const {getPrefixCls} = useDesign()

const prefixCls = getPrefixCls('layout')

const appStore = useAppStore()

const layout = computed(() => appStore.getLayout)

const renderLayout = () => {
  const {renderLanding} = useRenderLayout(false)
  return renderLanding()
}

export default defineComponent({
  name: 'Dashboard',
  setup() {
    return () => (
        <section class={[prefixCls, `${prefixCls}__${layout.value}`, 'w-[100%] h-[100%] relative']}>

          {renderLayout()}

          <Backtop></Backtop>

        </section>
    )
  }
})
</script>

<style lang="less" scoped>
@prefix-cls: ~'@{namespace}-layout';

.@{prefix-cls} {
  background-color: var(--app-content-bg-color);
  :deep(.@{elNamespace}-scrollbar__view) {
    height: 100% !important;
  }
}
</style>
