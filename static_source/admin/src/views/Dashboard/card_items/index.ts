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
import chart_custom from './chart_custom/index.vue';
import chartCustomEditor from './chart_custom/editor.vue';
import entityStorage from './entity_storage/index.vue';
import entityStorageEditor from './entity_storage/editor.vue';
import map from './map/index.vue';
import mapEditor from './map/editor.vue';
import slider from './slider/index.vue';
import sliderEditor from './slider/editor.vue';
import colorPicker from './color_picker/index.vue';
import colorPickerEditor from './color_picker/editor.vue';
import streamPlayer from './video/index.vue';
import streamPlayerEditor from './video/editor.vue';
import joystick from './joystick/index.vue';
import joystickEditor from './joystick/editor.vue';
import icon from './icon/index.vue';
import iconEditor from './icon/editor.vue';
import grid from './grid/index.vue';
import gridEditor from './grid/editor.vue';

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
    case 'chart_custom':
    case 'chartCustom':
      return chart_custom;
    case 'entityStorage':
      return entityStorage;
    case 'map':
      return map;
    case 'slider':
      return slider;
    case 'colorPicker':
      return colorPicker;
    case 'streamPlayer':
      return streamPlayer;
    case 'joystick':
      return joystick;
    case 'icon':
      return icon;
    case 'grid':
      return grid;
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
    case 'chart_custom':
    case 'chartCustom':
      return chartCustomEditor;
    case 'entityStorage':
      return entityStorageEditor;
    case 'map':
      return mapEditor;
    case 'slider':
      return sliderEditor;
    case 'colorPicker':
      return colorPickerEditor;
    case 'streamPlayer':
      return streamPlayerEditor;
    case 'joystick':
      return joystickEditor;
    case 'icon':
      return iconEditor;
    case 'grid':
      return gridEditor;
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
  {label: 'CHART_CUSTOM', value: 'chartCustom'},
  {label: 'ENTITY_STORAGE', value: 'entityStorage'},
  {label: 'MAP', value: 'map'},
  {label: 'SLIDER', value: 'slider'},
  {label: 'COLOR_PICKER', value: 'colorPicker'},
  {label: 'STREAM_PLAYER', value: 'streamPlayer'},
  {label: 'JOYSTICK', value: 'joystick'},
  {label: 'ICON', value: 'icon'},
  {label: 'GRID', value: 'grid'}
];
