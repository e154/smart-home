import {
  ApiDashboard,
  ApiDashboardCard,
  ApiDashboardCardItem,
  ApiDashboardTab,
  ApiEntity,
  ApiImage,
  ApiNewDashboardCardItemRequest,
  ApiNewDashboardCardRequest,
  ApiNewDashboardRequest,
  ApiNewDashboardTabRequest
} from '@/api/stub';
import api from '@/api/api';
import {randColor} from '@/utils/rans';
import {Attribute, EventStateChange, GetAttrValue} from '@/api/stream_types';
import {UUID} from 'uuid-generator-ts';
import {Compare, Resolve} from '@/views/Dashboard/render';
import stream from '@/api/stream';
import {useBus} from "@/views/Dashboard/bus";
import {debounce} from "lodash-es";
import {ref} from "vue";
import {ItemPayloadButton} from '@/views/Dashboard/card_items/button/types';
import {ItemPayloadText} from '@/views/Dashboard/card_items/text/types';
import {ItemPayloadState} from '@/views/Dashboard/card_items/state/types';
import {ItemPayloadLogs} from '@/views/Dashboard/card_items/logs/types';
import {ItemPayloadProgress} from '@/views/Dashboard/card_items/progress/types';
import {ItemPayloadChart} from '@/views/Dashboard/card_items/chart/types';
import {ItemPayloadChartCustom} from "@/views/Dashboard/card_items/chart_custom/types";
import {ItemPayloadMap, Marker} from '@/views/Dashboard/card_items/map/types';
import {ItemPayloadSlider} from "@/views/Dashboard/card_items/slider/types";
import {ItemPayloadColorPicker} from "@/views/Dashboard/card_items/color_picker/types";
import {ItemPayloadJoystick} from "@/views/Dashboard/card_items/joystick/types";
import {ItemPayloadVideo} from "@/views/Dashboard/card_items/video/types";
import {ItemPayloadEntityStorage} from "@/views/Dashboard/card_items/entity_storage/types";
import {prepareUrl} from "@/utils/serverId";
import {ItemPayloadTiles} from "@/views/Dashboard/card_items/tiles/types";

const {bus} = useBus()

export interface ButtonAction {
  entityId: string;
  entity?: { id?: string };
  action: string;
  image?: ApiImage | null;
  icon?: string;
  iconColor?: string;
  iconSize?: number;
}

export interface Position {
  width: string;
  height: string;
  transform: string;
}

export interface PositionInfo {
  left: string;
  top: string;
  width: string;
  height: string;
}

// eq: равно
// lt: меньше чем
// le: меньше или равно
// ne: не равно
// ge: больше или равно
// gt: больше чем
export enum comparisonType {
  EQ = 'eq',
  LT = 'lt',
  LE = 'le',
  NE = 'ne',
  GE = 'ge',
  GT = 'gt',
}

export interface CompareProp {
  key: string;
  comparison: comparisonType;
  value: string;
  entity?: { id?: string };
  entityId?: string;
}

export interface ItemPayloadImage {
  attrField?: string;
  image?: ApiImage;
}

export interface ItemPayloadIcon {
  attrField?: string;
  value?: string;
  iconColor?: string;
  iconSize?: number;
}

//todo: shouldn't be here, so will be optimize!!!
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
  tiles?: ItemPayloadTiles;
}

export interface ItemParams {
  style: object;
  payload: ItemPayload;
  type?: string;
  width: number;
  height: number;
  transform: string;
  showOn: CompareProp[];
  hideOn: CompareProp[];
  asButton: boolean;
  buttonActions: ButtonAction[];
}

export interface Action {
  value: string;
  label: string;
}

export interface State {
  value: string;
  label: string;
}

export class CardItem {
  readonly id: number;
  title: string;
  enabled: boolean;
  width = 0;
  weight = 0;
  height = 0;
  hidden = false;
  frozen = false;
  showOn: CompareProp[] = [];
  hideOn: CompareProp[] = [];
  transform = '';
  payload: ItemPayload = {} as ItemPayload;
  uuid: UUID = new UUID();
  asButton = false;
  buttonActions: ButtonAction[] = [];
  _target: any = null;

  private dashboardCardId: number;
  private styleObj: object = {};
  private styleString: string = serializedObject({});
  private _entityId: string;
  private _entity?: ApiEntity = {} as ApiEntity;
  private _type: string;
  private _entityActions: Action[] = [];
  private _entityStates: State[] = [];
  private _lastEvents?: Map<string, EventStateChange> = {} as Map<string, EventStateChange>;

