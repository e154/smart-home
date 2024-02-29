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
import {Three, ThreeEditor, ItemPayloadThree} from './three';
import {ItemPayloadJsonViewer, JsonViewer, JsonViewerEditor} from './json_viewer';
import {useI18n} from "@/hooks/web/useI18n";

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
    case 'three':
      return Three;
    case 'jsonViewer':
      return JsonViewer;
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
    case 'three':
      return ThreeEditor;
    case 'jsonViewer':
      return JsonViewerEditor;
    default:
      // console.error(`unknown card name "${name}"`);
      return DummyEditor;
  }
};

const {t} = useI18n()

export interface ItemsType {
  label: string;
  value: string;
}

export const CardItemList: ItemsType[] = [
  {
    value: 'general',
    label: 'General',
    children: [
      {label: t('dashboard.editor.TEXT'), value: 'text'},
      {label: t('dashboard.editor.IMAGE'), value: 'image'},
      {label: t('dashboard.editor.PROGRESS'), value: 'progress'},
      {label: t('dashboard.editor.ICON'), value: 'icon'},
      {label: t('dashboard.editor.STATE'), value: 'state'},
      {label: t('dashboard.editor.GRID'), value: 'grid'},
    ],
  },
  {
    value: 'dashboard',
    label: 'Dashboard',
    children: [
      {label: t('dashboard.editor.LOGS'), value: 'logs'},
      {label: t('dashboard.editor.ENTITY_STORAGE'), value: 'entityStorage'},
      {label: t('dashboard.editor.JSON_VIEWER'), value: 'jsonViewer'},
    ],
  },
  {
    value: 'video',
    label: 'Video',
    children: [
      {label: t('dashboard.editor.STREAM_PLAYER'), value: 'streamPlayer'},
    ],
  },
  {
    value: 'control',
    label: 'Control',
    children: [
      {label: t('dashboard.editor.BUTTON'), value: 'button'},
      {label: t('dashboard.editor.JOYSTICK'), value: 'joystick'},
      {label: t('dashboard.editor.SLIDER'), value: 'slider'},
      {label: t('dashboard.editor.COLOR_PICKER'), value: 'colorPicker'},
    ],
  },
  {
    value: 'charts',
    label: 'Charts',
    children: [
      {label: t('dashboard.editor.CHART'), value: 'chart'},
    ],
  },
  {
    value: 'geo',
    label: 'Geo',
    children: [
      {label: t('dashboard.editor.MAP'), value: 'map'},
    ],
  },
  {
    value: 'experimental',
    label: 'Experimental',
    children: [
      {label: t('dashboard.editor.THREE'), value: 'three'},
      {label: t('dashboard.editor.CHART_CUSTOM'), value: 'chartCustom'},
    ],
  }
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
  three?: ItemPayloadThree;
  jsonViewer?: ItemPayloadJsonViewer;
}
