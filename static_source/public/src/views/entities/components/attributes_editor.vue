<template>
  <el-main>
    <el-row>
      <el-col>
        <el-button
          align="left"
          @click='add()'
          v-if="customAttrs">
          <i class="el-icon-plus"/> {{ $t('entities.addAttribute') }}
        </el-button>

        <el-button
          align="left"
          v-if="attrs && Object.keys(attrs).length"
          @click='loadFromPlugin()'>
          {{ $t('entities.loadFromPlugin') }}
        </el-button>
      </el-col>
    </el-row>
    <el-row>
      <el-col>
        <el-table
          :data="value"
          style="width: 100%">
          <el-table-column
            prop="name"
            :label="$t('entities.table.name')"
            width="180px">

            <template slot-scope="{row}">
              <el-input type="string" v-model="row.name"/>
            </template>

          </el-table-column>
          <el-table-column
            prop="type"
            :label="$t('entities.table.type')"
            width="150px">

            <template slot-scope="{row}">
              <el-select v-model="row.type" placeholder="please select type">
                <el-option label="String" value="STRING"></el-option>
                <el-option label="Int" value="INT"></el-option>
                <el-option label="Float" value="FLOAT"></el-option>
                <el-option label="Bool" value="BOOL"></el-option>
                <el-option label="Array" value="ARRAY"></el-option>
                <el-option label="Time" value="TIME"></el-option>
                <el-option label="Map" value="MAP"></el-option>
                <el-option label="Image" value="IMAGE"></el-option>
              </el-select>
            </template>

          </el-table-column>
          <el-table-column
            width="auto"
            :label="$t('entities.table.value')"
          >

            <template slot-scope="{row}">
              <div v-if="row.type === 'STRING'">
                <el-input type="string" v-model="row.string"/>
              </div>
              <div v-if="row.type === 'IMAGE'">
                <el-input type="string" v-model="row.imageUrl"/>
              </div>
              <div v-if="row.type === 'INT'">
                <el-input type="number" v-model="row.int"/>
              </div>
              <div v-if="row.type === 'FLOAT'">
                <el-input type="number" v-model="row.float"/>
              </div>
              <div v-if="row.type === 'ARRAY'">
                <el-input type="string" v-model="row.array"/>
              </div>
              <el-select v-model="row.bool"
                         placeholder="please select value"
                         v-if="row.type === 'BOOL'"
              >
                <el-option label="TRUE" :value="true"></el-option>
                <el-option label="FALSE" :value="false"></el-option>
              </el-select>

              <div v-if="row.type === 'TIME'">
                <el-input type="string" v-model="row.time"/>
              </div>

              <div v-if="row.type === 'MAP'">
                <el-input type="string" v-model="row.map"/>
              </div>
            </template>

          </el-table-column>
          <el-table-column
            align="right"
            :label="$t('entities.table.operations')"
            v-if="customAttrs">
            <template slot-scope="{row, $index}">
              <el-button
                type="text"
                align="right"
                @click.prevent.stop="remove(row, $index)"
              >
                {{ $t('main.remove') }}
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-col>
    </el-row>
  </el-main>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'
import { Attribute } from '@/models/attributes'
import { ApiAttribute } from '@/api/stub'

@Component({
  name: 'AttributesEditor'
})
export default class extends Vue {
  @Prop({ default: [] }) private value!: ApiAttribute[];
  @Prop({ default: false }) private customAttrs?: boolean;
  @Prop({ default: {} }) private attrs?: Record<string, ApiAttribute>;

  private add() {
    this.value.push(new Attribute('new_value'))
  }

  private loadFromPlugin() {
    // ApiAttribute
    const attrs: ApiAttribute[] = []
    if (this.attrs) {
      for (const key in this.attrs) {
        attrs.push(this.attrs[key])
      }
    }
    this.$emit('update-value', attrs)
    setTimeout(() => {
      this.$forceUpdate()
    }, 0.5 * 1000)
  }

  private remove(row: Attribute, index: number) {
    this.value.splice(index, 1)
  }
}
</script>

<style>
.el-main {
  padding: 20px 0 0 0;
}
</style>
