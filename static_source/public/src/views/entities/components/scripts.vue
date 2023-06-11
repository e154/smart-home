<template>
  <div>
    <script-search
      style="width: 100%"
      v-model="value"
      @update-value="changedScript"/>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import { ApiScript } from '@/api/stub'
import ScriptSearch from '@/views/scripts/components/script_search.vue'

@Component({
  name: 'Scripts',
  components: { ScriptSearch }
})
export default class extends Vue {
  @Prop({ required: true }) private value!: ApiScript[];

  private changedScript(value: any, event: any) {
    this.$emit('update-value', value)
  }
}
</script>
