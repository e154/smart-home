<template>
  <div class="login-container">

    <div v-if="state==='REQUEST_SENDED'" class="login-form">
      <el-card class="box-card"
               style="background-color: #2d3a4b; color: #eee; border: 1px solid #515151;"
      >
        <span>
          Check your email for a link to reset your password. If it doesn’t appear within a few minutes, check your spam
          folder.
        </span>
        <el-button
          type="primary"
          style="width:100%; margin-top:20px;"
          @click.native.prevent="gotoLogin"
        >
          {{ $t('login.returnToSignIn') }}
        </el-button>

      </el-card>
    </div> <!-- /REQUEST_SENDED -->

    <div v-else-if="state==='NEW_REQUEST'">

      <el-form
        ref="newRequestForm"
        :model="newRequestForm"
        :rules="newRequestRules"
        class="login-form"
        autocomplete="on"
        label-position="left"
      >

        <div class="title-container">
          <h3 class="title">
            {{ $t('login.restore_title') }}
          </h3>
        </div>

        <el-form-item prop="email">
          <span class="svg-container">
            <svg-icon name="user"/>
          </span>
          <el-input
            ref="email"
            v-model="newRequestForm.email"
            :placeholder="$t('login.email')"
            name="email"
            type="text"
            tabindex="1"
            autocomplete="on"
          />
        </el-form-item>

        <el-row>
          <el-col>
            <span class="cursor-pointer"
                  @click="gotoLogin()">
                {{ $t('login.backLink') }}
              </span>
          </el-col>
        </el-row>

        <el-button
          :loading="loading"
          type="primary"
          style="width:100%; margin:20px 0 30px 0;"
          @click.native.prevent="handleNewRequest"
        >
          {{ $t('login.restore') }}
        </el-button>

      </el-form>
    </div><!-- /NEW_REQUEST -->

    <div v-else-if="state==='UPDATE_PASSWORD'">

      <el-form
        ref="updatePasswordForm"
        :model="updatePasswordForm"
        :rules="updatePasswordRules"
        class="login-form"
        autocomplete="on"
        label-position="left"
      >

        <div class="title-container">
          <h3 class="title">
            {{ $t('login.enter_new_password') }}
          </h3>
        </div>

        <el-tooltip
          v-model="capsTooltip"
          content="Caps lock is On"
          placement="right"
          manual
        >
          <el-form-item prop="password">
          <span class="svg-container">
            <svg-icon name="password"/>
          </span>
            <el-input
              :key="passwordType"
              ref="password"
              v-model="updatePasswordForm.password"
              :type="passwordType"
              :placeholder="$t('login.password')"
              name="password"
              tabindex="2"
              autocomplete="on"
              @keyup.native="checkCapslock"
              @blur="capsTooltip = false"
              @keyup.enter.native="handleUpdatePassword"
            />
            <span
              class="show-pwd"
              @click="showPwd"
            >
            <svg-icon :name="passwordType === 'password' ? 'eye-off' : 'eye-on'"/>
          </span>
          </el-form-item>
        </el-tooltip>

        <el-button
          :loading="loading"
          type="primary"
          style="width:100%; margin:20px 0 30px 0;"
          @click.native.prevent="handleUpdatePassword"
        >
          {{ $t('login.restore') }}
        </el-button>

      </el-form>
    </div> <!-- /UPDATE_PASSWORD -->

  </div>
</template>

<script lang="ts">
import {Component, Vue, Watch} from 'vue-property-decorator'
import {Route} from 'vue-router'
import {Dictionary} from 'vue-router/types/router'
import {Form as ElForm, Input} from 'element-ui'
import {UserModule} from '@/store/modules/user'
import router from "@/router";

@Component({
  name: 'Restore',
  components: {}
})
export default class extends Vue {

  private newRequestForm = {
    email: '',
  }

  private updatePasswordForm = {
    password: '',
  }
  private resetToken?: string = undefined;
  // NEW_REQUEST
  // REQUEST_SENDED
  // UPDATE_PASSWORD
  private state: string = 'NEW_REQUEST'

  private newRequestRules = {
    email: [
      { required: true, trigger: 'blur' },
    ],
  }

  private validatePassword = (rule: any, value: string, callback: Function) => {
    if (value.length < 4) {
      callback(new Error('The password can not be less than 6 digits'))
    } else {
      callback()
    }
  }

  private updatePasswordRules = {
    token: [
      { required: true, trigger: 'blur' },
    ],
    password: [
      { validator: this.validatePassword, required: true, trigger: 'blur' }
    ]
  }

  private loading = false
  private redirect?: string
  private otherQuery: Dictionary<string> = {}

