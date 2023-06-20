<template>
  <div class="app-container" v-if="!isLoading">
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
                       ref="currentArea"
                       :model="currentPlugin"
                       style="width: 100%">
                <el-form-item :label="$t('plugins.table.name')" prop="name">
                  <el-input disabled v-model.trim="currentPlugin.name"/>
                </el-form-item>

                <el-form-item :label="$t('plugins.table.version')" prop="version">
                  <el-input disabled v-model.trim="currentPlugin.version"/>
                </el-form-item>

                <el-form-item :label="$t('plugins.table.enabled')" prop="enabled">
                  <el-switch v-model="currentPlugin.enabled"
                             v-on:change="updateItem(currentPlugin)"
                             :disabled="currentPlugin.system"></el-switch>
                </el-form-item>

                <el-form-item :label="$t('plugins.table.system')" prop="system">
                  <el-switch v-model="currentPlugin.system" disabled></el-switch>
                </el-form-item>

                <el-form-item :label="$t('plugins.options.triggers')" prop="system">
                  <el-switch v-model="currentPlugin.options.triggers" disabled></el-switch>
                </el-form-item>

                <el-form-item :label="$t('plugins.options.actors')" prop="system">
                  <el-switch v-model="currentPlugin.options.actors" disabled></el-switch>
                </el-form-item>

                <el-form-item :label="$t('plugins.options.actorCustomAttrs')" prop="system">
                  <el-switch v-model="currentPlugin.options.actorCustomAttrs" disabled></el-switch>
                </el-form-item>

                <el-form-item :label="$t('plugins.options.actorCustomActions')" prop="system">
                  <el-switch v-model="currentPlugin.options.actorCustomActions" disabled></el-switch>
                </el-form-item>

                <el-form-item :label="$t('plugins.options.actorCustomStates')" prop="system">
                  <el-switch v-model="currentPlugin.options.actorCustomStates" disabled></el-switch>
                </el-form-item>

                <el-form-item :label="$t('plugins.options.actorCustomSetts')" prop="system">
                  <el-switch v-model="currentPlugin.options.actorCustomSetts" disabled></el-switch>
                </el-form-item>

              </el-form>

            </el-tab-pane>

            <el-tab-pane
              :label="$t('plugins.actor')"
              :disabled="!showActorTabIf()"
              name="actorAttrs"
            >
              <!-- Attributes -->
              <el-row style="margin-top: 20px"
                      v-if="internal.actorAttrs && Object.keys(internal.actorAttrs).length">
                <el-col>
                  {{ $t('plugins.actorAttrs') }}
                </el-col>
              </el-row>
              <el-row v-if="internal.actorAttrs && Object.keys(internal.actorAttrs).length">
                <el-col>
                  <attributes-viewer v-model="internal.actorAttrs"/>
                </el-col>
              </el-row>
              <!-- /Attributes -->

              <!-- Actions -->
              <el-row style="margin-top: 20px"
                      v-if="this.internal.actorActions && Object.keys(this.internal.actorActions).length">
                <el-col>
                  {{ $t('plugins.actorActions') }}
                </el-col>
              </el-row>
              <el-row
                v-if="this.internal.actorActions && Object.keys(this.internal.actorActions).length">
                <el-col>

                  <el-table
                    key="key"
                    :data="this.internal.actorActions"
                    style="width: 100%;"
                  >
                    <el-table-column
                      :label="$t('entities.table.name')"
                      prop="name"
                      align="left"
                    >
                      <template slot-scope="{row}">
                        <div>{{ row.name }}</div>
                      </template>
                    </el-table-column>

                    <el-table-column
                      :label="$t('entities.table.image')"
                      prop="image"
                      align="center"
                    >
                      <template slot-scope="{row}">
                        <i v-if="row.image" :class="'el-icon-check'"/>
                        <i v-if="!row.image" :class="'el-icon-minus'"/>
                      </template>
                    </el-table-column>

                    <el-table-column
                      :label="$t('entities.table.icon')"
                      prop="icon"
                      align="center"
                    >
                      <template slot-scope="{row}">
                        <i v-if="row.icon" :class="'el-icon-check'"/>
                        <i v-if="!row.icon" :class="'el-icon-minus'"/>
                      </template>
                    </el-table-column>

                    <el-table-column
                      :label="$t('entities.table.script')"
                      prop="script"
                      align="center"
                    >
                      <template slot-scope="{row}">
                        <i v-if="row.script" :class="'el-icon-check'"/>
                        <i v-if="!row.script" :class="'el-icon-minus'"/>
                      </template>
                    </el-table-column>

                    <el-table-column
                      :label="$t('entities.table.description')"
                      prop="description"
                      align="left"
                    >
                      <template slot-scope="{row}">
                        <div>{{ row.description }}</div>
                      </template>
                    </el-table-column>

                  </el-table>

                </el-col>
              </el-row>
              <!-- /Actions -->

              <!-- States -->
              <el-row style="margin-top: 20px"
                      v-if="this.internal.actorStates && Object.keys(this.internal.actorStates).length">
                <el-col>
                  {{ $t('plugins.actorStates') }}
                </el-col>
              </el-row>
              <el-row v-if="this.internal.actorStates && Object.keys(this.internal.actorStates).length">
                <el-col>

                  <el-table
                    key="key"
                    :data="this.internal.actorStates"
                    style="width: 100%;"
                  >
                    <el-table-column
                      :label="$t('entities.table.name')"
                      prop="name"
                      align="left"
                    >
                      <template slot-scope="{row}">
                        <div>{{ row.name }}</div>
                      </template>
                    </el-table-column>

                    <el-table-column
                      :label="$t('entities.table.image')"
                      prop="image"
                      align="center"
                    >
                      <template slot-scope="{row}">
                        <i v-if="row.image" :class="'el-icon-check'"/>
                        <i v-if="!row.image" :class="'el-icon-minus'"/>
                      </template>
                    </el-table-column>

                    <el-table-column
                      :label="$t('entities.table.icon')"
                      prop="icon"
                      align="center"
                    >
                      <template slot-scope="{row}">
                        <i v-if="row.icon" :class="'el-icon-check'"/>
                        <i v-if="!row.icon" :class="'el-icon-minus'"/>
                      </template>
                    </el-table-column>

                    <el-table-column
                      :label="$t('entities.table.description')"
                      prop="description"
                      align="left"
                    >
                      <template slot-scope="{row}">
                        <div>{{ row.description }}</div>
                      </template>
                    </el-table-column>

                  </el-table>

                </el-col>
              </el-row>
              <!-- /States -->

              <!-- Settings -->
              <el-row style="margin-top: 20px"
                      v-if="currentPlugin.options.actorSetts && Object.keys(currentPlugin.options.actorSetts).length">
                <el-col>
                  {{ $t('plugins.settings') }}
                </el-col>
              </el-row>
              <el-row v-if="currentPlugin.options.actorSetts && Object.keys(currentPlugin.options.actorSetts).length">
                <el-col>
                  <attributes-viewer v-model="currentPlugin.options.actorSetts"/>
                </el-col>
              </el-row>
              <!-- /Settings -->

            </el-tab-pane>

            <el-tab-pane
              :label="$t('plugins.settings')"
              name="settings"
              :disabled="!(currentPlugin.options.setts && Object.keys(currentPlugin.options.setts).length)"
            >

              <!-- plugin settings -->
              <el-row>
                <el-col>
                  <el-table
                    :data="internal.settings"
                    style="width: 100%">
                    <el-table-column
                      prop="name"
                      :label="$t('entities.table.name')"
                      width="180px">

                      <template slot-scope="{row}">
                        <el-input type="string" v-model="row.name" disabled/>
                      </template>

                    </el-table-column>
                    <el-table-column
                      prop="type"
                      :label="$t('entities.table.type')"
                      width="150px">

                      <template slot-scope="{row}">
                        <el-select v-model="row.type" placeholder="please select type" disabled>
                          <el-option label="String" value="STRING"></el-option>
                          <el-option label="Int" value="INT"></el-option>
                          <el-option label="Float" value="FLOAT"></el-option>
                          <el-option label="Bool" value="BOOL"></el-option>
                          <el-option label="Array" value="ARRAY"></el-option>
                          <el-option label="Time" value="TIME"></el-option>
                          <el-option label="Map" value="MAP"></el-option>
                          <el-option label="Image" value="IMAGE"></el-option>
                        </el-select>
                      </template>

                    </el-table-column>
                    <el-table-column
                      width="auto"
                      :label="$t('entities.table.value')"
                    >

                      <template slot-scope="{row}">
                        <div v-if="row.type === 'STRING'">
                          <el-input type="string" v-model="row.string"/>
                        </div>
                        <div v-if="row.type === 'IMAGE'">
                          <el-input type="string" v-model="row.imageUrl"/>
                        </div>
                        <div v-if="row.type === 'INT'">
                          <el-input type="number" v-model="row.int"/>
                        </div>
                        <div v-if="row.type === 'FLOAT'">
                          <el-input type="number" v-model="row.float"/>
                        </div>
                        <div v-if="row.type === 'ARRAY'">
                          <el-input type="string" v-model="row.array"/>
                        </div>
                        <el-select v-model="row.bool"
                                   placeholder="please select value"
                                   v-if="row.type === 'BOOL'"
                        >
                          <el-option label="TRUE" :value="true"></el-option>
                          <el-option label="FALSE" :value="false"></el-option>
                        </el-select>

                        <div v-if="row.type === 'TIME'">
                          <el-input type="string" v-model="row.time"/>
                        </div>

                        <div v-if="row.type === 'MAP'">
                          <el-input type="string" v-model="row.map"/>
                        </div>
                      </template>

                    </el-table-column>

                  </el-table>
                </el-col>
              </el-row>
              <!-- /plugin settings -->

              <el-row style="margin-top: 20px">
                <el-col align="right">
                  <el-button type="primary" @click="saveSetting()">{{ $t('main.update') }}</el-button>
                  <el-button @click="resetForm()">{{ $t('main.cancel') }}</el-button>
                </el-col>
              </el-row>

            </el-tab-pane>


          </el-tabs>

        </el-col>
      </el-row>


    </card-wrapper>
  </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator'
