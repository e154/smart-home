<script setup lang="ts">
import {onMounted, onUnmounted, ref} from "vue";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";
import {Version} from "@/components/ReloadPrompt/src/types";
import {ElButton} from 'element-plus'
import {useRegisterSW} from 'virtual:pwa-register/vue'
import {registerSW} from 'virtual:pwa-register'

const {updateServiceWorker, offlineReady} = useRegisterSW({
  immediate: true,
  onRegisteredSW(swUrl, r) {
    if (!r) return;
    console.log(`Service Worker at: ${swUrl}`)
    r.update()
  },
})

const needRefresh = ref<boolean>(false)
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

const update = async () => {

  useRegisterSW({
    immediate: true,
    async onRegisteredSW(swUrl, r) {
      if (!r) return;
      console.log(`Service Worker at: ${swUrl}`)
      await r.update()
    },
  })
}

function close() {
  needRefresh.value = false
}

const onVersion = (version: Version) => {
  if (!window?.app_settings?.server_version) {
    return
  }

  console.log('received server version', version)

  needRefresh.value = version.generated_string !== currentVersion.value?.generated_string
}

</script>

<template>
  <div
    v-if="needRefresh"
    class="pwa-toast"
    role="alert"
  >
    <div class="message">
      <span>
        {{ $t('main.newContentMessage') }}
      </span>
    </div>
    <ElButton type="primary" @click="update()" size="small">
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
