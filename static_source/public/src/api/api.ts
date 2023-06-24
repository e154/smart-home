import {Api} from '@/api/stub';
import {UserModule} from '@/store/modules/user';
import {Message, Notification} from 'element-ui';

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
    const res = response.data;
    if (res && res.meta) {
      res.meta.limit = +res.meta.limit;
      res.meta.page = +res.meta.page;
      res.meta.total = +res.meta.total;
    }
    return response;
  },
  (error) => {
    if (!error.response) {
      return
    }
    const response = error.response
    const res = response.data;

    Message({
      message: res.message || 'Error',
      type: 'error',
      duration: 5 * 1000
    });

    if (response.status === 401 || response.status === 400) {

      if (location.toString().includes('/login') || location.toString().includes('/password_reset')) {
        return
      }

      UserModule.ResetToken();
      location.reload() // To prevent bugs from vue-router
      return
    }

    if (res.details) {
      for (const i in res.details) {
        Notification({
          message: res.details[i].description || 'unknown error',
          title: 'Warning',
          type: 'warning',
          showClose: true,
          duration: 5 * 1000,
        });
      }
    }

    return Promise.reject(error);
  }
);

export default api;
