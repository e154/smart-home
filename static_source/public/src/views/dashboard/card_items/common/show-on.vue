<template>
  <el-row>
    <el-col>
      <div style="padding-bottom: 20px">
        <el-button type="default" @click.prevent.stop="addShowOnProp()"><i
          class="el-icon-plus"/>{{ $t('dashboard.editor.addNewProp') }}
        </el-button>
      </div>

      <!-- props -->
      <el-collapse>
        <el-collapse-item
          :name="index"
          :key="index"
          v-for="(prop, index) in value"
        >

          <template slot="title">
            <el-tag size="mini">{{ prop.key }}</el-tag>
            +
            <el-tag size="mini">{{ prop.comparison }}</el-tag>
            +
            <el-tag size="mini">{{ prop.value }}</el-tag>
          </template>

          <el-card shadow="never" class="item-card-editor">

            <el-form label-position="top"
                     :model="prop"
                     style="width: 100%"
                     ref="cardItemForm">

              <el-row :gutter="20">
                <el-col
                  :span="8"
                  :xs="8"
                >
                  <el-form-item :label="$t('dashboard.editor.text')" prop="text">
                    <el-input
                      placeholder="Please input"
                      v-model="prop.key">
                    </el-input>
                  </el-form-item>

                </el-col>

                <el-col
                  :span="8"
                  :xs="8"
                >
                  <el-form-item :label="$t('dashboard.editor.comparison')" prop="comparison">
                    <el-select
                      v-model="prop.comparison"
                      placeholder="please select type"
                      style="width: 100%"
                    >
                      <el-option label="==" value="eq"></el-option>
                      <el-option label="<" value="lt"></el-option>
                      <el-option label="<=" value="le"></el-option>
                      <el-option label="!=" value="ne"></el-option>
                      <el-option label=">=" value="ge"></el-option>
                      <el-option label=">" value="gt"></el-option>
                    </el-select>
                  </el-form-item>

                </el-col>

                <el-col
                  :span="8"
                  :xs="8"
                >

                  <el-form-item :label="$t('dashboard.editor.value')" prop="value">
                    <el-input
                      placeholder="Please input"
                      v-model="prop.value">
                    </el-input>
                  </el-form-item>

                </el-col>
              </el-row>

              <el-row>
                <el-col>
                  <div style="padding-bottom: 20px">
                    <div style="text-align: right;">
                      <el-popconfirm
                        :confirm-button-text="$t('main.ok')"
                        :cancel-button-text="$t('main.no')"
                        icon="el-icon-info"
                        icon-color="red"
                        style="margin-left: 10px;"
                        :title="$t('main.are_you_sure_to_do_want_this?')"
                        v-on:confirm="removeShowOnProp(index)"
                      >
                        <el-button type="danger" plain icon="el-icon-delete" slot="reference">{{
                            $t('main.remove')
                          }}
                        </el-button>
                      </el-popconfirm>
                    </div>
                  </div>
                </el-col>
              </el-row>

            </el-form>

          </el-card>

        </el-collapse-item>
      </el-collapse>
      <!-- /props -->

    </el-col>
  </el-row>
</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator'
import { CompareProp, comparisonType} from "@/views/dashboard/core";

@Component({
  name: 'ShowOn',
  components: {}
})
export default class extends Vue {
  @Prop({ default: [] }) private value!: CompareProp[];

  created() {

  }

  // ---------------------------------
  // common
  // ---------------------------------
  private addShowOnProp() {

    if (!this.value) {
      this.value = [];
    }

    let counter = 0;
    if (this.value.length) {
      counter = this.value.length;
    }

    this.value.push({
      key: 'new proper ' + counter,
      value: '',
      comparison: comparisonType.EQ
    });
  }


  private removeShowOnProp(index: number) {
    this.value.splice(index, 1);
  }
}
</script>

<style lang="scss">

</style>
