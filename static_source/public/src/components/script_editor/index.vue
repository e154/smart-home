<template>
  <div class="script-editor">
    <!--    <textarea ref="textarea" />-->
    <div ref="editor"></div>
  </div>
</template>

<script lang="ts">
import CodeMirror, {Editor} from 'codemirror'
import 'codemirror/addon/lint/lint.css'
import 'codemirror/lib/codemirror.css'
import 'codemirror/theme/mdn-like.css'
import 'codemirror/mode/coffeescript/coffeescript'
import 'codemirror/addon/lint/lint'
import 'codemirror/addon/lint/coffeescript-lint'
import 'codemirror/addon/hint/show-hint';
import 'codemirror/addon/hint/javascript-hint';
import 'codemirror/addon/hint/show-hint.css';
import 'codemirror/addon/search/search';
import 'codemirror/addon/search/searchcursor';
import 'codemirror/addon/dialog/dialog';
import 'codemirror/addon/dialog/dialog.css';
import 'codemirror/addon/fold/foldgutter';
import 'codemirror/addon/fold/foldgutter.css';
import 'codemirror/addon/fold/brace-fold';
import 'codemirror/mode/javascript/javascript';
import 'codemirror/addon/edit/closebrackets';
import 'codemirror/addon/edit/matchbrackets';
import 'codemirror/addon/comment/comment';
import 'codemirror/addon/comment/continuecomment';
import 'codemirror/addon/search/jump-to-line';
import 'codemirror/keymap/sublime';
import {Component, Prop, Vue, Watch} from 'vue-property-decorator'

// HACK: have to use script-loader to load jsonlint
/* eslint-disable import/no-webpack-loader-syntax */
require('script-loader!jsonlint')

@Component({
  name: 'ScriptEditor'
})
export default class extends Vue {
  @Prop({required: true}) private value!: string

  private jsonEditor?: Editor

  @Watch('value')
  private onValueChange(value: string) {
    if (this.jsonEditor) {
      const editorValue = this.jsonEditor.getValue()
      if (value !== editorValue) {
        this.jsonEditor.setValue(this.value)
      }
    }
  }

