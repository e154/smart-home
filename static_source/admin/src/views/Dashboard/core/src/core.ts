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
import {UUID} from 'uuid-generator-ts';
import {eventBus, RenderVar, requestCurrentState, Resolve, scriptService, stateService} from '@/views/Dashboard/core';
import {debounce} from "lodash-es";
import {ref} from "vue";
import {ItemPayload} from "@/views/Dashboard/card_items";
import {ButtonAction, Compare, CompareProp} from "./types"
import {AttributeValue, GetAttributeValue} from "@/components/Attributes"
import {FrameProp, KeysProp} from "@/views/Dashboard/components";
import {EventStateChange} from "@/api/types";
import {copyToClipboard, pasteFromClipboard} from "@/utils/clipboard";
import {generateName} from "@/utils/name";

export interface Position {
  width: string;
  height: string;
  transform: string;
}

export interface Action {
  value: string;
  label: string;
}

export interface State {
  value: string;
  label: string;
}

export interface CardItemPayload {
  style: object;
  payload: ItemPayload;
  type?: string;
  width: number;
  height: number;
  transform: string;
  showOn: CompareProp[];
  hideOn: CompareProp[];
  keysCapture?: KeysProp[];
  asButton: boolean;
  buttonActions: ButtonAction[];
  template?: boolean;
  templateFrame?: FrameProp;
}

export interface CardPayload {
  showOn?: CompareProp[];
  hideOn?: CompareProp[];
  keysCapture?: KeysProp[];
  template: boolean;
  templateFrame?: FrameProp;
  backgroundAdaptive?: boolean;
  modal?: boolean;
  modalHeader?: boolean;
}

