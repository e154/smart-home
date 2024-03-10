<script setup lang="ts">
import {computed, onMounted, onUnmounted, ref} from 'vue';
import {Terminal} from 'xterm';
import {FitAddon} from 'xterm-addon-fit';
import 'xterm/css/xterm.css';
import {useAppStore} from "@/store/modules/app";
import {DraggableContainer} from "@/components/DraggableContainer";
import {UUID} from "uuid-generator-ts";
import {ApiLog} from "@/api/stub";
import stream from "@/api/stream";
import {parseTime} from "@/utils";
import {useCache} from "@/hooks/web/useCache";
import api from "@/api/api";

const appStore = useAppStore()

const showTerminal = computed(() => appStore.getTerminal)
const terminalRef = ref<HTMLElement | null>(null);

let term: Terminal | null = null;
let fitAddon: FitAddon | null = null;

const shellprompt = "$ ";

const currentID = ref('')
onMounted(() => {

  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  if (terminalRef.value) {
    term = new Terminal({
      screenKeys: true,
      useStyle: true,
      cursorBlink: true,
      fullscreenWin: true,
      maximizeWin: true,
      screenReaderMode: true,
      cols: 128,
    });
    term.open(terminalRef.value);
    term.options.fontSize = 12;
    term.refresh(0, term.rows - 1);

    term._initialized = true;
    term.focus(); // Фокусируем терминал при монтировании
    term.scrollToBottom();

    // Добавляем обработчик событий для ввода текста
    term.onKey((event: { key: string }) => {
      handleInput(event);
    });
    // Сохраняем экземпляр Terminal в переменной

    fitAddon = new FitAddon();
    term.loadAddon(fitAddon);

    term.onResize(function (evt) {
      // console.log('onResize', evt)
    });

    term.prompt = function () {
      term.write("\r\n" + shellprompt);
    };
    term.prompt();
    term.write('Smart Home terminal initializing ...')
    term.prompt();
  }

  setTimeout(() => {
    if (fitAddon) {
      fitAddon.fit();
    }

    getList()

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

const getList = async () => {

  let params = {
    page: 0,
    limit: 200,
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

var currLine = '';
const handleInput = (e: KeyboardEvent) => {
  const printable = !e.altKey && !e.ctrlKey && !e.metaKey;
  if (term) {
    // console.log(e, currLine)
    if (e.domEvent.keyCode === 13) {
      if (currLine == 'clear') {
        currLine = '';
        term.clear()
        term.prompt();
        return
      }
      if (currLine != '') {
        sendCommand(currLine);
      }
      currLine = ''
      term.prompt();
    } else if (e.domEvent.keyCode === 8) {
      // Do not delete the prompt
      if (term['_core'].buffer.x > 2) {
        term.write("\b \b");
      }
    } else if (printable) {
      currLine += e.key;
      term.write(e.key);
    }
  }
};

const {wsCache} = useCache()
const updateAccessToken = (payload: any) => {
  const {access_token} = payload;
  wsCache.set("accessToken", access_token)
  appStore.SetToken(access_token);
  location.reload()
}

const serverResponse = (payload: any) => {
  const {body, type} = payload
  if (body == '') {
    return
  }
  const str = body.split("\n")
  if (term) {
    str.forEach((v) => {
      term.write(v + '\r\n');
    })
  }
}

const addLog = (log: ApiLog) => {
  if (term) {
    term.write(`${parseTime(log.createdAt || log.created_at)} [${log.level}] ${log.owner} -> ${log.body}\r\n`);
  }
}

const handleResize = () => {
  if (fitAddon) {
    fitAddon.fit();
  }
}

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
  <DraggableContainer
      @resize="handleResize"
      :name="'terminal'"
      v-show="showTerminal"
      :initial-width="800"
      :initial-height="400"
  >
    <template #header>
      <span>Terminal</span>
    </template>
    <template #default>
      <div ref="terminalRef" style="width: 100%; height: 100%;"></div>
    </template>
  </DraggableContainer>
</template>

<style lang="less">

.draggable-container.container-terminal {
  background-color: #000;

  .draggable-container-content {
    padding: 0;
    padding: 0 10px;
  }
}

</style>
