import {Completion, CompletionContext, snippetCompletion as snip} from '@codemirror/autocomplete';

export function Completions(context: CompletionContext) {
  const word = context.matchBefore(/\w*/);

  if (word === null || (word.from == word.to && !context.explicit)) return null;
  return {
    from: word.from,
    options: [
      // system events
      {label: 'command_enable_task', type: 'keyword'},
      {label: 'command_disable_task', type: 'keyword'},
      {label: 'command_enable_trigger', type: 'keyword'},
      {label: 'command_disable_trigger', type: 'keyword'},
      {label: 'event_call_trigger', type: 'keyword'},
      {label: 'event_call_action', type: 'keyword'},
      {label: 'command_load_entity', type: 'keyword'},
      {label: 'command_unload_entity', type: 'keyword'},
      // {label: 'true', type: 'keyword'},
      // {label: 'hello', type: 'variable', info: '(World)'},
      // {label: 'magic', type: 'text', apply: '⠁⭒*.✩.*⭒⠁', detail: 'macro'},
    ],
  };
}

export const snippets: readonly Completion[] = [
  snip("marshal(${params})", {
    label: "marshal",
    detail: "marshal value",
    type: "keyword"
  }),
  snip("unmarshal(${params})", {
    label: "unmarshal",
    detail: "unmarshal string",
    type: "keyword"
  }),
  snip("hex2arr(${params})", {
    label: "hex2arr",
    detail: "definition",
    type: "keyword"
  }),
  snip("ExecuteSync(${param1}, ${param2})", {
    label: "ExecuteSync",
    detail: "execute command",
    type: "keyword"
  }),
  snip("ExecuteAsync(${param1}, ${param2})", {
    label: "ExecuteAsync",
    detail: "execute command",
    type: "keyword"
  }),
  snip("parseFloat(${params})", {
    label: "parseFloat",
    detail: "parse string",
    type: "keyword"
  }),
  snip("Encrypt(${params})", {
    label: "Encrypt",
    detail: "encrypt string",
    type: "keyword"
  }),
  snip("Decrypt(${params})", {
    label: "Decrypt",
    detail: "decrypt string",
    type: "keyword"
  }),

  // notifr
  snip("notifr.newMessage()", {
    label: "notifr.newMessage",
    detail: "create new message",
    type: "keyword"
  }),
  snip("notifr.send(${params})", {
    label: "notifr.send",
    detail: "send message to notifr",
    type: "keyword"
  }),

  // print
  snip("print(${params})", {
    label: "print",
    detail: "print message",
    type: "keyword"
  }),

  // storage
  snip("Storage.push(${key}, ${value})", {
    label: "Storage.push",
    detail: "push value to storage",
    type: "keyword"
  }),
  snip("Storage.getByName(${key})", {
    label: "Storage.getByName",
    detail: "get value from storage",
    type: "keyword"
  }),
  snip("Storage.search(${key})", {
    label: "Storage.search",
    detail: "search in storage by key",
    type: "keyword"
  }),
  snip("Storage.pop(${key})", {
    label: "Storage.pop",
    detail: "remove value from storage by key",
    type: "keyword"
  }),

  // geo
  snip("GeoDistanceToArea(${areaId}, ${point})", {
    label: "GeoDistanceToArea",
    detail: "calculate distance between the point and the area",
    type: "keyword"
  }),
  snip("GeoPointInsideArea(${areaId}, ${point})", {
    label: "GeoPointInsideArea",
    detail: "check if point in the area",
    type: "keyword"
  }),
  snip("GeoDistanceBetweenPoints(${point1}, ${point2})", {
    label: "GeoDistanceBetweenPoints",
    detail: "calculate distance between points",
    type: "keyword"
  }),

  // http
  snip("http.get(${url})", {
    label: "http.get",
    detail: "get http request",
    type: "keyword"
  }),
  snip("http.post(${url}, ${body})", {
    label: "http.post",
    detail: "post http request",
    type: "keyword"
  }),
  snip("http.put(${url}, ${body})", {
    label: "http.put",
    detail: "put http request",
    type: "keyword"
  }),
  snip("http.delete(${url})", {
    label: "http.delete",
    detail: "delete http request",
    type: "keyword"
  }),
  snip("http.digestAuth(${username}, ${password})", {
    label: "http.digestAuth",
    detail: "set auth params",
    type: "keyword"
  }),
  snip("http.basicAuth(${username}, ${password})", {
    label: "http.basicAuth",
    detail: "set auth params",
    type: "keyword"
  }),
  snip("http.download(${url})", {
    label: "http.download",
    detail: "download request",
    type: "keyword"
  }),

  // zigbee2mqttEvent
  snip("function zigbee2mqttEvent(msg) {\n\t${}\n}", {
    label: "zigbee2mqttEvent",
    detail: "event handler",
    type: "keyword"
  }),

  // mqtt
  snip("Mqtt.publish(topic, payload, qos, retain)", {
    label: "Mqtt.publish",
    detail: "send message to mqtt server",
    type: "keyword"
  }),
  snip("function mqttEvent(ENTITY_ID, actionName) {\n\t${}\n}", {
    label: "mqttEvent",
    detail: "event handler",
    type: "keyword"
  }),

  // automation
  snip("function automationAction(entityId) {\n\t${}\n}", {
    label: "automationAction",
    detail: "event handler",
    type: "keyword"
  }),
  snip("function automationCondition(entityId) {\n\t${}\n}", {
    label: "automationCondition",
    detail: "event handler",
    type: "keyword"
  }),
  snip("function automationTriggerAlexa(msg) {\n\t${}\n}", {
    label: "automationTriggerAlexa",
    detail: "event handler",
    type: "keyword"
  }),
  snip("function automationTriggerTime(msg) {\n\t${}\n}", {
    label: "automationTriggerTime",
    detail: "event handler",
    type: "keyword"
  }),
  snip("function automationTriggerStateChanged(msg) {\n\t${}\n}", {
    label: "automationTriggerStateChanged",
    detail: "event handler",
    type: "keyword"
  }),
  snip("function automationTriggerSystem(msg) {\n\t${}\n}", {
    label: "automationTriggerSystem",
    detail: "event handler",
    type: "keyword"
  }),

  // supervisor
  snip("const entity = GetEntity(ENTITY_ID)", {
    label: "GetEntity",
    detail: "get entity",
    type: "keyword"
  }),
  snip("EntitySetState(ENTITY_ID, state)", {
    label: "EntitySetState",
    detail: "set entity state",
    type: "keyword"
  }),
  snip("EntitySetStateName(ENTITY_ID, ${params})", {
    label: "EntitySetStateName",
    detail: "set entity state",
    type: "keyword"
  }),
  snip("EntityGetState(ENTITY_ID)", {
    label: "EntityGetState",
    detail: "get entity state",
    type: "keyword"
  }),
  snip("EntitySetAttributes(ENTITY_ID, attr)", {
    label: "EntitySetAttributes",
    detail: "set attribute",
    type: "keyword"
  }),
  snip("EntitySetMetric(ENTITY_ID, name, value)", {
    label: "EntitySetMetric",
    detail: "set metric",
    type: "keyword"
  }),
  snip("EntityCallScene(ENTITY_ID, args)", {
    label: "EntityCallScene",
    detail: "set metric",
    type: "keyword"
  }),
  snip("EntityGetSettings(ENTITY_ID)", {
    label: "EntityGetSettings",
    detail: "get settings",
    type: "keyword"
  }),
  snip("EntityGetAttributes(ENTITY_ID)", {
    label: "EntityGetAttributes",
    detail: "get attributes",
    type: "keyword"
  }),

  // system events
  snip("PushSystemEvent(event, {id: 0})", {
    label: "PushSystemEvent",
    detail: "push system event",
    type: "keyword"
  }),

  // entity
  snip("function entityAction(ENTITY_ID, actionName, args) {\n\t${}\n}", {
    label: "entityAction",
    detail: "event handler",
    type: "keyword"
  }),

  // telegram
  snip("function telegramAction(ENTITY_ID, actionName) {\n\t${}\n}", {
    label: "telegramAction",
    detail: "event handler",
    type: "keyword"
  }),


  // camera
  snip("Camera.continuousMove(x, y)", {
    label: "Camera.continuousMove",
    detail: "control camera",
    type: "keyword"
  }),
  snip("Camera.stopContinuousMove(attr)", {
    label: "Camera.stopContinuousMove",
    detail: "control camera",
    type: "keyword"
  }),

  // alexa
  snip("function skillOnLaunch() {\n\t${}\n}", {
    label: "skillOnLaunch",
    detail: "event handler",
    type: "keyword"
  }),
  snip("function skillOnSessionEnd() {\n\t${}\n}", {
    label: "skillOnSessionEnd",
    detail: "event handler",
    type: "keyword"
  }),
  snip("function skillOnIntent() {\n\t${}\n}", {
    label: "skillOnIntent",
    detail: "event handler",
    type: "keyword"
  }),
  snip("Alexa.slots[\'place\']", {
    label: "Alexa",
    detail: "event handler",
    type: "keyword"
  }),
  snip("Alexa.sendMessage(\"#{place}_#{state}\")", {
    label: "Alexa.sendMessage",
    detail: "send message",
    type: "keyword"
  }),
  snip("Done(\"#{place}_#{state}\")", {
    label: "Done",
    detail: "send message",
    type: "keyword"
  }),

  // miner
  snip("Miner.stats()", {
    label: "Miner.stats",
    detail: "action",
    type: "keyword"
  }),
  snip("Miner.devs()", {
    label: "Miner.devs",
    detail: "action",
    type: "keyword"
  }),
  snip("Miner.summary()", {
    label: "Miner.summary",
    detail: "action",
    type: "keyword"
  }),
  snip("Miner.pools()", {
    label: "Miner.pools",
    detail: "action",
    type: "keyword"
  }),
  snip("Miner.addPool(url)", {
    label: "Miner.addPool",
    detail: "action",
    type: "keyword"
  }),
  snip("Miner.version()", {
    label: "Miner.version",
    detail: "action",
    type: "keyword"
  }),
  snip("Miner.enable(poolId)", {
    label: "Miner.enable",
    detail: "action",
    type: "keyword"
  }),
  snip("Miner.disable(poolId)", {
    label: "Miner.disable",
    detail: "action",
    type: "keyword"
  }),
  snip("Miner.delete(poolId)", {
    label: "Miner.delete",
    detail: "action",
    type: "keyword"
  }),
  snip("Miner.switchPool(poolId)", {
    label: "Miner.switchPool",
    detail: "action",
    type: "keyword"
  }),
  snip("Miner.restart()", {
    label: "Miner.restart",
    detail: "action",
    type: "keyword"
  }),
]