  constructor(item: ApiDashboardCardItem) {
    this.id = item.id;
    this.title = item.title;
    this._type = item.type;
    this.enabled = item.enabled;
    this.dashboardCardId = item.dashboardCardId;
    this._entityId = item.entityId;
    this.hidden = item.hidden;
    this.weight = item.weight;
    this.frozen = item.frozen;
    if (this._entityId) {
      this._entity = {id: this._entityId} as ApiEntity;
    }

    if (item.payload) {
      const result: any = parsedObject(decodeURIComponent(escape(atob(item.payload))));
      const payload = result as ItemParams;
      this.width = payload.width;
      this.height = payload.height;
      this.transform = payload.transform;
      this.payload = payload.payload;
      this.styleObj = payload.style;
      this.asButton = payload.asButton || false;
      this.buttonActions = payload.buttonActions || [];
      for (const i in this.buttonActions) {
        if (this.buttonActions[i].entityId) {
          this.buttonActions[i].entity = {id: this.buttonActions[i].entityId}
        }
      }
      this.styleString = serializedObject(payload.style || {});
      if (payload.showOn) {
        this.showOn = payload.showOn;
      }
      if (payload.hideOn) {
        this.hideOn = payload.hideOn;
      }
      if (!this.payload.image) {
        this.payload.image = {
          image: undefined,
          attrField: ''
        } as ItemPayloadImage;
      }
      if (!this.payload.icon) {
        this.payload.icon = {
          value: '',
          iconColor: '#000000',
          iconSize: 12
        } as ItemPayloadIcon;
      }
      if (this.payload.image.attrField == undefined) {
        this.payload.image.attrField = '';
      }
      if (!this.payload.image.image) {
        this.payload.image.image = undefined;
      }
      if (!this.payload.button) {
        this.payload.button = {} as ItemPayloadButton;
      }
      if (!this.payload.state) {
        this.payload.state = {
          items: [],
          default_image: undefined,
          defaultImage: undefined,
          defaultIcon: undefined,
          defaultIconColor: undefined,
          defaultIconSize: undefined,
        } as ItemPayloadState;
      }
      if (!this.payload.text) {
        this.payload.text = {
          items: [],
          default_text: 'default text',
          current_text: ''
        } as ItemPayloadText;
      }
      if (!this.payload.logs) {
        this.payload.logs = {
          limit: 20
        } as ItemPayloadLogs;
      }
      if (!this.payload.progress) {
        this.payload.progress = {
          type: '',
          textInside: true,
          strokeWidth: 26,
          width: 100
        } as ItemPayloadProgress;
      }
      if (!this.payload.chart) {
        this.payload.chart = {
          type: 'line',
          metric_index: 0,
          width: 400,
          height: 400,
          xAxis: false,
          yAxis: false,
          legend: false,
          range: '24h'
        } as ItemPayloadChart;
      }
      if (!this.payload.chartCustom) {
        this.payload.chartCustom = {} as ItemPayloadChartCustom;
      }
      if (!this.payload?.map) {
        this.payload.map = {
          markers: []
        } as ItemPayloadMap;
      } else {
        if (!this.payload.map?.markers) {
          this.payload.map.markers = [] as Marker[];
        }
        for (const index in this.payload.map?.markers) {
          const entityId = this.payload.map.markers[index].entityId;
          if (entityId) {
            this._lastEvents[entityId] = {} as EventStateChange;
          }
          requestCurrentState(entityId)
        }
      }
      if (!this.payload.slider) {
        this.payload.slider = {} as ItemPayloadSlider;
      }
      if (!this.payload.colorPicker) {
        this.payload.colorPicker = {} as ItemPayloadColorPicker;
      }
      if (!this.payload.joystick) {
        this.payload.joystick = {} as ItemPayloadJoystick;
      }
      if (!this.payload.video) {
        this.payload.video = {} as ItemPayloadVideo;
      }
      if (!this.payload.entityStorage) {
        this.payload.entityStorage = {} as ItemPayloadEntityStorage;
      }
      if (!this.payload.tiles) {
        this.payload.tiles = {
          items: [],
          defaultImage: undefined,
          columns: 5,
          rows: 5,
          tileHeight: 25,
          tileWidth: 25,
          attribute: '',
        } as ItemPayloadTiles;
      }
    }
  }

