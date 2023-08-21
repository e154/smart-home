import {useClipboard} from "@vueuse/core";
import {ElMessage} from "element-plus";
import {unref} from "vue";
import {useI18n} from "@/hooks/web/useI18n";
const {t} = useI18n()

export const copyToClipboard = async (sourceScript: string) => {
  const { copy, copied, isSupported } = useClipboard({source: sourceScript})
  if (!isSupported) {
    ElMessage.error(t('setting.copyFailed'))
  } else {
    await copy()
    if (unref(copied)) {
      ElMessage.success(t('setting.copySuccess'))
    }
  }
}
