<template>
  <div class="app-container">
    <el-empty :image-size="200"></el-empty>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import Pagination from '@/components/Pagination/index.vue'
import api from '@/api/api'
import { ApiLog } from '@/api/stub'

@Component({
  name: 'MqttList',
  components: {
    Pagination
  }
})
export default class extends Vue {

}
</script>

<style lang="scss">

</style>