  serialize(): ApiDashboardCardItem {
    const style = parsedObject(this.styleString || '{}');
    this.styleObj = style;
    const buttonActions: ButtonAction[] = [];
    for (const action of this.buttonActions) {
      let entity!: { id?: string };
      if (action.entity) {
        entity = {id: action.entity?.id};
      }
      buttonActions.push({
        entityId: action.entityId,
        entity: entity,
        action: action.action,
        image: action.image,
        icon: action.icon,
        iconColor: action.iconColor,
        iconSize: action.iconSize,
      });
    }
    const payload = {};
    payload[this._type] = this.payload[this._type];
    const cardItemPayload = btoa(unescape(encodeURIComponent(serializedObject({
      width: this.width,
      height: this.height,
      transform: this.transform,
      payload: payload,
      style: style,
      showOn: this.showOn,
      hideOn: this.hideOn,
      asButton: this.asButton,
      buttonActions: buttonActions
    }))));
    const item = {
      id: this.id,
      title: this.title,
      type: this._type,
      weight: this.weight,
      enabled: this.enabled,
      entityId: this._entityId || null,
      payload: cardItemPayload,
      hidden: this.hidden,
      frozen: this.frozen

    } as ApiDashboardCardItem;
    return item;
  }

  static async createNew(title: string, type: string,
                         dashboardCardId: number, weight: number): Promise<CardItem> {
    const request = {
      title: title,
      type: type,
      enabled: true,
      dashboardCardId: dashboardCardId,
      weight: weight,
      payload: btoa(serializedObject({
        style: {},
        width: 90,
        height: 50,
        payload: {
          text: {
            items: [],
            default_text: 'default text',
            current_text: ''
          } as ItemPayloadText
        },
        transform: 'matrix(1, 0, 0, 1, 0, 0) translate(10px, 10px)'
      }))
    } as ApiNewDashboardCardItemRequest;
    const {data} = await api.v1.dashboardCardItemServiceAddDashboardCardItem(request);

    return new CardItem(data);
  }

  async copy(): Promise<CardItem> {
    const serialized = this.serialize();
    serialized.title = serialized.title + ' [COPY]';
    // @ts-ignore
    delete serialized.id;
    const request = serialized as ApiNewDashboardCardItemRequest;
    request.dashboardCardId = this.dashboardCardId;

    const {data} = await api.v1.dashboardCardItemServiceAddDashboardCardItem(request);

    return new CardItem(data);
  }

  static async create(item: ApiDashboardCardItem): Promise<CardItem> {
    if (item.id) {
      // @ts-ignore
      delete item.id;
    }

    const request = item as ApiNewDashboardCardItemRequest;
    const {data} = await api.v1.dashboardCardItemServiceAddDashboardCardItem(request);

    return new CardItem(data);
  }

  getUrl(image: ApiImage | undefined): string {
    if (!image || !image.url) {
      return '';
    }
    return prepareUrl(import.meta.env.VITE_API_BASEPATH as string + image.url);
  }

  private clearActions() {
    this._entityActions = [];
    this._entityStates = [];
    this.showOn = [];
    this.hideOn = [];
    if (this.payload.button) {
      this.payload.button.action = undefined;
    }
  }

  // style
  get style(): object {
    return this.styleObj;
  }

  // entity
  get entity(): ApiEntity | undefined {
    // console.log('get entity', this._entity)
    return this._entity;
  }

  set entity(entity: ApiEntity | undefined) {
    // console.log('set entity', entity);
    // console.trace()
    this._entityId = entity?.id || '';
    if (entity?.id) {
      this._entity = entity;
    } else {
      this._entity = undefined;
      this.clearActions();
      return;
    }

    // update actions
    this._entityActions = [];
    if (this._entity.actions) {
      for (const item of this._entity.actions) {
        this._entityActions.push({label: item.description || item.name, value: item.name || 'no name'});
      }
    }

    // update states
    this._entityStates = [];
    if (this._entity.states) {
      for (const item of this._entity.states) {
        this._entityStates.push({label: item.description || item.name, value: item.name || 'no name'});
      }
    }
  }

  // entityId
  get entityId(): string {
    return this._entityId;
  }

  // entityActions
  get entityActions(): Action[] {
    return this._entityActions;
  }

  // entityStates
  get entityStates(): State[] {
    return this._entityStates;
  }

  // type
  set type(t: string) {
    this._type = t;
  }

  get type(): string {
    return this._type;
  }

  // target
  setTarget(e) {
    this._target = e;
  }

  get target(): any {
    return this._target;
  }

  // position
  get position(): Position {
    return {
      width: `${this.width}px`,
      height: `${this.height}px`,
      transform: this.transform
    };
  }

