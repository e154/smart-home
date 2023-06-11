<template>
  <div>
    <common-editor :item="item" :board="board"></common-editor>

    <el-divider content-position="left">{{ $t('dashboard.editor.stateOptions') }}</el-divider>

    <el-row>
      <el-col>
        <div style="padding-bottom: 20px">
          <el-button type="default" @click.prevent.stop="addProp()"><i
            class="el-icon-plus"/>{{ $t('dashboard.editor.addNewProp') }}
          </el-button>
        </div>
      </el-col>
    </el-row>

    <!-- props -->
    <el-collapse>
      <el-collapse-item
        :name="index"
        :key="index"
        v-for="(prop, index) in item.payload.state.items"
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
                <el-form-item :label="$t('dashboard.editor.image')" prop="image">
                  <image-preview
                    :image="prop.image"
                    @on-select="onSelectImageForState(index, ...arguments)"/>
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
      <el-collapse-item>
        <template slot="title">
          {{ $t('dashboard.editor.defaultImage') }}
        </template>
        <el-row>
          <el-col>
            <el-card shadow="never" class="item-card-editor">
              <image-preview
                :key="reloadKeyDefaultImage"
                :image="item.payload.state.default_image"
                @on-select="onSelectDefaultImage"/>
            </el-card>
          </el-col>
        </el-row>
      </el-collapse-item>
    </el-collapse>
    <!-- /props -->

    <el-row style="padding-bottom: 20px">
      <el-col>
        <event-viewer :item="item"></event-viewer>
      </el-col>
    </el-row>

  </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator';
import {ItemPayloadState} from '@/views/dashboard/card_items/state/types';
import {CardItem, comparisonType, Core} from '@/views/dashboard/core';
import CommonEditor from '@/views/dashboard/card_items/common/editor.vue';
import JsonEditor from '@/components/JsonEditor/index.vue';
import {ApiImage} from '@/api/stub';
import ImagePreview from '@/views/images/preview.vue';
import EventViewer from '@/views/dashboard/card_items/common/event_viewer.vue';

@Component({
  name: 'IStateEditor',
  components: {EventViewer, ImagePreview, JsonEditor, CommonEditor}
})
export default class extends Vue {
  @Prop() private item!: CardItem;
  @Prop() private board!: Core;

  private created() {
  }

  private mounted() {
  }

  private addProp() {
    // console.log('add prop');

    if (!this.item.payload.state?.items) {
      this.item.payload.state = {
        items: []
      };
    }

    let counter = 0;
    if (this.item.payload.state.items.length) {
      counter = this.item.payload.state.items.length;
    }

    this.item.payload.state.items.push({

      key: 'new proper ' + counter,
      value: '',
      comparison: comparisonType.EQ,
      image: undefined
    });
    this.item.update();
  }

  private removeProp(index: number) {
    if (!this.item.payload.state) {
      this.item.payload.state = {
        items: [],
        default_image: undefined
      } as ItemPayloadState;
    }

    this.item.payload.state.items!.splice(index, 1);
    this.item.update();
  }

  private onSelectImageForState(index: number, image: ApiImage) {
    console.log('select image', index, image);

    if (!this.item.payload.state) {
      this.item.payload.state = {
        items: [],
        default_image: undefined
      } as ItemPayloadState;
    }

    this.item.payload.state.items[index].image = image as ApiImage || undefined;
    this.item.update();
  }

  private reloadKeyDefaultImage = 0;

  private onSelectDefaultImage(image: ApiImage) {
    console.log('select image', image);

    if (!this.item.payload.state) {
      this.item.payload.state = {
        items: [],
        default_image: undefined
      } as ItemPayloadState;
    }

    this.item.payload.state.default_image = image as ApiImage || undefined;
    // this.reloadKeyDefaultImage += 1
    this.item.update();
  }
}
</script>

<style>

</style>
