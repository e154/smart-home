<template>
  <div class="app-container">
    <card-wrapper>
      <image-browser/>
    </card-wrapper>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import Pagination from '@/components/Pagination/index.vue'
import api from '@/api/api'
import { ApiLog } from '@/api/stub'
import CardWrapper from '@/components/card-wrapper/index.vue'
import ImageBrowser from '@/views/images/browser.vue'

@Component({
  name: 'MqttList',
  components: {
    ImageBrowser,
    CardWrapper,
    Pagination
  }
})
export default class extends Vue {

}
</script>

<style lang="scss">

</style>
