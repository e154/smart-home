/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

/**
 * Generic Error Response
 * Generic Error Response
 */
export interface GenericErrorResponse {
  message?: string;
  code?: string;
}

export interface AccessListListOfString {
  items: string[];
}

export interface GetImageFilterListResultfilter {
  date: string;
  /** @format int32 */
  count: number;
}

export interface UpdateDashboardCardRequestItem {
  /** @format int64 */
  id: number;
  title: string;
  type: string;
  /** @format int32 */
  weight: number;
  enabled: boolean;
  entityId?: string;
  /**
   * @format byte
   * @pattern ^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$
   */
  payload: string;
  hidden: boolean;
  frozen: boolean;
  showOn?: string[];
  hideOn?: string[];
}

export interface UpdateRoleAccessListRequestAccessListDiff {
  items: Record<string, boolean>;
}

export interface ApiAccessItem {
  actions: string[];
  method: string;
  description: string;
  roleName: string;
}

export interface ApiAccessLevels {
  items: Record<string, ApiAccessItem>;
}

export interface ApiAccessList {
  levels: Record<string, ApiAccessLevels>;
}

export interface ApiAccessListResponse {
  accessList?: ApiAccessList;
}

export interface ApiAction {
  /** @format int64 */
  id: number;
  name: string;
  description: string;
  /** @format int64 */
  scriptId?: number;
  script?: ApiScript;
  /** @format int64 */
  areaId?: number;
  area?: ApiArea;
  entity?: ApiEntity;
  entityId?: string;
  entityActionName?: string;
  completed?: boolean;
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
}

export interface ApiArea {
  /** @format int64 */
  id: number;
  name: string;
  description: string;
  polygon: ApiAreaLocation[];
  center?: ApiAreaLocation;
  /** @format float */
  zoom: number;
  /** @format float */
  resolution: number;
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
}

export interface ApiAreaLocation {
  /** @format double */
  lat: number;
  /** @format double */
  lon: number;
}

export interface ApiAttribute {
  name: string;
  type: ApiTypes;
  /** @format int64 */
  int?: number;
  string?: string;
  bool?: boolean;
  /** @format float */
  float?: number;
  array?: ApiAttribute[];
  map?: Record<string, ApiAttribute>;
  /** @format date-time */
  time?: string;
  imageUrl?: string;
  icon?: string;
  point?: string;
  encrypted?: string;
}

export interface ApiAutomationRequest {
  /** @format int64 */
  id: number;
  name: string;
}

export interface ApiBusStateItem {
  topic: string;
  /** @format int32 */
  subscribers: number;
  /** @format int64 */
  min: number;
  /** @format int64 */
  max: number;
  /** @format int64 */
  avg: number;
  /** @format double */
  rps: number;
}

export interface ApiClient {
  clientId: string;
  username: string;
  /** @format uint16 */
  keepAlive: number;
  /** @format int32 */
  version: number;
  willRetain: boolean;
  /** @format uint8 */
  willQos: number;
  willTopic: string;
  willPayload: string;
  remoteAddr: string;
  localAddr: string;
  /** @format uint32 */
  subscriptionsCurrent: number;
  /** @format uint32 */
  subscriptionsTotal: number;
  /** @format uint64 */
  packetsReceivedBytes: number;
  /** @format uint64 */
  packetsReceivedNums: number;
  /** @format uint64 */
  packetsSendBytes: number;
  /** @format uint64 */
  packetsSendNums: number;
  /** @format uint64 */
  messageDropped: number;
  /** @format uint32 */
  inflightLen: number;
  /** @format uint32 */
  queueLen: number;
  /** @format date-time */
  connectedAt: string;
  /** @format date-time */
  disconnectedAt?: string;
}

export interface ApiCondition {
  /** @format int64 */
  id: number;
  name: string;
  description: string;
  /** @format int64 */
  scriptId?: number;
  script?: ApiScript;
  area?: ApiArea;
  /** @format int64 */
  areaId?: number;
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
}

export interface ApiCurrentUser {
  /** @format int64 */
  id?: number;
  nickname?: string;
  firstName?: string;
  lastName?: string;
  email?: string;
  status?: string;
  history?: ApiUserHistory[];
  image?: ApiImage;
  /** @format int64 */
  signInCount?: number;
  meta?: ApiUserMeta[];
  role?: ApiRole;
  lang?: string;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
  /** @format date-time */
  currentSignInAt?: string;
  /** @format date-time */
  lastSignInAt?: string;
}

export interface ApiDashboard {
  /** @format int64 */
  id: number;
  name: string;
  description: string;
  enabled: boolean;
  /** @format int64 */
  areaId?: number;
  area?: ApiArea;
  tabs: ApiDashboardTab[];
  entities: Record<string, ApiEntity>;
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
}

export interface ApiDashboardCard {
  /** @format int64 */
  id: number;
  title: string;
  /** @format int32 */
  height: number;
  /** @format int32 */
  width: number;
  background?: string;
  /** @format int32 */
  weight: number;
  enabled: boolean;
  /** @format int64 */
  dashboardTabId: number;
  /**
   * @format byte
   * @pattern ^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$
   */
  payload: string;
  items: ApiDashboardCardItem[];
  entities: Record<string, ApiEntity>;
  hidden: boolean;
  entityId?: string;
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
}

export interface ApiDashboardCardItem {
  /** @format int64 */
  id: number;
  title: string;
  type: string;
  /** @format int32 */
  weight: number;
  enabled: boolean;
  /** @format int64 */
  dashboardCardId: number;
  entityId?: string;
  /**
   * @format byte
   * @pattern ^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$
   */
  payload: string;
  hidden: boolean;
  frozen: boolean;
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
}

export interface ApiDashboardShort {
  /** @format int64 */
  id: number;
  name: string;
  description: string;
  enabled: boolean;
  /** @format int64 */
  areaId?: number;
  area?: ApiArea;
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
}

export interface ApiDashboardTab {
  /** @format int64 */
  id: number;
  name: string;
  /** @format int32 */
  columnWidth: number;
  gap: boolean;
  background?: string;
  icon: string;
  enabled: boolean;
  /** @format int32 */
  weight: number;
  /** @format int64 */
  dashboardId: number;
  cards: ApiDashboardCard[];
  entities: Record<string, ApiEntity>;
  /**
   * @format byte
   * @pattern ^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$
   */
  payload: string;
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
}

