import {Action, getModule, Module, Mutation, VuexModule} from 'vuex-module-decorators';
import {removeToken, setToken} from '@/utils/cookies';
import router, {resetRouter} from '@/router';
import {PermissionModule} from './permission';
import {TagsViewModule} from './tags-view';
import store from '@/store';
import api from '@/api/api';
import {ApiMeta, ApiUserHistory} from '@/api/stub';
import stream from '@/api/stream';
import customNavigator from '@/navigator';
import registerServiceWorker from '@/pwa/register-service-worker';

export interface IUserState {
  id: number;
  first_name: string;
  last_name: string;
  lang: string;
  status: string;
  email: string;
  meta: ApiMeta[];
  history: ApiUserHistory[];
  // role: IRole;
  created_at: string;
  updated_at: string;
  current_sign_in_at: string;
  last_sign_in_at: string;
}

@Module({dynamic: true, store, name: 'user'})
class User extends VuexModule implements IUserState {
  public token = localStorage.getItem('token') || '';
  public id = +(localStorage.getItem('user_id') || 0);
  public first_name = '';
  public last_name = '';
  public lang = '';
  public status = '';
  public email = '';
  public meta = [];
  public history = [];
  // public role =
  public created_at = '';
  public updated_at = '';
  public current_sign_in_at = '';
  public last_sign_in_at = '';

  // todo remove
  public name = '';
  public avatar = '';
  public introduction = '';
  public roles: string[] = [];

  @Mutation
  private SET_TOKEN(token: string) {
    localStorage.token = token;
    this.token = token;
  }

  @Mutation
  private SET_NAME(name: string) {
    this.name = name;
  }

  @Mutation
  private SET_AVATAR(avatar: string) {
    this.avatar = avatar;
  }

  @Mutation
  private SET_INTRODUCTION(introduction: string) {
    this.introduction = introduction;
  }

  @Mutation
  private SET_ROLES(roles: string[]) {
    this.roles = roles;
  }

  @Mutation
  private SET_EMAIL(email: string) {
    this.email = email;
  }

  @Mutation
  private SET_ID(id: number) {
    localStorage.user_id = id;
    this.id = id;
  }

  @Action
  public async Signin(userInfo: { username: string, password: string }) {
    let {username, password} = userInfo;
    username = username.trim();

    const {data} = await api.v1.authServiceSignin({
      headers: {Authorization: 'Basic ' + btoa(username + ':' + password)}
    });
    setToken(data.accessToken);
    this.SET_TOKEN(data.accessToken);
    this.SET_ID(data.currentUser.id);
    // ws
    stream.connect(process.env.VUE_APP_BASE_API || window.location.origin, data.accessToken);
    // geo location
    customNavigator.watchPosition();
    // push service
    registerServiceWorker.start();
  }

  @Action
  public ResetToken() {
    removeToken();
    this.SET_TOKEN('');
    this.SET_ROLES([]);
    this.SET_ID(0);
    // push service
    registerServiceWorker.stop();
    // ws
    stream.disconnect();
  }

  @Action
  public async GetUserInfo() {
    if (this.token === '') {
      throw Error('GetUserInfo: token is undefined!');
    }
    const {data} = await api.v1.userServiceGetUserById(this.id);
    if (!data) {
      throw Error('Verification failed, please Login again.');
    }
    // const { roles, name, avatar, introduction, email } = user
    // roles must be a non-empty array
    // if (!roles || roles.length <= 0) {
    //   throw Error('GetUserInfo: roles must be a non-null array!')
    // }
    this.SET_ROLES(['admin']);
    this.SET_NAME(data.nickname);
    this.SET_AVATAR('https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif');
    this.SET_INTRODUCTION('I am a super administrator');
    this.SET_EMAIL(data.email);
  }

  @Action
  public async ChangeRoles(role: string) {
    // Dynamically modify permissions
    const token = role + '-token';
    this.SET_TOKEN(token);
    setToken(token);
    await this.GetUserInfo();
    resetRouter();
    // Generate dynamic accessible routes based on roles
    PermissionModule.GenerateRoutes(this.roles);
    // Add generated routes
    PermissionModule.dynamicRoutes.forEach(route => {
      router.addRoute(route);
    });
    // Reset visited views and cached views
    TagsViewModule.delAllViews();
  }

  @Action
  public async Signout() {
    if (this.token === '') {
      throw Error('LogOut: token is undefined!');
    }
    await api.v1.authServiceSignout();
    removeToken();
    resetRouter();

    // Reset visited views and cached views
    TagsViewModule.delAllViews();
    this.SET_TOKEN('');
    this.SET_ROLES([]);
  }
}

export const UserModule = getModule(User);