  // positionInfo
  get positionInfo(): PositionInfo {
    // todo optimize
    // let str = this.transform;
    //
    // const translate = str.split(') translate(');
    // const startItems = translate[0].split('matrix(')[1].split(',');
    // const startLeft = parseInt(startItems[4]);
    // const startTop = parseInt(startItems[5]);
    // const stag = translate[1].split('px,');
    // const left = startLeft + parseInt(stag[0]);
    // const top = startTop + parseInt(stag[1].split('px')[0]);

    // console.log('str', str)
    // console.log('left', left)
    // console.log('top', top)

    return {
      left: '0',
      top: '0',
      width: `${this.width}`,
      height: `${this.height}`
    };
  }

  update() {
    // console.log('update item', this.title)
    this.uuid = new UUID();
  }

  // ---------------------------------
  // events
  // ---------------------------------
  async onStateChanged(event: EventStateChange) {
    let updated: bool = false;

    // for common items
    if (this._lastEvents && this._lastEvents[event.entity_id]) {
      this._lastEvents[event.entity_id] = event;
      updated = true
    }

    // ...
    let exist: boolean
    for (const prop of this.hideOn) {
      if (prop.entityId == event.entity_id) {
        exist = true
        break;
      }
    }

    // for base entity
    if (!exist && (!this.entityId || event.entity_id != this.entityId)) {
      if (updated) {
        this.update();
      }
      return;
    }

    // console.log(event);

    this._lastEvents[this.entityId] = event;
    this.update();

    // hide
    for (const prop of this.hideOn) {
      let val = Resolve(prop.key, event);
      if (!val) {
        continue;
      }
      if (typeof val === 'object') {
        if (val && val.hasOwnProperty('type') && val.hasOwnProperty('name')) {
          val = GetAttrValue(val as Attribute);
        }
      }

      if (val == undefined) {
        val = '[NO VALUE]';
      }

      const tr = Compare(val, prop.value, prop.comparison);
      if (tr) {
        this.hidden = true;
        this.update();
        return;
      }
    }

    // show
    for (const prop of this.showOn) {
      let val = Resolve(prop.key, event);
      if (!val) {
        continue;
      }
      if (typeof val === 'object') {
        if (val && val.hasOwnProperty('type') && val.hasOwnProperty('name')) {
          val = GetAttrValue(val as Attribute);
        }
      }

      if (val == undefined) {
        val = '[NO VALUE]';
      }

      const tr = Compare(val, prop.value, prop.comparison);
      if (tr) {
        this.hidden = false;
        this.update();
        return;
      }
    }
  }

  // lastEvent
  get lastEvent(): EventStateChange | undefined {
    return this._lastEvents[this.entityId];
  }

  // lastEvents
  lastEvents(entityId: string): EventStateChange | undefined {
    return this._lastEvents[entityId];
  }
} // \CardItem

export class Card {
  id: number;
  title: string;
  height: number;
  width: number;
  background: string;
  weight: number;
  enabled: boolean;
  dashboardTabId: number;
  payload: ItemPayload = {} as ItemPayload;
  entities: Map<string, ApiEntity>;
  active = false;
  hidden: boolean
  showOn: CompareProp[] = [];
  hideOn: CompareProp[] = [];

  selectedItem = -1;

  items: CardItem[] = [];

  _document: any = null;
  currentID: string;

  private _entityId: string;
  private _entity?: ApiEntity = {} as ApiEntity;
  private _lastEvent?: EventStateChange = {} as EventStateChange;

  constructor(card: ApiDashboardCard) {
    this.id = card.id;
    this.title = card.title;
    this.height = card.height;
    this.width = card.width;
    this.background = card.background;
    this.weight = card.weight;
    this.enabled = card.enabled;
    this.dashboardTabId = card.dashboardTabId;
    this.entities = card.entities;
    this.items = [];
    this._entityId = card.entityId;
    this.hidden = card.hidden;
    if (this._entityId) {
      this._entity = {id: this._entityId} as ApiEntity;
    }
    if (card.payload) {
      const result: any = parsedObject(decodeURIComponent(escape(atob(card.payload))));
      const payload = result as ItemParams;
      if (payload.showOn) {
        this.showOn = payload.showOn;
      }
      if (payload.hideOn) {
        this.hideOn = payload.hideOn;
      }
    }

    for (const index in card.items) {
      this.items.push(new CardItem(card.items[index]));
    }

    const uuid = new UUID()
    this.currentID = uuid.getDashFreeUUID()

    this.sortItems();

    this.updateItemList()
  }