import api from '@/api/api'
import {
  ApiAttribute,
  ApiGetPluginOptionsResultEntityAction,
  ApiGetPluginOptionsResultEntityState,
  ApiPlugin,
  ApiPluginShort
} from '@/api/stub'
import CardWrapper from '@/components/card-wrapper/index.vue'
import AttributesEditor from "@/views/entities/components/attributes_editor.vue";
import AttributesViewer from "@/views/plugins/attributes_viewer.vue";


@Component({
  name: 'PluginEdit',
  components: {AttributesViewer, AttributesEditor, CardWrapper}
})
export default class extends Vue {
  @Prop({required: true}) private name!: string;

  private internal: {
    settings: ApiAttribute[],
    actorStates: ApiGetPluginOptionsResultEntityState[],
    actorActions: ApiGetPluginOptionsResultEntityAction[],
    actorAttrs: ApiAttribute[],
    activeTab: string
  } = {
    activeTab: 'main',
    settings: [],
    actorStates: [],
    actorActions: [],
    actorAttrs: []
  }

  private isLoading = true;
  private currentPlugin?: ApiPlugin;

  private async fetch() {
    this.isLoading = true
    const {data} = await api.v1.pluginServiceGetPlugin(this.name)
    this.currentPlugin = data

    this.internal.settings = [];
    this.internal.actorStates = [];

    // ApiAttribute
    if (data.options?.setts) {
      for (const key in data.options.setts) {
        let st = data.options.setts[key]
        if (data.settings[key]) {
          st = data.settings[key]
        }
        this.internal.settings.push(st)
      }
    }

    // actor states
    if (data.options?.actorStates) {
      for (const key in data.options.actorStates) {
        this.internal.actorStates.push(data.options.actorStates[key])
      }
    }

    // actor actions
    if (data.options?.actorActions) {
      for (const key in data.options.actorActions) {
        this.internal.actorActions.push(data.options.actorActions[key])
      }
    }

    // actor attributes
    if (data.options?.actorAttrs) {
      for (const key in data.options.actorAttrs) {
        this.internal.actorAttrs.push(data.options.actorAttrs[key])
      }
    }

    this.isLoading = false
  }

  created() {
    this.fetch()
  }

  private async saveSetting() {
    if (!this.currentPlugin) {
      return
    }
    let settings: { [key: string]: ApiAttribute } = {};
    for (const index in this.internal.settings) {
      settings[this.internal.settings[index].name] = this.internal.settings[index];
    }
    await api.v1.pluginServiceUpdatePluginSettings(this.currentPlugin.name, {settings: settings})
    this.fetch()
  }

  private async updateItem(plugin: ApiPluginShort) {
    if (plugin.enabled) {
      await api.v1.pluginServiceEnablePlugin(plugin.name)
    } else {
      await api.v1.pluginServiceDisablePlugin(plugin.name)
    }
    this.fetch()
  }

  private resetForm() {
    this.fetch()
  }

  private showActorTabIf() {
    if (Object.keys(this.internal.actorAttrs).length ||
      Object.keys(this.internal.actorActions).length ||
      Object.keys(this.internal.actorStates).length ||
      Object.keys(this.currentPlugin?.options?.actorSetts || {}).length) {
      return true
    }
    return false
  }

}
</script>

<style lang="scss" scoped>

.cursor-pointer {
  cursor: pointer;
}

.pagination-container {

}
</style>
