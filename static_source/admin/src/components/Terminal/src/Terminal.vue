<!--
  - This file is part of the Smart Home
  - Program complex distribution https://github.com/e154/smart-home
  - Copyright (C) 2023, Filippov Alex
  -
  - This library is free software: you can redistribute it and/or
  - modify it under the terms of the GNU Lesser General Public
  - License as published by the Free Software Foundation; either
  - version 3 of the License, or (at your option) any later version.
  -
  - This library is distributed in the hope that it will be useful,
  - but WITHOUT ANY WARRANTY; without even the implied warranty of
  - MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
  - Library General Public License for more details.
  -
  - You should have received a copy of the GNU Lesser General Public
  - License along with this library.  If not, see
  - <https://www.gnu.org/licenses/>.
  -->

<script setup lang="ts">
import Terminal from 'vue-web-terminal'
import {computed, onMounted, onUnmounted, ref} from "vue";
import {useAppStore} from "@/store/modules/app";
import {ApiLog} from "@/api/stub";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";
import api from "@/api/api";
import {useCache} from "@/hooks/web/useCache";

const appStore = useAppStore()

const context = ref("")
const initLog = ref([{type: 'normal',content: "Terminal Initializing ..."}])
const showTerminal = computed(() => appStore.getTerminal)

const onExecCmd = (key, command, success, failed) => {
  if (key === 'fail') {
    failed('Something wrong!!!')
  } else {
    // let allClass = ['success', 'error', 'system', 'info', 'warning'];
    // let clazz = allClass[Math.floor(Math.random() * allClass.length)];
    success({
      type: 'normal',
      class: 'info',
      content: command
    })

    sendCommand(command)
  }
}

const addLog = (log: ApiLog) => {
  // const t = parseTime(log.createdAt)
  const message = {
    class: log.level?.toLowerCase(),
    content: `${log.owner} -> ${log.body}`
  }
  Terminal.$api.pushMessage("terminal", [message])
}

const {wsCache} = useCache()
const updateAccessToken = (payload: any) => {
  const {access_token} = payload;
  wsCache.set("accessToken", access_token)
  appStore.SetToken(access_token);
  location.reload()
}

const serverResponse = (payload: any) => {
  const {body, type} = payload
  const str = body.split("\n")
  const message = []
  str.forEach((v) => {
    message.push({
      type: 'normal',
      class: type || 'info',
      content: v,
    } )
  })
  Terminal.$api.pushMessage("terminal", message)
}

const getList = async () => {

  let params = {
    page: 0,
    limit: 100,
  }

  const res = await api.v1.logServiceGetLogList(params)
      .catch(() => {
      })
      .finally(() => {
      })
  if (res) {
    const {items, meta} = res.data;
    for (var i=items.length-1; i>=0;i--) {
      // console.log(items[i])
      addLog(items[i])
    }
  }
}

const currentID = ref('')
onMounted(() => {
  getList()

  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  context.value = "/"

  setTimeout(() => {
    stream.subscribe('log', currentID.value, addLog)
    stream.subscribe('command_response', currentID.value, serverResponse)
    stream.subscribe('update_access_token', currentID.value, updateAccessToken)
  }, 1000)
})

onUnmounted(() => {
  stream.unsubscribe('log', currentID.value)
  stream.unsubscribe('command_response', currentID.value)
  stream.unsubscribe('update_access_token', currentID.value)
})

const sendCommand = (text?: string) => {
  if (!text) {
    return
  }
  stream.send({
    id: UUID.createUUID(),
    query: 'command_terminal',
    body: btoa(text)
  });
}

</script>

<template>
  <terminal
      v-show="showTerminal"
      show-header
      name="terminal"
      @exec-cmd="onExecCmd"
      :context="context"
      :auto-help="false"
      :enable-example-hint="false"
      :init-log="initLog"
      :drag-conf="{zIndex: 9999, width: 700, height: 500, init:{ x: 50, y: 50 }}">
    <template #header>
      <div class="terminal-header">
        Terminal
      </div>
    </template>

    <template #json="data">
      {{ data.message }}
    </template>

    <template #textEditor="{data}">
      <textarea
          name="editor" class="text-editor" v-model="data.value"
                @focus="data.onFocus" @blur="data.onBlur"></textarea>
<!--      <div class="text-editor-floor" align="center">-->
<!--        <button class="text-editor-floor-btn" @click="_textEditorClose">Save & Close(Ctrl + S)</button>-->
<!--      </div>-->
    </template>
  </terminal>
</template>

<style lang="less">

.terminal-header {
  background-color: #959598;
  text-align: center;
  padding: 2px;
}
.t-log-box {
  display: block;
  margin-block-start: 0;
  margin-block-end: 0;
}
.t-ask-input, .t-window, .t-window div, .t-window p {
  font-size: 11px;
  font-family: Monaco,Menlo,Consolas,monospace;
}
</style>