  @Watch('$route', {immediate: true})
  private onRouteChange(route: Route) {
    // TODO: remove the "as Dictionary<string>" hack after v4 release for vue-router
    // See https://github.com/vuejs/vue-router/pull/2050 for details
    const query = route.query as Dictionary<string>
    if (query) {
      this.redirect = query.redirect
      this.otherQuery = this.getOtherQuery(query)
    }
  }

  mounted() {
    if (this.newRequestForm.email === '') {
      (this.$refs.email as Input).focus()
    }
    if (this.$route.query.t) {
      this.resetToken = this.$route.query.t as string;
      this.state = 'UPDATE_PASSWORD'
    }
  }

  private handleUpdatePassword() {
    (this.$refs.updatePasswordForm as ElForm).validate(async (valid: boolean) => {
      if (valid) {

        this.loading = true
        await UserModule.PasswordReset({token: this.resetToken, new_password: this.updatePasswordForm.password})
        // Just to simulate the time of the request
        setTimeout(() => {
          this.loading = false
          this.gotoLogin()
        }, 0.5 * 1000)

      } else {
        return false
      }
    })
  }

  private handleNewRequest() {
    (this.$refs.newRequestForm as ElForm).validate(async (valid: boolean) => {
      if (valid) {

        this.loading = true
        UserModule.PasswordReset(this.newRequestForm).then(()=>{
          this.loading = false
          this.state = "REQUEST_SENDED"
        }).catch((e)=>{
          this.loading = false
        })
      } else {
        return false
      }
    })
  }

  private getOtherQuery(query: Dictionary<string>) {
    return Object.keys(query).reduce((acc, cur) => {
      if (cur !== 'redirect') {
        acc[cur] = query[cur]
      }
      return acc
    }, {} as Dictionary<string>)
  }

  private gotoLogin() {
    router.push({path: `/login`})
  }

  private capsTooltip = false
  private checkCapslock(e: KeyboardEvent) {
    const { key } = e
    this.capsTooltip = key !== null && key.length === 1 && (key >= 'A' && key <= 'Z')
  }

  private passwordType = 'password'
  private showPwd() {
    if (this.passwordType === 'password') {
      this.passwordType = ''
    } else {
      this.passwordType = 'password'
    }
    this.$nextTick(() => {
      (this.$refs.password as Input).focus()
    })
  }
}
</script>

<style lang="scss">
// References: https://www.zhangxinxu.com/wordpress/2018/01/css-caret-color-first-line/
@supports (-webkit-mask: none) and (not (cater-color: $loginCursorColor)) {
  .login-container .el-input {
  input { color: $loginCursorColor; }
  input::first-line { color: $lightGray; }
}
}

.login-container {
.el-input {
  display: inline-block;
  height: 47px;
  width: 85%;

input {
  height: 47px;
  background: transparent;
  border: 0px;
  border-radius: 0px;
  padding: 12px 5px 12px 15px;
  color: $lightGray;
  caret-color: $loginCursorColor;
  -webkit-appearance: none;

&:-webkit-autofill {
   box-shadow: 0 0 0px 1000px $loginBg inset !important;
   -webkit-text-fill-color: #fff !important;
 }
}
}

.el-form-item {
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(0, 0, 0, 0.1);
  border-radius: 5px;
  color: #454545;
}
}
</style>

<style lang="scss" scoped>
.login-container {
  height: 100%;
  width: 100%;
  overflow: hidden;
  background-color: $loginBg;

.login-form {
  position: relative;
  width: 520px;
  max-width: 100%;
  padding: 160px 35px 0;
  margin: 0 auto;
  overflow: hidden;
}

.tips {
  font-size: 14px;
  color: #fff;
  margin-bottom: 10px;

span {
&:first-of-type {
   margin-right: 16px;
 }
}
}

.svg-container {
  padding: 6px 5px 6px 15px;
  color: $darkGray;
  vertical-align: middle;
  width: 30px;
  display: inline-block;
}

.title-container {
  position: relative;

.title {
  font-size: 26px;
  color: $lightGray;
  margin: 0px auto 40px auto;
  text-align: center;
  font-weight: bold;
}

.set-language {
  color: #fff;
  position: absolute;
  top: 3px;
  font-size: 18px;
  right: 0px;
  cursor: pointer;
}
}

.show-pwd {
  position: absolute;
  right: 10px;
  top: 7px;
  font-size: 16px;
  color: $darkGray;
  cursor: pointer;
  user-select: none;
}

.thirdparty-button {
  position: absolute;
  right: 0;
  bottom: 6px;
}

@media only screen and (max-width: 470px) {
  .thirdparty-button {
    display: none;
  }
}

.cursor-pointer {
  cursor: pointer;
  color: $lightGray;
}
}
</style>
