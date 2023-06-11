<template>
  <div class="app-container">
    <card-wrapper>
      <el-row :gutter="20">
        <el-col
          :span="24"
          :xs="24"
        >

          <el-tabs v-model="internal.activeTab">
            <el-tab-pane
              label="Main"
              name="main"
            >
              <el-form label-position="top"
                       ref="currentEntity"
                       :model="currentEntity"
                       :rules="rules"
                       style="width: 100%">
                <el-form-item :label="$t('entities.table.name')" prop="name">
                  <el-input v-model.trim="currentEntity.name"/>
                </el-form-item>
                <el-form-item :label="$t('entities.table.pluginName')" prop="pluginName">
                  <plugin-search
                    v-model="currentEntity.pluginName"
                    @update-value="changedPlugin"/>
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
                <el-form-item :label="$t('entities.table.autoLoad')" prop="description">
                  <el-switch v-model="currentEntity.autoLoad"></el-switch>
                </el-form-item>
                <el-form-item :label="$t('entities.table.parent')" prop="parent">
                  <entity-search
                    v-model="currentEntity.parent"
                    @update-value="changedParent"
                  />
                </el-form-item>
              </el-form>
            </el-tab-pane>

          </el-tabs>

        </el-col>
      </el-row>
      <el-row>
        <el-col :span="24" align="right">
          <el-button type="primary" @click.prevent.stop="create">{{ $t('main.create') }}</el-button>
          <el-button @click.prevent.stop="cancel">{{ $t('main.cancel') }}</el-button>
        </el-col>
      </el-row>
    </card-wrapper>
  </div>
</template>

<script lang="ts">

import {Component, Vue} from 'vue-property-decorator';
import api from '@/api/api';
import {ApiAttribute, ApiEntityShort, ApiImage, ApiNewEntityRequest, ApiPlugin, ApiScript} from '@/api/stub';
import router from '@/router';
import AttributesEditor from './components/attributes_editor.vue';
import Scripts from './components/scripts.vue';
import Actions from './components/actions.vue';
import States from './components/states.vue';
import ScriptSearch from '@/views/scripts/components/script_search.vue';
import PluginSearch from '@/views/plugins/plugin_search.vue';
import EntitySearch from './components/entity_search.vue';
import Metrics from './components/metrics.vue';
import {Form} from 'element-ui';
import AreaSearch from '@/views/areas/components/areas_search.vue';
import ImagePreview from '@/views/images/preview.vue';
import CardWrapper from '@/components/card-wrapper/index.vue';
import {createCard} from '@/views/entities/common';

@Component({
  name: 'EntityEditor',
  components: {
    CardWrapper,
    AttributesEditor,
    Scripts,
    Actions,
    States,
    ScriptSearch,
    EntitySearch,
    PluginSearch,
    Metrics,
    AreaSearch,
    ImagePreview
  }
})
export default class extends Vue {
  private internal = {
    activeTab: 'main',
    pluginOptions: undefined
  };

  // entity params
  private currentEntity: ApiNewEntityRequest = {
    name: '',
    description: '',
    pluginName: '',
    autoLoad: true,
    parent: undefined,
    attributes: new Map<string, ApiAttribute>(),
    settings: new Map<string, ApiAttribute>(),
    metrics: [],
    actions: [],
    states: [],
    scripts: []
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
    pluginName: [
      {required: true, trigger: 'blur'},
      {max: 255, trigger: 'blur'}
    ]
  };

  private changedScript(values: ApiScript[], event: any) {
    if (values) {
      this.$set(this.currentEntity, 'scripts', values);
    } else {
      this.$set(this.currentEntity, 'scripts', undefined);
    }
  }

  private changedPlugin(value: ApiPlugin, event: any) {
    if (value) {
      this.$set(this.currentEntity, 'pluginName', value.name);
      this.fetchPlugin();
    } else {
      this.$set(this.internal, 'pluginOptions', undefined);
    }
  }

  private changedParent(values: ApiEntityShort, event?: any) {
    if (values) {
      this.$set(this.currentEntity, 'parent', values);
    } else {
      this.$set(this.currentEntity, 'parent', '');
    }
  }


  private async fetchPlugin() {
    if (!this.currentEntity || !this.currentEntity.pluginName) {
      return;
    }
    const {data} = await api.v1.pluginServiceGetPluginOptions(this.currentEntity.pluginName);
    this.$set(this.internal, 'pluginOptions', data);
  }

  private async create() {
    (this.$refs.currentEntity as Form).validate(async valid => {
      if (!valid) {
        return;
      }
      const {data} = await createCard(this.currentEntity);
      if (data) {
        router.push({path: `/entities/edit/${data.id}`});
      }
    });
  }

  private cancel() {
    router.go(-1)
  }

  private onSelectImage(value: ApiImage, event?: any) {
    this.$set(this.currentEntity, 'image', value);
  }
}
</script>
