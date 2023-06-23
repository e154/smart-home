<template>
  <div class="app-container" v-if="!listLoading">
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
              <el-input disabled v-model.trim="currentVariable.name"/>
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
import {ApiVariable} from '@/api/stub';
import router from '@/router';
import {Form} from 'element-ui';
import CardWrapper from '@/components/card-wrapper/index.vue';

@Component({
  name: 'Editor',
  components: {CardWrapper}
})
export default class extends Vue {
  @Prop({required: true}) private name!: string;

  private listLoading = true;
  private currentVariable: ApiVariable = {
    name: '',
    value: ''
  };

  private rules = {
    name: [
      {required: true, trigger: 'blur'},
      {min: 4, max: 255, trigger: 'blur'}
    ],
    value: [
      {required: false, trigger: 'blur'},
      {max: 255, trigger: 'blur'}
    ]
  };

  created() {
    this.fetch();
  }

  private async fetch() {
    this.listLoading = true;
    const {data} = await api.v1.variableServiceGetVariableByName(this.name);
    this.currentVariable = data;
    this.listLoading = false;
  }

  private async save() {
    (this.$refs.currentVariable as Form).validate(async valid => {
      if (!valid) {
        return;
      }
      const {data} = await api.v1.variableServiceUpdateVariable(this.name, {
        name: this.currentVariable.name!,
        value: this.currentVariable.value,
      });
    });
  }

  private cancel() {
    router.push({ path: '/etc/variables' })
  }

  private async remove() {
    const {data} = await api.v1.variableServiceDeleteVariable(this.name);
    this.$notify({
      title: 'Success',
      message: 'Delete Successfully',
      type: 'success',
      duration: 2000
    });
    router.go(-1);
  }
}
</script>
