import {useCache} from "@/hooks/web/useCache";

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
