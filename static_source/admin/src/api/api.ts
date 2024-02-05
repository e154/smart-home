import {Api} from '@/api/stub';
import {useCache} from "@/hooks/web/useCache";
import {ElMessage, ElNotification} from "element-plus";
import {useI18n} from "@/hooks/web/useI18n";
import {useAppStoreWithOut} from "@/store/modules/app";
import {resetRouter} from "@/router";
import {useTagsViewStore} from "@/store/modules/tagsView";

const {t} = useI18n()
const {wsCache} = useCache()
const appStore = useAppStoreWithOut()
const tagsViewStore = useTagsViewStore()

const api = new Api({
  baseURL: import.meta.env.VITE_API_BASEPATH as string || '/', // url = base url + request url
  timeout: 60000
  // withCredentials: true // send cookies when cross-domain requests
});

// Request interceptors
api.instance.interceptors.request.use(
  (config) => {
    // Add X-Access-Token header to every request, you can add other custom headers here
    if (wsCache.get('accessToken')) {
      config.headers.Authorization = wsCache.get('accessToken');
    }
    // Add X-SERVER-ID
    if (wsCache.get('serverId')) {
      config.headers['X-SERVER-ID'] = wsCache.get('serverId');
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
    return response;
  },
  (error) => {
    if (!error.response) {
      return
    }
    const response = error.response
    const res = response.data;

    if (response.status == 526 ) {
      ElMessage({
        message: 'No proxy available',
        type: 'error',
        duration: 5 * 1000
      });
      return
    }

    if (response.status == 403 ) {
      ElMessage({
        message: 'access forbidden: ' + res.error.message,
        type: 'error',
        duration: 5 * 1000
      });
      return
    }

    ElMessage({
      message: res.error.message || t('Error'),
      type: 'error',
      duration: 5 * 1000
    });

    if (response.status === 401 ) {

      if (location.toString().includes('/login') || location.toString().includes('/password_reset') ||
          location.toString().includes('/landing')) {
        return
      }

      appStore.RemoveToken()
      //wsCache.clear()

      tagsViewStore.delAllViews()
      resetRouter() // 重置静态路由表
      location.reload() // To prevent bugs from vue-router
      return
    }

    if (res.details) {
      for (const i in res.details) {
        ElNotification({
          message: res.details[i].description || t('message.unknownError'),
          title: t('Warning'),
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