  settings() {
    const selected = true;

    return {
      maxWidth: "auto",
      maxHeight: "auto",
      minWidth: "auto",
      minHeight: "auto",

      draggable: selected,
      throttleDrag: 1, // grid,

      keepRatio: false,
      resizable: selected,
      throttleResize: 1,

      scalable: false,
      throttleScale: 0,

      rotatable: selected,
      throttleRotate: 1,
      pinchable: selected,

      origin: false,

      snappable: selected,
      snapCenter: true,
      snapHorizontal: true,
      snapVertical: true,
      snapElement: true,
      snapThreshold: 5,
      maxSnapElementGuidelineDistance: null,
      elementGuidelines: this.itemList,
      snapDirections: {"top": true, "left": true, "bottom": true, "right": true, "center": true, "middle": true},
      elementSnapDirections: {"top": true, "left": true, "bottom": true, "right": true, "center": true, "middle": true},
      isDisplaySnapDigit: true,
      isDisplayInnerSnapDigit: false,
      snapGap: true,

      renderDirections: ["nw", "n", "ne", "w", "e", "sw", "s", "se"],
      snapDigit: 5,
      snapGridWidth: 5,
      snapGridHeight: 5,
    };
  }

  static async createNew(title: string, background: string, width: number,
                         height: number, dashboardTabId: number, weight: number): Promise<Card> {
    const request = {
      title: title,
      background: background,
      width: width,
      height: height,
      enabled: true,
      dashboardTabId: dashboardTabId,
      weight: weight,
      payload: btoa(JSON.stringify({}))
    } as ApiNewDashboardCardRequest;
    const {data} = await api.v1.dashboardCardServiceAddDashboardCard(request);

    return new Card(data);
  }

  serialize(): ApiDashboardCard {
    const items: ApiDashboardCardItem[] = [];

    for (const index in this.items) {
      items.push(this.items[index].serialize());
    }

    const payload = btoa(unescape(encodeURIComponent(JSON.stringify({
      showOn: this.showOn,
      hideOn: this.hideOn,
    }))));
    const card = {
      id: this.id,
      background: this.background,
      dashboardTabId: this.dashboardTabId,
      enabled: this.enabled,
      height: this.height,
      weight: this.weight,
      payload: payload,
      title: this.title,
      width: this.width,
      items: items,
      entityId: this._entityId || null,
      hidden: this.hidden,
    } as ApiDashboardCard;
    return card;
  }

  async update() {
    return await api.v1.dashboardCardServiceUpdateDashboardCard(this.id, this.serialize());
  }

  copy(): Card {
    return new Card(this.serialize());
  }

  static async import(card: ApiDashboardCard) {
    // todo ...
  }

  // entity
  get entity(): ApiEntity | undefined {
    return this._entity;
  }

  set entity(entity: ApiEntity | undefined) {
    this._entityId = entity?.id || '';
    if (entity?.id) {
      this._entity = entity;
    } else {
      this._entity = undefined;
      return;
    }
  }

  // entityId
  get entityId(): string {
    return this._entityId;
  }

  // lastEvent
  get lastEvent(): EventStateChange | undefined {
    return this._lastEvent;
  }

  set document(d) {
    this._document = d;
  }

  // ---------------------------------
  // items
  // ---------------------------------
  async createCardItem(): Promise<CardItem> {
    const item = await CardItem.createNew(
      'item' + this.items.length,
      'text',
      this.id,
      -1
    );

    this.items.push(item);
    this.selectedItem = this.items.length - 1;

    console.log('card item created, id:', item.id);

    this.updateItemList()

    return item;
  }

  addItem(item: CardItem) {
    this.items.push(item);
    this.updateItemList()
  }

  async removeItem(index: number) {
    console.log('remove card item id:', this.items[index].id);

    const {data} = await api.v1.dashboardCardItemServiceDeleteDashboardCardItem(this.items[index].id);
    if (data) {
      this.items.splice(index, 1);
      this.selectedItem = -1;
    }

    this.updateItemList()

    bus.emit('selected_card_item', -1);
  }

  async copyItem(index: number) {
    if (!this.items[index] && this.items[index].id) {
      return;
    }

    console.log('copy card item id:', this.items[index].id);

    const item = await this.items[index].copy();
    this.items.push(item);
    this.selectedItem = this.items.length - 1;

    this.updateItemList()
  }

  itemList = ref([])
  updateItemList = debounce(() => {
    if (!this._document) return;
    const container = this._document.querySelector('.class-' + this.currentID)
    if (!container) return;
    const cubeElements = container.querySelectorAll(".movable");
    this.itemList.value = Array.from(cubeElements)
  }, 100)

