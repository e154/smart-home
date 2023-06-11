<template>
  <el-row>
    <el-col v-bind="columnProps">
      <el-card>
        <slot></slot>
      </el-card>
    </el-col>
  </el-row>

</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import Pagination from '@/components/Pagination/index.vue'

@Component({
  name: 'CardWrapper',
  components: {
    Pagination
  }
})
export default class extends Vue {
  private columnProps = {
    xl: { span: 16, offset: 4 },
    lg: { span: 18, offset: 3 },
    md: { span: 24, offset: 0 },
    sm: { span: 24, offset: 0 }
  }
}
</script>

<style lang="scss">

</style>
