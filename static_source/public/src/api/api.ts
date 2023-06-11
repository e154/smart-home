import {Api} from '@/api/stub';
import {UserModule} from '@/store/modules/user';
import {Message, MessageBox, Notification} from 'element-ui';

const api = new Api({
  baseURL: process.env.VUE_APP_BASE_API || '/', // url = base url + request url
  timeout: 10000
  // withCredentials: true // send cookies when cross-domain requests
});

// Request interceptors
api.instance.interceptors.request.use(
  (config) => {
    // Add X-Access-Token header to every request, you can add other custom headers here
    if (UserModule.token) {
      config.headers.Authorization = UserModule.token;
    }
    return config;
  },
  (error) => {
    Promise.reject(error);
  }
);
// Response interceptors
api.instance.interceptors.response.use(
  (response) => {
    // Some example codes here:
    // code == 200: success
    // You can change this part for your own usage.
    const res = response.data;
    if (response.status !== 200) {
      Message({
        message: res.message || 'Error',
        type: 'error',
        duration: 5 * 1000
      });
      if (response.status === 401 || response.status === 400) {
        MessageBox.confirm(
          '你已被登出，可以取消继续留在该页面，或者重新登录',
          '确定登出',
          {
            confirmButtonText: '重新登录',
            cancelButtonText: '取消',
            type: 'warning'
          }
        ).then(() => {
          UserModule.ResetToken();
          location.reload(); // To prevent bugs from vue-router
        });
      }
      return Promise.reject(new Error(res.message || 'Error'));
    } else {
      if (response.data && response.data.meta) {
        response.data.meta.limit = +response.data.meta.limit;
        response.data.meta.page = +response.data.meta.page;
        response.data.meta.total = +response.data.meta.total;
      }
      return response;
    }
  },
  (error) => {
    const res = error.response.data;

    for (const i in res.details) {
      Notification({
        message: res.details[i].description || 'unknown error',
        title: 'Warning',
        type: 'warning',
        showClose: true,
      });
    }
    setTimeout(() => {
      Notification({
        message: res.message || 'Error',
        title: 'Error',
        type: 'error',
        duration: 5 * 1000,
        showClose: true,
      });
    }, 100)
    return Promise.reject(error);
  }
);

export default api;
