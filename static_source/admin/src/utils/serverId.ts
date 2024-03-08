import {useCache} from "@/hooks/web/useCache";
import {ApiImage} from "@/api/stub";

const {wsCache} = useCache()

export const prepareUrl = function (url: string | undefined): string {
  if (!url) {
    return ''
  }
  if (url?.includes('serverId')) {
    return url
  }
  const serverId = wsCache.get('serverId')
  if (!serverId) {
    return url;
  }
  if (url?.includes('?')) {
    return url + '&serverId=' + serverId;
  } else {
    return url + '?serverId=' + serverId;
  }

}

export const GetFullUrl = (uri: string | undefined): string => {
  if (!uri) {
    return '';
  }
  return prepareUrl(import.meta.env.VITE_API_BASEPATH as string + uri);
}

export const GetFullImageUrl = (image?: ApiImage | undefined): string => {
    if (!image) {
        return '';
    }
    if (image?.url?.includes(import.meta.env.VITE_API_BASEPATH)) {
        return prepareUrl(image?.url || '')
    }
    return prepareUrl(import.meta.env.VITE_API_BASEPATH as string + image?.url)
}
