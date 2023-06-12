<template>
  <div class="app-container">
    <card-wrapper>
      <el-row :gutter="20">
        <el-col
          :span="24"
          :xs="24"
        >

          <el-form label-position="top"
                   ref="currentUser"
                   :model="currentUser"
                   :rules="rules"
                   style="width: 100%">

            <el-form-item :label="$t('users.table.nickname')" prop="nickname">
              <el-input v-model.trim="currentUser.nickname"/>
            </el-form-item>

            <el-form-item :label="$t('users.table.firstName')" prop="firstName">
              <el-input v-model.trim="currentUser.firstName"/>
            </el-form-item>

            <el-form-item :label="$t('users.table.lastName')" prop="lastName">
              <el-input v-model.trim="currentUser.lastName"/>
            </el-form-item>

            <el-form-item :label="$t('users.table.email')" prop="email">
              <el-input v-model.trim="currentUser.email" type="email"/>
            </el-form-item>

            <el-form-item :label="$t('users.table.lang')" prop="lang">
              <el-select v-model="currentUser.lang" placeholder="Select" style="width: 100%">
                <el-option
                  v-for="item in options"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value">
                </el-option>
              </el-select>
            </el-form-item>

            <el-form-item :label="$t('users.table.role')" prop="role">
              <role-search
                :multiple=false
                v-model="currentUser.role"
                @update-value="changedRole"
              />
            </el-form-item>

            <el-form-item :label="$t('user.table.status')" prop="status">
              <el-select
                v-model="currentUser.status"
                placeholder="please select status"
                style="width: 100%"
                @change="changedStatus"
              >
                <el-option label="ACTIVE" value="active"></el-option>
                <el-option label="BLOCKED" value="blocked"></el-option>
              </el-select>
            </el-form-item>

            <el-form-item :label="$t('users.table.image')" prop="image">
              <image-preview :image="currentUser.image" @on-select="onSelectImage"/>
            </el-form-item>

            <el-form-item :label="$t('users.table.password')" prop="password">
              <el-input v-model="currentUser.password" type="password"/>
            </el-form-item>

            <el-form-item :label="$t('users.table.passwordRepeat')" prop="passwordRepeat">
              <el-input v-model="currentUser.passwordRepeat" type="password"/>
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
import {ApiImage, ApiNewtUserRequest, ApiRole, ApiUserMeta} from '@/api/stub';
import router from '@/router';
import {Form} from 'element-ui';
import CardWrapper from '@/components/card-wrapper/index.vue';
import ImagePreview from '@/views/images/preview.vue';
import RoleSearch from '@/views/users/components/role_search.vue';

// HACK: have to use script-loader to load jsonlint
/* eslint-disable import/no-webpack-loader-syntax */
require('script-loader!jsonlint');

@Component({
  name: 'UserEdit',
  components: {RoleSearch, ImagePreview, CardWrapper}
})
export default class extends Vue {
  @Prop({required: true}) private id!: number;

  created() {
    this.fetchUser();
  }

  private options: {
    label: string
    value: string
  }[] = [
    {
      label: 'RU',
      value: 'ru'
    },
    {
      label: 'EN',
      value: 'en'
    }
  ];

  private currentUser: {
    nickname: string;
    firstName: string;
    lastName: string;
    email: string;
    lang: string;
    password: string;
    passwordRepeat: string;
    status: string;
    image?: ApiImage;
    meta?: ApiUserMeta[];
    role?: ApiRole;
  } = {
    nickname: '',
    firstName: '',
    lastName: '',
    email: '',
    lang: '',
    password: '',
    passwordRepeat: '',
    meta: [],
    status: '',
  };

  private validatePasswordRepeat = (rule: any, value: string, callback: Function) => {
    if (!this.currentUser.password && !value) {
      callback();
      return;
    }
    if (value != this.currentUser.password) {
      callback(new Error('Please enter the correct password repeat'));
    } else {
      callback();
    }
  };

  private rules = {
    nickname: [
      {required: true, trigger: 'blur'},
      {min: 4, max: 255, trigger: 'blur'}
    ],
    firstName: [
      {required: false, trigger: 'blur'},
      {max: 255, trigger: 'blur'}
    ],
    lastName: [
      {required: false, trigger: 'blur'},
      {max: 255, trigger: 'blur'}
    ],
    role: [
      {required: false, trigger: 'blur'}
    ],
    email: [
      {required: true, trigger: 'blur'},
      {max: 255, trigger: 'blur'}
    ],
    lang: [
      {required: true, trigger: 'blur'},
      {max: 255, trigger: 'blur'}
    ],
    password: [
      {required: false, trigger: 'blur'},
      {max: 255, trigger: 'blur'}
    ],
    passwordRepeat: [
      {validator: this.validatePasswordRepeat, required: false, trigger: 'blur'}
    ]
  };

  private async save() {
    (this.$refs.currentUser as Form).validate(async valid => {
      if (!valid) {
        return;
      }
      const user: ApiNewtUserRequest = {
        nickname: this.currentUser.nickname,
        firstName: this.currentUser.firstName,
        lastName: this.currentUser.lastName,
        email: this.currentUser.email,
        lang: this.currentUser.lang,
        password: this.currentUser.password,
        passwordRepeat: this.currentUser.passwordRepeat,
        role: {name: this.currentUser?.role?.name},
        image: {id: this.currentUser?.image?.id},
        meta: this.currentUser.meta || [],
        status: this.currentUser.status
      };
      const {data} = await api.v1.userServiceUpdateUserById(this.id, user);
      if (data) {
        this.$notify({
          title: 'Success',
          message: 'user updated successfully',
          type: 'success',
          duration: 2000
        });
      }
    });
  }

  private cancel() {
    router.go(-1);
  }

  private onSelectImage(image: ApiImage, event?: any) {
    if (image) {
      this.$set(this.currentUser, 'image', image);
    } else {
      this.$set(this.currentUser, 'image', undefined);
    }
  }

  private changedRole(role: ApiRole, event?: any) {
    if (role) {
      this.$set(this.currentUser, 'role', {name: role.name});
    } else {
      this.$set(this.currentUser, 'role', undefined);
    }
  }

  private async fetchUser() {
    const {data} = await api.v1.userServiceGetUserById(this.id);
    this.currentUser.nickname = data.nickname;
    this.currentUser.firstName = data.firstName;
    this.currentUser.lastName = data.lastName;
    this.currentUser.email = data.email;
    this.currentUser.lang = data.lang;
    this.currentUser.meta = data.meta;
    this.currentUser.status = data.status;
    this.currentUser.role = data.role;
    this.currentUser.image = data.image;
  }

  private async remove() {
    const {data} = await api.v1.userServiceDeleteUserById(this.id);
    this.$notify({
      title: 'Success',
      message: 'Delete Successfully',
      type: 'success',
      duration: 2000
    });
    router.go(-1);
  }

  private changedStatus(status: string) {
    console.log('-----');
    console.log(status);
    console.log('-----');
  }
}
</script>