  sortItems() {
    this.items.sort(sortCardItems);
  }

  // ---------------------------------
  // events
  // ---------------------------------
  onStateChanged(event: EventStateChange) {
    for (const index in this.items) {
      this.items[index].onStateChanged(event);
    }

    if (!this.entityId || event.entity_id != this.entityId) {
      return;
    }

    // console.log(event);

    this._lastEvent = event;
    // this.update();

    // hide
    for (const prop of this.hideOn) {
      let val = Resolve(prop.key, event);
      if (!val) {
        continue;
      }
      if (typeof val === 'object') {
        if (val && val.hasOwnProperty('type') && val.hasOwnProperty('name')) {
          val = GetAttrValue(val as Attribute);
        }
      }

      if (val == undefined) {
        val = '[NO VALUE]';
      }

      const tr = Compare(val, prop.value, prop.comparison);
      if (tr) {
        this.hidden = true;
        // this.update();
        return;
      }
    }

    // show
    for (const prop of this.showOn) {
      let val = Resolve(prop.key, event);
      if (!val) {
        continue;
      }
      if (typeof val === 'object') {
        if (val && val.hasOwnProperty('type') && val.hasOwnProperty('name')) {
          val = GetAttrValue(val as Attribute);
        }
      }

      if (val == undefined) {
        val = '[NO VALUE]';
      }

      const tr = Compare(val, prop.value, prop.comparison);
      if (tr) {
        this.hidden = false;
        this.update();
        return;
      }
    }

  }
} // \Card

export class Tab {
  background: string;
  cards: Card[] = [];
  columnWidth: number;
  dashboardId: number;
  enabled: boolean;
  entities: Map<string, ApiEntity>;
  gap: boolean;
  icon: string;
  id: number;
  name: string;
  weight: number;

  constructor(tab: ApiDashboardTab) {
    this.background = tab.background;
    this.cards = [];
    this.columnWidth = tab.columnWidth;
    this.dashboardId = tab.dashboardId;
    this.enabled = tab.enabled;
    this.entities = tab.entities;
    this.gap = tab.gap;
    this.icon = tab.icon;
    this.id = tab.id;
    this.name = tab.name;
    this.weight = tab.weight;

    for (const index in tab.cards) {
      this.cards.push(new Card(tab.cards[index]));
    }

    this.sortCards();
  }

  static async createNew(boardId: number, name: string, columnWidth: number, weight: number): Promise<Tab> {
    const request: ApiNewDashboardTabRequest = {
      name: name,
      icon: '',
      columnWidth: columnWidth,
      gap: false,
      background: 'inherit',
      enabled: true,
      weight: weight,
      dashboardId: boardId
    };
    const {data} = await api.v1.dashboardTabServiceAddDashboardTab(request);

    return new Tab(data);
  }

  async update() {
    const request = {
      name: this.name || '',
      icon: this.icon,
      columnWidth: this.columnWidth,
      gap: this.gap,
      background: this.background,
      enabled: this.enabled,
      weight: this.weight,
      dashboardId: this.dashboardId
    };
    return api.v1.dashboardTabServiceUpdateDashboardTab(this.id, request);
  }

  serialize(): ApiDashboardTab {
    const cards: ApiDashboardCard[] = [];

    for (const index in this.cards) {
      cards.push(this.cards[index].serialize());
    }

    return {
      id: this.id,
      name: this.name,
      columnWidth: this.columnWidth,
      gap: this.gap,
      background: this.background,
      icon: this.icon,
      enabled: this.enabled,
      weight: this.weight,
      dashboardId: this.dashboardId,
      cards: cards,
      entities: this.entities
    } as ApiDashboardTab;
  }

  copy(): Tab {
    return new Tab(this.serialize());
  }

  sortCards() {
    this.cards.sort(sortCards);
  }

  get cards2(): Card[] {
    //todo fix items sort
    const cards: Card[] = [];
    for (const card of this.cards) {
      if (card.hidden) {
        continue
      }
      cards.push(card)
    }
    return cards
  }

  // ---------------------------------
  // events
  // ---------------------------------
  onStateChanged(event: EventStateChange) {
    for (const index in this.cards) {
      this.cards[index].onStateChanged(event);
    }
  }
} // \Tab

export class Core {
  current: ApiDashboard = {} as ApiDashboard;

  activeTab = 0; // index
  currentTabId: number | undefined;

  activeCard: number | undefined = undefined; // index
  currentCardId: number | undefined;

  tabs: Tab[] = [];

  mainTab = 'cards';
  secondTab = '1';
  editorDisabled = false;

  constructor() {
  }

