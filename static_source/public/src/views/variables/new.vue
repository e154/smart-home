<template>
  <div class="app-container">
    <card-wrapper>
      <el-row :gutter="20">
        <el-col
          :span="24"
          :xs="24"
        >

          <el-form label-position="top"
                   ref="currentVariable"
                   :model="currentVariable"
                   :rules="rules"
                   style="width: 100%">
            <el-form-item :label="$t('variables.table.name')" prop="name">
              <el-input v-model.trim="currentVariable.name"/>
            </el-form-item>

            <el-form-item :label="$t('variables.table.value')" prop="value">
              <el-input v-model.trim="currentVariable.value"/>
            </el-form-item>

          </el-form>

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
import {ApiArea, ApiVariable} from '@/api/stub'
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
  private currentVariable: ApiVariable = {
    name: '',
    value: ''
  };

  private rules = {
    name: [
      { required: true, trigger: 'blur' },
      { min: 4, max: 255, trigger: 'blur' }
    ],
    value: [
      { required: false, trigger: 'blur' },
      { max: 255, trigger: 'blur' }
    ]
  };

  private async save() {
    (this.$refs.currentVariable as Form).validate(async valid => {
      if (!valid) {
        return
      }
      const { data } = await api.v1.variableServiceAddVariable(this.currentVariable)
      router.push({ path: `/etc/variables/edit/${data.name}` })
    })
  }

  private cancel() {
    router.go(-1)
  }
}
</script>
