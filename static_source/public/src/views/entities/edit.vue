<template>
  <div class="app-container" v-if="!listLoading">
    <card-wrapper>
      <el-row :gutter="20">
        <el-col
          :span="24"
          :xs="24"
        >

          <el-tabs v-model="internal.activeTab" @tab-click="handleTabChange">
            <el-tab-pane
              label="Main"
              name="main"
            >
              <el-form label-position="top"
                       ref="currentEntity"
                       :model="currentEntity"
                       :rules="rules"
                       style="width: 100%">
                <el-form-item :label="$t('entities.table.id')" prop="id">
                  <el-input disabled v-model.trim="currentEntity.id"/>
                </el-form-item>
                <el-form-item :label="$t('entities.table.pluginName')" prop="pluginName">
                  <el-input disabled v-model.trim="currentEntity.pluginName"/>
                </el-form-item>
                <el-form-item :label="$t('entities.table.scripts')" prop="scripts">
                  <scripts
                    v-model="currentEntity.scripts"
                    @update-value="changedScript"
                  />
                </el-form-item>
                <el-form-item :label="$t('entities.table.description')" prop="description">
                  <el-input v-model.trim="currentEntity.description"/>
                </el-form-item>
                <el-form-item :label="$t('entities.table.icon')" prop="icon">
                  <el-input v-model.trim="currentEntity.icon"/>
                </el-form-item>
                <el-form-item :label="$t('entities.table.image')" prop="image">
                  <image-preview :image="currentEntity.image" @on-select="onSelectImage"/>
                </el-form-item>
                <el-form-item :label="$t('entities.table.autoLoad')" prop="autoLoad">
                  <el-switch v-model="currentEntity.autoLoad"></el-switch>
                </el-form-item>
                <el-form-item :label="$t('entities.table.parent')" prop="parent">
                  <entity-search
                    v-model="currentEntity.parent"
                    @update-value="changedParent"
                  />
                </el-form-item>
                <el-form-item :label="$t('entities.table.area')" prop="area">
                  <area-search
                    :multiple=false
                    v-model="currentEntity.area"
                    @update-value="changedArea"
                  />
                </el-form-item>
                <el-form-item :label="$t('entities.table.createdAt')">
                  <div>{{ currentEntity.createdAt | parseTime }}</div>
                </el-form-item>
                <el-form-item :label="$t('entities.table.updatedAt')">
                  <div>{{ currentEntity.updatedAt | parseTime }}</div>
                </el-form-item>
              </el-form>
            </el-tab-pane>

            <el-tab-pane
              label="Actions"
              name="actions"
              v-if="internal.pluginOptions.actorCustomActions || Object.keys(internal.pluginOptions.actorActions).length"
            >
              <actions
                v-model="currentEntity.actions"
                :settings="internal.pluginOptions.actorActions"
                :customActions="internal.pluginOptions.actorCustomActions"
                v-on:call-action="callAction"
              />
            </el-tab-pane>

            <el-tab-pane
              label="States"
              name="states"
              v-if="internal.pluginOptions.actorCustomStates || Object.keys(internal.pluginOptions.actorStates).length"
            >
              <states
                v-model="currentEntity.states"
                :settings="internal.pluginOptions.actorStates"
                :customStates="internal.pluginOptions.actorCustomStates"
                @update-value="changedStates"
                v-on:set-state="setState"
              />
            </el-tab-pane>

            <el-tab-pane
              label="Attributes"
              name="attributes"
              v-if="internal.pluginOptions.actorCustomAttrs || Object.keys(internal.pluginOptions.actorAttrs).length"
            >
              <AttributesEditor
                v-model="internal.attributes"
                :attrs="internal.pluginOptions.actorAttrs"
                :customAttrs="internal.pluginOptions.actorCustomAttrs"
                @update-value="changedAttributes($event, 'attributes')"
              />
            </el-tab-pane>

            <el-tab-pane
              label="Settings"
              name="settings"
              v-if="internal.pluginOptions.actorCustomSetts || Object.keys(internal.pluginOptions.actorSetts).length"
            >
              <AttributesEditor
                v-model="internal.settings"
                :attrs="internal.pluginOptions.actorSetts"
                :customAttrs="internal.pluginOptions.actorCustomSetts"
                @update-value="changedAttributes($event, 'settings')"
              />
            </el-tab-pane>

            <el-tab-pane
              label="Metrics"
              name="metrics"
            >
              <metrics
                v-model="internal.metrics"
                @update-value="changedMetrics"
              />
            </el-tab-pane>

            <el-tab-pane
              label="Storage"
              name="storage"
            >
              <storage :entity="currentEntity"/>
            </el-tab-pane>

            <el-tab-pane
              label="Current event"
              name="current_event"
            >
              <el-button type="default" @click.prevent.stop="requestCurrentState()" style="margin-bottom: 20px"><i
                class="el-icon-refresh"/> {{ $t('dashboard.editor.getEvent') }}
              </el-button>

              <json-editor
                ref="jsoneditor"
                :value="lastEvent"
              />
            </el-tab-pane>

          </el-tabs>

        </el-col>
      </el-row>
      <el-row style="margin-top: 20px">
        <el-col :span="24" align="right">

          <export-tool
            :title="$t('main.export')"
            :visible="showExport"
            :value="internal.exportValue"
            @on-close="showExport=false"/>

          <el-button type="primary" @click.prevent.stop="save">{{ $t('main.save') }}</el-button>
          <el-button type="primary" icon="el-icon-refresh" @click.prevent.stop='restart'>{{
              $t('main.restart')
            }}
          </el-button>
          <el-button type="primary" icon="el-icon-document" @click.prevent.stop='_export'>{{
              $t('main.export')
            }}
          </el-button>
          <el-button @click.prevent.stop="fetchEntity">{{ $t('main.load_from_server') }}</el-button>
          <el-button @click.prevent.stop="cancel">{{ $t('main.cancel') }}</el-button>
          <el-popconfirm
            :confirm-button-text="$t('main.ok')"
            :cancel-button-text="$t('main.no')"
            icon="el-icon-info"
            icon-color="red"
            style="margin-left: 10px;"
            :title="$t('main.are_you_sure_to_do_want_this?')"
            v-on:confirm="remove"
          >
            <el-button type="danger" icon="el-icon-delete" slot="reference">{{ $t('main.remove') }}</el-button>
          </el-popconfirm>
        </el-col>
      </el-row>
    </card-wrapper>
  </div>
