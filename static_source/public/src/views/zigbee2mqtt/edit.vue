<template>
  <div class="app-container" v-if="!listLoading">
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
                       ref="currentBridge"
                       :model="currentBridge"
                       :rules="rules"
                       style="width: 100%">
                <el-form-item :label="$t('zigbee2mqtt.table.name')" prop="name">
                  <el-input disabled v-model="currentBridge.name"/>
                </el-form-item>
                <el-form-item :label="$t('zigbee2mqtt.table.login')" prop="login">
                  <el-input v-model="currentBridge.login"/>
                </el-form-item>
                <el-form-item :label="$t('zigbee2mqtt.table.password')" prop="password">
                  <el-input v-model="currentBridge.password"
                            placeholder="Please input password"
                            show-password/>
                </el-form-item>
                <el-form-item :label="$t('zigbee2mqtt.table.permitJoin')" prop="permitJoin">
                  <el-switch v-model="currentBridge.permitJoin"></el-switch>
                </el-form-item>

                <el-form-item :label="$t('zigbee2mqtt.table.createdAt')">
                  <div>{{ currentBridge.createdAt | parseTime }}</div>
                </el-form-item>
                <el-form-item :label="$t('zigbee2mqtt.table.updatedAt')">
                  <div>{{ currentBridge.updatedAt | parseTime }}</div>
                </el-form-item>
              </el-form>

            </el-tab-pane>

            <el-tab-pane
              label="Devices"
              name="devices"
            >

              <devices
                :id="id"
              />

            </el-tab-pane>
          </el-tabs>

        </el-col>
      </el-row>
      <el-row style="margin-top: 20px">
        <el-col :span="24" align="right">
          <el-button type="primary" @click.prevent.stop="save">{{ $t('main.save') }}</el-button>
          <el-button @click.prevent.stop="fetchBridge">{{ $t('main.load_from_server') }}</el-button>
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

import { Component, Prop, Vue } from 'vue-property-decorator'
import api from '@/api/api'
import { ApiZigbee2Mqtt } from '@/api/stub'
import router from '@/router'
import { Form } from 'element-ui'
import Devices from '@/views/zigbee2mqtt/components/devices.vue'
import CardWrapper from '@/components/card-wrapper/index.vue'

@Component({
  name: 'Editor',
  components: {
    CardWrapper,
    Devices
  }
})
export default class extends Vue {
  @Prop({ required: true }) private id!: number;

  private listLoading = true;

  private internal = {
    activeTab: 'devices'
  };

  // entity params
  private currentBridge: ApiZigbee2Mqtt = {
    name: '',
    password: ''
  };

  private rules = {
    name: [
      { required: true, trigger: 'blur' },
      { min: 4, max: 255, trigger: 'blur' }
    ],
    login: [
      { required: false, trigger: 'blur' },
      { max: 255, trigger: 'blur' }
    ]
  };

  created() {
    this.fetchBridge()
  }

  private async fetchBridge() {
    this.listLoading = true
    const { data } = await api.v1.zigbee2MqttServiceGetZigbee2MqttBridge(this.id)
    this.currentBridge = data
    this.listLoading = false
  }

  private async save() {
    (this.$refs.currentBridge as Form).validate(async valid => {
      if (!valid) {
        return
      }
      const bridge = {
        name: this.currentBridge.name,
        login: this.currentBridge.login,
        password: this.currentBridge.password,
        permitJoin: this.currentBridge.permitJoin,
        baseTopic: this.currentBridge.baseTopic
      }
      const { data } = await api.v1.zigbee2MqttServiceUpdateBridgeById(this.id, bridge)
      if (data) {
        this.$notify({
          title: 'Success',
          message: 'bridge updated successfully',
          type: 'success',
          duration: 2000
        })
      }
    })
  }

  private async remove() {
    const { data } = await api.v1.zigbee2MqttServiceDeleteBridgeById(this.id)
    this.$notify({
      title: 'Success',
      message: 'Delete Successfully',
      type: 'success',
      duration: 2000
    })
    router.push({ path: '/zigbee2mqtt' })
  }

  private cancel() {
    router.go(-1)
  }
}
</script>
