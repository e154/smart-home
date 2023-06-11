<template>
  <div>
    <common-editor :item="item" :board="board"></common-editor>

    <el-divider content-position="left">{{ $t('dashboard.editor.textOptions') }}</el-divider>

    <el-row style="padding-bottom: 20px">
      <el-col>
        <div style="padding-bottom: 20px">
          <el-button type="default" @click.prevent.stop="addProp()"><i
            class="el-icon-plus"/>{{ $t('dashboard.editor.addProp') }}
          </el-button>
        </div>

        <!-- props -->
        <el-collapse>
          <el-collapse-item
            :name="index"
            :key="index"
            v-for="(prop, index) in item.payload.text.items"
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
                    <el-form-item :label="$t('dashboard.editor.html')" prop="enabled">
                      <el-switch
                        v-model="defaultTextHtml"></el-switch>
                    </el-form-item>

                    <el-form-item :label="$t('dashboard.editor.text')" prop="text">
                      <el-input
                        v-if="defaultTextHtml"
                        type="textarea"
                        :autosize="{minRows: 10}"
                        placeholder="Please input"
                        v-model="prop.text"
                        @change="propTextUpdated(prop)"
                      >
                      </el-input>
                      <tinymce
                        v-else
                        v-model="prop.text"
                        :height="400"
                        @text-change="propTextUpdated(prop)"
                      />
                    </el-form-item>
                  </el-col>
                </el-row>

                <el-row>
                  <el-col>
                    <el-form-item :label="$t('dashboard.editor.tokens')">
                      <el-tag size="small" v-for="(token, index) in prop.tokens">{{ token }}</el-tag>
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

      </el-col>
    </el-row>

    <el-row>
      <el-col>
        <el-form-item :label="$t('dashboard.editor.html')" prop="enabled">
          <el-switch
            v-model="defaultTextHtml"></el-switch>
        </el-form-item>

        <el-form-item :label="$t('dashboard.editor.defaultText')" prop="text">
          <el-input
            v-if="defaultTextHtml"
            type="textarea"
            :autosize="{minRows: 10}"
            placeholder="Please input"
            v-model="item.payload.text.default_text"
            @change="defaultTextUpdated"
          >
          </el-input>
          <tinymce
            v-else
            v-model="item.payload.text.default_text"
            :height="400"
            @text-change="defaultTextUpdated"
          />
        </el-form-item>
      </el-col>
    </el-row>

    <el-row>
      <el-col>
        <el-form-item :label="$t('dashboard.editor.tokens')">
          <el-tag size="small" v-for="(token, index) in tokens">{{ token }}</el-tag>
        </el-form-item>
      </el-col>
    </el-row>

    <el-row style="padding-bottom: 20px">
      <el-col>
        <event-viewer :item="item"></event-viewer>
      </el-col>
    </el-row>

  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'
import { CardItem, comparisonType, Core } from '@/views/dashboard/core'
import { Cache, GetTokens } from '@/views/dashboard/render'
import CommonEditor from '@/views/dashboard/card_items/common/editor.vue'
import ScriptEditor from '@/components/script_editor/index.vue'
import JsonEditor from '@/components/JsonEditor/index.vue'
import EventViewer from '@/views/dashboard/card_items/common/event_viewer.vue'
import Tinymce from '@/components/Tinymce/index.vue'
import {TextProp} from '@/views/dashboard/card_items/text/types';

@Component({
  name: 'ITextEditor',
  components: { Tinymce, EventViewer, JsonEditor, ScriptEditor, CommonEditor }
})
export default class extends Vue {
  @Prop() private item!: CardItem;
  @Prop() private board!: Core;

  private _cache!: Cache;
  private tokens: string[] = [];
  private customToolbar = [
    [{ header: [false, 1, 2, 3, 4, 5, 6] }],
    ['bold', 'italic', 'underline', 'strike'], // toggled buttons
    [{ align: '' }, { align: 'center' }, { align: 'right' }, { align: 'justify' }],
    ['blockquote', 'code-block'], [{ list: 'ordered' }, { list: 'bullet' }, { list: 'check' }],
    [{ indent: '-1' }, { indent: '+1' }], // outdent/indent
    [{ color: [] }, { background: [] }], // dropdown with defaults from theme
    /* ["link", "image", "video"], */ ['clean'] // remove formatting button
  ];

  private defaultTextHtml = false;

  private created() {
    this._cache = new Cache()
    this.update()
  }

  private mounted() {
  }

  get lastEvent() {
    return this.item.lastEvent || {}
  }

  private update() {
    this.updateTokensDefaultText()

    if (this.item?.payload?.text?.default_text) {
      for (const prop of this.item.payload.text.items) {
        this.updateTokensPropText(prop)
      }
    }
  }

  private updateTokensDefaultText() {
    if (!this.item?.payload?.text?.default_text) {
      this.tokens = []
      return
    }

    const tokens = GetTokens(this.item.payload.text.default_text, this._cache)
    this.tokens = tokens || []
  }

  private defaultTextUpdated() {
    this.updateTokensDefaultText()
  }

  private updateTokensPropText(prop: TextProp) {
    prop.tokens = GetTokens(prop.text, this._cache) || []
  }

  private propTextUpdated(prop: TextProp) {
    this.updateTokensPropText(prop)
  }

  @Watch('item.payload.text')
  private onUpdateText(item: CardItem) {
    this.update()
  }

  private addProp() {
    console.log('addProp')

    if (!this.item.payload.text?.items) {
      this.item.payload.text!.items = []
    }

    let counter = 0
    if (this.item.payload.text!.items.length) {
      counter = this.item.payload.text!.items.length
    }

    this.item.payload.text!.items.push({

      key: 'new proper ' + counter,
      value: '',
      comparison: comparisonType.EQ,
      text: ''
    })
  }

  private removeProp(index: number) {
    if (!this.item.payload.text?.items) {
      return
    }

    this.item.payload.text?.items.splice(index, 1)
  }
}
</script>

<style>

</style>