</template>

<script lang="ts">

import {Component, Prop, Vue} from 'vue-property-decorator';
import api from '@/api/api';
import {
  ApiArea,
  ApiAttribute,
  ApiEntity,
  ApiEntityShort,
  ApiEntityState,
  ApiGetPluginOptionsResult,
  ApiImage,
  ApiMetric,
  ApiScript,
} from '@/api/stub';
import router from '@/router';
import AttributesEditor from './components/attributes_editor.vue';
import Scripts from './components/scripts.vue';
import Actions from './components/actions.vue';
import States from './components/states.vue';
import ScriptSearch from '@/views/scripts/components/script_search.vue';
import AreaSearch from '@/views/areas/components/areas_search.vue';
import EntitySearch from './components/entity_search.vue';
import Metrics from './components/metrics.vue';
import {Form} from 'element-ui';
import ImagePreview from '@/views/images/preview.vue';
import CardWrapper from '@/components/card-wrapper/index.vue';
import ExportTool from '@/components/export-tool/index.vue';
import stream from '@/api/stream';
import {UUID} from 'uuid-generator-ts';
import JsonEditor from '@/components/JsonEditor/index.vue';
import {EventStateChange} from '@/api/stream_types';
import Storage from '@/views/entities/components/storage.vue';

@Component({
  name: 'EntityEditor',
  components: {
    Storage,
    JsonEditor,
    ExportTool,
    CardWrapper,
    AttributesEditor,
    Scripts,
    Actions,
    States,
    ScriptSearch,
    EntitySearch,
    Metrics,
    AreaSearch,
    ImagePreview
  }
})
export default class extends Vue {
  @Prop({required: true}) private id!: string;

  private listLoading: boolean = true;

  // id for streaming subscribe
  private currentID: string = '';

