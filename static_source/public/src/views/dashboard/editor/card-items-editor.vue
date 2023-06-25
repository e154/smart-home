<template>
  <el-row :gutter="20">

    <!-- item detail -->
    <el-col :span="15" :xs="15">

      <el-card>
        <div slot="header" class="clearfix">
          <span>{{ $t('dashboard.editor.itemDetail') }}</span>
        </div>

        <div style="padding: 0 0 20px 0" v-if="core.activeCard !== undefined">
          <el-button type="default" @click.prevent.stop="addCardItem"><i
            class="el-icon-plus"/>{{ $t('dashboard.editor.addCardItem') }}
          </el-button>
        </div>

        <div
          v-if="card && card.items && card.selectedItem > -1"
          style="padding: 0 0 20px;"
        >

          <el-card shadow="never" class="item-card-editor">

            <el-form label-position="top"
                     :model="card.items[card.selectedItem]"
                     style="width: 100%"
                     ref="cardItemForm">

              <el-row :gutter="20">
                <el-col
                  :span="8"
                  :xs="8"
                >
                  <el-form-item :label="$t('dashboard.editor.type')" prop="type">
                    <el-select
                      v-model="card.items[card.selectedItem].type"
                      :placeholder="$t('dashboard.editor.pleaseSelectType')"
                      style="width: 100%"
                    >
                      <el-option
                        v-for="item in itemTypes"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value">
                      </el-option>
                    </el-select>
                  </el-form-item>

                </el-col>
                <el-col
                  :span="8"
                  :xs="8"
                >
                  <el-form-item :label="$t('dashboard.editor.title')" prop="title">
                    <el-input size="small"
                              v-model="card.items[card.selectedItem].title"></el-input>
                  </el-form-item>
                </el-col>
              </el-row>

              <component
                :is="getCardEditorName(card.items[card.selectedItem].type)"
                :core="core"
                :item="card.items[card.selectedItem]"
                :index="card.selectedItem"
              />

            </el-form>

            <div style="text-align: right;">

              <el-button type="primary" @click.prevent.stop="updateCard">{{ $t('main.update') }}</el-button>
              <el-button
                @click.prevent.stop="copyCardItem(card.selectedItem)">
                {{ $t('main.copy') }}
              </el-button>
              <el-popconfirm
                :confirm-button-text="$t('main.ok')"
                :cancel-button-text="$t('main.no')"
                icon="el-icon-info"
                icon-color="red"
                style="margin-left: 10px;"
                :title="$t('main.are_you_sure_to_do_want_this?')"
                v-on:confirm="cancel()"
              >
                <el-button type="default" slot="reference">{{ $t('main.cancel') }}</el-button>
              </el-popconfirm>
              <el-popconfirm
                :confirm-button-text="$t('main.ok')"
                :cancel-button-text="$t('main.no')"
                icon="el-icon-info"
                icon-color="red"
                style="margin-left: 10px;"
                :title="$t('main.are_you_sure_to_do_want_this?')"
                v-on:confirm="removeCardItem(card.selectedItem)"
              >
                <el-button type="danger" icon="el-icon-delete" slot="reference">{{
                    $t('main.remove')
                  }}
                </el-button>
              </el-popconfirm>
            </div>

          </el-card>

        </div>

      </el-card>
    </el-col>
    <!-- /item detail -->

    <!-- item list -->
    <el-col :span="8" :xs="12">
      <el-card>
        <div slot="header" class="clearfix">
          <span>{{ $t('dashboard.editor.itemList') }}</span>
        </div>

        <el-menu
          v-if="card && card.id"
          ref="tabMenu"
          :default-active="card.selectedItem + ''"
          v-model="card.selectedItem"
          class="el-menu-vertical-demo"
        >
          <el-menu-item :index="index + ''"
                        v-for="(item, index) in card.items"
                        @click="menuCardItemClick(index)">
                    <span>
                      {{ item.title }}
                      <el-tag size="mini">{{ item.type }}</el-tag>
                      <el-tag v-if="item.hidden" size="mini" type="info" effect="plain">
                        {{ $t('dashboard.editor.hidden') }}
                      </el-tag>
                      <el-tag v-if="!item.enabled" size="mini" type="info" effect="plain">
                        {{ $t('dashboard.editor.disabled') }}
                      </el-tag>
                      <el-tag v-if="item.frozen" size="mini" type="info" effect="plain">
                        {{ $t('dashboard.editor.frozen') }}
                      </el-tag>
                      <el-button style="float: right" size="mini" @click.prevent.stop="sortItemUp(item, index)"><i
                        class="el-icon-upload2"/></el-button>
                          <el-button style="float: right" size="mini" @click.prevent.stop="sortItemDown(item, index)"><i
                            class="el-icon-download"/></el-button>
                    </span>
          </el-menu-item>
        </el-menu>

      </el-card>


      <!-- TODO: fix -->
      <table style="margin-top: 20px" class="filter-list">
        <thead style="background: #d7d7d7">
        <tr>
          <td>name</td>
          <td>description</td>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(filter, index) in getFilterList()">
          <td><strong>{{ filter.name }}</strong></td>
          <td>{{ filter.description }}</td>
        </tr>
        </tbody>
      </table>

    </el-col>
    <!-- /item list -->
  </el-row>