  currentBoard(current: ApiDashboard) {
    this.current = current;
    this.tabs = [];
    this.activeTab = 0;
    this.currentTabId = undefined;
    if (current.tabs && current.tabs.length > 0) {
      for (const index in current.tabs) {
        this.tabs.push(new Tab(current.tabs[index]));
      }
      this.currentTabId = this.tabs[this.activeTab].id;

      if (this.tabs[this.activeTab].cards.length > 0) {
        this.activeCard = 0;
      }
    }

    // get entity for card item
    for (const t in this.tabs) {
      for (const c in this.tabs[t].cards) {
        for (const i in this.tabs[t].cards[c].items) {
          (async (id: string) => {
            if (!id) {
              return;
            }
            this.tabs[t].cards[c].items[i].entity = await this.fetchEntity(id);
          })(this.tabs[t].cards[c].items[i].entityId);
        }
      }
    }

    this.sortTabs();

    this.updateCurrentTab();
  }

  sortTabs() {
    this.tabs.sort(sortTabs);
  }

  async fetchEntity(id: string): Promise<ApiEntity> {
    if (this.current.entities && this.current.entities[id]) {
      return this.current.entities[id];
    }
    const {data} = await api.v1.entityServiceGetEntity(id);
    return data;
  }

  // ---------------------------------
  // dashboard
  // ---------------------------------
  static async createNew(name: string): Promise<ApiDashboard> {
    const request: ApiNewDashboardRequest = {
      name: name
    };
    const {data} = await api.v1.dashboardServiceAddDashboard(request);
    return data;
  }

  async update() {
    const request = {
      name: this.current.name,
      description: this.current.description,
      enabled: this.current.enabled,
      areaId: this.current.areaId || undefined,
    };
    return await api.v1.dashboardServiceUpdateDashboard(this.current.id, request);
  }

  serialize(): ApiDashboard {
    const tabs: ApiDashboardTab[] = [];
    for (const index in this.tabs) {
      tabs.push(this.tabs[index].serialize());
    }

    return {
      id: this.current.id,
      name: this.current.name,
      description: this.current.description,
      enabled: this.current.enabled,
      areaId: this.current.areaId,
      tabs: tabs,
      createdAt: this.current.createdAt,
      updatedAt: this.current.updatedAt
    } as ApiDashboard;
  }

  static async _import(dashboard: ApiDashboard): Promise<ApiDashboard> {
    const {data} = await api.v1.dashboardServiceImportDashboard(dashboard);

    return data;
  }

  async removeBoard() {
    if (!this.current || !this.current.id) {
      return;
    }
    return await api.v1.dashboardServiceDeleteDashboard(this.current.id);
  }

  // ---------------------------------
  // events
  // ---------------------------------
  onStateChanged(event: EventStateChange) {
    // console.log('onStateChanged', event.entity_id);
    for (const index in this.tabs) {
      this.tabs[index].onStateChanged(event);
    }
  }

  // ---------------------------------
  // tabs
  // ---------------------------------
  async createTab() {
    const tab = await Tab.createNew(this.current.id, 'NEW_TAB' + (this.tabs.length + 1), 300, this.tabs.length);
    this.tabs.push(tab);
    this.activeTab = (this.tabs.length - 1);
    this.currentTabId = tab.id;
    this.currentCardId = undefined;
  }

  async updateTab() {
    if (this.activeTab < 0) {
      return;
    }

    bus.emit('update_tab', this.currentTabId)
    if (this.tabs[this.activeTab]) {
      return this.tabs[this.activeTab].update();
    }
  }

  async removeTab() {
    if (!this.currentTabId) {
      return;
    }

    this.currentCardId = undefined;
    this.activeCard = undefined;

    bus.emit('remove_tab', this.currentTabId)
    return api.v1.dashboardTabServiceDeleteDashboardTab(this.currentTabId);
  }

  updateCurrentTab() {
    if (!this.tabs.length) {
      this.currentTabId = undefined;
      return;
    }

    this.currentTabId = this.tabs[this.activeTab].id;
    console.log(`select tab id:${this.currentTabId}`);
    // this.activeCard = -1
    // this.currentCardId = undefined

    bus.emit('update_tab', this.currentTabId)
  }

  // ---------------------------------
  // cards
  // ---------------------------------
  onSelectedCard(id: number) {
    if (this.activeTab < 0) {
      return;
    }
    // console.log(`select card id:${id}`);
    for (const index in this.tabs[this.activeTab].cards) {
      const cardId = this.tabs[this.activeTab].cards[index].id;
      if (cardId === id) {
        this.activeCard = index as unknown as number;
        this.currentCardId = id;
        this.tabs[this.activeTab].cards[index].active = true
      } else {
        // console.log(`disable id:${cardId}`)
        this.tabs[this.activeTab].cards[index].active = false
      }
    }
  }

