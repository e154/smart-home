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
    <show-on v-model="item.showOn"/>
    <!-- /show on -->

    <!-- hide on-->
    <el-divider content-position="left">{{ $t('dashboard.editor.hideOn') }}</el-divider>
    <show-on v-model="item.hideOn"/>
    <!-- /hide on-->

    <!-- button options -->
    <div v-if="item.type !== 'button' && item.type !== 'chart'">
      <el-divider content-position="left">{{
          $t('dashboard.editor.buttonOptions')
        }}
      </el-divider>
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
    </div>
    <!-- /button options -->

  </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator';
import {CardItem, Core} from '@/views/dashboard/core';
import {ApiEntity, ApiImage} from '@/api/stub';
import EntitySearch from '@/views/entities/components/entity_search.vue';
import ImagePreview from '@/views/images/preview.vue';
import ShowOn from "@/views/dashboard/card_items/common/show-on.vue";

interface Action {
  value: string;
  label: string;
}

@Component({
  name: 'CommonEditor',
  components: {ShowOn, ImagePreview, EntitySearch}
})
export default class extends Vue {
  @Prop() private item!: CardItem;
  @Prop() private core!: Core;

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
    const entity = await this.core.fetchEntity(id);
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
      this.item.buttonActions[index].entity = await this.core.fetchEntity(entity.id);
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