  private internal: {
    activeTab: string,
    pluginOptions?: ApiGetPluginOptionsResult,
    exportValue?: ApiEntity,
    attributes: ApiAttribute[],
    settings: ApiAttribute[],
    metrics: ApiMetric[],
  } = {
    activeTab: 'main',
    pluginOptions: {} as ApiGetPluginOptionsResult,
    exportValue: {} as ApiEntity,
    attributes: [],
    settings: [],
    metrics: [],
  };

  // entity params
  private currentEntity: ApiEntity = {
    pluginName: '',
    autoLoad: true,
    parent: {},
    metrics: []
  };

  private rules = {
    name: [
      {required: true, trigger: 'blur'},
      {min: 4, max: 255, trigger: 'blur'}
    ],
    description: [
      {required: false, trigger: 'blur'},
      {max: 255, trigger: 'blur'}
    ],
    plugin: [
      {required: false, trigger: 'blur'},
      {max: 255, trigger: 'blur'}
    ]
  };

  created() {
    let uuid = new UUID();
    this.currentID = uuid.getDashFreeUUID();

    this.fetchEntity();

    setTimeout(() => {
      stream.subscribe('state_changed', this.currentID, this.onStateChanged);
    }, 1000);
  }

  private destroyed() {
    stream.unsubscribe('state_changed', this.currentID);
  }

  private onStateChanged(event: EventStateChange) {
    if (event.entity_id != this.currentEntity.id) {
      return;
    }

    this.lastEvent = event;
  }

  private changedAttributes(attributes: ApiAttribute[], event: any) {
    if (event == 'attributes') {
      if (attributes) {
        this.$set(this.internal, 'attributes', attributes);
        // this.internal = Object.assign({}, this.internal, {attributes: attributes})
      } else {
        this.$set(this.internal, 'attributes', undefined);
      }
    } else if (event == 'settings') {
      if (attributes) {
        this.$set(this.internal, 'settings', attributes);
        // this.internal = Object.assign({}, this.internal, {settings: value})
      } else {
        this.$set(this.internal, 'settings', undefined);
      }
    }
  }

  private changedScript(values: ApiScript[], event: any) {
    if (values) {
      this.$set(this.currentEntity, 'scripts', values);
    } else {
      this.$set(this.currentEntity, 'scripts', undefined);
    }
  }

  private changedParent(values: ApiEntityShort, event?: any) {
    console.log('changedParent', values);
    if (values) {
      this.$set(this.currentEntity, 'parent', values);
    } else {
      this.$set(this.currentEntity, 'parent', undefined);
    }
  }

  private changedArea(values: ApiArea, event?: any) {
    if (values) {
      this.$set(this.currentEntity, 'area', values);
    } else {
      this.$set(this.currentEntity, 'area', undefined);
    }
  }

  private changedStates(values: ApiEntityState[], event?: any) {
    if (values) {
      this.$set(this.currentEntity, 'states', values);
    } else {
      this.$set(this.currentEntity, 'states', undefined);
    }
  }

  private changedMetrics(values: ApiMetric[], event?: any) {
    if (values) {
      this.$set(this.internal, 'metrics', values);
    } else {
      this.$set(this.internal, 'metrics', undefined);
    }
  }

  private async fetchEntity() {
    this.listLoading = true;
    const {data} = await api.v1.entityServiceGetEntity(this.id);
    this.currentEntity = data;

    // attributes
    let attr: ApiAttribute[] = [];
    if (this.currentEntity.attributes) {
      for (const key in this.currentEntity.attributes) {
        attr.push(this.currentEntity.attributes[key]);
      }
    }
    this.$set(this.internal, 'attributes', attr);

    // settings
    let setts: ApiAttribute[] = [];
    if (this.currentEntity.settings) {
      for (const key in this.currentEntity.settings) {
        setts.push(this.currentEntity.settings[key]);
      }
    }
    this.$set(this.internal, 'settings', setts);

    // metrics
    let metrics: ApiMetric[] = [];
    if (this.currentEntity.metrics) {
      for (const key in this.currentEntity.metrics) {
        metrics.push(this.currentEntity.metrics[key]);
      }
    }
    this.$set(this.internal, 'metrics', metrics);

    await this.fetchPlugin();
    this.listLoading = false;
  }

  private async fetchPlugin() {
    this.listLoading = true;
    const {data} = await api.v1.pluginServiceGetPluginOptions(this.currentEntity.pluginName);
    this.$set(this.internal, 'pluginOptions', data);
    this.listLoading = false;
  }

