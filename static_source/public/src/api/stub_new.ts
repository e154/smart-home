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

export interface AccessListListOfString {
  items?: string[];
}

export interface GetImageFilterListResultfilter {
  date?: string;
  /** @format int32 */
  count?: number;
}

export interface NewEntityRequestActionScript {
  /** @format int64 */
  id?: string;
}

export interface UpdateDashboardCardRequestItem {
  /** @format int64 */
  id?: string;
  title?: string;
  type?: string;
  /** @format int32 */
  weight?: number;
  enabled?: boolean;
  entityId?: string;
  /** @format byte */
  payload?: string;
  hidden?: boolean;
  frozen?: boolean;
  showOn?: string[];
  hideOn?: string[];
}

export interface UpdateRoleAccessListRequestAccessListDiff {
  items?: Record<string, boolean>;
}

export interface ApiAccessItem {
  actions?: string[];
  method?: string;
  description?: string;
  roleName?: string;
}

export interface ApiAccessLevels {
  items?: Record<string, ApiAccessItem>;
}

export interface ApiAccessList {
  levels?: Record<string, ApiAccessLevels>;
}

export interface ApiAccessListResponse {
  accessList?: ApiAccessList;
}

export interface ApiAction {
  /** @format int64 */
  id?: string;
  name?: string;
  script?: ApiScript;
}

export interface ApiArea {
  /** @format int64 */
  id?: string;
  name?: string;
  description?: string;
}

export interface ApiAttribute {
  name?: string;
  type?: ApiTypes;
  /** @format int64 */
  int?: string;
  string?: string;
  bool?: boolean;
  /** @format float */
  float?: number;
  array?: ApiAttribute[];
  map?: Record<string, ApiAttribute>;
  /** @format date-time */
  time?: string;
  imageUrl?: string;
}

export interface ApiAutomationRequest {
  /** @format int64 */
  id?: string;
  name?: string;
}

export interface ApiCondition {
  /** @format int64 */
  id?: string;
  name?: string;
  script?: ApiScript;
}