export interface ApiDashboardTabShort {
  /** @format int64 */
  id?: number;
  name?: string;
  /** @format int32 */
  columnWidth?: number;
  gap?: boolean;
  background?: string;
  icon?: string;
  enabled?: boolean;
  /** @format int32 */
  weight?: number;
  /** @format int64 */
  dashboardId?: number;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface ApiDeviceBanRequest {
  /** @format int64 */
  id: number;
  friendlyName: string;
}

export interface ApiDeviceListResult {
  items: ApiZigbee2MqttDevice[];
  meta?: ApiMeta;
}

export interface ApiDeviceRenameRequest {
  friendlyName: string;
  newName: string;
}

export interface ApiDeviceWhitelistRequest {
  /** @format int64 */
  id: number;
  friendlyName: string;
}

export type ApiDisablePluginResult = object;

export type ApiEnablePluginResult = object;

export interface ApiEntity {
  id: string;
  pluginName: string;
  description: string;
  area?: ApiArea;
  image?: ApiImage;
  icon?: string;
  autoLoad: boolean;
  restoreState: boolean;
  parent?: ApiEntityParent;
  actions: ApiEntityAction[];
  states: ApiEntityState[];
  scripts: ApiScript[];
  scriptIds: number[];
  attributes: Record<string, ApiAttribute>;
  settings: Record<string, ApiAttribute>;
  metrics: ApiMetric[];
  isLoaded?: boolean;
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
  tags: string[];
}

export interface ApiEntityAction {
  name: string;
  description: string;
  icon?: string;
  image?: ApiImage;
  script?: ApiScript;
  /** @format int64 */
  scriptId?: number;
  type: string;
}

export interface ApiEntityCallActionRequest {
  id?: string;
  name: string;
  attributes: Record<string, ApiAttribute>;
  tags: string[];
  /** @format int64 */
  areaId?: number;
}

export interface ApiEntityParent {
  id: string;
}

export interface ApiEntityRequest {
  id: string;
  name: string;
}

export interface ApiEntityShort {
  id: string;
  pluginName: string;
  description: string;
  area?: ApiArea;
  icon?: string;
  autoLoad: boolean;
  restoreState: boolean;
  parentId?: string;
  isLoaded?: boolean;
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
  tags: string[];
}

export interface ApiEntityState {
  name: string;
  description: string;
  icon?: string;
  image?: ApiImage;
  style: string;
}

export interface ApiEntityStorageFilter {
  entityId: string;
  description: string;
}

export interface ApiEntityStorage {
  /** @format int64 */
  id: number;
  entityId: string;
  entity_description: string;
  state: string;
  state_description: string;
  attributes: Record<string, ApiAttribute>;
  /** @format date-time */
  createdAt: string;
}

export interface ApiEventBusStateListResult {
  items: ApiBusStateItem[];
  meta?: ApiMeta;
}

export interface ApiExecScriptResult {
  result: string;
}

export interface ApiExecSrcScriptRequest {
  lang: string;
  name: string;
  source: string;
  description: string;
}

export interface ApiGetActionListResult {
  items: ApiAction[];
  meta?: ApiMeta;
}

export interface ApiGetAreaListResult {
  items: ApiArea[];
  meta?: ApiMeta;
}

export interface ApiGetBackupListResult {
  items: string[];
  meta?: ApiMeta;
}

export interface ApiGetBridgeListResult {
  items: ApiZigbee2MqttShort[];
  meta?: ApiMeta;
}

export interface ApiGetClientListResult {
  items: ApiClient[];
  meta?: ApiMeta;
}

export interface ApiGetConditionListResult {
  items: ApiCondition[];
  meta?: ApiMeta;
}

export interface ApiGetDashboardCardItemListResult {
  items: ApiDashboardCardItem[];
  meta?: ApiMeta;
}

export interface ApiGetDashboardCardListResult {
  items: ApiDashboardCard[];
  meta?: ApiMeta;
}

export interface ApiGetDashboardListResult {
  items: ApiDashboardShort[];
  meta?: ApiMeta;
}

export interface ApiGetDashboardTabListResult {
  items: ApiDashboardTabShort[];
  meta?: ApiMeta;
}

export interface ApiGetEntityListResult {
  items: ApiEntityShort[];
  meta?: ApiMeta;
}

export interface ApiGetEntityStorageResult {
  items: ApiEntityStorage[];
  filter: ApiEntityStorageFilter[];
  meta: ApiMeta;
}

export interface ApiGetImageFilterListResult {
  items: GetImageFilterListResultfilter[];
}

export interface ApiGetImageListByDateResult {
  items: ApiImage[];
}

export interface ApiGetImageListResult {
  items: ApiImage[];
  meta?: ApiMeta;
}

export interface ApiGetLogListResult {
  items: ApiLog[];
  meta?: ApiMeta;
}

export interface ApiGetMessageDeliveryListResult {
  items: ApiMessageDelivery[];
  meta?: ApiMeta;
}

export interface ApiGetPluginListResult {
  items: ApiPluginShort[];
  meta?: ApiMeta;
}

export interface ApiGetRoleListResult {
  items: ApiRole[];
  meta?: ApiMeta;
}

export interface ApiGetScriptListResult {
  items: ApiScript[];
  meta?: ApiMeta;
}

export interface ApiGetSubscriptionListResult {
  items: ApiSubscription[];
  meta?: ApiMeta;
}

export interface ApiGetTaskListResult {
  items: ApiTask[];
  meta?: ApiMeta;
}

export interface ApiGetTriggerListResult {
  items: ApiTrigger[];
  meta?: ApiMeta;
}

export interface ApiGetUserListResult {
  items: ApiUserShot[];
  meta?: ApiMeta;
}

export interface ApiGetVariableListResult {
  items: ApiVariable[];
  meta?: ApiMeta;
}

export interface ApiGetTagListResult {
  items: ApiTag[];
  meta?: ApiMeta;
}

export interface ApiBackup {
  name: string;
  /** @format int64 */
  size: number;
  /** @format uint32 */
  fileMode: number;
  /** @format date-time */
  modTime: string;
}

export interface ApiImage {
  /** @format int64 */
  id: number;
  thumb: string;
  url: string;
  image: string;
  mimeType: string;
  title: string;
  /** @format int64 */
  size: number;
  name: string;
  /** @format date-time */
  createdAt: string;
}

export interface ApiLog {
  /** @format int64 */
  id: number;
  level: string;
  body: string;
  owner: string;
  /** @format date-time */
  createdAt: string;
}

export interface ApiMessage {
  /** @format int64 */
  id: number;
  type: string;
  entityId?: string;
  attributes: Record<string, string>;
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
}

export interface ApiMessageDelivery {
  /** @format int64 */
  id: number;
  message: ApiMessage;
  address: string;
  status: string;
  errorMessageStatus?: string;
  errorMessageBody?: string;
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
}

export interface ApiPagination {
  /** @format uint64 */
  limit: number;
  /** @format uint64 */
  page: number;
  /** @format uint64 */
  total: number;
}

export interface ApiMeta {
  pagination: ApiPagination;
  sort: string;
}

export interface ApiMetric {
  /** @format int64 */
  id: number;
  name: string;
  description: string;
  options?: ApiMetricOption;
  data: ApiMetricOptionData[];
  type: string;
  ranges: string[];
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
}

export interface ApiMetricOption {
  items: ApiMetricOptionItem[];
}

export interface ApiMetricOptionData {
  value: Record<string, any>;
  /** @format int64 */
  metricId?: number;
  /** @format date-time */
  time: string;
}

export interface ApiMetricOptionItem {
  name: string;
  description: string;
  color: string;
  translate: string;
  label: string;
}

export interface ApiNetworkmapResponse {
  networkmap: string;
}

export interface ApiNewActionRequest {
  name: string;
  description: string;
  /** @format int64 */
  scriptId?: number;
  /** @format int64 */
  areaId?: number;
  entityId?: string;
  entityActionName?: string;
}

export interface ApiNewAreaRequest {
  name: string;
  description: string;
  polygon: ApiAreaLocation[];
  center?: ApiAreaLocation;
  /** @format float */
  zoom: number;
  /** @format float */
  resolution: number;
}

export interface ApiNewConditionRequest {
  name: string;
  description: string;
  /** @format int64 */
  scriptId?: number;
  /** @format int64 */
  areaId?: number;
}

export interface ApiNewDashboardCardItemRequest {
  title: string;
  type: string;
  /** @format int32 */
  weight: number;
  enabled: boolean;
  /** @format int64 */
  dashboardCardId: number;
  entityId?: string;
  /**
   * @format byte
   * @pattern ^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$
   */
  payload: string;
  hidden: boolean;
  frozen: boolean;
}

export interface ApiNewDashboardCardRequest {
  title: string;
  /** @format int32 */
  height: number;
  /** @format int32 */
  width: number;
  background?: string;
  /** @format int32 */
  weight: number;
  enabled: boolean;
  /** @format int64 */
  dashboardTabId: number;
  /**
   * @format byte
   * @pattern ^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$
   */
  payload: string;
  hidden: boolean;
  entityId?: string;
}

export interface ApiNewDashboardRequest {
  name: string;
  description: string;
  enabled: boolean;
  /** @format int64 */
  areaId?: number;
}

export interface ApiNewDashboardTabRequest {
  name: string;
  /** @format int32 */
  columnWidth: number;
  gap: boolean;
  background?: string;
  icon: string;
  enabled: boolean;
  /** @format int32 */
  weight: number;
  /** @format int64 */
  dashboardId: number;
  /**
   * @format byte
   * @pattern ^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$
   */
  payload: string;
}

export interface ApiNewEntityRequest {
  name: string;
  pluginName: string;
  description: string;
  /** @format int64 */
  areaId?: number;
  icon?: string;
  /** @format int64 */
  imageId?: number;
  autoLoad: boolean;
  restoreState: boolean;
  parentId?: string;
  actions: ApiNewEntityRequestAction[];
  states: ApiNewEntityRequestState[];
  attributes: Record<string, ApiAttribute>;
  settings: Record<string, ApiAttribute>;
  metrics: ApiMetric[];
  scriptIds: number[];
  tags: string[];
}

export interface ApiNewEntityRequestAction {
  name: string;
  description: string;
  icon?: string;
  /** @format int64 */
  imageId?: number;
  /** @format int64 */
  scriptId?: number;
  type: string;
}

export interface ApiNewEntityRequestState {
  name: string;
  description: string;
  icon?: string;
  /** @format int64 */
  imageId?: number;
  style: string;
}

export interface ApiNewImageRequest {
  thumb: string;
  image: string;
  mimeType: string;
  title: string;
  /** @format int64 */
  size: number;
  name: string;
}

export interface ApiNewRoleRequest {
  name: string;
  description: string;
  parent?: string;
}

export interface ApiNewScriptRequest {
  lang: string;
  name: string;
  source: string;
  description: string;
}

export interface ApiNewTaskRequest {
  name: string;
  description: string;
  enabled: boolean;
  condition: string;
  triggerIds: number[];
  conditionIds: number[];
  actionIds: number[];
  /** @format int64 */
  areaId?: number;
}

export interface ApiNewTriggerRequest {
  name: string;
  description: string;
  entityIds: string[];
  script?: ApiScript;
  /** @format int64 */
  scriptId?: number;
  pluginName: string;
  attributes: Record<string, ApiAttribute>;
  enabled: boolean;
  /** @format int64 */
  areaId?: number;
}

export interface ApiNewVariableRequest {
  name: string;
  value: string;
  tags: string[];
}

export interface ApiNewZigbee2MqttRequest {
  name: string;
  login: string;
  password?: string;
  permitJoin: boolean;
  baseTopic: string;
}

export interface ApiNewtUserRequest {
  nickname: string;
  firstName?: string;
  lastName?: string;
  password: string;
  passwordRepeat: string;
  email: string;
  status?: string;
  lang?: string;
  /** @format int64 */
  imageId?: number;
  roleName: string;
  meta?: ApiUserMeta[];
}

export interface ApiPasswordResetRequest {
  email: string;
  token?: string;
  newPassword?: string;
}

export interface ApiPlugin {
  name: string;
  version: string;
  enabled: boolean;
  system: boolean;
  actor: boolean;
  settings: Record<string, ApiAttribute>;
  options?: ApiPluginOptionsResult;
  isLoaded?: boolean;
}

export interface ApiPluginOptionsResult {
  triggers: boolean;
  actors: boolean;
  actorCustomAttrs: boolean;
  actorAttrs: Record<string, ApiAttribute>;
  actorCustomActions: boolean;
  actorActions: Record<string, ApiPluginOptionsResultEntityAction>;
  actorCustomStates: boolean;
  actorStates: Record<string, ApiPluginOptionsResultEntityState>;
  actorCustomSetts: boolean;
  actorSetts: Record<string, ApiAttribute>;
  setts: Record<string, ApiAttribute>;
}

export interface ApiPluginOptionsResultEntityAction {
  name: string;
  description: string;
  imageUrl: string;
  icon: string;
}

export interface ApiPluginOptionsResultEntityState {
  name: string;
  description: string;
  imageUrl: string;
  icon: string;
}

export interface ApiPluginShort {
  name: string;
  version: string;
  enabled: boolean;
  system: boolean;
  actor?: boolean;
  isLoaded?: boolean;
}

export interface ApiReloadRequest {
  id: string;
}

export interface ApiResponse {
  id?: string;
  query?: string;
  /**
   * @format byte
   * @pattern ^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$
   */
  body?: string;
}

export interface ApiRestoreBackupRequest {
  name: string;
}

export interface ApiRole {
  parent?: ApiRole;
  name: string;
  description: string;
  children: ApiRole[];
  accessList?: ApiRoleAccessList;
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
}

export interface ApiRoleAccessList {
  levels: Record<string, AccessListListOfString>;
}

export interface ApiRoleAccessListResult {
  levels: Record<string, ApiAccessLevels>;
}

export interface ApiScriptVersion {
  /** @format int64 */
  id: number;
  lang: string;
  source: string;
  /** @format date-time */
  createdAt: string;
}

export interface ApiTag {
  /** @format int64 */
  id: number;
  name: string;
}

export interface ApiScript {
  /** @format int64 */
  id: number;
  lang: string;
  name: string;
  source: string;
  description: string;
  scriptInfo?: ApiScriptInfo;
  versions: ApiScriptVersion[];
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
}

export interface ApiScriptInfo {
  /** @format int32 */
  alexaIntents: number;
  /** @format int32 */
  entityActions: number;
  /** @format int32 */
  entityScripts: number;
  /** @format int32 */
  automationTriggers: number;
  /** @format int32 */
  automationConditions: number;
  /** @format int32 */
  automationActions: number;
}

export interface ApiSearchActionResult {
  items: ApiAction[];
}

export interface ApiSearchAreaResult {
  items: ApiArea[];
}

export interface ApiSearchConditionResult {
  items: ApiCondition[];
}

export interface ApiSearchDashboardResult {
  items: ApiDashboard[];
}

export interface ApiSearchDeviceResult {
  items: ApiZigbee2MqttDevice[];
}

export interface ApiSearchEntityResult {
  items: ApiEntityShort[];
}

export interface ApiSearchPluginResult {
  items: ApiPluginShort[];
}

export interface ApiSearchRoleListResult {
  items: ApiRole[];
}

export interface ApiSearchScriptListResult {
  items: ApiScript[];
}

export interface ApiSearchTagListResult {
  items: ApiTag[];
}

export interface ApiSearchTriggerResult {
  items: ApiTrigger[];
}

export interface ApiSearchVariableResult {
  items: ApiVariable[];
}

export interface ApiSigninResponse {
  currentUser?: ApiCurrentUser;
  accessToken: string;
}

export interface ApiStatistic {
  name: string;
  description: string;
  /** @format int32 */
  value: number;
  /** @format int32 */
  diff: number;
}

export interface ApiStatistics {
  items: ApiStatistic[];
}

export interface ApiSubscription {
  /** @format uint32 */
  id: number;
  clientId: string;
  topicName: string;
  name: string;
  /** @format uint32 */
  qos: number;
  noLocal: boolean;
  retainAsPublished: boolean;
  /** @format uint32 */
  retainHandling: number;
}

export interface ApiTask {
  /** @format int64 */
  id: number;
  name: string;
  description: string;
  enabled: boolean;
  condition: string;
  triggers: ApiTrigger[];
  conditions: ApiCondition[];
  actions: ApiAction[];
  area?: ApiArea;
  /** @format int64 */
  areaId?: number;
  isLoaded?: boolean;
  triggerIds: number[];
  conditionIds: number[];
  actionIds: number[];
  completed?: boolean;
  telemetry: ApiTelemetryItem[];
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
}

export interface ApiTelemetryItem {
  name: string;
  /** @format int32 */
  num: number;
  /** @format int64 */
  start: number;
  /** @format int64 */
  end?: number;
  /** @format int64 */
  timeEstimate: number;
  attributes: Record<string, string>;
  status: string;
  /** @format int32 */
  level: number;
}

export interface ApiTrigger {
  /** @format int64 */
  id: number;
  name: string;
  description: string;
  entities: ApiEntityShort[];
  entityIds: string[];
  script?: ApiScript;
  /** @format int64 */
  scriptId?: number;
  area?: ApiArea;
  /** @format int64 */
  areaId?: number;
  pluginName: string;
  attributes: Record<string, ApiAttribute>;
  enabled: boolean;
  isLoaded?: boolean;
  completed?: boolean;
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
}

/** @default "INT" */
export enum ApiTypes {
  INT = "INT",
  STRING = "STRING",
  FLOAT = "FLOAT",
  BOOL = "BOOL",
  ARRAY = "ARRAY",
  MAP = "MAP",
  TIME = "TIME",
  IMAGE = "IMAGE",
  ICON = "ICON",
  POINT = "POINT",
  ENCRYPTED = "ENCRYPTED",
}

export interface ApiUpdateEntityRequestAction {
  name: string;
  description: string;
  icon?: string;
  /** @format int64 */
  imageId?: number;
  /** @format int64 */
  scriptId?: number;
  type: string;
}

export interface ApiUpdateEntityRequestState {
  name: string;
  description: string;
  icon?: string;
  /** @format int64 */
  imageId?: number;
  style: string;
}

export interface ApiUserFull {
  /** @format int64 */
  id: number;
  nickname: string;
  firstName?: string;
  lastName?: string;
  email: string;
  status: string;
  history: ApiUserHistory[];
  image?: ApiImage;
  /** @format int64 */
  signInCount: number;
  meta: ApiUserMeta[];
  role: ApiRole;
  roleName: string;
  lang: string;
  authenticationToken: string;
  currentSignInIp?: string;
  lastSignInIp?: string;
  user?: ApiUserFullParent;
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
  /** @format date-time */
  currentSignInAt?: string;
  /** @format date-time */
  lastSignInAt?: string;
  /** @format date-time */
  resetPasswordSentAt?: string;
  /** @format date-time */
  deletedAt?: string;
}

export interface ApiUserFullParent {
  /** @format int64 */
  id: number;
  nickname: string;
}

export interface ApiUserHistory {
  ip: string;
  /** @format date-time */
  time: string;
}

export interface ApiUserMeta {
  key: string;
  value: string;
}

export interface ApiUserShot {
  /** @format int64 */
  id: number;
  nickname: string;
  firstName?: string;
  lastName?: string;
  email: string;
  status: string;
  image?: ApiImage;
  lang: string;
  role: ApiRole;
  roleName: string;
  user?: ApiUserShotParent;
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
}

export interface ApiUserShotParent {
  /** @format int64 */
  id: number;
  nickname: string;
}

export interface ApiVariable {
  name: string;
  value: string;
  system: boolean;
  tags: string[];
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
}

export interface ApiZigbee2Mqtt {
  scanInProcess: boolean;
  /** @format date-time */
  lastScan?: string;
  networkmap: string;
  status: string;
  /** @format int64 */
  id: number;
  name: string;
  login: string;
  permitJoin: boolean;
  baseTopic: string;
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
}

export interface ApiZigbee2MqttDevice {
  id: string;
  /** @format int64 */
  zigbee2mqttId: number;
  name: string;
  type: string;
  model: string;
  description: string;
  manufacturer: string;
  functions: string[];
  imageUrl: string;
  icon: string;
  status: string;
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
}

export interface ApiZigbee2MqttShort {
  /** @format int64 */
  id: number;
  name: string;
  login: string;
  permitJoin: boolean;
  baseTopic: string;
  /** @format date-time */
  createdAt: string;
  /** @format date-time */
  updatedAt: string;
}

/**
 * `Any` contains an arbitrary serialized protocol buffer message along with a
 * URL that describes the type of the serialized message.
 *
 * Protobuf library provides support to pack/unpack Any values in the form
 * of utility functions or additional generated methods of the Any type.
 *
 * Example 1: Pack and unpack a message in C++.
 *
 *     Foo foo = ...;
 *     Any any;
 *     any.PackFrom(foo);
 *     ...
 *     if (any.UnpackTo(&foo)) {
 *       ...
 *     }
 *
 * Example 2: Pack and unpack a message in Java.
 *
 *     Foo foo = ...;
 *     Any any = Any.pack(foo);
 *     ...
 *     if (any.is(Foo.class)) {
 *       foo = any.unpack(Foo.class);
 *     }
 *
 * Example 3: Pack and unpack a message in Python.
 *
 *     foo = Foo(...)
 *     any = Any()
 *     any.Pack(foo)
 *     ...
 *     if any.Is(Foo.DESCRIPTOR):
 *       any.Unpack(foo)
 *       ...
 *
 * Example 4: Pack and unpack a message in Go
 *
 *      foo := &pb.Foo{...}
 *      any, err := anypb.New(foo)
 *      if err != nil {
 *        ...
 *      }
 *      ...
 *      foo := &pb.Foo{}
 *      if err := any.UnmarshalTo(foo); err != nil {
 *        ...
 *      }
 *
 * The pack methods provided by protobuf library will by default use
 * 'type.googleapis.com/full.type.name' as the type URL and the unpack
 * methods only use the fully qualified type name after the last '/'
 * in the type URL, for example "foo.bar.com/x/y.z" will yield type
 * name "y.z".
 *
 *
 * JSON
 *
 * The JSON representation of an `Any` value uses the regular
 * representation of the deserialized, embedded message, with an
 * additional field `@type` which contains the type URL. Example:
 *
 *     package google.profile;
 *     message Person {
 *       string first_name = 1;
 *       string last_name = 2;
 *     }
 *
 *     {
 *       "@type": "type.googleapis.com/google.profile.Person",
 *       "firstName": <string>,
 *       "lastName": <string>
 *     }
 *
 * If the embedded message type is well-known and has a custom JSON
 * representation, that representation will be embedded adding a field
 * `value` which holds the custom JSON in addition to the `@type`
 * field. Example (for message [google.protobuf.Duration][]):
 *
 *     {
 *       "@type": "type.googleapis.com/google.protobuf.Duration",
 *       "value": "1.212s"
 *     }
 */
export interface ProtobufAny {
  /**
   * A URL/resource name that uniquely identifies the type of the serialized
   * protocol buffer message. This string must contain at least
   * one "/" character. The last segment of the URL's path must represent
   * the fully qualified name of the type (as in
   * `path/google.protobuf.Duration`). The name should be in a canonical form
   * (e.g., leading "." is not accepted).
   *
   * In practice, teams usually precompile into the binary all types that they
   * expect it to use in the context of Any. However, for URLs which use the
   * scheme `http`, `https`, or no scheme, one can optionally set up a type
   * server that maps type URLs to message definitions as follows:
   *
   * * If no scheme is provided, `https` is assumed.
   * * An HTTP GET on the URL must yield a [google.protobuf.Type][]
   *   value in binary format, or produce an error.
   * * Applications are allowed to cache lookup results based on the
   *   URL, or have them precompiled into a binary to avoid any
   *   lookup. Therefore, binary compatibility needs to be preserved
   *   on changes to types. (Use versioned type names to manage
   *   breaking changes.)
   *
   * Note: this functionality is not currently available in the official
   * protobuf release, and it is not used for type URLs beginning with
   * type.googleapis.com.
   *
   * Schemes other than `http`, `https` (or the empty scheme) might be
   * used with implementation specific semantics.
   */
  "@type"?: string;
  [key: string]: any;
}

export interface RpcStatus {
  /** @format int32 */
  code?: number;
  message?: string;
  details?: ProtobufAny[];
}

import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, HeadersDefaults, ResponseType } from "axios";

