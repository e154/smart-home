import dummy from './dummy/index.vue';
import dummyEditor from './dummy/editor.vue';
import button from './button/index.vue';
import buttonEditor from './button/editor.vue';
import text from './text/index.vue'
import textEditor from './text/editor.vue'
import state from './state/index.vue';
import stateEditor from './state/editor.vue';
import image from './image/index.vue';
import imageEditor from './image/editor.vue';
import logs from './logs/index.vue';
import logsEditor from './logs/editor.vue';
import progress from './progress/index.vue';
import progressEditor from './progress/editor.vue';
import chart from './chart/index.vue';
import chartEditor from './chart/editor.vue';
import entityStorage from './entity_storage/index.vue';
import entityStorageEditor from './entity_storage/editor.vue';

export const CardItemName = (name: string): any => {
  switch (name) {
    case 'button':
      return button;
    case 'text':
      return text;
    case 'state':
      return state;
    case 'image':
      return image;
    case 'logs':
      return logs;
    case 'progress':
      return progress;
    case 'chart':
      return chart;
    case 'entityStorage':
      return entityStorage;
    default:
      // console.error(`unknown card name "${name}"`);
      return dummy;
  }
};

export const CardEditorName = (name: string): any => {
  switch (name) {
    case 'button':
      return buttonEditor;
    case 'text':
      return textEditor;
    case 'state':
      return stateEditor;
    case 'image':
      return imageEditor;
    case 'logs':
      return logsEditor;
    case 'progress':
      return progressEditor;
    case 'chart':
      return chartEditor;
    case 'entityStorage':
      return entityStorageEditor;
    default:
      // console.error(`unknown card name "${name}"`);
      return dummyEditor;
  }
};

export interface ItemsType {
  label: string;
  value: string;
}

export const CardItemList: ItemsType[] = [
  {label: 'TEXT', value: 'text'},
  {label: 'IMAGE', value: 'image'},
  {label: 'BUTTON', value: 'button'},
  {label: 'STATE', value: 'state'},
  {label: 'LOGS', value: 'logs'},
  {label: 'PROGRESS', value: 'progress'},
  {label: 'CHART', value: 'chart'},
  {label: 'ENTITY_STORAGE', value: 'entityStorage'}
];
