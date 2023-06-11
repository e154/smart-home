<template>

  <el-drawer
    :visible.sync="visible"
    :direction="direction"
    :append-to-body="true"
    title="Smart Home terminal"
    size="100%"
    :withHeader="true"
    :destroy-on-close="true"
    :modal="false"
    :close-on-press-escape="true"
  >

    <div class="terminal-viewer">
      <ul class="list-unstyled">
        <li v-for="i in items" class="infinite-list-item">{{ i }}</li>
      </ul>
    </div>

  </el-drawer>

</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import Pagination from '@/components/Pagination/index.vue'
import stream from '@/api/stream'
import { LogObject } from '@/models'
import { parseTime } from '@/utils'

@Component({
  name: 'Terminal',
  components: {
    Pagination
  }
})
export default class extends Vue {
  @Prop() private show = false;

  get visible(): boolean {
    return this.show
  }

  set visible(val: boolean) {
    this.$emit('on-hidden', false)
  }

  private direction = 'btt';

  private items: string[] = [];
  private load = true;

  private created() {
    // todo id
    setTimeout(() => {
      stream.subscribe('log', '1', this.onMessage)
    }, 1000)
  }

  private destroyed() {
    stream.unsubscribe('log', '1')
  }

  private onMessage(m: LogObject) {
    const t = `${parseTime(m.created_at)}        `.substring(0, 33)
    const l = `[${m.level}]       `.substring(0, 9)
    const b = `${m.body}`
    this.items.unshift(t + l + b)
  }
}
</script>

<style lang="scss">

.list-unstyled {
  list-style: none;
}

.el-drawer__header {
  padding: 5px 20px;
  margin-bottom: 0;
  background-color: #0c0c0c;
  color: #c5c5c5;
}

.terminal-viewer {
  background-color: #202020;
  color: #c5c5c5;
  padding: 20px;
  height: 100%;
  min-height: 100%;
  overflow-y: scroll;
  font-size: 0.8em;
  line-height: 1.6;

ul {
  margin: 0;
  padding: 0;
}

}
</style>