export type QueryParamsType = Record<string | number, any>;

export interface FullRequestParams extends Omit<AxiosRequestConfig, "data" | "params" | "url" | "responseType"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseType;
  /** request body */
  body?: unknown;
}

export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;

export interface ApiConfig<SecurityDataType = unknown> extends Omit<AxiosRequestConfig, "data" | "cancelToken"> {
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<AxiosRequestConfig | void> | AxiosRequestConfig | void;
  secure?: boolean;
  format?: ResponseType;
}

export enum ContentType {
  Json = "application/json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
  Text = "text/plain",
}

export class HttpClient<SecurityDataType = unknown> {
  public instance: AxiosInstance;
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private secure?: boolean;
  private format?: ResponseType;

  constructor({ securityWorker, secure, format, ...axiosConfig }: ApiConfig<SecurityDataType> = {}) {
    this.instance = axios.create({ ...axiosConfig, baseURL: axiosConfig.baseURL || "/" });
    this.secure = secure;
    this.format = format;
    this.securityWorker = securityWorker;
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected mergeRequestParams(params1: AxiosRequestConfig, params2?: AxiosRequestConfig): AxiosRequestConfig {
    const method = params1.method || (params2 && params2.method);

    return {
      ...this.instance.defaults,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...((method && this.instance.defaults.headers[method.toLowerCase() as keyof HeadersDefaults]) || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected stringifyFormItem(formItem: unknown) {
    if (typeof formItem === "object" && formItem !== null) {
      return JSON.stringify(formItem);
    } else {
      return `${formItem}`;
    }
  }

  protected createFormData(input: Record<string, unknown>): FormData {
    return Object.keys(input || {}).reduce((formData, key) => {
      const property = input[key];
      const propertyContent: any[] = property instanceof Array ? property : [property];

      for (const formItem of propertyContent) {
        const isFileType = formItem instanceof Blob || formItem instanceof File;
        formData.append(key, isFileType ? formItem : this.stringifyFormItem(formItem));
      }

      return formData;
    }, new FormData());
  }

  public request = async <T = any, _E = any>({
    secure,
    path,
    type,
    query,
    format,
    body,
    ...params
  }: FullRequestParams): Promise<AxiosResponse<T>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const responseFormat = format || this.format || undefined;

    if (type === ContentType.FormData && body && body !== null && typeof body === "object") {
      body = this.createFormData(body as Record<string, unknown>);
    }

    if (type === ContentType.Text && body && body !== null && typeof body !== "string") {
      body = JSON.stringify(body);
    }

    return this.instance.request({
      ...requestParams,
      headers: {
        ...(requestParams.headers || {}),
        ...(type && type !== ContentType.FormData ? { "Content-Type": type } : {}),
      },
      params: query,
      responseType: responseFormat,
      data: body,
      url: path,
    });
  };
}

/**
 * @title Smart home api
 * @version 1.0
 * @baseUrl /
 * @contact Alex Filippov <support@e154.ru> (https://e154.github.io/smart-home/)
 *
 * This documentation describes APIs found under https://github.com/e154/smart-home
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  v1 = {
    /**
     * No description
     *
     * @tags AuthService
     * @name AuthServiceAccessList
     * @summary get user access list object
     * @request GET:/v1/access_list
     * @secure
     */
    authServiceAccessList: (params: RequestParams = {}) =>
      this.request<
        ApiAccessListResponse,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/access_list`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ActionService
     * @name ActionServiceAddAction
     * @summary add new action
     * @request POST:/v1/action
     * @secure
     */
    actionServiceAddAction: (data: ApiNewActionRequest, params: RequestParams = {}) =>
      this.request<
        ApiAction,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/action`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ActionService
     * @name ActionServiceGetActionById
     * @summary get action by id
     * @request GET:/v1/action/{id}
     * @secure
     */
    actionServiceGetActionById: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiAction,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/action/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ActionService
     * @name ActionServiceUpdateAction
     * @summary update action
     * @request PUT:/v1/action/{id}
     * @secure
     */
    actionServiceUpdateAction: (
      id: number,
      data: {
        name: string;
        description: string;
        /** @format int64 */
        scriptId?: number;
        /** @format int64 */
        areaId?: number;
        entityId?: string;
        entityActionName?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiAction,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/action/${id}`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ActionService
     * @name ActionServiceDeleteAction
     * @summary delete action
     * @request DELETE:/v1/action/{id}
     * @secure
     */
    actionServiceDeleteAction: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/action/${id}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ActionService
     * @name ActionServiceGetActionList
     * @summary get action list
     * @request GET:/v1/actions
     * @secure
     */
    actionServiceGetActionList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
        /** The number of results returned on a page */
        "ids[]"?: number[];
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetActionListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/actions`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ActionService
     * @name ActionServiceSearchAction
     * @summary search action
     * @request GET:/v1/actions/search
     * @secure
     */
    actionServiceSearchAction: (
      query?: {
        query?: string;
        /** @format int64 */
        offset?: number;
        /** @format int64 */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiSearchActionResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/actions/search`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags AreaService
     * @name AreaServiceAddArea
     * @summary add new area
     * @request POST:/v1/area
     * @secure
     */
    areaServiceAddArea: (data: ApiNewAreaRequest, params: RequestParams = {}) =>
      this.request<
        ApiArea,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/area`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags AreaService
     * @name AreaServiceGetAreaById
     * @summary get area by id
     * @request GET:/v1/area/{id}
     * @secure
     */
    areaServiceGetAreaById: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiArea,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/area/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags AreaService
     * @name AreaServiceUpdateArea
     * @summary update area
     * @request PUT:/v1/area/{id}
     * @secure
     */
    areaServiceUpdateArea: (
      id: number,
      data: {
        name: string;
        description: string;
        polygon: ApiAreaLocation[];
        center?: ApiAreaLocation;
        /** @format float */
        zoom: number;
        /** @format float */
        resolution: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiArea,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/area/${id}`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags AreaService
     * @name AreaServiceDeleteArea
     * @summary delete area
     * @request DELETE:/v1/area/{id}
     * @secure
     */
    areaServiceDeleteArea: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/area/${id}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags AreaService
     * @name AreaServiceGetAreaList
     * @summary get area list
     * @request GET:/v1/areas
     * @secure
     */
    areaServiceGetAreaList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetAreaListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/areas`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags AreaService
     * @name AreaServiceSearchArea
     * @summary search area
     * @request GET:/v1/areas/search
     * @secure
     */
    areaServiceSearchArea: (
      query?: {
        query?: string;
        /** @format int64 */
        offset?: number;
        /** @format int64 */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiSearchAreaResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/areas/search`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags BackupService
     * @name BackupServiceGetBackupList
     * @summary get backup list
     * @request GET:/v1/backups
     * @secure
     */
    backupServiceGetBackupList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetBackupListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/backups`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags BackupService
     * @name BackupServiceNewBackup
     * @summary new backup
     * @request POST:/v1/backups
     * @secure
     */
    backupServiceNewBackup: (data: ApiDisablePluginResult, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/backups`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags BackupService
     * @name BackupServiceUploadBackup
     * @summary upload backup file
     * @request POST:/v1/backup/upload
     * @secure
     */
    backupServiceUploadBackup: (
      data: {
        filename?: File[];
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiBackup,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/backup/upload`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.FormData,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags BackupService
     * @name BackupServiceApplyState
     * @summary apply state
     * @request POST:/v1/backup/apply
     * @secure
     */
    backupServiceApplyState: (params: RequestParams = {}) =>
      this.request<
        void,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/backup/apply`,
        method: "POST",
        secure: true,
        ...params,
      }),

    /**
     * No description
     *
     * @tags BackupService
     * @name BackupServiceRevertState
     * @summary revert state
     * @request POST:/v1/backup/rollback
     * @secure
     */
    backupServiceRevertState: (params: RequestParams = {}) =>
      this.request<
        void,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/backup/rollback`,
        method: "POST",
        secure: true,
        ...params,
      }),

    /**
     * No description
     *
     * @tags BackupService
     * @name BackupServiceRestoreBackup
     * @summary restore backup
     * @request PUT:/v1/backup/{name}
     * @secure
     */
    backupServiceRestoreBackup: (name: string, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/backup/${name}`,
        method: "PUT",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags BackupService
     * @name BackupServiceDeleteBackup
     * @summary delete backup
     * @request DELETE:/v1/backup/{name}
     * @secure
     */
    backupServiceDeleteBackup: (name: string, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/backup/${name}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ConditionService
     * @name ConditionServiceAddCondition
     * @summary add new condition
     * @request POST:/v1/condition
     * @secure
     */
    conditionServiceAddCondition: (data: ApiNewConditionRequest, params: RequestParams = {}) =>
      this.request<
        ApiCondition,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/condition`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ConditionService
     * @name ConditionServiceGetConditionById
     * @summary get condition by id
     * @request GET:/v1/condition/{id}
     * @secure
     */
    conditionServiceGetConditionById: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiCondition,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/condition/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ConditionService
     * @name ConditionServiceUpdateCondition
     * @summary update condition
     * @request PUT:/v1/condition/{id}
     * @secure
     */
    conditionServiceUpdateCondition: (
      id: number,
      data: {
        name: string;
        description: string;
        /** @format int64 */
        scriptId?: number;
        /** @format int64 */
        areaId?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiCondition,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/condition/${id}`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ConditionService
     * @name ConditionServiceDeleteCondition
     * @summary delete condition
     * @request DELETE:/v1/condition/{id}
     * @secure
     */
    conditionServiceDeleteCondition: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/condition/${id}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ConditionService
     * @name ConditionServiceGetConditionList
     * @summary get condition list
     * @request GET:/v1/conditions
     * @secure
     */
    conditionServiceGetConditionList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
        /** The number of results returned on a page */
        "ids[]"?: number[];
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetConditionListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/conditions`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ConditionService
     * @name ConditionServiceSearchCondition
     * @summary search condition
     * @request GET:/v1/conditions/search
     * @secure
     */
    conditionServiceSearchCondition: (
      query?: {
        query?: string;
        /** @format int64 */
        offset?: number;
        /** @format int64 */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiSearchConditionResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/conditions/search`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardService
     * @name DashboardServiceAddDashboard
     * @summary add new dashboard
     * @request POST:/v1/dashboard
     * @secure
     */
    dashboardServiceAddDashboard: (data: ApiNewDashboardRequest, params: RequestParams = {}) =>
      this.request<
        ApiDashboard,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/dashboard`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardService
     * @name DashboardServiceGetDashboardById
     * @summary get dashboard by id
     * @request GET:/v1/dashboard/{id}
     * @secure
     */
    dashboardServiceGetDashboardById: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDashboard,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/dashboard/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardService
     * @name DashboardServiceUpdateDashboard
     * @summary update dashboard
     * @request PUT:/v1/dashboard/{id}
     * @secure
     */
    dashboardServiceUpdateDashboard: (
      id: number,
      data: {
        name: string;
        description: string;
        enabled: boolean;
        /** @format int64 */
        areaId?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiDashboard,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/dashboard/${id}`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardService
     * @name DashboardServiceDeleteDashboard
     * @summary delete dashboard
     * @request DELETE:/v1/dashboard/{id}
     * @secure
     */
    dashboardServiceDeleteDashboard: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/dashboard/${id}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardCardService
     * @name DashboardCardServiceAddDashboardCard
     * @summary add new dashboard_card
     * @request POST:/v1/dashboard_card
     * @secure
     */
    dashboardCardServiceAddDashboardCard: (data: ApiNewDashboardCardRequest, params: RequestParams = {}) =>
      this.request<
        ApiDashboardCard,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/dashboard_card`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardCardService
     * @name DashboardCardServiceImportDashboardCard
     * @summary import dashboard_card
     * @request POST:/v1/dashboard_card/import
     * @secure
     */
    dashboardCardServiceImportDashboardCard: (data: ApiDashboardCard, params: RequestParams = {}) =>
      this.request<
        ApiDashboardCard,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/dashboard_card/import`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardCardService
     * @name DashboardCardServiceGetDashboardCardById
     * @summary get dashboard_card by id
     * @request GET:/v1/dashboard_card/{id}
     * @secure
     */
    dashboardCardServiceGetDashboardCardById: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDashboardCard,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/dashboard_card/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardCardService
     * @name DashboardCardServiceUpdateDashboardCard
     * @summary update dashboard_card
     * @request PUT:/v1/dashboard_card/{id}
     * @secure
     */
    dashboardCardServiceUpdateDashboardCard: (
      id: number,
      data: {
        title: string;
        /** @format int32 */
        height: number;
        /** @format int32 */
        width: number;
        background?: string;
        /** @format int32 */
        weight: number;
        enabled: boolean;
        /** @format int64 */
        dashboardTabId: number;
        /**
         * @format byte
         * @pattern ^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$
         */
        payload: string;
        items: UpdateDashboardCardRequestItem[];
        hidden: boolean;
        entityId?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiDashboardCard,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/dashboard_card/${id}`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardCardService
     * @name DashboardCardServiceDeleteDashboardCard
     * @summary delete dashboard_card
     * @request DELETE:/v1/dashboard_card/{id}
     * @secure
     */
    dashboardCardServiceDeleteDashboardCard: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/dashboard_card/${id}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardCardItemService
     * @name DashboardCardItemServiceAddDashboardCardItem
     * @summary add new dashboard_card
     * @request POST:/v1/dashboard_card_item
     * @secure
     */
    dashboardCardItemServiceAddDashboardCardItem: (data: ApiNewDashboardCardItemRequest, params: RequestParams = {}) =>
      this.request<
        ApiDashboardCardItem,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/dashboard_card_item`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardCardItemService
     * @name DashboardCardItemServiceGetDashboardCardItemById
     * @summary get dashboard_card_item by id
     * @request GET:/v1/dashboard_card_item/{id}
     * @secure
     */
    dashboardCardItemServiceGetDashboardCardItemById: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDashboardCardItem,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/dashboard_card_item/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardCardItemService
     * @name DashboardCardItemServiceUpdateDashboardCardItem
     * @summary update dashboard_card_item
     * @request PUT:/v1/dashboard_card_item/{id}
     * @secure
     */
    dashboardCardItemServiceUpdateDashboardCardItem: (
      id: number,
      data: {
        title: string;
        type: string;
        /** @format int32 */
        weight: number;
        enabled: boolean;
        /** @format int64 */
        dashboardCardId: number;
        entityId?: string;
        /**
         * @format byte
         * @pattern ^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$
         */
        payload: string;
        hidden: boolean;
        frozen: boolean;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiDashboardCardItem,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/dashboard_card_item/${id}`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardCardItemService
     * @name DashboardCardItemServiceDeleteDashboardCardItem
     * @summary delete dashboard_card_item
     * @request DELETE:/v1/dashboard_card_item/{id}
     * @secure
     */
    dashboardCardItemServiceDeleteDashboardCardItem: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/dashboard_card_item/${id}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardCardItemService
     * @name DashboardCardItemServiceGetDashboardCardItemList
     * @summary get dashboard_card_item list
     * @request GET:/v1/dashboard_card_items
     * @secure
     */
    dashboardCardItemServiceGetDashboardCardItemList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetDashboardCardItemListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/dashboard_card_items`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardCardService
     * @name DashboardCardServiceGetDashboardCardList
     * @summary get dashboard_card list
     * @request GET:/v1/dashboard_cards
     * @secure
     */
    dashboardCardServiceGetDashboardCardList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetDashboardCardListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/dashboard_cards`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardTabService
     * @name DashboardTabServiceAddDashboardTab
     * @summary add new dashboard_tab
     * @request POST:/v1/dashboard_tab
     * @secure
     */
    dashboardTabServiceAddDashboardTab: (data: ApiNewDashboardTabRequest, params: RequestParams = {}) =>
      this.request<
        ApiDashboardTab,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/dashboard_tab`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardTabService
     * @name DashboardTabServiceGetDashboardTabById
     * @summary get dashboard_tab by id
     * @request GET:/v1/dashboard_tab/{id}
     * @secure
     */
    dashboardTabServiceGetDashboardTabById: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDashboardTab,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/dashboard_tab/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardTabService
     * @name DashboardTabServiceUpdateDashboardTab
     * @summary update dashboard
     * @request PUT:/v1/dashboard_tab/{id}
     * @secure
     */
    dashboardTabServiceUpdateDashboardTab: (
      id: number,
      data: {
        name: string;
        /** @format int32 */
        columnWidth: number;
        gap: boolean;
        background?: string;
        icon: string;
        enabled: boolean;
        /** @format int32 */
        weight: number;
        /** @format int64 */
        dashboardId: number;
        /**
         * @format byte
         * @pattern ^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$
         */
        payload: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiDashboardTab,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/dashboard_tab/${id}`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardTabService
     * @name DashboardTabServiceDeleteDashboardTab
     * @summary delete dashboard_tab
     * @request DELETE:/v1/dashboard_tab/{id}
     * @secure
     */
    dashboardTabServiceDeleteDashboardTab: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/dashboard_tab/${id}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardTabService
     * @name DashboardTabServiceGetDashboardTabList
     * @summary get dashboard_tab list
     * @request GET:/v1/dashboard_tabs
     * @secure
     */
    dashboardTabServiceGetDashboardTabList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetDashboardTabListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/dashboard_tabs`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardTabService
     * @name DashboardTabServiceImportDashboardTab
     * @summary import dashboard_tab
     * @request POST:/v1/dashboard_tabs/import
     * @secure
     */
    dashboardTabServiceImportDashboardTab: (data: ApiDashboardTab, params: RequestParams = {}) =>
      this.request<
        ApiDashboardTab,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/dashboard_tabs/import`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardService
     * @name DashboardServiceGetDashboardList
     * @summary get dashboard list
     * @request GET:/v1/dashboards
     * @secure
     */
    dashboardServiceGetDashboardList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetDashboardListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/dashboards`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardService
     * @name DashboardServiceImportDashboard
     * @summary import dashboard
     * @request POST:/v1/dashboards/import
     * @secure
     */
    dashboardServiceImportDashboard: (data: ApiDashboard, params: RequestParams = {}) =>
      this.request<
        ApiDashboard,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/dashboards/import`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardService
     * @name DashboardServiceSearchDashboard
     * @summary search area
     * @request GET:/v1/dashboards/search
     * @secure
     */
    dashboardServiceSearchDashboard: (
      query?: {
        query?: string;
        /** @format int64 */
        offset?: number;
        /** @format int64 */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiSearchDashboardResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/dashboards/search`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DeveloperToolsService
     * @name DeveloperToolsServiceCallAction
     * @summary call action
     * @request POST:/v1/developer_tools/automation/call_action
     * @secure
     */
    developerToolsServiceCallAction: (data: ApiAutomationRequest, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/developer_tools/automation/call_action`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DeveloperToolsService
     * @name DeveloperToolsServiceCallTrigger
     * @summary call trigger
     * @request POST:/v1/developer_tools/automation/call_trigger
     * @secure
     */
    developerToolsServiceCallTrigger: (data: ApiAutomationRequest, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/developer_tools/automation/call_trigger`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DeveloperToolsService
     * @name DeveloperToolsServiceGetEventBusStateList
     * @summary bas state
     * @request GET:/v1/developer_tools/bus/state
     * @secure
     */
    developerToolsServiceGetEventBusStateList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiEventBusStateListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/developer_tools/bus/state`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DeveloperToolsService
     * @name DeveloperToolsServiceReloadEntity
     * @summary reload entity
     * @request POST:/v1/developer_tools/entity/reload
     * @secure
     */
    developerToolsServiceReloadEntity: (data: ApiReloadRequest, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/developer_tools/entity/reload`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DeveloperToolsService
     * @name DeveloperToolsServiceEntitySetState
     * @summary entity set state
     * @request POST:/v1/developer_tools/entity/set_state
     * @secure
     */
    developerToolsServiceEntitySetState: (data: ApiEntityRequest, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/developer_tools/entity/set_state`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags EntityService
     * @name EntityServiceGetEntityList
     * @summary get entity list
     * @request GET:/v1/entities
     * @secure
     */
    entityServiceGetEntityList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
        query?: string;
        "tags[]"?: string[];
        plugin?: string;
        /** @format int64 */
        area?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetEntityListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/entities`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags EntityService
     * @name EntityServiceImportEntity
     * @summary import entity
     * @request POST:/v1/entities/import
     * @secure
     */
    entityServiceImportEntity: (data: ApiEntity, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/entities/import`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags EntityService
     * @name EntityServiceAddEntity
     * @summary add new entity
     * @request POST:/v1/entity
     * @secure
     */
    entityServiceAddEntity: (data: ApiNewEntityRequest, params: RequestParams = {}) =>
      this.request<
        ApiEntity,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/entity`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags EntityService
     * @name EntityServiceSearchEntity
     * @summary search entity
     * @request GET:/v1/entity/search
     * @secure
     */
    entityServiceSearchEntity: (
      query?: {
        query?: string;
        /** @format int64 */
        offset?: number;
        /** @format int64 */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiSearchEntityResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/entity/search`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags EntityService
     * @name EntityServiceGetEntity
     * @summary get entity
     * @request GET:/v1/entity/{id}
     * @secure
     */
    entityServiceGetEntity: (id: string, params: RequestParams = {}) =>
      this.request<
        ApiEntity,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/entity/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags EntityService
     * @name EntityServiceUpdateEntity
     * @summary update entity
     * @request PUT:/v1/entity/{id}
     * @secure
     */
    entityServiceUpdateEntity: (
      id: string,
      data: {
        id: string;
        name?: string;
        pluginName: string;
        description: string;
        /** @format int64 */
        areaId?: number;
        icon?: string;
        /** @format int64 */
        imageId?: number;
        autoLoad: boolean;
        restoreState: boolean;
        parentId?: string;
        actions: ApiUpdateEntityRequestAction[];
        states: ApiUpdateEntityRequestState[];
        attributes: Record<string, ApiAttribute>;
        settings: Record<string, ApiAttribute>;
        scriptIds: number[];
        metrics: ApiMetric[];
        tags: string[];
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiEntity,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/entity/${id}`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags EntityService
     * @name EntityServiceDeleteEntity
     * @summary delete entity
     * @request DELETE:/v1/entity/{id}
     * @secure
     */
    entityServiceDeleteEntity: (id: string, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/entity/${id}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags EntityService
     * @name EntityServiceDisabledEntity
     * @summary disable entity
     * @request POST:/v1/entity/{id}/disable
     * @secure
     */
    entityServiceDisabledEntity: (id: string, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/entity/${id}/disable`,
        method: "POST",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags EntityService
     * @name EntityServiceEnabledEntity
     * @summary enabled entity
     * @request POST:/v1/entity/{id}/enable
     * @secure
     */
    entityServiceEnabledEntity: (id: string, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/entity/${id}/enable`,
        method: "POST",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags EntityService
     * @name EntityServiceGetStatistic
     * @summary get statistic
     * @request GET:/v1/entities/statistic
     * @secure
     */
    entityServiceGetStatistic: (params: RequestParams = {}) =>
      this.request<
        ApiStatistics,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/entities/statistic`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags EntityStorageService
     * @name EntityStorageServiceGetEntityStorageList
     * @request GET:/v1/entity_storage
     * @secure
     */
    entityStorageServiceGetEntityStorageList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
        /** @format date-time */
        startDate?: string;
        /** @format date-time */
        endDate?: string;
        "entityId[]"?: string[];
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetEntityStorageResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/entity_storage`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ImageService
     * @name ImageServiceAddImage
     * @summary add new image
     * @request POST:/v1/image
     * @secure
     */
    imageServiceAddImage: (data: ApiNewImageRequest, params: RequestParams = {}) =>
      this.request<
        ApiImage,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/image`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ImageService
     * @name ImageServiceUploadImage
     * @summary upload image
     * @request POST:/v1/image/upload
     * @secure
     */
    imageServiceUploadImage: (
      data: {
        filename?: File[];
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiImage,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/image/upload`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.FormData,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ImageService
     * @name ImageServiceGetImageById
     * @summary get image by id
     * @request GET:/v1/image/{id}
     * @secure
     */
    imageServiceGetImageById: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiImage,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/image/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ImageService
     * @name ImageServiceUpdateImageById
     * @summary update image
     * @request PUT:/v1/image/{id}
     * @secure
     */
    imageServiceUpdateImageById: (
      id: number,
      data: {
        thumb: string;
        image: string;
        mimeType: string;
        title: string;
        /** @format int64 */
        size: number;
        name: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiImage,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/image/${id}`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ImageService
     * @name ImageServiceDeleteImageById
     * @summary delete image by id
     * @request DELETE:/v1/image/{id}
     * @secure
     */
    imageServiceDeleteImageById: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/image/${id}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ImageService
     * @name ImageServiceGetImageList
     * @summary get image list
     * @request GET:/v1/images
     * @secure
     */
    imageServiceGetImageList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetImageListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/images`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ImageService
     * @name ImageServiceGetImageFilterList
     * @summary get image filter list
     * @request GET:/v1/images/filter_list
     * @secure
     */
    imageServiceGetImageFilterList: (params: RequestParams = {}) =>
      this.request<
        ApiGetImageFilterListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/images/filter_list`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ImageService
     * @name ImageServiceGetImageListByDate
     * @summary get image list by date
     * @request GET:/v1/images/filtered
     * @secure
     */
    imageServiceGetImageListByDate: (
      query?: {
        filter?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetImageListByDateResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/images/filtered`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags InteractService
     * @name InteractServiceEntityCallAction
     * @summary entity call action
     * @request POST:/v1/interact/entity/call_action
     * @secure
     */
    interactServiceEntityCallAction: (data: ApiEntityCallActionRequest, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/interact/entity/call_action`,
        method: "POST",
        body: data,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags LogService
     * @name LogServiceGetLogList
     * @request GET:/v1/logs
     * @secure
     */
    logServiceGetLogList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
        /** @format date-time */
        startDate?: string;
        /** @format date-time */
        endDate?: string;
        query?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetLogListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/logs`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags MessageDeliveryService
     * @name MessageDeliveryServiceGetMessageDeliveryList
     * @summary get list
     * @request GET:/v1/message_delivery
     * @secure
     */
    messageDeliveryServiceGetMessageDeliveryList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
        /** @format date-time */
        startDate?: string;
        /** @format date-time */
        endDate?: string;
        messageType?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetMessageDeliveryListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/message_delivery`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags MetricService
     * @name MetricServiceGetMetric
     * @summary get metric
     * @request GET:/v1/metric
     * @secure
     */
    metricServiceGetMetric: (
      query: {
        /** @format int64 */
        id: number;
        range?: "6h" | "12h" | "24h" | "7d" | "30d" | "1m";
        /** @format date-time */
        startDate?: string;
        /** @format date-time */
        endDate?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiMetric,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/metric`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags MqttService
     * @name MqttServiceGetClientById
     * @summary get client by id
     * @request GET:/v1/mqtt/client/{id}
     * @secure
     */
    mqttServiceGetClientById: (id: string, params: RequestParams = {}) =>
      this.request<
        ApiClient,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/mqtt/client/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags MqttService
     * @name MqttServiceGetClientList
     * @summary get client list
     * @request GET:/v1/mqtt/clients
     * @secure
     */
    mqttServiceGetClientList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetClientListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/mqtt/clients`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags MqttService
     * @name MqttServiceGetSubscriptionList
     * @summary get subscription list
     * @request GET:/v1/mqtt/subscriptions
     * @secure
     */
    mqttServiceGetSubscriptionList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
        clientId?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetSubscriptionListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/mqtt/subscriptions`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags AuthService
     * @name AuthServicePasswordReset
     * @summary sign out user
     * @request POST:/v1/password_reset
     * @secure
     */
    authServicePasswordReset: (data: ApiPasswordResetRequest, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/password_reset`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags PluginService
     * @name PluginServiceGetPlugin
     * @summary get plugin
     * @request GET:/v1/plugin/{name}
     * @secure
     */
    pluginServiceGetPlugin: (name: string, params: RequestParams = {}) =>
      this.request<
        ApiPlugin,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/plugin/${name}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags PluginService
     * @name PluginServiceDisablePlugin
     * @summary disable plugin
     * @request POST:/v1/plugin/{name}/disable
     * @secure
     */
    pluginServiceDisablePlugin: (name: string, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/plugin/${name}/disable`,
        method: "POST",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags PluginService
     * @name PluginServiceEnablePlugin
     * @summary enable plugin
     * @request POST:/v1/plugin/{name}/enable
     * @secure
     */
    pluginServiceEnablePlugin: (name: string, params: RequestParams = {}) =>
      this.request<
        ApiEnablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/plugin/${name}/enable`,
        method: "POST",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags PluginService
     * @name PluginServiceUpdatePluginSettings
     * @summary update plugin settings
     * @request PUT:/v1/plugin/{name}/settings
     * @secure
     */
    pluginServiceUpdatePluginSettings: (
      name: string,
      data: {
        settings: Record<string, ApiAttribute>;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/plugin/${name}/settings`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags PluginService
     * @name PluginServiceGetPluginReadme
     * @summary get plugin readme
     * @request GET:/v1/plugin/{name}/readme
     * @secure
     */
    pluginServiceGetPluginReadme: (
      name: string,
      query?: {
        lang?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        string,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/plugin/${name}/readme`,
        method: "GET",
        query: query,
        secure: true,
        ...params,
      }),

    /**
     * No description
     *
     * @tags PluginService
     * @name PluginServiceGetPluginList
     * @summary get plugin list
     * @request GET:/v1/plugins
     * @secure
     */
    pluginServiceGetPluginList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetPluginListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/plugins`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags PluginService
     * @name PluginServiceSearchPlugin
     * @summary search plugin
     * @request GET:/v1/plugins/search
     * @secure
     */
    pluginServiceSearchPlugin: (
      query?: {
        query?: string;
        /** @format int64 */
        offset?: number;
        /** @format int64 */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiSearchPluginResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/plugins/search`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags RoleService
     * @name RoleServiceAddRole
     * @summary add new role
     * @request POST:/v1/role
     * @secure
     */
    roleServiceAddRole: (data: ApiNewRoleRequest, params: RequestParams = {}) =>
      this.request<
        ApiRole,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/role`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags RoleService
     * @name RoleServiceGetRoleByName
     * @summary get role by name
     * @request GET:/v1/role/{name}
     * @secure
     */
    roleServiceGetRoleByName: (name: string, params: RequestParams = {}) =>
      this.request<
        ApiRole,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/role/${name}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags RoleService
     * @name RoleServiceUpdateRoleByName
     * @summary update role
     * @request PUT:/v1/role/{name}
     * @secure
     */
    roleServiceUpdateRoleByName: (
      name: string,
      data: {
        description: string;
        parent?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiRole,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/role/${name}`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags RoleService
     * @name RoleServiceDeleteRoleByName
     * @summary delete role by name
     * @request DELETE:/v1/role/{name}
     * @secure
     */
    roleServiceDeleteRoleByName: (name: string, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/role/${name}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags RoleService
     * @name RoleServiceGetRoleAccessList
     * @summary get role access list
     * @request GET:/v1/role/{name}/access_list
     * @secure
     */
    roleServiceGetRoleAccessList: (name: string, params: RequestParams = {}) =>
      this.request<
        ApiRoleAccessListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/role/${name}/access_list`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags RoleService
     * @name RoleServiceUpdateRoleAccessList
     * @summary update role access list
     * @request PUT:/v1/role/{name}/access_list
     * @secure
     */
    roleServiceUpdateRoleAccessList: (
      name: string,
      data: {
        accessListDiff: Record<string, UpdateRoleAccessListRequestAccessListDiff>;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiRoleAccessListResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/role/${name}/access_list`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags RoleService
     * @name RoleServiceGetRoleList
     * @summary get role list
     * @request GET:/v1/roles
     * @secure
     */
    roleServiceGetRoleList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetRoleListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/roles`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags RoleService
     * @name RoleServiceSearchRoleByName
     * @summary delete role by name
     * @request GET:/v1/roles/search
     * @secure
     */
    roleServiceSearchRoleByName: (
      query?: {
        query?: string;
        /** @format int64 */
        offset?: number;
        /** @format int64 */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiSearchRoleListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/roles/search`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ScriptService
     * @name ScriptServiceAddScript
     * @summary add new script
     * @request POST:/v1/script
     * @secure
     */
    scriptServiceAddScript: (data: ApiNewScriptRequest, params: RequestParams = {}) =>
      this.request<
        ApiScript,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/script`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ScriptService
     * @name ScriptServiceExecSrcScriptById
     * @summary exec src script by id
     * @request POST:/v1/script/exec_src
     * @secure
     */
    scriptServiceExecSrcScriptById: (data: ApiExecSrcScriptRequest, params: RequestParams = {}) =>
      this.request<
        ApiExecScriptResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/script/exec_src`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ScriptService
     * @name ScriptServiceGetScriptById
     * @summary get script by id
     * @request GET:/v1/script/{id}
     * @secure
     */
    scriptServiceGetScriptById: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiScript,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/script/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ScriptService
     * @name ScriptServiceUpdateScriptById
     * @summary update script
     * @request PUT:/v1/script/{id}
     * @secure
     */
    scriptServiceUpdateScriptById: (
      id: number,
      data: {
        lang: string;
        name: string;
        source: string;
        description: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiScript,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/script/${id}`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ScriptService
     * @name ScriptServiceDeleteScriptById
     * @summary delete script by id
     * @request DELETE:/v1/script/{id}
     * @secure
     */
    scriptServiceDeleteScriptById: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/script/${id}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ScriptService
     * @name ScriptServiceCopyScriptById
     * @summary copy script by id
     * @request POST:/v1/script/{id}/copy
     * @secure
     */
    scriptServiceCopyScriptById: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiScript,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/script/${id}/copy`,
        method: "POST",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ScriptService
     * @name ScriptServiceGetCompiledScriptById
     * @summary get compiled script by id
     * @request GET:/v1/script/{id}/compiled
     * @secure
     */
    scriptServiceGetCompiledScriptById: (id: number, params: RequestParams = {}) =>
      this.request<
        string,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/script/${id}/compiled`,
        method: "GET",
        secure: true,
        ...params,
      }),

    /**
     * No description
     *
     * @tags ScriptService
     * @name ScriptServiceExecScriptById
     * @summary exec script by id
     * @request POST:/v1/script/{id}/exec
     * @secure
     */
    scriptServiceExecScriptById: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiExecScriptResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/script/${id}/exec`,
        method: "POST",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ScriptService
     * @name ScriptServiceGetScriptList
     * @summary get script list
     * @request GET:/v1/scripts
     * @secure
     */
    scriptServiceGetScriptList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
        /** The number of results returned on a page */
        "ids[]"?: number[];
        query?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetScriptListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/scripts`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ScriptService
     * @name ScriptServiceSearchScript
     * @summary search script by name
     * @request GET:/v1/scripts/search
     * @secure
     */
    scriptServiceSearchScript: (
      query?: {
        query?: string;
        /** @format int64 */
        offset?: number;
        /** @format int64 */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiSearchScriptListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/scripts/search`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags ScriptService
     * @name ScriptServiceGetStatistic
     * @summary get statistic
     * @request GET:/v1/scripts/statistic
     * @secure
     */
    scriptServiceGetStatistic: (params: RequestParams = {}) =>
      this.request<
        ApiStatistics,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/scripts/statistic`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags TagService
     * @name TagServiceSearchTag
     * @summary search tag by name
     * @request GET:/v1/tags/search
     * @secure
     */
    tagServiceSearchTag: (
      query?: {
        query?: string;
        /** @format int64 */
        offset?: number;
        /** @format int64 */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiSearchTagListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/tags/search`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags TagService
     * @name TagServiceGetTagList
     * @summary get tag list
     * @request GET:/v1/tags
     * @secure
     */
    tagServiceGetTagList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
        query?: string;
        "tags[]"?: string[];
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetTagListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/tags`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags TagService
     * @name TagServiceGetTagById
     * @summary get tag by id
     * @request GET:/v1/tag/{id}
     * @secure
     */
    tagServiceGetTagById: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiTag,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/tag/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags TagService
     * @name TagServiceUpdateTagById
     * @summary update tag
     * @request PUT:/v1/tag/{id}
     * @secure
     */
    tagServiceUpdateTagById: (
      id: number,
      data: {
        name: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiTag,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/tag/${id}`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags TagService
     * @name TagServiceDeleteTagById
     * @summary delete tag by id
     * @request DELETE:/v1/tag/{id}
     * @secure
     */
    tagServiceDeleteTagById: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/tag/${id}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags AuthService
     * @name AuthServiceSignin
     * @summary sign in user
     * @request POST:/v1/signin
     * @secure
     */
    authServiceSignin: (params: RequestParams = {}) =>
      this.request<
        ApiSigninResponse,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/signin`,
        method: "POST",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags AuthService
     * @name AuthServiceSignout
     * @summary sign out user
     * @request POST:/v1/signout
     * @secure
     */
    authServiceSignout: (params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/signout`,
        method: "POST",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags AutomationService
     * @name AutomationServiceAddTask
     * @summary add new task
     * @request POST:/v1/task
     * @secure
     */
    automationServiceAddTask: (data: ApiNewTaskRequest, params: RequestParams = {}) =>
      this.request<
        ApiTask,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/task`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags AutomationService
     * @name AutomationServiceGetTask
     * @summary get task
     * @request GET:/v1/task/{id}
     * @secure
     */
    automationServiceGetTask: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiTask,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/task/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags AutomationService
     * @name AutomationServiceUpdateTask
     * @summary update task
     * @request PUT:/v1/task/{id}
     * @secure
     */
    automationServiceUpdateTask: (
      id: number,
      data: {
        name: string;
        description: string;
        enabled: boolean;
        condition: string;
        triggerIds: number[];
        conditionIds: number[];
        actionIds: number[];
        /** @format int64 */
        areaId?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiTask,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/task/${id}`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags AutomationService
     * @name AutomationServiceDeleteTask
     * @summary delete task
     * @request DELETE:/v1/task/{id}
     * @secure
     */
    automationServiceDeleteTask: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/task/${id}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags AutomationService
     * @name AutomationServiceDisableTask
     * @summary disable task
     * @request POST:/v1/task/{id}/disable
     * @secure
     */
    automationServiceDisableTask: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/task/${id}/disable`,
        method: "POST",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags AutomationService
     * @name AutomationServiceEnableTask
     * @summary enable task
     * @request POST:/v1/task/{id}/enable
     * @secure
     */
    automationServiceEnableTask: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/task/${id}/enable`,
        method: "POST",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags AutomationService
     * @name AutomationServiceGetTaskList
     * @summary get task list
     * @request GET:/v1/tasks
     * @secure
     */
    automationServiceGetTaskList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetTaskListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/tasks`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags AutomationService
     * @name AutomationServiceImportTask
     * @summary import task
     * @request POST:/v1/tasks/import
     * @secure
     */
    automationServiceImportTask: (data: ApiTask, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/tasks/import`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags TriggerService
     * @name TriggerServiceAddTrigger
     * @summary add new trigger
     * @request POST:/v1/trigger
     * @secure
     */
    triggerServiceAddTrigger: (data: ApiNewTriggerRequest, params: RequestParams = {}) =>
      this.request<
        ApiTrigger,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/trigger`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags TriggerService
     * @name TriggerServiceGetTriggerById
     * @summary get trigger by id
     * @request GET:/v1/trigger/{id}
     * @secure
     */
    triggerServiceGetTriggerById: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiTrigger,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/trigger/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags TriggerService
     * @name TriggerServiceUpdateTrigger
     * @summary update trigger
     * @request PUT:/v1/trigger/{id}
     * @secure
     */
    triggerServiceUpdateTrigger: (
      id: number,
      data: {
        name: string;
        description: string;
        entityIds: string[];
        script?: ApiScript;
        /** @format int64 */
        scriptId?: number;
        /** @format int64 */
        areaId?: number;
        pluginName: string;
        attributes: Record<string, ApiAttribute>;
        enabled: boolean;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiTrigger,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/trigger/${id}`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags TriggerService
     * @name TriggerServiceDeleteTrigger
     * @summary delete trigger
     * @request DELETE:/v1/trigger/{id}
     * @secure
     */
    triggerServiceDeleteTrigger: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/trigger/${id}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags TriggerService
     * @name TriggerServiceGetTriggerList
     * @summary get trigger list
     * @request GET:/v1/triggers
     * @secure
     */
    triggerServiceGetTriggerList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
        /** The number of results returned on a page */
        "ids[]"?: number[];
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetTriggerListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/triggers`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags TriggerService
     * @name TriggerServiceSearchTrigger
     * @summary search trigger
     * @request GET:/v1/triggers/search
     * @secure
     */
    triggerServiceSearchTrigger: (
      query?: {
        query?: string;
        /** @format int64 */
        offset?: number;
        /** @format int64 */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiSearchTriggerResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/triggers/search`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags TriggerService
     * @name TriggerServiceDisableTrigger
     * @summary disable triggers
     * @request POST:/v1/triggers/{id}/disable
     * @secure
     */
    triggerServiceDisableTrigger: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/triggers/${id}/disable`,
        method: "POST",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags TriggerService
     * @name TriggerServiceEnableTrigger
     * @summary enable triggers
     * @request POST:/v1/triggers/{id}/enable
     * @secure
     */
    triggerServiceEnableTrigger: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/triggers/${id}/enable`,
        method: "POST",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags UserService
     * @name UserServiceAddUser
     * @summary add new user
     * @request POST:/v1/user
     * @secure
     */
    userServiceAddUser: (data: ApiNewtUserRequest, params: RequestParams = {}) =>
      this.request<
        ApiUserFull,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/user`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags UserService
     * @name UserServiceGetUserById
     * @summary get user by id
     * @request GET:/v1/user/{id}
     * @secure
     */
    userServiceGetUserById: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiUserFull,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/user/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags UserService
     * @name UserServiceUpdateUserById
     * @summary update user by id
     * @request PUT:/v1/user/{id}
     * @secure
     */
    userServiceUpdateUserById: (
      id: number,
      data: {
        nickname: string;
        firstName: string;
        lastName?: string;
        password: string;
        passwordRepeat: string;
        email: string;
        status: string;
        lang: string;
        /** @format int64 */
        imageId?: number;
        roleName: string;
        meta?: ApiUserMeta[];
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiUserFull,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/user/${id}`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags UserService
     * @name UserServiceDeleteUserById
     * @summary delete user by id
     * @request DELETE:/v1/user/{id}
     * @secure
     */
    userServiceDeleteUserById: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/user/${id}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags UserService
     * @name UserServiceGetUserList
     * @summary get user list
     * @request GET:/v1/users
     * @secure
     */
    userServiceGetUserList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetUserListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/users`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags VariableService
     * @name VariableServiceAddVariable
     * @summary add new variable
     * @request POST:/v1/variable
     * @secure
     */
    variableServiceAddVariable: (data: ApiNewVariableRequest, params: RequestParams = {}) =>
      this.request<
        ApiVariable,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/variable`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags VariableService
     * @name VariableServiceGetVariableByName
     * @summary get variable by name
     * @request GET:/v1/variable/{name}
     * @secure
     */
    variableServiceGetVariableByName: (name: string, params: RequestParams = {}) =>
      this.request<
        ApiVariable,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/variable/${name}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags VariableService
     * @name VariableServiceUpdateVariable
     * @summary update variable
     * @request PUT:/v1/variable/{name}
     * @secure
     */
    variableServiceUpdateVariable: (
      name: string,
      data: {
        value: string;
        tags: string[];
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiVariable,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/variable/${name}`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags VariableService
     * @name VariableServiceDeleteVariable
     * @summary delete variable
     * @request DELETE:/v1/variable/{name}
     * @secure
     */
    variableServiceDeleteVariable: (name: string, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/variable/${name}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags VariableService
     * @name VariableServiceGetVariableList
     * @summary get variable list
     * @request GET:/v1/variables
     * @secure
     */
    variableServiceGetVariableList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetVariableListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/variables`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags VariableService
     * @name VariableServiceSearchVariable
     * @summary search trigger
     * @request GET:/v1/variables/search
     * @secure
     */
    variableServiceSearchVariable: (
      query?: {
        query?: string;
        /** @format int64 */
        offset?: number;
        /** @format int64 */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiSearchVariableResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/variables/search`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Zigbee2mqttService
     * @name Zigbee2MqttServiceGetBridgeList
     * @summary get bridge list
     * @request GET:/v1/zigbee2mqtt/bridge
     * @secure
     */
    zigbee2MqttServiceGetBridgeList: (
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiGetBridgeListResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/zigbee2mqtt/bridge`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Zigbee2mqttService
     * @name Zigbee2MqttServiceAddZigbee2MqttBridge
     * @summary add new bridge
     * @request POST:/v1/zigbee2mqtt/bridge
     * @secure
     */
    zigbee2MqttServiceAddZigbee2MqttBridge: (data: ApiNewZigbee2MqttRequest, params: RequestParams = {}) =>
      this.request<
        ApiZigbee2Mqtt,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/zigbee2mqtt/bridge`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Zigbee2mqttService
     * @name Zigbee2MqttServiceGetZigbee2MqttBridge
     * @summary get bridge
     * @request GET:/v1/zigbee2mqtt/bridge/{id}
     * @secure
     */
    zigbee2MqttServiceGetZigbee2MqttBridge: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiZigbee2Mqtt,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/zigbee2mqtt/bridge/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Zigbee2mqttService
     * @name Zigbee2MqttServiceDeleteBridgeById
     * @summary delete bridge by id
     * @request DELETE:/v1/zigbee2mqtt/bridge/{id}
     * @secure
     */
    zigbee2MqttServiceDeleteBridgeById: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/zigbee2mqtt/bridge/${id}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Zigbee2mqttService
     * @name Zigbee2MqttServiceUpdateBridgeById
     * @summary update bridge by id
     * @request PUT:/v1/zigbee2mqtt/bridge/{id}/bridge
     * @secure
     */
    zigbee2MqttServiceUpdateBridgeById: (
      id: number,
      data: {
        name: string;
        login: string;
        password?: string;
        permitJoin: boolean;
        baseTopic: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiZigbee2Mqtt,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/zigbee2mqtt/bridge/${id}/bridge`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Zigbee2mqttService
     * @name Zigbee2MqttServiceDeviceList
     * @summary list device
     * @request GET:/v1/zigbee2mqtt/bridge/{id}/devices
     * @secure
     */
    zigbee2MqttServiceDeviceList: (
      id: number,
      query?: {
        /**
         * Field on which to sort and its direction
         * @example "-created_at"
         */
        sort?: string;
        /**
         * Page number of the requested result set
         * @format uint64
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * The number of results returned on a page
         * @format uint64
         */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiDeviceListResult,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/zigbee2mqtt/bridge/${id}/devices`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Zigbee2mqttService
     * @name Zigbee2MqttServiceNetworkmap
     * @summary networkmap
     * @request GET:/v1/zigbee2mqtt/bridge/{id}/networkmap
     * @secure
     */
    zigbee2MqttServiceNetworkmap: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiNetworkmapResponse,
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
        | {
            error?: GenericErrorResponse;
          }
      >({
        path: `/v1/zigbee2mqtt/bridge/${id}/networkmap`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Zigbee2mqttService
     * @name Zigbee2MqttServiceUpdateNetworkmap
     * @summary update networkmap
     * @request POST:/v1/zigbee2mqtt/bridge/{id}/networkmap
     * @secure
     */
    zigbee2MqttServiceUpdateNetworkmap: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/zigbee2mqtt/bridge/${id}/networkmap`,
        method: "POST",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Zigbee2mqttService
     * @name Zigbee2MqttServiceResetBridgeById
     * @summary reset bridge by id
     * @request POST:/v1/zigbee2mqtt/bridge/{id}/reset
     * @secure
     */
    zigbee2MqttServiceResetBridgeById: (id: number, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/zigbee2mqtt/bridge/${id}/reset`,
        method: "POST",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Zigbee2mqttService
     * @name Zigbee2MqttServiceDeviceBan
     * @summary device ban
     * @request POST:/v1/zigbee2mqtt/device_ban
     * @secure
     */
    zigbee2MqttServiceDeviceBan: (data: ApiDeviceBanRequest, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/zigbee2mqtt/device_ban`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Zigbee2mqttService
     * @name Zigbee2MqttServiceDeviceRename
     * @summary device rename
     * @request POST:/v1/zigbee2mqtt/device_rename
     * @secure
     */
    zigbee2MqttServiceDeviceRename: (data: ApiDeviceRenameRequest, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/zigbee2mqtt/device_rename`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Zigbee2mqttService
     * @name Zigbee2MqttServiceDeviceWhitelist
     * @summary device whitelist
     * @request POST:/v1/zigbee2mqtt/device_whitelist
     * @secure
     */
    zigbee2MqttServiceDeviceWhitelist: (data: ApiDeviceWhitelistRequest, params: RequestParams = {}) =>
      this.request<
        ApiDisablePluginResult,
        | {
            error?: GenericErrorResponse;
          }
        | {
            error?: GenericErrorResponse & {
              code?: "UNAUTHORIZED";
            };
          }
      >({
        path: `/v1/zigbee2mqtt/device_whitelist`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Zigbee2mqttService
     * @name Zigbee2MqttServiceSearchDevice
     * @summary search device
     * @request GET:/v1/zigbee2mqtt/search_device
     * @secure
     */
    zigbee2MqttServiceSearchDevice: (
      query?: {
        query?: string;
        /** @format int64 */
        offset?: number;
        /** @format int64 */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        ApiSearchDeviceResult,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/zigbee2mqtt/search_device`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags StreamService
     * @name StreamServiceSubscribe
     * @request GET:/v1/ws
     * @secure
     */
    streamServiceSubscribe: (params: RequestParams = {}) =>
      this.request<
        ApiResponse,
        {
          error?: GenericErrorResponse & {
            code?: "UNAUTHORIZED";
          };
        }
      >({
        path: `/v1/ws`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),
  };
}
