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
import {StreamPlayer, StreamPlayerEditor, ItemPayloadVideo} from './stream_player';
import {Joystick, JoystickEditor, ItemPayloadJoystick} from './joystick';
import {Icon, IconEditor, ItemPayloadIcon} from './icon';
import {Grid, GridEditor, ItemPayloadGrid} from './grid';
import {Three, ThreeEditor, ItemPayloadThree} from './three';
import {Modal, ModalEditor, ItemPayloadModal} from './modal';
import {IFrame, IFrameEditor, ItemPayloadIFrame} from './iframe';
import {ItemPayloadJsonViewer, JsonViewer, JsonViewerEditor} from './json_viewer';
import {useI18n} from "@/hooks/web/useI18n";
import {stateService} from "@/views/Dashboard/core";

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
    case 'modal':
      return Modal;
    case 'iframe':
      return IFrame;
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
    case 'video':
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
    case 'modal':
      return ModalEditor;
    case 'iframe':
      return IFrameEditor;
    default:
      // console.error(`unknown card name "${name}"`);
      return DummyEditor;
  }
};

const {t} = useI18n()

export interface ItemsType {
  label: string;
  value: string;
  children?: ItemsType[];
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
      // {label: t('dashboard.editor.MODAL'), value: 'modal'},
    ],
  },
  {
    value: 'dashboard',
    label: 'Dashboard',
    children: [
      {label: t('dashboard.editor.LOGS'), value: 'logs'},
      {label: t('dashboard.editor.ENTITY_STORAGE'), value: 'entityStorage'},
      {label: t('dashboard.editor.JSON_VIEWER'), value: 'jsonViewer'},
      {label: t('dashboard.editor.IFRAME'), value: 'iframe'},
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
      // {label: t('dashboard.editor.THREE'), value: 'three'},
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
  streamPlayer?: ItemPayloadVideo;
  entityStorage?: ItemPayloadEntityStorage;
  grid?: ItemPayloadGrid;
  three?: ItemPayloadThree;
  jsonViewer?: ItemPayloadJsonViewer;
  modal?: ItemPayloadModal;
  iframe?: ItemPayloadIFrame;
}

export const initFunc = function(payload :any) {
    if (!payload.image) {
        payload.image = {
            image: undefined,
            attrField: ''
        };
    }
    if (!payload.icon) {
        payload.icon = {
            value: '',
            iconColor: '#000000',
        };
    }
    if (payload.image.attrField == undefined) {
        payload.image.attrField = '';
    }
    if (!payload.image.image) {
        payload.image.image = undefined;
    }
    if (!payload.button) {
        payload.button = {};
    }
    if (!payload.state) {
        payload.state = {
            items: [],
            default_image: undefined,
            defaultImage: undefined,
            defaultIcon: undefined,
            defaultIconColor: undefined,
            defaultIconSize: undefined,
        };
    }
    if (!payload.text) {
        payload.text = {
            items: [],
            default_text: '<div>default text</div>',
            current_text: ''
        };
    }
    if (!payload.logs) {
        payload.logs = {
            limit: 20
        };
    }
    if (!payload.progress) {
        payload.progress = {
            items: [],
            type: '',
            showText: false,
            textInside: false,
            strokeWidth: 26,
            width: 100
        };
    }
    if (!payload.chart) {
        payload.chart = {
            type: 'line',
            metric_index: 0,
            width: 400,
            height: 400,
            xAxis: false,
            yAxis: false,
            legend: false,
            range: '24h'
        };
    }
    if (!payload.chartCustom) {
        payload.chartCustom = {};
    }
    if (!payload?.map) {
        payload.map = {
            markers: []
        };
    } else {
        if (!payload.map?.markers) {
            payload.map.markers = [];
        }
        for (const index in payload.map?.markers) {
            const entityId = payload.map.markers[index].entityId;
            if (entityId) {
                stateService.requestCurrentState(entityId)
            }
        }
    }
    if (!payload.slider) {
        payload.slider = {};
    }
    if (!payload.colorPicker) {
        payload.colorPicker = {};
    }
    if (!payload.joystick) {
        payload.joystick = {};
    }
    if (!payload.streamPlayer) {
        payload.streamPlayer = {};
    }
    if (!payload.entityStorage) {
        payload.entityStorage = {};
    }
    if (!payload.grid) {
        payload.grid = {
            items: [],
            tooltip: false,
            gap: false,
            gapSize: 5,
            defaultImage: undefined,
            cellHeight: 25,
            cellWidth: 25,
            attribute: '',
        };
    }
}

