<template>
  <div class="app-container">
    <card-wrapper>
      <el-row :gutter="20">
        <el-col
          :span="24"
          :xs="24"
        >

          <el-form label-position="top"
                   ref="currentBridge"
                   :model="currentBridge"
                   :rules="rules"
                   style="width: 100%">
            <el-form-item :label="$t('zigbee2mqtt.table.name')" prop="name">
              <el-input v-model="currentBridge.name"/>
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

          </el-form>

        </el-col>
        <el-col
          :span="18"
          :xs="24"
        >

        </el-col>
      </el-row>
      <el-row>
        <el-col :span="24" align="right">
          <el-button type="primary" @click.prevent.stop="save">{{ $t('main.save') }}</el-button>
          <el-button @click.prevent.stop="cancel">{{ $t('main.cancel') }}</el-button>
        </el-col>
      </el-row>
    </card-wrapper>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import api from '@/api/api'
import { ApiNewZigbee2MqttRequest } from '@/api/stub'
import router from '@/router'
import { Form } from 'element-ui'
import CardWrapper from '@/components/card-wrapper/index.vue'

// HACK: have to use script-loader to load jsonlint
/* eslint-disable import/no-webpack-loader-syntax */
require('script-loader!jsonlint')

@Component({
  name: 'New',
  components: { CardWrapper }
})
export default class extends Vue {
  private currentBridge: ApiNewZigbee2MqttRequest = {
    name: '',
    login: '',
    password: '',
    permitJoin: true
  };

  private rules = {
    name: [
      { required: true, trigger: 'blur' },
      { min: 4, max: 255, trigger: 'blur' }
    ],
    description: [
      { required: false, trigger: 'blur' },
      { max: 255, trigger: 'blur' }
    ]
  };

  private async save() {
    (this.$refs.currentBridge as Form).validate(async valid => {
      if (!valid) {
        return
      }
      const { data } = await api.v1.zigbee2MqttServiceAddZigbee2MqttBridge(this.currentBridge)
      if (data) {
        this.$notify({
          title: 'Success',
          message: 'bridge created successfully',
          type: 'success',
          duration: 2000
        })
        router.push({ path: `/zigbee2mqtt/edit/${data.id}` })
      }
    })
  }

  private cancel() {
    router.go(-1)
  }
}
</script>
