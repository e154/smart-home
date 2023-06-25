<template>
  <el-row :gutter="20">
    <!-- card details -->
    <el-col :span="15" :xs="15">
      <el-card>
        <div slot="header" class="clearfix">
          <span>{{ $t('dashboard.editor.cardSettings') }}</span>
        </div>

        <!-- main options -->
        <el-row style="padding-bottom: 20px" v-if="tab.id">
          <el-col>
            <el-button type="default" @click.prevent.stop="addCard"><i class="el-icon-plus"/>
              {{ $t('dashboard.editor.addNewCard') }}
            </el-button>

            <el-button type="default" @click.prevent.stop="showCardImport = true">{{
                $t('main.import')
              }}
            </el-button>

            <export-tool
              :title="$t('main.import')"
              :visible="showCardImport"
              :value="importCardValue"
              @on-close="showCardImport=false"
              @on-import="importCard"
              :import-dialog="true"/>
          </el-col>
        </el-row>

        <div v-if="tab.cards[core.activeCard] && tab.cards[core.activeCard].id">
          <el-form label-position="top"
                   :model="tab.cards[core.activeCard]"
                   style="width: 100%"
                   :rules="cardRules"
                   ref="cardForm">
            <el-form-item :label="$t('dashboard.editor.title')" prop="title">
              <el-input size="small"
                        v-model="tab.cards[core.activeCard].title"></el-input>
            </el-form-item>

            <el-form-item :label="$t('dashboard.editor.height')" prop="height">
              <el-input-number size="small"
                               v-model="tab.cards[core.activeCard].height"></el-input-number>
            </el-form-item>

            <el-form-item :label="$t('dashboard.editor.width')" prop="width">
              <el-input-number size="small" v-model="tab.cards[core.activeCard].width"></el-input-number>
            </el-form-item>

            <el-form-item :label="$t('dashboard.editor.background')" prop="background">
              <el-color-picker
                v-model="tab.cards[core.activeCard].background"></el-color-picker>
            </el-form-item>
            <el-form-item :label="$t('dashboard.editor.enabled')" prop="enabled">
              <el-switch v-model="tab.cards[core.activeCard].enabled"></el-switch>
            </el-form-item>
            <el-form-item :label="$t('dashboard.editor.hidden')" prop="hidden">
              <el-switch v-model="tab.cards[core.activeCard].hidden"></el-switch>
            </el-form-item>

            <el-form-item :label="$t('dashboard.editor.entity')" prop="entity">
              <entity-search
                v-model="tab.cards[core.activeCard].entity"
                @update-value="changedCardEntity"
              />
            </el-form-item>

            <!-- show on -->
            <el-divider content-position="left">{{ $t('dashboard.editor.showOn') }}</el-divider>
            <show-on v-model="tab.cards[core.activeCard].showOn" />
            <!-- /show on -->

            <!-- hide on-->
            <el-divider content-position="left">{{ $t('dashboard.editor.hideOn') }}</el-divider>
            <show-on v-model="tab.cards[core.activeCard].hideOn" />
            <!-- /hide on-->

            <el-row style="padding-bottom: 20px">
              <el-col>
                <event-viewer :item="tab.cards[core.activeCard]"></event-viewer>
              </el-col>
            </el-row>

          </el-form>

        </div>
        <!-- /main options -->

        <el-row v-if="tab.cards[core.activeCard] && tab.cards[core.activeCard].id" style="margin-top: 20px">
          <el-col style="text-align: right">

            <export-tool
              :title="$t('main.export')"
              :visible="showExportCard"
              :value="exportCardValue"
              @on-close="showExportCard=false"/>
            <el-button type="primary" icon="el-icon-document" @click.prevent.stop='exportCard'>{{
                $t('main.export')
              }}
            </el-button>

            <el-button type="primary" @click.prevent.stop="copyCard">{{ $t('main.copy') }}</el-button>

            <el-button type="primary" @click.prevent.stop="updateCard">{{ $t('main.update') }}</el-button>

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
              v-on:confirm="removeCard"
            >
              <el-button type="danger" icon="el-icon-delete" slot="reference">{{
                  $t('main.remove')
                }}
              </el-button>
            </el-popconfirm>
          </el-col>
        </el-row>

      </el-card>
    </el-col>
    <!-- card details -->

    <!-- card list -->
    <el-col :span="8" :xs="12">

      <el-card>
        <div slot="header" class="clearfix">
          <span>{{ $t('dashboard.editor.cardList') }}</span>
        </div>

        <el-menu
          v-if="core.activeTab && tab.cards"
          ref="tabMenu"
          :default-active="core.activeCard + ''"
          v-model="core.activeCard"
          class="el-menu-vertical-demo"
        >
          <el-menu-item :index="index + ''" v-for="(card, index) in tab.cards"
                        @click="menuCardsClick(card)">
            <div style="clear: both">
              <span>{{ card.title }}</span>
              <span>
                            <el-button style="float: right" size="mini" @click.prevent.stop="sortCardUp(card, index)"><i
                              class="el-icon-upload2"/></el-button>
                          <el-button style="float: right" size="mini" @click.prevent.stop="sortCardDown(card, index)"><i
                            class="el-icon-download"/></el-button>
                          </span>
            </div>

          </el-menu-item>
        </el-menu>

      </el-card>
    </el-col>
    <!-- /card list -->

  </el-row>