</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator'
import {Card, CardItem, Core} from "@/views/dashboard/core";
import {ApiDashboard} from "@/api/stub";
import {UUID} from "uuid-generator-ts";
import {filterInfo, filterList} from "@/views/dashboard/filters";
import {
  CardEditorName,
  CardItemList,
  IButtonEditor,
  IChartEditor,
  IImageEditor,
  ILogsEditor,
  IProgressEditor,
  IStateEditor,
  ITextEditor
} from "@/views/dashboard/card_items";
import EntitySearch from "@/views/entities/components/entity_search.vue";
import ExportTool from "@/components/export-tool/index.vue";
import AreaSearch from "@/views/areas/components/areas_search.vue";
import CardEditor from "@/views/dashboard/editor/card-editor.vue";
import TabEditor from "@/views/dashboard/editor/tab-editor.vue";
import TabDashboardSettings from "@/views/dashboard/editor/tab-dashboard-settings.vue";
import CardWrapper from "@/components/card-wrapper/index.vue";
import Editor from "@/views/automation/new.vue";
import EditorTabMuu from "@/views/dashboard/editor/tab-muu.vue";
import ImagePreview from "@/views/images/preview.vue";

@Component({
  name: 'CardItemsEditor',
  components: {
    CardEditor,
    TabEditor,
    TabDashboardSettings,
    AreaSearch,
    ExportTool,
    EntitySearch,
    CardWrapper,
    Editor,
    EditorTabMuu,
    ImagePreview,
    IButtonEditor,
    ITextEditor,
    IImageEditor,
    IStateEditor,
    ILogsEditor,
    IProgressEditor,
    IChartEditor,
  }
})
export default class extends Vue {
  @Prop() private card!: Card;
  @Prop() private board!: ApiDashboard;
  @Prop() private core!: Core;
  @Prop() private bus!: Vue;

  // id for streaming subscribe
  private currentID = '';

  private itemTypes = CardItemList;


  private mounted() {
  }

  created() {
    const uuid = new UUID();
    this.currentID = uuid.getDashFreeUUID();

  }

  private destroyed() {

  }

  // ---------------------------------
  // common
  // ---------------------------------
  private addCardItem() {
    this.core.createCardItem();
  }

  private removeCardItem(index: number) {
    this.core.removeCardItem(index);
  }

  private copyCardItem(index: number) {
    this.card.copyItem(index);
  }

  private menuCardItemClick(index: number) {
    if (!this.core.activeTab || this.core.activeCard == undefined) {
      return;
    }

    this.card.selectedItem = index;
  }

  private getFilterList(): filterInfo[] {
    return filterList();
  }

  private getCardEditorName(name: string): string {
    return CardEditorName(name);
  }

  private cancel() {
    this.$router.go(-1);
  }

  private async updateCard() {
    const {data} = await this.core.updateCard();

    if (data) {
      this.$notify({
        title: 'Success',
        message: 'card updated successfully',
        type: 'success',
        duration: 2000
      });
    }
  }

  private sortItemUp(item: CardItem, index: number) {
    // console.log('up', item, index)

    if (!this.card.items[index - 1]) {
      return;
    }

    let rows = [this.card.items[index - 1], this.card.items[index]];
    this.card.items.splice(index - 1, 2, rows[1], rows[0]);

    let counter = 0
    for (const index in this.card.items) {
      this.card.items[index].weight = counter;
      counter++;
    }

    this.core.updateCard();
  }

  private sortItemDown(item: CardItem, index: number) {
    // console.log('down', item, index)

    if (!this.card.items[index + 1]) {
      return;
    }

    let rows = [this.card.items[index], this.card.items[index + 1]];
    this.card.items.splice(index, 2, rows[1], rows[0]);

    let counter = 0
    for (const index in this.card.items) {
      this.card.items[index].weight = counter;
      counter++;
    }

    this.core.updateCard();
  }
}
</script>

<style lang="scss">

</style>
