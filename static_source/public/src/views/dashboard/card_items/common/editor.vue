<template>
  <div>
    <!-- main options -->
    <el-divider content-position="left">{{ $t('dashboard.editor.mainOptions') }}</el-divider>

    <el-row :gutter="20">
      <el-col
        :span="8"
        :xs="8"
      >
        <el-form-item :label="$t('dashboard.editor.entity')" prop="entity">
          <entity-search
            v-model="item.entity"
            @update-value="changedEntity"
          />
        </el-form-item>

        <el-form-item :label="$t('dashboard.editor.frozen')" prop="frozen">
          <el-switch
            v-model="item.frozen"></el-switch>
        </el-form-item>
      </el-col>

      <el-col
        :span="8"
        :xs="8">

      </el-col>

      <el-col
        :span="8"
        :xs="8">
        <el-form-item :label="$t('dashboard.editor.enabled')" prop="enabled">
          <el-switch
            v-model="item.enabled"></el-switch>
        </el-form-item>
        <el-form-item :label="$t('dashboard.editor.hidden')" prop="hidden">
          <el-switch
            v-model="item.hidden"></el-switch>
        </el-form-item>
      </el-col>
    </el-row>
    <!-- /main options -->

    <!-- show on -->
    <el-divider content-position="left">{{ $t('dashboard.editor.showOn') }}</el-divider>

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
            v-for="(prop, index) in item.showOn"
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
    <!-- /show on -->

    <!-- hide on-->
    <el-divider content-position="left">{{ $t('dashboard.editor.hideOn') }}</el-divider>

    <el-row>
      <el-col>
        <div style="padding-bottom: 20px">
          <el-button type="default" @click.prevent.stop="addHideOnProp()"><i
            class="el-icon-plus"/>{{ $t('dashboard.editor.addNewProp') }}
          </el-button>
        </div>

        <!-- props -->
        <el-collapse>
          <el-collapse-item
            :name="index"
            :key="index"
            v-for="(prop, index) in item.hideOn"
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
                          v-on:confirm="removeHideOnProp(index)"
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
    <!-- /hide on-->

    <!-- button options -->
    <el-divider content-position="left">{{ $t('dashboard.editor.buttonOptions') }}</el-divider>
    <el-row :gutter="20">
      <el-col
        :span="8"
        :xs="8"
      >
        <el-form-item :label="$t('dashboard.editor.asButton')" prop="enabled">
          <el-switch
            v-model="item.asButton"></el-switch>
        </el-form-item>
      </el-col>
    </el-row>

    <el-row>
      <el-col>
        <div style="padding-bottom: 20px">
          <el-button type="default" @click.prevent.stop="addAction()"><i
            class="el-icon-plus"/>{{ $t('dashboard.editor.addAction') }}
          </el-button>
        </div>

        <!-- props -->
        <el-collapse>
          <el-collapse-item
            v-for="(prop, index) in item.buttonActions"
            :name="index"
            :key="index"
          >

            <template slot="title">
              {{ prop.entityId }} - {{ prop.action }}
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
                    <el-form-item :label="$t('dashboard.editor.entity')" prop="entity">
                      <entity-search
                        disabled=""
                        v-model="prop.entity"
                        @update-value="changedForActionButton($event, index)"
                      />
                    </el-form-item>

                  </el-col>

                  <el-col
                    :span="8"
                    :xs="8"
                  >

                    <el-form-item :label="$t('dashboard.editor.action')" prop="action" :aria-disabled="!item.entity">

                      <el-select
                        v-model="prop.action"
                        clearable
                        :placeholder="$t('dashboard.editor.selectAction')"
                        style="width: 100%"
                      >
                        <el-option
                          v-for="item in getButtonAction(prop.entity)"
                          :key="item.name"
                          :label="item.name"
                          :value="item.name">
                        </el-option>
                      </el-select>

                    </el-form-item>

                  </el-col>

                  <el-col
                    :span="8"
                    :xs="8"
                  >

                    <el-form-item :label="$t('dashboard.editor.image')" prop="image">
                      <image-preview
                        :image="prop.image"
                        @on-select="onSelectImageForAction(index, ...arguments)"/>
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
                          v-on:confirm="removeAction(index)"
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
    <!-- /button options -->

  </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator';
import {CardItem, comparisonType, Core} from '@/views/dashboard/core';
import {ApiEntity, ApiImage} from '@/api/stub';
import EntitySearch from '@/views/entities/components/entity_search.vue';
import ImagePreview from '@/views/images/preview.vue';

interface Action {
  value: string;
  label: string;
}

@Component({
  name: 'CommonEditor',
  components: {ImagePreview, EntitySearch}
})
export default class extends Vue {
  @Prop() private item!: CardItem;
  @Prop() private board!: Core;

  private created() {
    setTimeout(() => {
      if (this.item.entityId) {
        this.fetchEntity(this.item.entityId);
      }
    }, 1000);
  }

  private mounted() {
  }

  private async fetchEntity(id: string) {
    const entity = await this.board.fetchEntity(id);
    this.item.entity = entity;
  }

  private changedEntity(entity: ApiEntity, event?: any) {
    if (!entity?.id) {
      this.item.entity = undefined;
      return;
    }
    this.fetchEntity(entity.id);
  }

  private async changedForActionButton(entity: ApiEntity, index: number) {
    if (entity?.id) {
      this.item.buttonActions[index].entity = await this.board.fetchEntity(entity.id);
      this.item.buttonActions[index].entityId = entity.id;
    } else {
      this.item.buttonActions[index].entity = undefined;
      this.item.buttonActions[index].entityId = '';
      this.item.buttonActions[index].action = '';
    }
  }

  private getButtonAction(entity?: ApiEntity) {
    if (!entity) {
      return [];
    }
    return entity.actions;
  }

  private addShowOnProp() {
    // console.log('addShowOnProp');

    if (!this.item?.showOn) {
      this.item.showOn = [];
    }

    let counter = 0;
    if (this.item.showOn.length) {
      counter = this.item.showOn.length;
    }

    this.item.showOn.push({
      key: 'new proper ' + counter,
      value: '',
      comparison: comparisonType.EQ
    });
  }

  private addHideOnProp() {
    // console.log('addHideOnProp');

    if (!this.item?.hideOn) {
      this.item.hideOn = [];
    }

    let counter = 0;
    if (this.item.hideOn.length) {
      counter = this.item.hideOn.length;
    }

    this.item.hideOn.push({
      key: 'new proper ' + counter,
      value: '',
      comparison: comparisonType.EQ
    });
  }

  private removeShowOnProp(index: number) {
    if (!this.item.showOn) {
      return;
    }

    this.item.showOn.splice(index, 1);
  }

  private removeHideOnProp(index: number) {
    if (!this.item.hideOn) {
      return;
    }

    this.item.hideOn.splice(index, 1);
  }

  private onSelectImageForAction(index: number, image: ApiImage) {
    // console.log('select image', index, image);
    if (!this.item.buttonActions[index]) {
      return;
    }

    this.item.buttonActions[index].image = image as ApiImage || undefined;
  }

  private addAction() {
    this.item.buttonActions.push({
      entity: undefined,
      entityId: this.item.entityId,
      action: '',
      image: null,
    });
  }

  private removeAction(index: number) {
    if (!this.item.buttonActions) {
      return;
    }

    this.item.buttonActions.splice(index, 1);
  }
}
</script>

<style scoped>

.clearfix:before,
.clearfix:after {
  display: table;
  content: "";
}

.clearfix:after {
  clear: both
}

</style>
