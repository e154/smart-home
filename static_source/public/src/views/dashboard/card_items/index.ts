export {default as Dummy} from './dummy/index.vue';
export {default as IButton} from './button/index.vue';
export {default as IButtonEditor} from './button/editor.vue';
export {default as IText} from './text/index.vue';
export {default as ITextEditor} from './text/editor.vue';
export {default as IState} from './state/index.vue';
export {default as IStateEditor} from './state/editor.vue';
export {default as IImage} from './image/index.vue';
export {default as IImageEditor} from './image/editor.vue';
export {default as ILogs} from './logs/index.vue';
export {default as ILogsEditor} from './logs/editor.vue';
export {default as IProgress} from './progress/index.vue';
export {default as IProgressEditor} from './progress/editor.vue';
export {default as IChart} from './chart/index.vue';
export {default as IChartEditor} from './chart/editor.vue';
export {default as CommonEditor} from './common/editor.vue';

export const CardItemName = (name: string): string => {
  switch (name) {
    case 'button':
      return 'IButton';
    case 'text':
      return 'IText';
    case 'state':
      return 'IState';
    case 'image':
      return 'IImage';
    case 'logs':
      return 'ILogs';
    case 'progress':
      return 'IProgress';
    case 'chart':
      return 'IChart';
    default:
      console.error(`unknown card name "${name}"`);
      return 'Dummy';
  }
};

export const CardEditorName = (name: string): string => {
  switch (name) {
    case 'button':
      return 'IButtonEditor';
    case 'text':
      return 'ITextEditor';
    case 'state':
      return 'IStateEditor';
    case 'image':
      return 'IImageEditor';
    case 'logs':
      return 'ILogsEditor';
    case 'progress':
      return 'IProgressEditor';
    case 'chart':
      return 'IChartEditor';
    default:
      console.error(`unknown card name "${name}"`);
      return 'DummyEditor';
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
  {label: 'CHART', value: 'chart'}
];