export interface TabPayload {
  backgroundImage?: ApiImage;
  backgroundAdaptive?: boolean;
  fonts?: string[];
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
  private dashboardCardId: number;
  private styleObj: object = {};
  private styleString: string = serializedObject({});

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
      const payload = result as CardItemPayload;
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
        };
      }
      if (!this.payload.icon) {
        this.payload.icon = {
          value: '',
          iconColor: '#000000',
        };
      }
      if (this.payload.image.attrField == undefined) {
        this.payload.image.attrField = '';
      }
      if (!this.payload.image.image) {
        this.payload.image.image = undefined;
      }
      if (!this.payload.button) {
        this.payload.button = {};
      }
      if (!this.payload.state) {
        this.payload.state = {
          items: [],
          default_image: undefined,
          defaultImage: undefined,
          defaultIcon: undefined,
          defaultIconColor: undefined,
          defaultIconSize: undefined,
        };
      }
      if (!this.payload.text) {
        this.payload.text = {
          items: [],
          default_text: '<div>default text</div>',
          current_text: ''
        };
      }
      if (!this.payload.logs) {
        this.payload.logs = {
          limit: 20
        };
      }
      if (!this.payload.progress) {
        this.payload.progress = {
          items: [],
          type: '',
          showText: false,
          textInside: false,
          strokeWidth: 26,
          width: 100
        };
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
        };
      }
      if (!this.payload.chartCustom) {
        this.payload.chartCustom = {};
      }
      if (!this.payload?.map) {
        this.payload.map = {
          markers: []
        };
      } else {
        if (!this.payload.map?.markers) {
          this.payload.map.markers = [];
        }
        for (const index in this.payload.map?.markers) {
          const entityId = this.payload.map.markers[index].entityId;
          if (entityId) {
            stateService.requestCurrentState(entityId)
          }
        }
      }
      if (!this.payload.slider) {
        this.payload.slider = {};
      }
      if (!this.payload.colorPicker) {
        this.payload.colorPicker = {};
      }
      if (!this.payload.joystick) {
        this.payload.joystick = {};
      }
      if (!this.payload.video) {
        this.payload.video = {};
      }
      if (!this.payload.entityStorage) {
        this.payload.entityStorage = {};
      }
      if (!this.payload.grid) {
        this.payload.grid = {
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

    // get state
    for (const prop of [...this.hideOn, ...this.showOn]) {
      if (prop.entityId) {
        requestCurrentState(prop.entityId)
      }
    }
    if (this._entityId) {
      requestCurrentState(this._entityId)
    }
  }

  private _entityId: string;

  // entityId
  get entityId(): string {
    return this._entityId;
  }

  private _entity?: ApiEntity = {} as ApiEntity;

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

  private _type: string;

  get type(): string {
    return this._type;
  }

  // type
  set type(t: string) {
    this._type = t;
  }

  private _entityActions: Action[] = [];

  // entityActions
  get entityActions(): Action[] {
    return this._entityActions;
  }

  private _entityStates: State[] = [];

  // entityStates
  get entityStates(): State[] {
    return this._entityStates;
  }

  // style
  get style(): object {
    return this.styleObj;
  }

  // position
  get position(): Position {
    return {
      width: `${this.width}px`,
      height: `${this.height}px`,
      transform: this.transform
    };
  }

  // lastEvent
  get lastEvent(): EventStateChange | undefined {
    return stateService.lastEvent(this.entityId);
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
            default_text: '<div>default text</div>',
            current_text: ''
          }
        },
        transform: 'matrix(1, 0, 0, 1, 0, 0) translate(10px, 10px)'
      }))
    } as ApiNewDashboardCardItemRequest;
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

  serialize(): ApiDashboardCardItem {
    const style = parsedObject(this.styleString || '{}');
    this.styleObj = style;
    const buttonActions: ButtonAction[] = [];
    for (const action of this.buttonActions) {
      let entity!: {
        id?: string
      };
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
        tags: action.tags,
        areaId: action.areaId,
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

  async copyToClipboard() {
    const serialized = this.serialize();
    // @ts-ignore
    delete serialized.id;
    const request = serialized as ApiNewDashboardCardItemRequest;
    request.dashboardCardId = this.dashboardCardId;

    copyToClipboard(JSON.stringify(request, null, 2))
  }

  async copy(): Promise<CardItem> {
    const serialized = this.serialize();
    serialized.title = generateName()
    // @ts-ignore
    delete serialized.id;
    const request = serialized as ApiNewDashboardCardItemRequest;
    request.dashboardCardId = this.dashboardCardId;

    const {data} = await api.v1.dashboardCardItemServiceAddDashboardCardItem(request);

    return new CardItem(data);
  }


  update() {
    // console.log('update item', this.title)
    this.uuid = new UUID();
  }

  // ---------------------------------
  // events
  // ---------------------------------
  checkPropEntity(entityId: string): boolean {
    for (const prop of [...this.hideOn, ...this.showOn]) {
      if (prop.entityId == entityId || prop.entity?.id == entityId) {
        return true
      }
    }
    return false
  }

  async onStateChanged(event: EventStateChange) {

    if (event.entity_id != this.entityId && !this.checkPropEntity(event.entity_id)) {
      return;
    }

    this.update();

    // hide
    for (const prop of this.hideOn) {
      const val: any = await RenderVar(prop.key || '', event)
      if ('[NO VALUE]' == val) {
        continue
      }
      const tr = Compare(val, prop.value, prop.comparison);
      if (tr) {
        this.hidden = true;
        this.update();
      }
    }

    // show
    for (const prop of this.showOn) {
      const val: any = await RenderVar(prop.key || '', event)
      if ('[NO VALUE]' == val) {
        continue
      }
      const tr = Compare(val, prop.value, prop.comparison);
      if (tr) {
        this.hidden = false;
        this.update();
      }
    }
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

  eventBusHandler(event: string, args: any[]) {

    // hide
    if (this.hideOn) {
      for (const prop of this.hideOn) {
        if (prop?.eventName == event && prop?.eventArgs == args) {
          this.hidden = true;
          this.update();
        }
      }
    }

    // show
    if (this.showOn) {
      for (const prop of this.showOn) {
        if (prop?.eventName == event && prop?.eventArgs == args) {
          this.hidden = false;
          this.update();
        }
      }
    }
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
  payload: CardPayload = {} as CardPayload;
  entities: Map<string, ApiEntity>;
  active = false;
  hidden: boolean
  showOn: CompareProp[] = [];
  hideOn: CompareProp[] = [];
  keysCapture: KeysProp[] = [];
  template = false;
  templateFrame: FrameProp;
  backgroundAdaptive = true;
  modal = false;
  modalHeader = true;

  selectedItem = -1;

  items: CardItem[] = [];
  currentID: string;
  itemList = ref([])

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
      const payload = result as CardPayload;
      this.showOn = payload?.showOn || [];
      this.hideOn = payload?.hideOn || [];
      this.keysCapture = payload?.keysCapture;
      this.template = payload?.template || false;
      this.templateFrame = payload?.templateFrame || {};
      this.backgroundAdaptive = payload?.backgroundAdaptive || false;
      this.modal = payload?.modal || false;
      this.modalHeader = payload?.modalHeader == undefined? true : payload.modalHeader;
    }

    for (const index in card.items) {
      this.items.push(new CardItem(card.items[index]));
    }

    const uuid = new UUID()
    this.currentID = uuid.getDashFreeUUID()

    this.sortItems();

    this.updateItemList()

    // get state
    for (const prop of [...this.hideOn, ...this.showOn]) {
      if (prop.entityId) {
        requestCurrentState(prop.entityId)
      }
    }
  }

  _document: ElRef = null;

  set document(e: ElRef) {
    this._document = e;
  }

  updateItemList = debounce(() => {
    if (!this._document) return;
    const items = this._document.querySelectorAll(".movable");
    this.itemList.value = Array.from(items)
  }, 100)

  private _entityId: string;

  // entityId
  get entityId(): string {
    return this._entityId;
  }

  private _entity?: ApiEntity = {} as ApiEntity;

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

  static async createNew(
    title: string,
    background: string,
    backgroundAdaptive: boolean,
    template: boolean,
    templateFrame: FrameProp,
    width: number,
    height: number,
    dashboardTabId: number,
    weight: number): Promise<Card> {
    const request = {
      title: title,
      background: background,
      width: width,
      height: height,
      enabled: true,
      dashboardTabId: dashboardTabId,
      weight: weight,
      payload: btoa(JSON.stringify({
        backgroundAdaptive: backgroundAdaptive,
        template: template,
        templateFrame: templateFrame,
      }))
    } as ApiNewDashboardCardRequest;
    const {data} = await api.v1.dashboardCardServiceAddDashboardCard(request);

    return new Card(data);
  }

  static async import(card: ApiDashboardCard) {
    // todo ...
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

  serialize(): ApiDashboardCard {
    const items: ApiDashboardCardItem[] = [];

    for (const index in this.items) {
      items.push(this.items[index].serialize());
    }

    const payload = btoa(unescape(encodeURIComponent(serializedObject({
      showOn: this.showOn,
      hideOn: this.hideOn,
      keysCapture: this.keysCapture,
      template: this.template,
      templateFrame: this.templateFrame,
      backgroundAdaptive: this.backgroundAdaptive,
      modal: this.modal,
      modalHeader: this.modalHeader,
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

  // ---------------------------------
  // items

  async update() {
    return await api.v1.dashboardCardServiceUpdateDashboardCard(this.id, this.serialize());
  }

  copy(): Card {
    return new Card(this.serialize());
  }

  // ---------------------------------
  async createCardItem(type: string): Promise<CardItem> {
    const item = await CardItem.createNew(
      'item' + this.items.length,
      type,
      this.id,
      -1
    );

    item.weight = 1
    if (this.items && this.items.length > 1) {
      item.weight = this.items[this.items.length - 1].weight + 1
    }

    this.items.push(item);
    this.selectedItem = this.items.length - 1;

    // console.log('card item created, id:', item.id);

    this.updateItemList()

    return item;
  }

  addItem(item: CardItem) {
    this.items.push(item);
    this.updateItemList()
  }

  async removeItem(index: number) {
    // console.log('remove card item id:', this.items[index].id);

    const {data} = await api.v1.dashboardCardItemServiceDeleteDashboardCardItem(this.items[index].id);
    if (data) {
      this.items.splice(index, 1);
      this.selectedItem = -1;
    }

    this.updateItemList()

    eventBus.emit('selectedCardItem', -1);
  }

  async copyItem(index: number) {
    if (!this.items[index] && this.items[index].id) {
      return;
    }

    // console.log('copy card item id:', this.items[index].id);

    const item = await this.items[index].copy();
    this.items.push(item);
    this.selectedItem = this.items.length - 1;

    this.updateItemList()
  }

  async pasteCardItem() {
    const request = JSON.parse(await pasteFromClipboard())
    this.importCardItem(request)
  }

  async importCardItem(request) {

    if (!request) {
      return
    }

    request.dashboardCardId = this.id
    // console.log(request)

    try {
      const {data} = await api.v1.dashboardCardItemServiceAddDashboardCardItem(request);

      const item = new CardItem(data);

      this.items.push(item);
      this.selectedItem = this.items.length - 1;

      this.updateItemList()
    } catch (e) {

    }
  }

  sortItems() {
    this.items.sort(sortCardItems);
  }

  sortCardItemUp(item: CardItem, index: number) {
    if (!this.items[index - 1]) {
      return;
    }

    const rows = [this.items[index - 1], this.items[index]];
    this.items.splice(index - 1, 2, rows[1], rows[0]);

    let counter = 0
    for (const index in this.items) {
      this.items[index].weight = counter;
      counter++;
    }
    this.updateItemList()
  }

  sortCardItemDown(item: CardItem, index: number) {
    if (!this.items[index + 1]) {
      return;
    }

    const rows = [this.items[index], this.items[index + 1]];
    this.items.splice(index, 2, rows[1], rows[0]);

    let counter = 0
    for (const index in this.items) {
      this.items[index].weight = counter;
      counter++;
    }
    this.updateItemList()
  }

  // ---------------------------------
  // events
  // ---------------------------------
  checkPropEntity(entityId: string): boolean {
    for (const prop of [...this.hideOn, ...this.showOn]) {
      if (prop.entityId == entityId || prop.entity?.id == entityId) {
        return true
      }
    }
    return false
  }

  onStateChanged(event: EventStateChange) {
    for (const index in this.items) {
      this.items[index].onStateChanged(event);
    }

    if (event.entity_id != this.entityId && !this.checkPropEntity(event.entity_id)) {
      return;
    }

    // hide
    for (const prop of this.hideOn) {
      let val = Resolve(prop.key, event);
      if (!val) {
        continue;
      }
      if (typeof val === 'object') {
        if (val && val.hasOwnProperty('type') && val.hasOwnProperty('name')) {
          val = GetAttributeValue(val as AttributeValue);
        }
      }

      if (val == undefined) {
        val = '[NO VALUE]';
      }

      const tr = Compare(val, prop.value, prop.comparison);
      if (tr) {
        this.hidden = true;
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
          val = GetAttributeValue(val as AttributeValue);
        }
      }

      if (val == undefined) {
        val = '[NO VALUE]';
      }

      const tr = Compare(val, prop.value, prop.comparison);
      if (tr) {
        this.hidden = false;
        return;
      }
    }

  }

  eventBusHandler(event: string, args: any[]) {
    for (const index in this.items) {
      this.items[index].eventBusHandler(event, args);
    }

    // hide
    if (this.hideOn) {
      for (const prop of this.hideOn) {
        if (prop?.eventName == event && prop?.eventArgs == args) {
          this.hidden = true;
          return
        }
      }
    }

    // show
    if (this.showOn) {
      for (const prop of this.showOn) {
        if (prop?.eventName == event && prop?.eventArgs == args) {
          this.hidden = false;
          return
        }
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
  backgroundImage: ApiImage = undefined;
  backgroundAdaptive: true;
  fonts: string[];

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
    if (tab.payload) {
      const payload = parsedObject(decodeURIComponent(escape(atob(tab.payload)))) as TabPayload;
      this.backgroundImage = payload.backgroundImage || undefined;
      this.backgroundAdaptive = payload.backgroundAdaptive
      this.fonts = payload.fonts || []
    } else {
      this.backgroundImage = undefined
      this.backgroundAdaptive = true
      this.fonts = []
    }

    this.sortCards();
  }

  get cards2(): Card[] {
    return this.cards ? this.cards.filter(c => !c.hidden && !c.modal) : []
  }

  get modalCards(): Card[] {
    return this.cards ? this.cards.filter(c => c.modal) : []
  }

  static async createNew(boardId: number, name: string, columnWidth: number, weight: number): Promise<Tab> {
    const request: ApiNewDashboardTabRequest = {
      name: name,
      icon: '',
      columnWidth: columnWidth,
      gap: false,
      background: '',
      enabled: true,
      weight: weight,
      dashboardId: boardId,
      payload: btoa(JSON.stringify({
        backgroundAdaptive: true,
        fonts: [],
      }))
    };
    const {data} = await api.v1.dashboardTabServiceAddDashboardTab(request);

    return new Tab(data);
  }

  async update() {
    const request = this.serialize()
    return api.v1.dashboardTabServiceUpdateDashboardTab(this.id, request);
  }

  serialize(): ApiDashboardTab {
    const cards: ApiDashboardCard[] = [];

    for (const index in this.cards) {
      cards.push(this.cards[index].serialize());
    }

    const payload = btoa(unescape(encodeURIComponent(serializedObject({
      backgroundAdaptive: this.backgroundAdaptive,
      backgroundImage: this.backgroundImage,
      fonts: this.fonts,
    }))));
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
      entities: this.entities,
      payload: payload,
    } as ApiDashboardTab;
  }

  copy(): Tab {
    return new Tab(this.serialize());
  }

  sortCards() {
    this.cards.sort(sortCards);
  }

  sortCardUp(card: Card, index: number) {
    if (!this.cards[index - 1]) {
      return;
    }

    const rows = [this.cards[index - 1], this.cards[index]];
    this.cards.splice(index - 1, 2, rows[1], rows[0]);

    let counter = 0
    for (const index in this.cards) {
      this.cards[index].weight = counter;
      this.cards[index].update();
      counter++;
    }

  }

  sortCardDown(card: Card, index: number) {
    if (!this.cards[index + 1]) {
      return;
    }

    const rows = [this.cards[index], this.cards[index + 1]];
    this.cards.splice(index, 2, rows[1], rows[0]);

    let counter = 0
    for (const index in this.cards) {
      this.cards[index].weight = counter;
      this.cards[index].update();
      counter++;
    }
  }

  // ---------------------------------
  // events

  // ---------------------------------
  onStateChanged(event: EventStateChange) {
    for (const index in this.cards) {
      this.cards[index].onStateChanged(event);
    }
  }

  eventBusHandler(event: string, args: any[]) {
    for (const index in this.cards) {
      this.cards[index].eventBusHandler(event, args);
    }
  }
} // \Tab

export class Core {
  current: ApiDashboard = {} as ApiDashboard;
  activeCard: number | undefined = undefined; // index
  currentCardId: number | undefined;
  tabs: Tab[] = [];
  mainTab = 'cards';

  constructor() {
    //todo: move to global scope
    scriptService.start()
  }

  private _activeTabIdx = 0; // index

  get activeTabIdx(): number {
    return this._activeTabIdx;
  }

  set activeTabIdx(idx: number) {
    if (this._activeTabIdx == idx) {
      return
    }
    this._activeTabIdx = idx;
  }

  get getActiveTab(): Tab | undefined {
    if (this._activeTabIdx === undefined || this._activeTabIdx < 0) {
      this._activeTabIdx = 0
    }
    return this.tabs[this._activeTabIdx] || undefined;
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

  static async _import(dashboard: ApiDashboard): Promise<ApiDashboard> {
    const {data} = await api.v1.dashboardServiceImportDashboard(dashboard);

    return data;
  }

  currentBoard(current: ApiDashboard) {
    this.current = current;
    this.tabs = [];
    this._activeTabIdx = 0;
    if (current.tabs && current.tabs.length > 0) {
      for (const index in current.tabs) {
        this.tabs.push(new Tab(current.tabs[index]));
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

    if (!this.getActiveTab) {
      return
    }

    this.sortTabs();

    if (this.getActiveTab.cards.length > 0) {
      this.onSelectedCard(this.getActiveTab.cards[0].id)
    }

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
  // events

  async update() {
    const request = {
      name: this.current.name,
      description: this.current.description,
      enabled: this.current.enabled,
      areaId: this.current.areaId || undefined,
    };
    return await api.v1.dashboardServiceUpdateDashboard(this.current.id, request);
  }

  // ---------------------------------
  // tabs
  // ---------------------------------

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

  async removeBoard() {
    if (!this.current || !this.current.id) {
      return;
    }
    return await api.v1.dashboardServiceDeleteDashboard(this.current.id);
  }

  async importTab(tab: ApiDashboardTab) {
    if (!tab) {
      return;
    }

    tab.dashboardId = this.current.id;
    tab.id = undefined
    if (this.tabs && this.tabs.length) {
      tab.weight = this.tabs[this.tabs.length - 1].weight + 1
    }
    for (const index in this.tabs) {
      if (this.tabs[index].name == tab.name) {
        tab.name = generateName()
      }
    }

    const {data} = await api.v1.dashboardTabServiceImportDashboardTab(tab);
    if (data) {
      const tab = new Tab(data);
      this.tabs.push(tab);
      this._activeTabIdx = (this.tabs.length - 1);
      this.currentCardId = undefined;
      this.activeCard = -1
    }
  }

  // ---------------------------------
  onStateChanged(event: EventStateChange) {
    // console.log('onStateChanged', event.entity_id);
    for (const index in this.tabs) {
      this.tabs[index].onStateChanged(event);
    }
  }

  eventBusHandler(eventName: string, args: any[]) {
    if (typeof eventName != 'string') return
    for (const index in this.tabs) {
      this.tabs[index].eventBusHandler(eventName, args);
    }
  }

  selectTabInMenu(idx: number) {
    if (this._activeTabIdx === idx) return;
    this._activeTabIdx = idx;
    this.updateCurrentTab();
  }

  async createTab() {
    const tab = await Tab.createNew(this.current.id, generateName(), 300, this.tabs.length);
    this.tabs.push(tab);
    this._activeTabIdx = (this.tabs.length - 1);
    this.currentCardId = undefined;
    this.activeCard = -1
  }

  async updateTab() {
    const tab = this.getActiveTab
    if (!tab) {
      return;
    }

    // eventBus.emit('updateTab', tab.id)
    if (this.getActiveTab) {
      return this.getActiveTab.update();
    }
  }

  async removeTab() {
    const tab = this.getActiveTab
    if (!tab) {
      return;
    }

    this.tabs.splice(this._activeTabIdx, 1);
    this._activeTabIdx = this.tabs.length - 1;

    this.currentCardId = undefined;
    this.activeCard = undefined;

    this.updateCurrentTab();

    return api.v1.dashboardTabServiceDeleteDashboardTab(tab.id);
  }

  updateCurrentTab() {
    const tab = this.getActiveTab
    if (!tab) {
      return;
    }

    // console.log(`select tab id:${tab.id}`);
    eventBus.emit('updateTab', tab.id)
  }

  // ---------------------------------
  // cards
  // ---------------------------------
  onSelectedCard(id: number) {
    const tab = this.getActiveTab
    if (!tab) {
      return;
    }

    // console.log(`select card id:${id}`);
    for (const index in tab.cards) {
      const cardId = tab.cards[index].id;
      if (cardId == id) {
        this.activeCard = index as unknown as number;
        this.currentCardId = id;
        tab.cards[index].active = true
      } else {
        // console.log(`disable id:${cardId}`)
        tab.cards[index].active = false
      }
    }
  }

  async createCard() {
    const tab = this.getActiveTab
    if (!tab) {
      return;
    }

    let width: number = tab.columnWidth;
    let height: number = getSize()

    let background = ''
    let backgroundAdaptive = true
    let template = false
    let templateFrame: FrameProp
    if (tab.cards && tab.cards.length) {
      const card = tab.cards[tab.cards.length - 1];
      background = card.background
      backgroundAdaptive = card.backgroundAdaptive
      template = card.template
      templateFrame = parsedObject(serializedObject(card.templateFrame))
      width = card.width
      height = card.height
    }

    const card = await Card.createNew(
      generateName(),
      background,
      backgroundAdaptive,
      template,
      templateFrame,
      width,
      height,
      tab.id,
      10 * tab.cards.length || 0
    );

    for (const t in this.tabs) {
      for (const c in this.tabs[t].cards) {
        this.tabs[t].cards[c].active = false
      }
    }

    card.active = true
    tab.cards.push(card);
    this.activeCard = tab.cards.length - 1;
    this.currentCardId = card.id;
  }

  async updateCard() {
    const tab = this.getActiveTab
    if (!tab) {
      return;
    }

    if (this.activeCard == undefined) {
      return;
    }

    return tab.cards[this.activeCard].update();
  }

  async removeCard(cardId?: number) {
    const tab = this.getActiveTab
    if (!tab) {
      return;
    }

    let _cardId = this.currentCardId
    if (cardId != undefined) {
      _cardId = cardId
    }

    if (!_cardId) {
      return;
    }

    // console.log('remove card id:', this.currentCardId);

    const {data} = await api.v1.dashboardCardServiceDeleteDashboardCard(_cardId);
    if (data) {
      for (const index in tab.cards) {
        if (tab.cards[index].id == _cardId) {
          tab.cards.splice(parseInt(index), 1);
          break
        }
      }

      this.currentCardId = undefined;
      this.activeCard = undefined;
    }
  }

  async importCard(card: ApiDashboardCard) {
    const tab = this.getActiveTab
    if (!tab) {
      return;
    }

    card.dashboardTabId = tab.id;
    card.id = undefined
    if (tab.cards && tab.cards.length) {
      card.weight = tab.cards[tab.cards.length - 1].weight + 1
    }

    const {data} = await api.v1.dashboardCardServiceImportDashboardCard(card);
    if (data) {
      this.getActiveTab.cards.push(new Card(data));
    }

    return data;
  }

  serializeCard(cardId: number): ApiDashboardCard | undefined {
    for (const t in this.tabs) {
      for (const c in this.tabs[t].cards) {
        if (cardId == this.tabs[t].cards[c].id) {
          return this.tabs[t].cards[c].serialize()
        }
      }
    }
    return
  }

  copyCard(cardId?: number) {
    if (!cardId) return

    const source = this.serializeCard(cardId)
    source.title = generateName()
    source.weight += 10
    if (!source) return;
    this.importCard(source)
  }

  // ---------------------------------
  // Card item
  // ---------------------------------
  async createCardItem(cardId?: number, type: string) {
    const tab = this.getActiveTab
    if (!tab) {
      return;
    }

    let card: Card
    if (cardId) {
      for (const t in this.tabs) {
        for (const c in this.tabs[t].cards) {
          if (this.tabs[t].cards[c].id == cardId) {
            this.activeCard = parseInt(c);
            card = this.tabs[t].cards[c]
            break
          }
        }
      }
    } else {
      if (this.activeCard == undefined) {
        return;
      }

      card = tab.cards[this.activeCard];
    }

    if (card == undefined) {
      return;
    }

    for (const t in this.tabs) {
      for (const c in this.tabs[t].cards) {
        this.tabs[t].cards[c].active = false
      }
    }

    card.active = true
    this.currentCardId = card.id;
    await card.createCardItem(type);
  }

  importCardItem(cardId: number, request) {
    if (!cardId) return

    for (const t in this.tabs) {
      for (const c in this.tabs[t].cards) {
        if (cardId == this.tabs[t].cards[c].id) {
          return this.tabs[t].cards[c].importCardItem(request)
        }
      }
    }
  }

  serializeCardItem(cardItemId: number): ApiDashboardCardItem | undefined {
    if (!cardItemId) return

    for (const t in this.tabs) {
      for (const c in this.tabs[t].cards) {
        for (const i in this.tabs[t].cards[c].items) {
          if (cardItemId == this.tabs[t].cards[c].items[i].id) {
            return this.tabs[t].cards[c].items[i].serialize()
          }
        }
      }
    }
    return
  }

  copyCardItem(cardItemId?: number) {
    if (!cardItemId) return

    for (const t in this.tabs) {
      for (const c in this.tabs[t].cards) {
        for (const i in this.tabs[t].cards[c].items) {
          if (this.tabs[t].cards[c].items[i].id == cardItemId) {
            this.tabs[t].cards[c].copyItem(parseInt(i))
          }
        }
      }
    }
  }

  async removeCardItemById(cardItemId?: number) {
    if (!cardItemId) return
    for (const t in this.tabs) {
      for (const c in this.tabs[t].cards) {
        for (const i in this.tabs[t].cards[c].items) {
          if (this.tabs[t].cards[c].items[i].id == cardItemId) {
            await this.tabs[t].cards[c].removeItem(parseInt(i));
            return
          }
        }
      }
    }
  }

  async removeCardItem(index: number) {
    const tab = this.getActiveTab
    if (!tab) {
      return;
    }

    if (this.activeCard == undefined) {
      return;
    }

    await tab.cards[this.activeCard].removeItem(index);
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

export function serializedObject(obj: any): string {
  return JSON.stringify(obj, function (key, value) {
    if (typeof value === 'function') {
      return value.toString(); // Convert function to string
    }
    if (value instanceof Map) {
      return {
        dataType: 'Map',
        value: Array.from(value.entries()), // or with spread: value: [...value]
      };
    }
    return value;
  });
}

export function parsedObject(str): any {
  return JSON.parse(str, function (key, value) {
    if (typeof value === 'string' && value.indexOf('function') === 0) {
      return new Function('return ' + value)(); // Create a function using Function constructor
    }
    if (typeof value === 'object' && value !== null) {
      if (value.dataType === 'Map') {
        return new Map(value.value);
      }
    }
    return value;
  });
}