</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator';
import {
  Dummy,
  IButton,
  IChart,
  IImage,
  ILogs,
  IProgress,
  IState,
  IText
} from '@/views/dashboard/card_items';
import {Card, Core, Tab} from '@/views/dashboard/core';
import Moveable from 'vue-moveable';
import {ApiDashboard, ApiEntity} from "@/api/stub";
import {UUID} from "uuid-generator-ts";
import AreaSearch from "@/views/areas/components/areas_search.vue";
import ExportTool from "@/components/export-tool/index.vue";
import {Form} from "element-ui";
import EntitySearch from "@/views/entities/components/entity_search.vue";
import ShowOn from "@/views/dashboard/card_items/common/show-on.vue";
import EventViewer from "@/views/dashboard/card_items/common/event_viewer.vue";

class elementOption {
  public value = '';
  public label = '';
}

@Component({
  name: 'CardEditor',
  components: {
    EventViewer,
    ShowOn,
    EntitySearch,
    ExportTool,
    AreaSearch,
    Moveable,
    Dummy,
    IText,
    IImage,
    IButton,
    IState,
    ILogs,
    IProgress,
    IChart
  }
})
export default class extends Vue {
  @Prop() private tab!: Tab;
  @Prop() private board!: ApiDashboard;
  @Prop() private core!: Core;
  @Prop() private bus!: Vue;

  // id for streaming subscribe
  private currentID = '';

  private mounted() {
  }

  created() {
    const uuid = new UUID();
    this.currentID = uuid.getDashFreeUUID();

    this.bus.$on('selected_card', (m: any) => {
      this.onSelectedCard(m);
    });
  }

  private destroyed() {

  }

  private cardRules = {
    title: [
      {required: true, trigger: 'blur'},
      {min: 4, max: 255, trigger: 'blur'}
    ]
  };

  // ---------------------------------
  // common
  // ---------------------------------

  private onSelectedCard(id: number) {
    this.core.onSelectedCard(id);
  }

  private async addCard() {
    return this.core.createCard();
  }

  private async updateCard() {
    (this.$refs.cardForm as Form).validate(async valid => {
      if (!valid) {
        return;
      }

      const {data} = await this.core.updateCard();

      if (data) {
        this.$notify({
          title: 'Success',
          message: 'card updated successfully',
          type: 'success',
          duration: 2000
        });
      }
    });
  }

  private async removeCard() {
    await this.core.removeCard();

    this.$notify({
      title: 'Success',
      message: 'card removed successfully',
      type: 'success',
      duration: 2000
    });
  }

  private menuCardsClick(card: Card) {
    this.bus.$emit('selected_card', card.id);
  }

  // export card
  private showExportCard = false;
  private exportCardValue = {};

  private exportCard() {
    if (this.core.activeTab == undefined || this.core.activeCard == undefined) {
      return;
    }
    this.exportCardValue = this.board.tabs[this.core.activeTab].cards[this.core.activeCard].serialize();
    this.showExportCard = true;
  }

  // import card
  private showCardImport = false;
  private importCardValue = '';

  private async importCard(value: string, event?: any) {
    const card = JSON.parse(value);
    const data = await this.core.importCard(card);
    if (data) {
      this.$notify({
        title: 'Success',
        message: 'card imported successfully',
        type: 'success',
        duration: 2000
      });
    }
  }

  private async copyCard() {
    if (this.core.activeTab == undefined || this.core.activeCard == undefined) {
      return;
    }

    const card = this.tab.cards[this.core.activeCard].serialize();

    const data = await this.core.importCard(card);
    if (data) {
      this.$notify({
        title: 'Success',
        message: 'card copied successfully',
        type: 'success',
        duration: 2000
      });
    }
  }

  private sortCardUp(card: Card, index: number) {
    if (!this.core.tabs || !this.core.activeTab) {
      return;
    }

    if (!this.tab.cards[index - 1]) {
      return;
    }

    let rows = [this.tab.cards[index - 1], this.tab.cards[index]];
    this.tab.cards.splice(index - 1, 2, rows[1], rows[0]);

    let counter = 0
    for (const index in this.tab.cards) {
      this.tab.cards[index].weight = counter;
      this.tab.cards[index].update();
      counter++;
    }

    this.core.updateCurrentTab();
  }

  private sortCardDown(card: Card, index: number) {
    if (!this.core.tabs || !this.core.activeTab) {
      return;
    }

    if (!this.tab.cards[index + 1]) {
      return;
    }

    let rows = [this.tab.cards[index], this.tab.cards[index + 1]];
    this.tab.cards.splice(index, 2, rows[1], rows[0]);

    let counter = 0
    for (const index in this.tab.cards) {
      this.tab.cards[index].weight = counter;
      this.tab.cards[index].update();
      counter++;
    }

    this.core.updateCurrentTab();
  }

  private changedCardEntity(entity: ApiEntity, event?: any) {
    if (!this.core.activeTab || this.core.activeCard == undefined) {
      return;
    }

    if (!entity?.id) {
      this.tab.cards[this.core.activeCard].entity = undefined;
      return;
    }
    this.fetchEntity(entity.id);
  }

  private async fetchEntity(id: string) {
    if (!this.core.activeCard) {
      return
    }

    const entity = await this.core.fetchEntity(id);
    this.tab.cards[this.core.activeCard].entity = entity;
  }

  private cancel() {
    this.$router.go(-1);
  }
}
</script>

<style lang="scss">



</style>