  async createCard() {
    if (this.activeTab < 0 || !this.currentTabId) {
      return;
    }

    const card = await Card.createNew(
      'new card' + this.tabs[this.activeTab].cards.length,
      randColor(),
      this.tabs[this.activeTab].columnWidth,
      getSize(),
      this.tabs[this.activeTab].id,
      10 * this.tabs[this.activeTab].cards.length || 0
    );

    this.tabs[this.activeTab].cards.push(card);
    this.activeCard = this.tabs[this.activeTab].cards.length - 1;
    this.currentCardId = card.id;

    bus.emit('update_tab', this.currentTabId);
  }

  async updateCard() {
    if (this.activeTab < 0 || this.activeCard == undefined) {
      return;
    }

    // move to direct call
    // bus.emit('update_tab', this.currentTabId);

    return this.tabs[this.activeTab].cards[this.activeCard].update();
  }

  async removeCard() {
    if (this.activeTab < 0 || !this.currentCardId) {
      return;
    }

    console.log('remove card id:', this.currentCardId);

    const {data} = await api.v1.dashboardCardServiceDeleteDashboardCard(this.currentCardId);
    if (data) {
      for (const index: number in this.tabs[this.activeTab].cards) {
        if (this.tabs[this.activeTab].cards[index].id == this.currentCardId) {
          this.tabs[this.activeTab].cards.splice(index, 1);
        }
      }

      this.currentCardId = undefined;
      this.activeCard = undefined;
    }

    if (this.currentTabId) {
      bus.emit('update_tab', this.currentTabId);
    }
  }

  async importCard(card: ApiDashboardCard) {
    if (this.activeTab < 0 || !this.currentTabId) {
      return;
    }
    card.dashboardTabId = this.currentTabId;
    const {data} = await api.v1.dashboardCardServiceImportDashboardCard(card);
    if (data) {
      this.tabs[this.activeTab].cards.push(new Card(data));
    }

    bus.emit('update_tab', this.currentTabId);

    return data;
  }

  // ---------------------------------
  // Card item
  // ---------------------------------
  async createCardItem() {
    if (this.activeTab < 0 || this.activeCard == undefined) {
      return;
    }

    const card = await this.tabs[this.activeTab].cards[this.activeCard];
    if (!card) {
      return;
    }

    await this.tabs[this.activeTab].cards[this.activeCard].createCardItem();

    // bus.emit('update_tab', this.currentTabId);
  }

  async removeCardItem(index: number) {
    if (this.activeTab < 0 || this.activeCard == undefined) {
      return;
    }

    await this.tabs[this.activeTab].cards[this.activeCard].removeItem(index);

    // bus.emit('update_tab', this.currentTabId);
  }
} // \Core

function sortTabs(t1: Tab, t2: Tab) {
  if (t1.weight > t2.weight) {
    return 1;
  }

  if (t1.weight < t2.weight) {
    return -1;
  }

  return 0;
}

function sortCards(n1: Card, n2: Card) {
  if (n1.weight > n2.weight) {
    return 1;
  }

  if (n1.weight < n2.weight) {
    return -1;
  }

  return 0;
}

function sortCardItems(n1: CardItem, n2: CardItem) {
  if (n1.weight > n2.weight) {
    return 1;
  }

  if (n1.weight < n2.weight) {
    return -1;
  }

  return 0;
}

function getSize(): number {
  const number = Math.random();
  if (number < 0.333) {
    return 100;
  }

  if (number < 0.666) {
    return 150;
  }

  return 200;
}

export function requestCurrentState(entityId?: string) {
  if (!entityId) {
    return;
  }
  // console.log('requestCurrentState', entityId);
  stream.send({
    id: UUID.createUUID(),
    query: 'event_get_last_state',
    body: btoa(JSON.stringify({entity_id: entityId}))
  });
}

export function serializedObject(obj: any): string {
  return JSON.stringify(obj, function(key, value) {
    if (typeof value === 'function') {
      return value.toString(); // Convert function to string
    }
    return value;
  });
}

export function parsedObject(str): any {
  return JSON.parse(str, function(key, value) {
    if (typeof value === 'string' && value.indexOf('function') === 0) {
      return new Function('return ' + value)(); // Create a function using Function constructor
    }
    return value;
  });
}