export interface ApiCurrentUser {
  /** @format int64 */
  id?: string;
  nickname?: string;
  firstName?: string;
  lastName?: string;
  email?: string;
  status?: string;
  history?: ApiUserHistory[];
  image?: ApiImage;
  /** @format int64 */
  signInCount?: string;
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
  id?: string;
  name?: string;
  description?: string;
  enabled?: boolean;
  /** @format int64 */
  areaId?: string;
  area?: ApiArea;
  tabs?: ApiDashboardTab[];
  entities?: Record<string, ApiEntity>;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface ApiDashboardCard {
  /** @format int64 */
  id?: string;
  title?: string;
  /** @format int32 */
  height?: number;
  /** @format int32 */
  width?: number;
  background?: string;
  /** @format int32 */
  weight?: number;
  enabled?: boolean;
  /** @format int64 */
  dashboardTabId?: string;
  /** @format byte */
  payload?: string;
  items?: ApiDashboardCardItem[];
  entities?: Record<string, ApiEntity>;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface ApiDashboardCardItem {
  /** @format int64 */
  id?: string;
  title?: string;
  type?: string;
  /** @format int32 */
  weight?: number;
  enabled?: boolean;
  /** @format int64 */
  dashboardCardId?: string;
  entityId?: string;
  /** @format byte */
  payload?: string;
  hidden?: boolean;
  frozen?: boolean;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface ApiDashboardShort {
  /** @format int64 */
  id?: string;
  name?: string;
  description?: string;
  enabled?: boolean;
  /** @format int64 */
  areaId?: string;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface ApiDashboardTab {
  /** @format int64 */
  id?: string;
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
  dashboardId?: string;
  cards?: ApiDashboardCard[];
  entities?: Record<string, ApiEntity>;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface ApiDashboardTabShort {
  /** @format int64 */
  id?: string;
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
  dashboardId?: string;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface ApiDeviceBanRequest {
  /** @format int64 */
  id?: string;
  friendlyName?: string;
}

export interface ApiDeviceListResult {
  items?: ApiZigbee2MqttDevice[];
  meta?: ApiMeta;
}

export interface ApiDeviceRenameRequest {
  friendlyName?: string;
  newName?: string;
}

export interface ApiDeviceWhitelistRequest {
  /** @format int64 */
  id?: string;
  friendlyName?: string;
}

export type ApiDisablePluginResult = object;

export type ApiEnablePluginResult = object;

export interface ApiEntity {
  id?: string;
  pluginName?: string;
  description?: string;
  area?: ApiArea;
  image?: ApiImage;
  icon?: string;
  autoLoad?: boolean;
  parent?: ApiEntityParent;
  actions?: ApiEntityAction[];
  states?: ApiEntityState[];
  scripts?: ApiScript[];
  attributes?: Record<string, ApiAttribute>;
  settings?: Record<string, ApiAttribute>;
  metrics?: ApiMetric[];
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface ApiEntityAction {
  name?: string;
  description?: string;
  icon?: string;
  image?: ApiImage;
  script?: ApiScript;
  type?: string;
}

export interface ApiEntityCallActionRequest {
  id?: string;
  name?: string;
}

export interface ApiEntityParent {
  id?: string;
}

export interface ApiEntityRequest {
  id?: string;
  name?: string;
}

export interface ApiEntityShort {
  id?: string;
  pluginName?: string;
  description?: string;
  area?: ApiArea;
  image?: ApiImage;
  icon?: string;
  autoLoad?: boolean;
  parent?: ApiEntityParent;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface ApiEntityState {
  name?: string;
  description?: string;
  icon?: string;
  image?: ApiImage;
  style?: string;
}

export interface ApiEntityStorage {
  /** @format int64 */
  id?: string;
  entityId?: string;
  state?: string;
  attributes?: Record<string, ApiAttribute>;
  /**
   * map<string, google.protobuf.Any> attributes = 4;
   * @format date-time
   */
  createdAt?: string;
}

export interface ApiExecScriptResult {
  result?: string;
}

export interface ApiExecSrcScriptRequest {
  lang?: string;
  name?: string;
  source?: string;
  description?: string;
}

export interface ApiGetAreaListResult {
  items?: ApiArea[];
  meta?: ApiMeta;
}

export interface ApiGetBackupListResult {
  items?: string[];
}

export interface ApiGetBridgeListResult {
  items?: ApiZigbee2MqttShort[];
  meta?: ApiMeta;
}

export interface ApiGetDashboardCardItemListResult {
  items?: ApiDashboardCardItem[];
  meta?: ApiMeta;
}

export interface ApiGetDashboardCardListResult {
  items?: ApiDashboardCard[];
  meta?: ApiMeta;
}

export interface ApiGetDashboardListResult {
  items?: ApiDashboardShort[];
  meta?: ApiMeta;
}

export interface ApiGetDashboardTabListResult {
  items?: ApiDashboardTabShort[];
  meta?: ApiMeta;
}

export interface ApiGetEntityListResult {
  items?: ApiEntity[];
  meta?: ApiMeta;
}

export interface ApiGetEntityStorageResult {
  items?: ApiEntityStorage[];
  meta?: ApiMeta;
}

export interface ApiGetImageFilterListResult {
  items?: GetImageFilterListResultfilter[];
}

export interface ApiGetImageListByDateResult {
  items?: ApiImage[];
}

export interface ApiGetImageListResult {
  items?: ApiImage[];
  meta?: ApiMeta;
}

export interface ApiGetLogListResult {
  items?: ApiLog[];
  meta?: ApiMeta;
}

export interface ApiGetMessageDeliveryListResult {
  items?: ApiMessageDelivery[];
  meta?: ApiMeta;
}

export interface ApiGetPluginListResult {
  items?: ApiPlugin[];
  meta?: ApiMeta;
}

export interface ApiGetPluginOptionsResult {
  triggers?: boolean;
  actors?: boolean;
  actorCustomAttrs?: boolean;
  actorAttrs?: Record<string, ApiAttribute>;
  actorCustomActions?: boolean;
  actorActions?: Record<string, ApiGetPluginOptionsResultEntityAction>;
  actorCustomStates?: boolean;
  actorStates?: Record<string, ApiGetPluginOptionsResultEntityState>;
  actorCustomSetts?: boolean;
  actorSetts?: Record<string, ApiAttribute>;
  setts?: Record<string, ApiAttribute>;
}

export interface ApiGetPluginOptionsResultEntityAction {
  name?: string;
  description?: string;
  imageUrl?: string;
  icon?: string;
}

export interface ApiGetPluginOptionsResultEntityState {
  name?: string;
  description?: string;
  imageUrl?: string;
  icon?: string;
}

export interface ApiGetRoleListResult {
  items?: ApiRole[];
  meta?: ApiMeta;
}

export interface ApiGetScriptListResult {
  items?: ApiScript[];
  meta?: ApiMeta;
}

export interface ApiGetTaskListResult {
  items?: ApiTask[];
  meta?: ApiMeta;
}

export interface ApiGetUserListResult {
  items?: ApiUserShot[];
  meta?: ApiMeta;
}

export interface ApiGetVariableListResult {
  items?: ApiVariable[];
  meta?: ApiMeta;
}

export interface ApiImage {
  /** @format int64 */
  id?: string;
  thumb?: string;
  url?: string;
  image?: string;
  mimeType?: string;
  title?: string;
  /** @format int64 */
  size?: string;
  name?: string;
  /** @format date-time */
  createdAt?: string;
}

export interface ApiLog {
  /** @format int64 */
  id?: string;
  level?: string;
  body?: string;
  owner?: string;
  /** @format date-time */
  createdAt?: string;
}

export interface ApiMessage {
  /** @format int64 */
  id?: string;
  type?: string;
  attributes?: Record<string, string>;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface ApiMessageDelivery {
  /** @format int64 */
  id?: string;
  message?: ApiMessage;
  address?: string;
  status?: string;
  errorMessageStatus?: string;
  errorMessageBody?: string;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface ApiMeta {
  /** @format uint64 */
  limit?: string;
  /** @format uint64 */
  page?: string;
  /** @format uint64 */
  total?: string;
  sort?: string;
}

export interface ApiMetric {
  /** @format int64 */
  id?: string;
  name?: string;
  description?: string;
  options?: ApiMetricOption;
  data?: ApiMetricOptionData[];
  type?: string;
  ranges?: string[];
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface ApiMetricOption {
  items?: ApiMetricOptionItem[];
}

export interface ApiMetricOptionData {
  value?: Record<string, number>;
  /** @format int64 */
  metricId?: string;
  /** @format date-time */
  time?: string;
}

export interface ApiMetricOptionItem {
  name?: string;
  description?: string;
  color?: string;
  translate?: string;
  label?: string;
}

export interface ApiNetworkmapResponse {
  networkmap?: string;
}

export interface ApiNewAreaRequest {
  name?: string;
  description?: string;
}

export interface ApiNewDashboardCardItemRequest {
  title?: string;
  type?: string;
  /** @format int32 */
  weight?: number;
  enabled?: boolean;
  /** @format int64 */
  dashboardCardId?: string;
  entityId?: string;
  /** @format byte */
  payload?: string;
  hidden?: boolean;
  frozen?: boolean;
}

export interface ApiNewDashboardCardRequest {
  title?: string;
  /** @format int32 */
  height?: number;
  /** @format int32 */
  width?: number;
  background?: string;
  /** @format int32 */
  weight?: number;
  enabled?: boolean;
  /** @format int64 */
  dashboardTabId?: string;
  /** @format byte */
  payload?: string;
}

export interface ApiNewDashboardRequest {
  name?: string;
  description?: string;
  enabled?: boolean;
  /** @format int64 */
  areaId?: string;
}

export interface ApiNewDashboardTabRequest {
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
  dashboardId?: string;
}

export interface ApiNewEntityRequest {
  name?: string;
  pluginName?: string;
  description?: string;
  area?: ApiArea;
  icon?: string;
  image?: ApiNewEntityRequestImage;
  autoLoad?: boolean;
  parent?: ApiEntityParent;
  actions?: ApiNewEntityRequestAction[];
  states?: ApiNewEntityRequestState[];
  attributes?: Record<string, ApiAttribute>;
  settings?: Record<string, ApiAttribute>;
  metrics?: ApiMetric[];
  scripts?: ApiScript[];
}

export interface ApiNewEntityRequestAction {
  name?: string;
  description?: string;
  icon?: string;
  image?: ApiNewEntityRequestImage;
  script?: NewEntityRequestActionScript;
  type?: string;
}

export interface ApiNewEntityRequestImage {
  /** @format int64 */
  id?: string;
}

export interface ApiNewEntityRequestState {
  name?: string;
  description?: string;
  icon?: string;
  image?: ApiNewEntityRequestImage;
  style?: string;
}

export interface ApiNewImageRequest {
  thumb?: string;
  image?: string;
  mimeType?: string;
  title?: string;
  /** @format int64 */
  size?: string;
  name?: string;
}

export interface ApiNewRoleRequest {
  name?: string;
  description?: string;
  parent?: string;
}

export interface ApiNewScriptRequest {
  lang?: string;
  name?: string;
  source?: string;
  description?: string;
}

export interface ApiNewTaskRequest {
  name?: string;
  description?: string;
  enabled?: boolean;
  condition?: string;
  triggers?: ApiTrigger[];
  conditions?: ApiCondition[];
  actions?: ApiAction[];
  area?: ApiArea;
}

export interface ApiNewVariableRequest {
  name?: string;
  value?: string;
}

export interface ApiNewZigbee2MqttRequest {
  name?: string;
  login?: string;
  password?: string;
  permitJoin?: boolean;
  baseTopic?: string;
}

export interface ApiNewtUserRequest {
  nickname?: string;
  firstName?: string;
  lastName?: string;
  password?: string;
  passwordRepeat?: string;
  email?: string;
  status?: string;
  lang?: string;
  image?: ApiNewtUserRequestImage;
  role?: ApiNewtUserRequestRole;
  meta?: ApiUserMeta[];
}

export interface ApiNewtUserRequestImage {
  /** @format int64 */
  id?: string;
}

export interface ApiNewtUserRequestRole {
  name?: string;
}

export interface ApiPlugin {
  name?: string;
  version?: string;
  enabled?: boolean;
  system?: boolean;
  actor?: boolean;
  settings?: Record<string, ApiAttribute>;
}

export interface ApiReloadRequest {
  id?: string;
}

export interface ApiResponse {
  id?: string;
  query?: string;
  /** @format byte */
  body?: string;
}

export interface ApiRestoreBackupRequest {
  name?: string;
}

export interface ApiRole {
  parent?: ApiRole;
  name?: string;
  description?: string;
  children?: ApiRole[];
  accessList?: ApiRoleAccessList;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface ApiRoleAccessList {
  levels?: Record<string, AccessListListOfString>;
}

export interface ApiRoleAccessListResult {
  levels?: Record<string, ApiAccessLevels>;
}

export interface ApiScript {
  /** @format int64 */
  id?: string;
  lang?: string;
  name?: string;
  source?: string;
  description?: string;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface ApiSearchAreaResult {
  items?: ApiArea[];
}

export interface ApiSearchDashboardResult {
  items?: ApiDashboard[];
}

export interface ApiSearchDeviceResult {
  items?: ApiZigbee2MqttDevice[];
}

export interface ApiSearchEntityResult {
  items?: ApiEntityShort[];
}

export interface ApiSearchPluginResult {
  items?: ApiPlugin[];
}

export interface ApiSearchRoleListResult {
  items?: ApiRole[];
}

export interface ApiSearchScriptListResult {
  items?: ApiScript[];
}

export interface ApiSigninResponse {
  currentUser?: ApiCurrentUser;
  accessToken?: string;
}

export interface ApiTask {
  /** @format int64 */
  id?: string;
  name?: string;
  description?: string;
  enabled?: boolean;
  condition?: string;
  triggers?: ApiTrigger[];
  conditions?: ApiCondition[];
  actions?: ApiAction[];
  area?: ApiArea;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface ApiTrigger {
  /** @format int64 */
  id?: string;
  name?: string;
  entity?: ApiTriggerEntity;
  script?: ApiScript;
  pluginName?: string;
  attributes?: Record<string, ApiAttribute>;
}

export interface ApiTriggerEntity {
  id?: string;
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
}

export interface ApiUpdateEntityRequestAction {
  name?: string;
  description?: string;
  icon?: string;
  image?: ApiUpdateEntityRequestImage;
  script?: ApiScript;
  type?: string;
}

export interface ApiUpdateEntityRequestImage {
  /** @format int64 */
  id?: string;
}

export interface ApiUpdateEntityRequestState {
  name?: string;
  description?: string;
  icon?: string;
  image?: ApiUpdateEntityRequestImage;
  style?: string;
}

export interface ApiUpdateUserRequestImage {
  /** @format int64 */
  id?: string;
}

export interface ApiUpdateUserRequestRole {
  name?: string;
}

export interface ApiUploadImageRequest {
  /** @format byte */
  body?: string;
}

export interface ApiUserFull {
  /** @format int64 */
  id?: string;
  nickname?: string;
  firstName?: string;
  lastName?: string;
  email?: string;
  status?: string;
  history?: ApiUserHistory[];
  image?: ApiImage;
  /** @format int64 */
  signInCount?: string;
  meta?: ApiUserMeta[];
  role?: ApiRole;
  roleName?: string;
  lang?: string;
  authenticationToken?: string;
  currentSignInIp?: string;
  lastSignInIp?: string;
  user?: ApiUserFullParent;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
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
  id?: string;
  nickname?: string;
}

export interface ApiUserHistory {
  ip?: string;
  /** @format date-time */
  time?: string;
}

export interface ApiUserMeta {
  key?: string;
  value?: string;
}

export interface ApiUserShot {
  /** @format int64 */
  id?: string;
  nickname?: string;
  firstName?: string;
  lastName?: string;
  email?: string;
  status?: string;
  image?: ApiImage;
  lang?: string;
  role?: ApiRole;
  roleName?: string;
  user?: ApiUserShotParent;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface ApiUserShotParent {
  /** @format int64 */
  id?: string;
  nickname?: string;
}

export interface ApiVariable {
  name?: string;
  value?: string;
  system?: boolean;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface ApiZigbee2Mqtt {
  scanInProcess?: boolean;
  /** @format date-time */
  lastScan?: string;
  networkmap?: string;
  status?: string;
  /** @format int64 */
  id?: string;
  name?: string;
  login?: string;
  permitJoin?: boolean;
  baseTopic?: string;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface ApiZigbee2MqttDevice {
  id?: string;
  /** @format int64 */
  zigbee2mqttId?: string;
  name?: string;
  type?: string;
  model?: string;
  description?: string;
  manufacturer?: string;
  functions?: string[];
  imageUrl?: string;
  status?: string;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface ApiZigbee2MqttShort {
  /** @format int64 */
  id?: string;
  name?: string;
  login?: string;
  permitJoin?: boolean;
  baseTopic?: string;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
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
  typeUrl?: string;
  /**
   * Must be a valid serialized protocol buffer of the above specified type.
   * @format byte
   */
  value?: string;
}

export interface RpcStatus {
  /** @format int32 */
  code?: number;
  message?: string;
  details?: ProtobufAny[];
}

import HeadersDefaults from "axios"
import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse,  ResponseType } from "axios";

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
    this.instance = axios.create({ ...axiosConfig, baseURL: axiosConfig.baseURL || "" });
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
        // ...((method && this.instance.defaults.headers[method.toLowerCase() as keyof HeadersDefaults]) || {}),
        ...(this.instance.defaults.headers || {}),
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
      this.request<ApiAccessListResponse, RpcStatus>({
        path: `/v1/access_list`,
        method: "GET",
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
    areaServiceAddArea: (body: ApiNewAreaRequest, params: RequestParams = {}) =>
      this.request<ApiArea, RpcStatus>({
        path: `/v1/area`,
        method: "POST",
        body: body,
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
    areaServiceGetAreaById: (id: string, params: RequestParams = {}) =>
      this.request<ApiArea, RpcStatus>({
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
     * @name AreaServiceDeleteArea
     * @summary delete area
     * @request DELETE:/v1/area/{id}
     * @secure
     */
    areaServiceDeleteArea: (id: string, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
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
     * @name AreaServiceUpdateArea
     * @summary update area
     * @request PUT:/v1/area/{id}
     * @secure
     */
    areaServiceUpdateArea: (
      id: string,
      body: {
        name?: string;
        description?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiArea, RpcStatus>({
        path: `/v1/area/${id}`,
        method: "PUT",
        body: body,
        secure: true,
        type: ContentType.Json,
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
        /** @format uint64 */
        page?: string;
        /** @format uint64 */
        limit?: string;
        sort?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiGetAreaListResult, RpcStatus>({
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
        limit?: string;
        /** @format int64 */
        offset?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiSearchAreaResult, RpcStatus>({
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
     * @name BackupServiceNewBackup
     * @summary new backup
     * @request POST:/v1/backup
     * @secure
     */
    backupServiceNewBackup: (body: any, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
        path: `/v1/backup`,
        method: "POST",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags BackupService
     * @name BackupServiceRestoreBackup
     * @summary restore backup
     * @request PUT:/v1/backup/restore
     * @secure
     */
    backupServiceRestoreBackup: (body: ApiRestoreBackupRequest, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
        path: `/v1/backup/restore`,
        method: "PUT",
        body: body,
        secure: true,
        type: ContentType.Json,
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
    backupServiceGetBackupList: (params: RequestParams = {}) =>
      this.request<ApiGetBackupListResult, RpcStatus>({
        path: `/v1/backups`,
        method: "GET",
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
    dashboardServiceAddDashboard: (body: ApiNewDashboardRequest, params: RequestParams = {}) =>
      this.request<ApiDashboard, RpcStatus>({
        path: `/v1/dashboard`,
        method: "POST",
        body: body,
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
    dashboardServiceGetDashboardById: (id: string, params: RequestParams = {}) =>
      this.request<ApiDashboard, RpcStatus>({
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
     * @name DashboardServiceDeleteDashboard
     * @summary delete dashboard
     * @request DELETE:/v1/dashboard/{id}
     * @secure
     */
    dashboardServiceDeleteDashboard: (id: string, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
        path: `/v1/dashboard/${id}`,
        method: "DELETE",
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
      id: string,
      body: {
        name?: string;
        description?: string;
        enabled?: boolean;
        /** @format int64 */
        areaId?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiDashboard, RpcStatus>({
        path: `/v1/dashboard/${id}`,
        method: "PUT",
        body: body,
        secure: true,
        type: ContentType.Json,
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
    dashboardCardServiceAddDashboardCard: (body: ApiNewDashboardCardRequest, params: RequestParams = {}) =>
      this.request<ApiDashboardCard, RpcStatus>({
        path: `/v1/dashboard_card`,
        method: "POST",
        body: body,
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
    dashboardCardServiceImportDashboardCard: (body: ApiDashboardCard, params: RequestParams = {}) =>
      this.request<ApiDashboardCard, RpcStatus>({
        path: `/v1/dashboard_card/import`,
        method: "POST",
        body: body,
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
    dashboardCardServiceGetDashboardCardById: (id: string, params: RequestParams = {}) =>
      this.request<ApiDashboardCard, RpcStatus>({
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
     * @name DashboardCardServiceDeleteDashboardCard
     * @summary delete dashboard_card
     * @request DELETE:/v1/dashboard_card/{id}
     * @secure
     */
    dashboardCardServiceDeleteDashboardCard: (id: string, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
        path: `/v1/dashboard_card/${id}`,
        method: "DELETE",
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
      id: string,
      body: {
        title?: string;
        /** @format int32 */
        height?: number;
        /** @format int32 */
        width?: number;
        background?: string;
        /** @format int32 */
        weight?: number;
        enabled?: boolean;
        /** @format int64 */
        dashboardTabId?: string;
        /** @format byte */
        payload?: string;
        items?: UpdateDashboardCardRequestItem[];
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiDashboardCard, RpcStatus>({
        path: `/v1/dashboard_card/${id}`,
        method: "PUT",
        body: body,
        secure: true,
        type: ContentType.Json,
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
    dashboardCardItemServiceAddDashboardCardItem: (body: ApiNewDashboardCardItemRequest, params: RequestParams = {}) =>
      this.request<ApiDashboardCardItem, RpcStatus>({
        path: `/v1/dashboard_card_item`,
        method: "POST",
        body: body,
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
    dashboardCardItemServiceGetDashboardCardItemById: (id: string, params: RequestParams = {}) =>
      this.request<ApiDashboardCardItem, RpcStatus>({
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
     * @name DashboardCardItemServiceDeleteDashboardCardItem
     * @summary delete dashboard_card_item
     * @request DELETE:/v1/dashboard_card_item/{id}
     * @secure
     */
    dashboardCardItemServiceDeleteDashboardCardItem: (id: string, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
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
     * @name DashboardCardItemServiceUpdateDashboardCardItem
     * @summary update dashboard_card_item
     * @request PUT:/v1/dashboard_card_item/{id}
     * @secure
     */
    dashboardCardItemServiceUpdateDashboardCardItem: (
      id: string,
      body: {
        title?: string;
        type?: string;
        /** @format int32 */
        weight?: number;
        enabled?: boolean;
        /** @format int64 */
        dashboardCardId?: string;
        entityId?: string;
        /** @format byte */
        payload?: string;
        hidden?: boolean;
        frozen?: boolean;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiDashboardCardItem, RpcStatus>({
        path: `/v1/dashboard_card_item/${id}`,
        method: "PUT",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DashboardCardItemService
     * @name DashboardCardItemServiceGetDashboardCardItemList
     * @summary get dashboard_card_item list
     * @request GET:/v1/dashboard_cards
     * @secure
     */
    dashboardCardItemServiceGetDashboardCardItemList: (
      query?: {
        /** @format uint64 */
        page?: string;
        /** @format uint64 */
        limit?: string;
        sort?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiGetDashboardCardItemListResult, RpcStatus>({
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
    dashboardTabServiceAddDashboardTab: (body: ApiNewDashboardTabRequest, params: RequestParams = {}) =>
      this.request<ApiDashboardTab, RpcStatus>({
        path: `/v1/dashboard_tab`,
        method: "POST",
        body: body,
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
    dashboardTabServiceGetDashboardTabById: (id: string, params: RequestParams = {}) =>
      this.request<ApiDashboardTab, RpcStatus>({
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
     * @name DashboardTabServiceDeleteDashboardTab
     * @summary delete dashboard_tab
     * @request DELETE:/v1/dashboard_tab/{id}
     * @secure
     */
    dashboardTabServiceDeleteDashboardTab: (id: string, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
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
     * @name DashboardTabServiceUpdateDashboardTab
     * @summary update dashboard
     * @request PUT:/v1/dashboard_tab/{id}
     * @secure
     */
    dashboardTabServiceUpdateDashboardTab: (
      id: string,
      body: {
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
        dashboardId?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiDashboardTab, RpcStatus>({
        path: `/v1/dashboard_tab/${id}`,
        method: "PUT",
        body: body,
        secure: true,
        type: ContentType.Json,
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
        /** @format uint64 */
        page?: string;
        /** @format uint64 */
        limit?: string;
        sort?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiGetDashboardTabListResult, RpcStatus>({
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
     * @tags DashboardService
     * @name DashboardServiceGetDashboardList
     * @summary get dashboard list
     * @request GET:/v1/dashboards
     * @secure
     */
    dashboardServiceGetDashboardList: (
      query?: {
        /** @format uint64 */
        page?: string;
        /** @format uint64 */
        limit?: string;
        sort?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiGetDashboardListResult, RpcStatus>({
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
    dashboardServiceImportDashboard: (body: ApiDashboard, params: RequestParams = {}) =>
      this.request<ApiDashboard, RpcStatus>({
        path: `/v1/dashboards/import`,
        method: "POST",
        body: body,
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
        limit?: string;
        /** @format int64 */
        offset?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiSearchDashboardResult, RpcStatus>({
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
     * @name DeveloperToolsServiceReloadEntity
     * @summary reload entity
     * @request POST:/v1/developer_tools/entity_reload
     * @secure
     */
    developerToolsServiceReloadEntity: (body: ApiReloadRequest, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
        path: `/v1/developer_tools/entity_reload`,
        method: "POST",
        body: body,
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
     * @request POST:/v1/developer_tools/entity_set_state
     * @secure
     */
    developerToolsServiceEntitySetState: (body: ApiEntityRequest, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
        path: `/v1/developer_tools/entity_set_state`,
        method: "POST",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DeveloperToolsService
     * @name DeveloperToolsServiceTaskCallAction
     * @summary task call action
     * @request POST:/v1/developer_tools/task_call_action
     * @secure
     */
    developerToolsServiceTaskCallAction: (body: ApiAutomationRequest, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
        path: `/v1/developer_tools/task_call_action`,
        method: "POST",
        body: body,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags DeveloperToolsService
     * @name DeveloperToolsServiceTaskCallTrigger
     * @summary task call trigger
     * @request POST:/v1/developer_tools/task_call_trigger
     * @secure
     */
    developerToolsServiceTaskCallTrigger: (body: ApiAutomationRequest, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
        path: `/v1/developer_tools/task_call_trigger`,
        method: "POST",
        body: body,
        secure: true,
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
        /** @format uint64 */
        page?: string;
        /** @format uint64 */
        limit?: string;
        sort?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiGetEntityListResult, RpcStatus>({
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
    entityServiceImportEntity: (body: ApiEntity, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
        path: `/v1/entities/import`,
        method: "POST",
        body: body,
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
    entityServiceAddEntity: (body: ApiNewEntityRequest, params: RequestParams = {}) =>
      this.request<ApiEntity, RpcStatus>({
        path: `/v1/entity`,
        method: "POST",
        body: body,
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
        limit?: string;
        /** @format int64 */
        offset?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiSearchEntityResult, RpcStatus>({
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
      this.request<ApiEntity, RpcStatus>({
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
     * @name EntityServiceDeleteEntity
     * @summary delete entity
     * @request DELETE:/v1/entity/{id}
     * @secure
     */
    entityServiceDeleteEntity: (id: string, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
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
     * @name EntityServiceUpdateEntity
     * @summary update entity
     * @request PUT:/v1/entity/{id}
     * @secure
     */
    entityServiceUpdateEntity: (
      id: string,
      body: {
        name?: string;
        pluginName?: string;
        description?: string;
        area?: ApiArea;
        icon?: string;
        image?: ApiUpdateEntityRequestImage;
        autoLoad?: boolean;
        parent?: ApiEntityParent;
        actions?: ApiUpdateEntityRequestAction[];
        states?: ApiUpdateEntityRequestState[];
        attributes?: Record<string, ApiAttribute>;
        settings?: Record<string, ApiAttribute>;
        scripts?: ApiScript[];
        metrics?: ApiMetric[];
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiEntity, RpcStatus>({
        path: `/v1/entity/${id}`,
        method: "PUT",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags EntityStorageService
     * @name EntityStorageServiceGetEntityStorageList
     * @request GET:/v1/entity_storage/{entityId}
     * @secure
     */
    entityStorageServiceGetEntityStorageList: (
      entityId: string,
      query?: {
        /** @format uint64 */
        page?: string;
        /** @format uint64 */
        limit?: string;
        sort?: string;
        startDate?: string;
        endDate?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiGetEntityStorageResult, RpcStatus>({
        path: `/v1/entity_storage/${entityId}`,
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
    imageServiceAddImage: (body: ApiNewImageRequest, params: RequestParams = {}) =>
      this.request<ApiImage, RpcStatus>({
        path: `/v1/image`,
        method: "POST",
        body: body,
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
    imageServiceUploadImage: (body: ApiUploadImageRequest, params: RequestParams = {}) =>
      this.request<ApiImage, RpcStatus>({
        path: `/v1/image/upload`,
        method: "POST",
        body: body,
        secure: true,
        type: ContentType.Json,
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
    imageServiceGetImageById: (id: string, params: RequestParams = {}) =>
      this.request<ApiImage, RpcStatus>({
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
     * @name ImageServiceDeleteImageById
     * @summary delete image by id
     * @request DELETE:/v1/image/{id}
     * @secure
     */
    imageServiceDeleteImageById: (id: string, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
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
     * @name ImageServiceUpdateImageById
     * @summary update image
     * @request PUT:/v1/image/{id}
     * @secure
     */
    imageServiceUpdateImageById: (
      id: string,
      body: {
        thumb?: string;
        image?: string;
        mimeType?: string;
        title?: string;
        /** @format int64 */
        size?: string;
        name?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiImage, RpcStatus>({
        path: `/v1/image/${id}`,
        method: "PUT",
        body: body,
        secure: true,
        type: ContentType.Json,
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
        /** @format uint64 */
        page?: string;
        /** @format uint64 */
        limit?: string;
        sort?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiGetImageListResult, RpcStatus>({
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
      this.request<ApiGetImageFilterListResult, RpcStatus>({
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
      this.request<ApiGetImageListByDateResult, RpcStatus>({
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
     * @request POST:/v1/interact/entity_call_action
     * @secure
     */
    interactServiceEntityCallAction: (body: ApiEntityCallActionRequest, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
        path: `/v1/interact/entity_call_action`,
        method: "POST",
        body: body,
        secure: true,
        type: ContentType.Json,
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
        /** @format uint64 */
        page?: string;
        /** @format uint64 */
        limit?: string;
        sort?: string;
        query?: string;
        startDate?: string;
        endDate?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiGetLogListResult, RpcStatus>({
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
        /** @format uint64 */
        page?: string;
        /** @format uint64 */
        limit?: string;
        sort?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiGetMessageDeliveryListResult, RpcStatus>({
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
      query?: {
        /** @format int64 */
        id?: string;
        range?: string;
        /** @format date-time */
        from?: string;
        /** @format date-time */
        to?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiMetric, RpcStatus>({
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
     * @tags PluginService
     * @name PluginServiceDisablePlugin
     * @summary disable plugin
     * @request POST:/v1/plugin/{name}/disable
     * @secure
     */
    pluginServiceDisablePlugin: (name: string, params: RequestParams = {}) =>
      this.request<ApiDisablePluginResult, RpcStatus>({
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
      this.request<ApiEnablePluginResult, RpcStatus>({
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
     * @name PluginServiceGetPluginOptions
     * @summary get plugin options
     * @request GET:/v1/plugin/{name}/options
     * @secure
     */
    pluginServiceGetPluginOptions: (name: string, params: RequestParams = {}) =>
      this.request<ApiGetPluginOptionsResult, RpcStatus>({
        path: `/v1/plugin/${name}/options`,
        method: "GET",
        secure: true,
        format: "json",
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
        /** @format uint64 */
        page?: string;
        /** @format uint64 */
        limit?: string;
        sort?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiGetPluginListResult, RpcStatus>({
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
        limit?: string;
        /** @format int64 */
        offset?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiSearchPluginResult, RpcStatus>({
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
    roleServiceAddRole: (body: ApiNewRoleRequest, params: RequestParams = {}) =>
      this.request<ApiRole, RpcStatus>({
        path: `/v1/role`,
        method: "POST",
        body: body,
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
      this.request<ApiRole, RpcStatus>({
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
     * @name RoleServiceDeleteRoleByName
     * @summary delete role by name
     * @request DELETE:/v1/role/{name}
     * @secure
     */
    roleServiceDeleteRoleByName: (name: string, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
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
     * @name RoleServiceUpdateRoleByName
     * @summary update role
     * @request PUT:/v1/role/{name}
     * @secure
     */
    roleServiceUpdateRoleByName: (
      name: string,
      body: {
        description?: string;
        parent?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiRole, RpcStatus>({
        path: `/v1/role/${name}`,
        method: "PUT",
        body: body,
        secure: true,
        type: ContentType.Json,
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
      this.request<ApiRoleAccessListResult, RpcStatus>({
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
      body: {
        accessListDiff?: Record<string, UpdateRoleAccessListRequestAccessListDiff>;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiRoleAccessListResult, RpcStatus>({
        path: `/v1/role/${name}/access_list`,
        method: "PUT",
        body: body,
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
        /** @format uint64 */
        page?: string;
        /** @format uint64 */
        limit?: string;
        sort?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiGetRoleListResult, RpcStatus>({
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
        limit?: string;
        /** @format int64 */
        offset?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiSearchRoleListResult, RpcStatus>({
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
    scriptServiceAddScript: (body: ApiNewScriptRequest, params: RequestParams = {}) =>
      this.request<ApiScript, RpcStatus>({
        path: `/v1/script`,
        method: "POST",
        body: body,
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
    scriptServiceExecSrcScriptById: (body: ApiExecSrcScriptRequest, params: RequestParams = {}) =>
      this.request<ApiExecScriptResult, RpcStatus>({
        path: `/v1/script/exec_src`,
        method: "POST",
        body: body,
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
    scriptServiceGetScriptById: (id: string, params: RequestParams = {}) =>
      this.request<ApiScript, RpcStatus>({
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
     * @name ScriptServiceDeleteScriptById
     * @summary delete script by id
     * @request DELETE:/v1/script/{id}
     * @secure
     */
    scriptServiceDeleteScriptById: (id: string, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
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
     * @name ScriptServiceUpdateScriptById
     * @summary update script
     * @request PUT:/v1/script/{id}
     * @secure
     */
    scriptServiceUpdateScriptById: (
      id: string,
      body: {
        lang?: string;
        name?: string;
        source?: string;
        description?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiScript, RpcStatus>({
        path: `/v1/script/${id}`,
        method: "PUT",
        body: body,
        secure: true,
        type: ContentType.Json,
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
    scriptServiceCopyScriptById: (id: string, params: RequestParams = {}) =>
      this.request<ApiScript, RpcStatus>({
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
     * @name ScriptServiceExecScriptById
     * @summary exec script by id
     * @request POST:/v1/script/{id}/exec
     * @secure
     */
    scriptServiceExecScriptById: (id: string, params: RequestParams = {}) =>
      this.request<ApiExecScriptResult, RpcStatus>({
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
        /** @format uint64 */
        page?: string;
        /** @format uint64 */
        limit?: string;
        sort?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiGetScriptListResult, RpcStatus>({
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
     * @summary delete script by id
     * @request GET:/v1/scripts/search
     * @secure
     */
    scriptServiceSearchScript: (
      query?: {
        query?: string;
        /** @format int64 */
        limit?: string;
        /** @format int64 */
        offset?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiSearchScriptListResult, RpcStatus>({
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
     * @tags AuthService
     * @name AuthServiceSignin
     * @summary sign in user
     * @request POST:/v1/signin
     * @secure
     */
    authServiceSignin: (params: RequestParams = {}) =>
      this.request<ApiSigninResponse, RpcStatus>({
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
      this.request<any, RpcStatus>({
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
    automationServiceAddTask: (body: ApiNewTaskRequest, params: RequestParams = {}) =>
      this.request<ApiTask, RpcStatus>({
        path: `/v1/task`,
        method: "POST",
        body: body,
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
    automationServiceGetTask: (id: string, params: RequestParams = {}) =>
      this.request<ApiTask, RpcStatus>({
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
     * @name AutomationServiceDeleteTask
     * @summary delete task
     * @request DELETE:/v1/task/{id}
     * @secure
     */
    automationServiceDeleteTask: (id: string, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
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
     * @name AutomationServiceUpdateTask
     * @summary update task
     * @request PUT:/v1/task/{id}
     * @secure
     */
    automationServiceUpdateTask: (
      id: string,
      body: {
        name?: string;
        description?: string;
        enabled?: boolean;
        condition?: string;
        triggers?: ApiTrigger[];
        conditions?: ApiCondition[];
        actions?: ApiAction[];
        area?: ApiArea;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiTask, RpcStatus>({
        path: `/v1/task/${id}`,
        method: "PUT",
        body: body,
        secure: true,
        type: ContentType.Json,
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
    automationServiceDisableTask: (id: string, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
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
    automationServiceEnableTask: (id: string, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
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
        /** @format uint64 */
        page?: string;
        /** @format uint64 */
        limit?: string;
        sort?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiGetTaskListResult, RpcStatus>({
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
     * @tags UserService
     * @name UserServiceAddUser
     * @summary add new user
     * @request POST:/v1/user
     * @secure
     */
    userServiceAddUser: (body: ApiNewtUserRequest, params: RequestParams = {}) =>
      this.request<ApiUserFull, RpcStatus>({
        path: `/v1/user`,
        method: "POST",
        body: body,
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
    userServiceGetUserById: (id: string, params: RequestParams = {}) =>
      this.request<ApiUserFull, RpcStatus>({
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
     * @name UserServiceDeleteUserById
     * @summary delete user by id
     * @request DELETE:/v1/user/{id}
     * @secure
     */
    userServiceDeleteUserById: (id: string, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
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
     * @name UserServiceUpdateUserById
     * @summary update user by id
     * @request PUT:/v1/user/{id}
     * @secure
     */
    userServiceUpdateUserById: (
      id: string,
      body: {
        nickname?: string;
        firstName?: string;
        lastName?: string;
        password?: string;
        passwordRepeat?: string;
        email?: string;
        status?: string;
        lang?: string;
        image?: ApiUpdateUserRequestImage;
        role?: ApiUpdateUserRequestRole;
        meta?: ApiUserMeta[];
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiUserFull, RpcStatus>({
        path: `/v1/user/${id}`,
        method: "PUT",
        body: body,
        secure: true,
        type: ContentType.Json,
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
        /** @format uint64 */
        page?: string;
        /** @format uint64 */
        limit?: string;
        sort?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiGetUserListResult, RpcStatus>({
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
    variableServiceAddVariable: (body: ApiNewVariableRequest, params: RequestParams = {}) =>
      this.request<ApiVariable, RpcStatus>({
        path: `/v1/variable`,
        method: "POST",
        body: body,
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
      this.request<ApiVariable, RpcStatus>({
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
     * @name VariableServiceDeleteVariable
     * @summary delete variable
     * @request DELETE:/v1/variable/{name}
     * @secure
     */
    variableServiceDeleteVariable: (name: string, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
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
     * @name VariableServiceUpdateVariable
     * @summary update variable
     * @request PUT:/v1/variable/{name}
     * @secure
     */
    variableServiceUpdateVariable: (
      name: string,
      body: {
        value?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiVariable, RpcStatus>({
        path: `/v1/variable/${name}`,
        method: "PUT",
        body: body,
        secure: true,
        type: ContentType.Json,
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
        /** @format uint64 */
        page?: string;
        /** @format uint64 */
        limit?: string;
        sort?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiGetVariableListResult, RpcStatus>({
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
     * @tags Zigbee2mqttService
     * @name Zigbee2MqttServiceGetBridgeList
     * @summary get bridge list
     * @request GET:/v1/zigbee2mqtt/bridge
     * @secure
     */
    zigbee2MqttServiceGetBridgeList: (
      query?: {
        /** @format uint64 */
        page?: string;
        /** @format uint64 */
        limit?: string;
        sort?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiGetBridgeListResult, RpcStatus>({
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
    zigbee2MqttServiceAddZigbee2MqttBridge: (body: ApiNewZigbee2MqttRequest, params: RequestParams = {}) =>
      this.request<ApiZigbee2Mqtt, RpcStatus>({
        path: `/v1/zigbee2mqtt/bridge`,
        method: "POST",
        body: body,
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
    zigbee2MqttServiceGetZigbee2MqttBridge: (id: string, params: RequestParams = {}) =>
      this.request<ApiZigbee2Mqtt, RpcStatus>({
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
    zigbee2MqttServiceDeleteBridgeById: (id: string, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
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
      id: string,
      body: {
        name?: string;
        login?: string;
        password?: string;
        permitJoin?: boolean;
        baseTopic?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiZigbee2Mqtt, RpcStatus>({
        path: `/v1/zigbee2mqtt/bridge/${id}/bridge`,
        method: "PUT",
        body: body,
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
      id: string,
      query?: {
        /** @format uint64 */
        page?: string;
        /** @format uint64 */
        limit?: string;
        sort?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiDeviceListResult, RpcStatus>({
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
    zigbee2MqttServiceNetworkmap: (id: string, params: RequestParams = {}) =>
      this.request<ApiNetworkmapResponse, RpcStatus>({
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
    zigbee2MqttServiceUpdateNetworkmap: (id: string, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
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
    zigbee2MqttServiceResetBridgeById: (id: string, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
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
    zigbee2MqttServiceDeviceBan: (body: ApiDeviceBanRequest, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
        path: `/v1/zigbee2mqtt/device_ban`,
        method: "POST",
        body: body,
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
    zigbee2MqttServiceDeviceRename: (body: ApiDeviceRenameRequest, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
        path: `/v1/zigbee2mqtt/device_rename`,
        method: "POST",
        body: body,
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
    zigbee2MqttServiceDeviceWhitelist: (body: ApiDeviceWhitelistRequest, params: RequestParams = {}) =>
      this.request<any, RpcStatus>({
        path: `/v1/zigbee2mqtt/device_whitelist`,
        method: "POST",
        body: body,
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
        limit?: string;
        /** @format int64 */
        offset?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<ApiSearchDeviceResult, RpcStatus>({
        path: `/v1/zigbee2mqtt/search_device`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),
  };
  ws = {
    /**
     * No description
     *
     * @tags StreamService
     * @name StreamServiceSubscribe
     * @request GET:/ws
     * @secure
     */
    streamServiceSubscribe: (
      query?: {
        id?: string;
        query?: string;
        /** @format byte */
        body?: string;
        accessToken?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        {
          result?: ApiResponse;
          error?: RpcStatus;
        },
        RpcStatus
      >({
        path: `/ws`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),
  };
}
