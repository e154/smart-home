<script setup lang="ts">
import {onMounted, onUnmounted, ref} from "vue";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";
import {Version} from "@/components/ReloadPrompt/src/types";
import {ElButton} from 'element-plus'
import {useRegisterSW} from 'virtual:pwa-register/vue'

const {offlineReady, needRefresh, updateServiceWorker} = useRegisterSW({immediate: true})

const currentVersion = ref<Version>()

const currentID = ref('')
onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  try {
    currentVersion.value = JSON.parse(window?.app_settings?.server_version)
  } catch (e) {

  }

  console.log('current server version', currentVersion.value)

  setTimeout(() => {
    stream.subscribe('event_server_version', currentID.value, onVersion)
    stream.send({
      id: UUID.createUUID(),
      query: 'event_get_server_version',
    });
  }, 2000)
})

onUnmounted(() => stream.unsubscribe('event_server_version', currentID.value))

const isNewAppFetching = ref(false)
const update = async () => {

  if (isNewAppFetching.value) {
    return
  }

  isNewAppFetching.value = true
  updateServiceWorker(true)

}

function close() {
  offlineReady.value = false
  needRefresh.value = false
}

const onVersion = (version: Version) => {
  if (!window?.app_settings?.server_version) {
    return
  }

  console.log('received server version', version)

  needRefresh.value = version.generated_string !== currentVersion.value?.generated_string

  console.log("need refresh", needRefresh.value)
}

</script>

<template>
  <div
    v-if="offlineReady || needRefresh"
    class="pwa-toast"
    role="alert"
  >
    <div class="message">
      <span v-if="isNewAppFetching">
        Loading...
      </span>
      <span v-else>
        <span v-if="offlineReady">
        {{ $t('main.offlineReady') }}
        </span>
        <span v-else>
          {{ $t('main.newContentMessage') }}
        </span>
      </span>
    </div>
    <ElButton v-if="needRefresh" type="primary" @click="update()" :disabled="isNewAppFetching" size="small">
      {{ $t('main.reload') }}
    </ElButton>
    <ElButton @click="close" size="small">
      {{ $t('main.closeDialog') }}
    </ElButton>
  </div>
</template>

<style>
.pwa-toast {
  position: fixed;
  right: 0;
  bottom: 0;
  margin: 16px;
  padding: 12px;
  border: 1px solid #8885;
  border-radius: 4px;
  z-index: 9999;
  text-align: left;
  box-shadow: 3px 4px 5px 0px #8885;
  background: var(--el-bg-color);
}

.pwa-toast .message {
  margin-bottom: 8px;
}

.pwa-toast button {
  border: 1px solid #8885;
  outline: none;
  margin-right: 5px;
  border-radius: 2px;
  padding: 3px 10px;
}
</style>
