import {Dummy, DummyEditor} from './dummy';
import {Button, ButtonEditor, ItemPayloadButton} from './button';
import {ItemPayloadText, Text, TextEditor} from './text';
import {State, StateEditor, ItemPayloadState} from './state';
import {Image, ImageEditor, ItemPayloadImage} from './image';
import {Logs, LogsEditor, ItemPayloadLogs} from './logs';
import {Progress, ProgressEditor, ItemPayloadProgress} from './progress';
import {Chart, ChartEditor, ItemPayloadChart} from './chart';
import {ChartCustom, ChartCustomEditor, ItemPayloadChartCustom} from './chart_custom';
import {EntityStorage, EntityStorageEditor, ItemPayloadEntityStorage} from './entity_storage';
import {Map, MapEditor, ItemPayloadMap} from './map';
import {Slider, SliderEditor, ItemPayloadSlider} from './slider';
import {ColorPicker, ColorPickerEditor, ItemPayloadColorPicker} from './color_picker';
import {StreamPlayer, StreamPlayerEditor, ItemPayloadVideo} from './video';
import {Joystick, JoystickEditor, ItemPayloadJoystick} from './joystick';
import {Icon, IconEditor, ItemPayloadIcon} from './icon';
import {Grid, GridEditor, ItemPayloadGrid} from './grid';

export const CardItemName = (name: string): any => {
  switch (name) {
    case 'button':
      return Button;
    case 'text':
      return Text;
    case 'state':
      return State;
    case 'image':
      return Image;
    case 'logs':
      return Logs;
    case 'progress':
      return Progress;
    case 'chart':
      return Chart;
    case 'chart_custom':
    case 'chartCustom':
      return ChartCustom;
    case 'entityStorage':
      return EntityStorage;
    case 'map':
      return Map;
    case 'slider':
      return Slider;
    case 'colorPicker':
      return ColorPicker;
    case 'streamPlayer':
      return StreamPlayer;
    case 'joystick':
      return Joystick;
    case 'icon':
      return Icon;
    case 'grid':
      return Grid;
    default:
      // console.error(`unknown card name "${name}"`);
      return Dummy;
  }
};

export const CardEditorName = (name: string): any => {
  switch (name) {
    case 'button':
      return ButtonEditor;
    case 'text':
      return TextEditor;
    case 'state':
      return StateEditor;
    case 'image':
      return ImageEditor;
    case 'logs':
      return LogsEditor;
    case 'progress':
      return ProgressEditor;
    case 'chart':
      return ChartEditor;
    case 'chart_custom':
    case 'chartCustom':
      return ChartCustomEditor;
    case 'entityStorage':
      return EntityStorageEditor;
    case 'map':
      return MapEditor;
    case 'slider':
      return SliderEditor;
    case 'colorPicker':
      return ColorPickerEditor;
    case 'streamPlayer':
      return StreamPlayerEditor;
    case 'joystick':
      return JoystickEditor;
    case 'icon':
      return IconEditor;
    case 'grid':
      return GridEditor;
    default:
      // console.error(`unknown card name "${name}"`);
      return DummyEditor;
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

export interface ItemPayload {
  text?: ItemPayloadText;
  image?: ItemPayloadImage;
  icon?: ItemPayloadIcon;
  button?: ItemPayloadButton;
  state?: ItemPayloadState;
  logs?: ItemPayloadLogs;
  progress?: ItemPayloadProgress;
  chart?: ItemPayloadChart;
  chartCustom?: ItemPayloadChartCustom;
  map?: ItemPayloadMap;
  slider?: ItemPayloadSlider;
  colorPicker?: ItemPayloadColorPicker;
  joystick?: ItemPayloadJoystick;
  video?: ItemPayloadVideo;
  entityStorage?: ItemPayloadEntityStorage;
  grid?: ItemPayloadGrid;
}