  private prepareSave(): ApiEntity {
    let entity: ApiEntity = {
      pluginName: this.currentEntity.pluginName,
      description: this.currentEntity.description,
      area: this.currentEntity.area,
      icon: this.currentEntity.icon,
      image: this.currentEntity.image,
      autoLoad: this.currentEntity.autoLoad,
      parent: this.currentEntity.parent || undefined,
      actions: [],
      metrics: [],
      states: [],
      scripts: this.currentEntity.scripts,
    };

    // attributes
    let attributes: { [key: string]: ApiAttribute } = {};
    for (const index in this.internal.attributes) {
      attributes[this.internal.attributes[index].name] = this.internal.attributes[index];
    }
    this.$set(entity, 'attributes', attributes);

    // settings
    let settings: { [key: string]: ApiAttribute } = {};
    for (const index in this.internal.settings) {
      settings[this.internal.settings[index].name] = this.internal.settings[index];
    }
    this.$set(entity, 'settings', settings);

    // update image
    if (entity.image) {
      entity.image = {id: entity.image.id};
    }

    // update actions
    for (const i in this.currentEntity.actions) {
      let action = Object.assign({}, this.currentEntity.actions[<any> i]);
      if (action.image?.id) {
        action.image = {id: action.image?.id};
      }
      entity.actions?.push(action);
    }

    // update states
    for (const i in this.currentEntity.states) {
      let state = Object.assign({}, this.currentEntity.states[<any> i]);
      if (state.image?.id) {
        state.image = {id: state.image?.id};
      }
      entity.states?.push(state);
    }

    // metrics
    let metrics: ApiMetric[] = [];
    if (this.internal.metrics) {
      for (const key in this.internal.metrics) {
        metrics.push(this.internal.metrics[key]);
      }
    }
    entity.metrics = metrics;

    return entity;
  }

  private async save() {
    (this.$refs.currentEntity as Form).validate(async valid => {
      if (!valid) {
        return;
      }
      const entity = this.prepareSave();
      const {data} = await api.v1.entityServiceUpdateEntity(this.id, entity);
      if (data) {
        this.$notify({
          title: 'Success',
          message: 'entity updated successfully',
          type: 'success',
          duration: 2000
        });
      }
    });
  }

  private async remove() {
    const {data} = await api.v1.entityServiceDeleteEntity(this.id);
    this.$notify({
      title: 'Success',
      message: 'Delete Successfully',
      type: 'success',
      duration: 2000
    });
    router.push({path: `/entities`});
  }

  private async callAction(name: string) {
    await api.v1.interactServiceEntityCallAction({id: this.currentEntity.id, name: name});
    this.$notify({
      title: 'Success',
      message: 'Call Successfully',
      type: 'success',
      duration: 2000
    });
  }

  private async setState(name: string) {
    await api.v1.developerToolsServiceEntitySetState({id: this.currentEntity.id, name: name});
    this.$notify({
      title: 'Success',
      message: 'Call Successfully',
      type: 'success',
      duration: 2000
    });
  }

  private onSelectImage(value: ApiImage, event?: any) {
    this.$set(this.currentEntity, 'image', value);
  }

  private async restart() {
    await api.v1.developerToolsServiceReloadEntity({id: this.currentEntity.id});
    this.$notify({
      title: 'Success',
      message: 'entity reloaded successfully',
      type: 'success',
      duration: 2000
    });
  }

  private cancel() {
    router.go(-1);
  }

  private showExport: boolean = false;

  private _export() {
    let entity: any;
    entity = this.prepareSave();
    entity.id = this.currentEntity.id;
    entity.name = this.currentEntity.id?.replaceAll('.' + entity.pluginName, '');
    this.internal.exportValue = entity;
    this.showExport = true;
  }

  private lastEvent: EventStateChange = {} as EventStateChange;

  private handleTabChange(tab: any, event: any) {
    if (this.internal.activeTab == 'current_event') {
      this.requestCurrentState();
    }
  }

  private requestCurrentState() {
    stream.send({
      id: UUID.createUUID(),
      query: 'event_get_last_state',
      body: btoa(JSON.stringify({'entity_id': this.currentEntity.id}))
    });
  }
}
</script>