  mounted() {
    this.jsonEditor = CodeMirror(this.$refs.editor as HTMLDivElement, {
      lineNumbers: true,
      mode: 'application/vnd.coffeescript',
      gutters: ['CodeMirror-lint-markers'],
      theme: 'mdn-like',
      lint: false,
      indentWithTabs: true,
      smartIndent: true,
      autoCloseBrackets: true, // Автоматическое закрытие скобок
      matchBrackets: true, // Подсветка соответствующих скобок
      extraKeys: {
        'Ctrl-Space': 'autocomplete' // Комбинация клавиш для активации автодополнения
      },
      autofocus: true,
      hintOptions: {
        completeSingle: false,
      },
      lineWrapping: true
    })

    // Пример установки начального содержимого редактора
    this.jsonEditor.setValue('console.log("Hello, World!");');

    this.jsonEditor.setValue(this.value)
    this.jsonEditor.on('change', jsonEditor => {
      this.$emit('changed', jsonEditor.getValue())
      this.$emit('input', jsonEditor.getValue())
    })

    // Добавление параметров для поиска и замены
    this.jsonEditor.setOption('lineWrapping', true); // Опционально: переносить строки для лучшего отображения текста в диалоговом окне поиска и замены

    const searchWrapper = document.createElement('div');
    searchWrapper.style.position = 'absolute';
    searchWrapper.style.zIndex = '999';
    searchWrapper.style.right = '10px';
    searchWrapper.style.top = '10px';
    this.jsonEditor.getWrapperElement().appendChild(searchWrapper);

    const searchButton = document.createElement('button');
    searchButton.innerHTML = '🔍';
    searchWrapper.appendChild(searchButton);

    const replaceButton = document.createElement('button');
    replaceButton.innerHTML = '🔄';
    searchWrapper.appendChild(replaceButton);

    searchButton.addEventListener('click', () => {
      if (this.jsonEditor) {
        this.jsonEditor.execCommand('find');
      }
    });

    replaceButton.addEventListener('click', () => {
      if (this.jsonEditor) {
        this.jsonEditor.execCommand('replace');
      }
    });

    // Добавление параметров для сворачивания
    this.jsonEditor.setOption('gutters', ['CodeMirror-foldgutter']);
    this.jsonEditor.setOption('foldGutter', true);


    const dictionary = {
      words: [
        // common
        {text: 'main', displayText: 'main'},
        {text: 'unmarshal', displayText: 'unmarshal'},
        {text: 'marshal', displayText: 'marshal'},
        {text: 'hex2arr(hexString)', displayText: 'hex2arr'},
        {text: 'ExecuteSync(file, args)', displayText: 'ExecuteSync'},
        {text: 'ExecuteAsync(file, args)', displayText: 'ExecuteAsync'},
        {text: 'for v, i in items', displayText: 'for'},
        {text: 'parseFloat', displayText: 'parseFloat'},
        {text: 'indexOf', displayText: 'indexOf'},
        {text: 'substring', displayText: 'substring'},

        // notifr
        {text: 'notifr.newMessage()', displayText: 'notifr.newMessage'},
        {text: 'notifr.send(msg)', displayText: 'notifr.send'},

        // logging
        {text: 'print', displayText: 'print'},
        {text: 'Log.info', displayText: 'Log.info'},
        {text: 'Log.debug', displayText: 'Log.debug'},
        {text: 'Log.warn', displayText: 'Log.warn'},
        {text: 'Log.error', displayText: 'Log.error'},

        // storage
        {text: 'Storage.push(key, value)', displayText: 'Storage.push'},
        {text: 'Storage.getByName(key)', displayText: 'Storage.getByName'},
        {text: 'Storage.search(key)', displayText: 'Storage.search'},
        {text: 'Storage.pop(key)', displayText: 'Storage.pop'},
        {text: 'push(key, value)', displayText: 'push'},
        {text: 'getByName(key)', displayText: 'getByName'},
        {text: 'search(key)', displayText: 'search'},
        {text: 'pop(key)', displayText: 'Storage.pop'},

        // http
        {text: 'http.get(url)', displayText: 'http.get'},
        {text: 'http.post(url, body)', displayText: 'http.post'},
        {text: 'http.put(url, body)', displayText: 'http.put'},
        {text: 'http.delete(url)', displayText: 'http.delete'},

        // mqtt
        {text: 'Mqtt.publish(topic, payload, qos, retain)', displayText: 'Mqtt.publish'},
        {text: 'mqttEvent = (entityId, actionName) ->', displayText: 'mqttEvent'},
        {text: 'message', displayText: 'message'},
        {text: 'message.payload', displayText: 'message.payload'},
        {text: 'message.topic', displayText: 'message.topic'},
        {text: 'message.qos', displayText: 'message.qos'},
        {text: 'message.duplicate', displayText: 'message.duplicate'},
        {text: 'message.storage', displayText: 'message.storage'},
        {text: 'message.error', displayText: 'message.error'},
        {text: 'message.success', displayText: 'message.success'},
        {text: 'message.new_state', displayText: 'message.new_state'},

        // automation
        {text: 'automationAction = (entityId)->', displayText: 'automationAction'},
        {text: 'Action', displayText: 'Action'},
        {text: 'Action.callAction(id, action, args)', displayText: 'Action.callAction'},
        {text: 'automationCondition = (entityId)->', displayText: 'automationCondition'},
        {text: 'Condition', displayText: 'Condition'},
        {text: 'Trigger', displayText: 'Trigger'},
        {text: 'automationTriggerAlexa = (msg) ->', displayText: 'automationTriggerAlexa'},
        {text: 'automationTriggerTime = (msg) ->', displayText: 'automationTriggerTime'},
        {text: 'automationTriggerStateChanged = (msg)->', displayText: 'automationTriggerStateChanged'},
        {text: 'automationTriggerSystem = (msg)->', displayText: 'automationTriggerSystem'},

        // entity manager
        {text: 'entityManager.getEntity(id)', displayText: 'entityManager.getEntity'},
        {text: 'entityManager.setState(id, state)', displayText: 'entityManager.setState'},
        {text: 'entityManager.setAttributes(id, attr)', displayText: 'entityManager.setAttributes'},
        {text: 'entityManager.setMetric(id, name, value)', displayText: 'entityManager.setMetric'},
        {text: 'entityManager.callAction(id, action, args)', displayText: 'entityManager.callAction'},
        {text: 'entityManager.callScene(id, args)', displayText: 'entityManager.callScene'},

        // actor
        {text: 'Actor.setState(attr)', displayText: 'Actor.setState'},
        {text: 'Actor.getSettings()', displayText: 'Actor.getSettings'},

        // entity
        {text: 'entity.setState(state)', displayText: 'entity.setState'},
        {text: 'entity.setAttributes(attr)', displayText: 'entity.setAttributes'},
        {text: 'entity.getAttributes()', displayText: 'entity.getAttributes'},
        {text: 'entity.getSettings()', displayText: 'entity.getSettings'},
        {text: 'entity.setMetric(name, value)', displayText: 'entity.setMetric'},
        {text: 'entity.callAction(name, args)', displayText: 'entity.callAction'},
        {text: 'entityAction = (entityId, actionName)->', displayText: 'entityAction'},

        // telegram
        {text: 'telegramAction = (entityId, actionName)->', displayText: 'telegramAction'},

        // alexa
        {text: 'skillOnLaunch = ()->', displayText: 'skillOnLaunch'},
        {text: 'skillOnSessionEnd = ()->', displayText: 'skillOnSessionEnd'},
        {text: 'skillOnIntent = ()->', displayText: 'skillOnIntent'},
        {text: 'Alexa.slots[\'place\']', displayText: 'Alexa.slots'},
        {text: 'Alexa.sendMessage("#{place}_#{state}")', displayText: 'Alexa.sendMessage'},
        {text: 'Done("#{place}_#{state}")', displayText: 'Done'},

        // miner
        {text: 'Miner.stats()', displayText: 'Miner.stats'},
        {text: 'Miner.devs()', displayText: 'Miner.devs'},
        {text: 'Miner.summary()', displayText: 'Miner.summary'},
        {text: 'Miner.pools()', displayText: 'Miner.pools'},
        {text: 'Miner.addPool(url)', displayText: 'Miner.addPool'},
        {text: 'Miner.version()', displayText: 'Miner.version'},
        {text: 'Miner.enable(poolId)', displayText: 'Miner.enable'},
        {text: 'Miner.disable(poolId)', displayText: 'Miner.disable'},
        {text: 'Miner.delete(poolId)', displayText: 'Miner.delete'},
        {text: 'Miner.switchPool(poolId)', displayText: 'Miner.switchPool'},
        {text: 'Miner.restart()', displayText: 'Miner.restart'},

      ]
    };

    // Register our custom Codemirror hint plugin.
    CodeMirror.registerHelper('hint', 'smartHome', function (editor: CodeMirror.Editor, options: object) {
      var cur = editor.getCursor();
      var curLine = editor.getLine(cur.line);
      var start = cur.ch;
      var end = start;
      while (end < curLine.length && /[\w$]/.test(curLine.charAt(end))) ++end;
      while (start && /[\w$]/.test(curLine.charAt(start - 1))) --start;
      var curWord = start !== end && curLine.slice(start, end);
      var regex = new RegExp('^' + curWord, 'i');
      return {
        list: (!curWord ? [] : dictionary.words.filter(function (item) {
          return item.text.match(regex);
        })).sort(),
        from: CodeMirror.Pos(cur.line, start),
        to: CodeMirror.Pos(cur.line, end)
      }
    });

    this.jsonEditor.setOption('hintOptions', {
      // @ts-ignore
      hint: CodeMirror.hint.smartHome,
    })

    // Включение автодополнения
    this.jsonEditor.on('keyup', (cm, event) => {
      if (!cm.state.completionActive && event.keyCode !== 13) {
        // @ts-ignore
        CodeMirror.commands.autocomplete(cm, null, { completeSingle: false });
      }
    });

    /* eslint-disable */
    /* tslint:disable */
// @ts-ignore
//     CodeMirror.commands.autocomplete = function (cm: any) {
//
//       var doc = cm.getDoc();
//       var POS = doc.getCursor();
//       var mode = CodeMirror.innerMode(cm.getMode(), cm.getTokenAt(POS).state).mode.name;
//
//       console.log(mode)
//       if (mode == 'xml') { //html depends on xml
//         /* eslint-disable */
//         /* tslint:disable */
// // @ts-ignore
//         CodeMirror.showHint(cm, CodeMirror.hint.smartHome);
//       } else if (mode == 'javascript') {
//         /* eslint-disable */
//         /* tslint:disable */
// // @ts-ignore
//         CodeMirror.showHint(cm, CodeMirror.hint.javascript);
//       } else if (mode == 'coffeescript') {
//         /* eslint-disable */
//         /* tslint:disable */
// // @ts-ignore
//         CodeMirror.showHint(cm, CodeMirror.hint.coffeescript);
//       }
//     };

    // extend exist hint
//     let extendHints = function (hints, cm, pos) {
//       CodeMirror.showHint(cm, CodeMirror.hint.coffeescript);
//     }
//
//     this.jsonEditor.setOption('hintOptions', {
//       // @ts-ignore
// //       hint: CodeMirror.hint.javascript,
//       hint: {
//         // @ts-ignore
//         resolve: function (cm, pos) {
//           // @ts-ignore
//           var resolved = CodeMirror.hint.auto.resolve(cm, pos);
//           var result = function () {
//             if (resolved.async) {
//               var callback = arguments[1];
//               // @ts-ignore
//               arguments[1] = function (hints) {
//                 callback(extendHints(hints, cm, pos));
//               };
//               // @ts-ignore
//               resolved.apply(this, arguments);
//             } else {
//               // @ts-ignore
//               var hints = resolved.apply(this, arguments);
//               return extendHints(hints, cm, pos);
//             }
//           };
//           for (var k in resolved) {
//             result[k] = resolved[k]; // copy attributes like async & supportsSelection
//           }
//           return result;
//         }
//       },
//       completeSingle: false,
//
//       // Дополнительные параметры для настройки автодополнения
//     });
  }

  public setValue(value: string) {
    if (this.jsonEditor) {
      this.jsonEditor.setValue(value)
    }
  }

  public getValue() {
    if (this.jsonEditor) {
      return this.jsonEditor.getValue()
    }
    return ''
  }
}
</script>

<style lang="scss">
.CodeMirror {
  height: auto;
  min-height: 300px;
  font-family: inherit;
}

.CodeMirror-scroll {
  min-height: 300px;
}

.CodeMirror-foldgutter .CodeMirror-guttermarker {
  background-color: #888;
  background-repeat: no-repeat;
  color: #fff;
  font-family: arial, sans-serif;
}

.CodeMirror-foldgutter .CodeMirror-foldgutter-open:after {
  content: "▼";
}

.CodeMirror-foldgutter .CodeMirror-foldgutter-folded:after {
  content: "►";
}

.cm span.cm-string {
  color: #F08047;
}
</style>

<style lang="scss" scoped>
.script-editor {
  height: 100%;
  position: relative;
}
</style>
