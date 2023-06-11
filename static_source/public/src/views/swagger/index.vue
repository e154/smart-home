<template>
  <div>
    <iframe
      :src="src"
      style="position: absolute; border: none; width: 100%; height: 100%"
      width="100%"
      height="100%"
      frameborder="0"
    >
    </iframe>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'

@Component({
  name: 'Swagger',
  components: {
  }
})
export default class extends Vue {
  private myIframe?: string;
  private src: string = (process.env.VUE_APP_BASE_API || window.location.origin) + '/swagger-ui'

  private onLoad(frame: any) {
    this.myIframe = frame.contentWindow
  }
}
</script>

<style lang="scss" scoped>

</style>
