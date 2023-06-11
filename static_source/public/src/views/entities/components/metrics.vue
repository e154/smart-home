<template>
  <el-row :gutter="20">

    <el-col :span="16" :xs="12">
      <div style="padding-bottom: 20px">
        <el-button type="default" @click.prevent.stop="add"><i class="el-icon-plus"/> {{
            $t('entities.metric.addMetric')
          }}
        </el-button>
      </div>

      <el-form v-if="value[selectedMetric]" label-position="top" label-width="100px"
               ref="currentItem"
               :model="value[selectedMetric]"
               :rules="rules2"
               style="width: 100%">
        <el-form-item :label="$t('entities.metric.name')" prop="name">
          <el-input size="small" v-model="value[selectedMetric].name"></el-input>
        </el-form-item>
        <el-form-item :label="$t('entities.metric.description')" prop="description">
          <el-input size="small" v-model="value[selectedMetric].description"></el-input>
        </el-form-item>
<!--        <el-form-item :label="$t('entities.metric.type')" prop="type">-->
<!--          <el-select v-model="value[selectedMetric].type" placeholder="please select type" default-first-option-->
<!--                     style="width: 100%">-->
<!--            <el-option label="LINE" value="line"></el-option>-->
<!--            <el-option label="PIE" value="pie"></el-option>-->
<!--          </el-select>-->
<!--        </el-form-item>-->
<!--        <el-form-item :label="$t('entities.metric.ranges')" prop="ranges">-->
<!--          <el-select-->
<!--            v-model="value[selectedMetric].ranges"-->
<!--            placeholder="please select type"-->
<!--            :multiple="true"-->
<!--            remote-->
<!--            clearable-->
<!--            style="width: 100%"-->
<!--          >-->
<!--            <el-option label="1h" value="1h"></el-option>-->
<!--            <el-option label="6h" value="6h"></el-option>-->
<!--            <el-option label="12h" value="12h"></el-option>-->
<!--            <el-option label="24h" value="24h"></el-option>-->
<!--            <el-option label="7d" value="7d"></el-option>-->
<!--            <el-option label="30d" value="30d"></el-option>-->
<!--          </el-select>-->
<!--        </el-form-item>-->
      </el-form>

      <el-divider v-if="selectedMetric >= 0" content-position="left">Properties</el-divider>

      <div style="padding-bottom: 20px" v-if="selectedMetric >= 0">
        <el-button type="default" @click.prevent.stop="addProp()"><i
          class="el-icon-plus"/>{{ $t('entities.metric.addProp') }}
        </el-button>
      </div>

      <!-- props -->
      <el-collapse v-if="value[selectedMetric] && value[selectedMetric].options">
        <el-collapse-item
          :name="index"
          :key="index"
          v-for="(prop, index) in value[selectedMetric].options.items"
        >

          <template slot="title">
            {{ prop.name }}
          </template>

          <el-card shadow="never" class="item-card-editor">
            <el-form label-position="top"
                     :model="prop"
                     style="width: 100%"
                     ref="cardItemForm">

              <el-row :gutter="20">
                <el-col>
                  <el-form-item :label="$t('entities.metric.name')" prop="name">
                    <el-input size="small" v-model="prop.name"></el-input>
                  </el-form-item>

                  <el-form-item :label="$t('entities.metric.description')" prop="description">
                    <el-input size="small" v-model="prop.description"></el-input>
                  </el-form-item>

                  <el-form-item :label="$t('entities.metric.color')" prop="background">
                    <el-color-picker v-model="prop.color"></el-color-picker>
                  </el-form-item>

                  <el-form-item :label="$t('entities.metric.translate')" prop="translate">
                    <el-input size="small" v-model="prop.translate"></el-input>
                  </el-form-item>

                  <el-form-item :label="$t('entities.metric.label')" prop="label">
                    <el-input size="small" v-model="prop.label"></el-input>
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
                        v-on:confirm="removeProp(index)"
                      >
                        <el-button type="danger" icon="el-icon-delete" slot="reference">{{
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

      <div style="padding: 20px 0; text-align: right"
           v-if="value[selectedMetric]"
      >
        <el-popconfirm
          :confirm-button-text="$t('main.ok')"
          :cancel-button-text="$t('main.no')"
          icon="el-icon-info"
          icon-color="red"
          :title="$t('main.are_you_sure_to_do_want_this?')"
          v-on:confirm="remove(index)"
        >
          <el-button type="danger" icon="el-icon-delete" slot="reference">{{
              $t('main.remove')
            }}
          </el-button>
        </el-popconfirm>
      </div>

    </el-col>
    <el-col :span="8" :xs="12">
      <el-card>
        <div slot="header" class="clearfix">
          <span>{{ $t('entities.metric.list') }}</span>
        </div>

        <el-menu
          v-if="value"
          :default-active="selectedMetric + ''"
          v-model="selectedMetric"
          class="el-menu-vertical-demo"
        >
          <el-menu-item :index="index + ''" v-for="(metric, index) in value"
                        @click="menuClick(index)">
            <span>{{ metric.name }}</span>
          </el-menu-item>
        </el-menu>

      </el-card>
    </el-col>
  </el-row>
</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator';
import {ApiMetric} from '@/api/stub';

@Component({
  name: 'Metrics',
  components: {}
})
export default class extends Vue {
  @Prop({default: [], required: true}) private value!: ApiMetric[];

  private selectedMetric = -1;

  private rules = {
    name: [
      {required: true, trigger: 'blur'},
      {min: 4, max: 255, trigger: 'blur'}
    ],
    description: [
      {required: false, trigger: 'blur'},
      {max: 255, trigger: 'blur'}
    ],
    type: [
      {required: false, trigger: 'blur'},
      {max: 255, trigger: 'blur'}
    ]
  };

  private rules2 = {
    name: [
      {required: true, trigger: 'blur'},
      {min: 4, max: 255, trigger: 'blur'}
    ],
    description: [
      {required: false, trigger: 'blur'},
      {max: 255, trigger: 'blur'}
    ],
    color: [
      {required: false, trigger: 'blur'},
      {max: 255, trigger: 'blur'}
    ],
    translate: [
      {required: false, trigger: 'blur'},
      {max: 255, trigger: 'blur'}
    ],
    label: [
      {required: false, trigger: 'blur'},
      {max: 255, trigger: 'blur'}
    ]
  };

  created() {
    if (this.value && this.value.length) {
      this.selectedMetric = 0;
    }
  }

  private add() {
    this.value.push({
      description: undefined,
      name: `new metric ${this.value.length}`,
      ranges: [],
      type: 'LINE',
      options: {
        items: []
      }
    } as ApiMetric);
    this.selectedMetric = this.value.length - 1 || 0;
  }

  private remove() {
    if (!this.value || !this.value.length || this.selectedMetric < 0) {
      return;
    }
    this.value.splice(this.selectedMetric, 1);
    this.selectedMetric = this.value.length - 1;
  }

  private addProp() {


    this.value[this.selectedMetric].options!.items!.push({
      name: 'new label',
      description: '',
      color: '#FF0000',
      translate: '',
      label: ''
    });
  }

  private removeProp(index: number) {
    if (this.selectedMetric < 0 || !this.value[this.selectedMetric]) {
      return;
    }

    this.value[this.selectedMetric].options!.items!.splice(index, 1);
    this.selectedMetric = 0;
  }

  private menuClick(index: number) {
    this.selectedMetric = index;
  }
}
</script>
