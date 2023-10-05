
export const HintDictionaryCoffee = {
    words: [
        // common
        {text: 'main', displayText: 'main'},
        {text: 'unmarshal', displayText: 'unmarshal'},
        {text: 'marshal', displayText: 'marshal'},
        {text: 'hex2arr(hexString)', displayText: 'hex2arr'},
        {text: 'ExecuteSync(file, args)', displayText: 'ExecuteSync'},
        {text: 'ExecuteAsync(file, args)', displayText: 'ExecuteAsync'},
        {text: 'for v, i in items', displayText: 'for'},
        {text: 'parseFloat', displayText: 'parseFloat'},
        {text: 'indexOf', displayText: 'indexOf'},
        {text: 'substring', displayText: 'substring'},
        {text: 'Encrypt(val)', displayText: 'Encrypt'},
        {text: 'Decrypt(val)', displayText: 'Decrypt'},

        // notifr
        {text: 'notifr.newMessage()', displayText: 'notifr.newMessage'},
        {text: 'notifr.send(msg)', displayText: 'notifr.send'},
        {text: 'template.render(name, params)', displayText: 'template.render'},

        // logging
        {text: 'print', displayText: 'print'},
        {text: 'Log.info', displayText: 'Log.info'},
        {text: 'Log.debug', displayText: 'Log.debug'},
        {text: 'Log.warn', displayText: 'Log.warn'},
        {text: 'Log.error', displayText: 'Log.error'},

        // storage
        {text: 'Storage.push(key, value)', displayText: 'Storage.push'},
        {text: 'Storage.getByName(key)', displayText: 'Storage.getByName'},
        {text: 'Storage.search(key)', displayText: 'Storage.search'},
        {text: 'Storage.pop(key)', displayText: 'Storage.pop'},
        {text: 'push(key, value)', displayText: 'push'},
        {text: 'getByName(key)', displayText: 'getByName'},
        {text: 'search(key)', displayText: 'search'},

        // geo
        {text: 'GeoDistanceToArea(areaId, point)', displayText: 'GeoDistanceToArea'},
        {text: 'GeoPointInsideAria(areaId, point)', displayText: 'GeoPointInsideAria'},
        {text: 'GeoDistanceBetweenPoints(point1, point2)', displayText: 'GeoDistanceBetweenPoints'},

        // http
        {text: 'http.get(url)', displayText: 'http.get'},
        {text: 'http.post(url, body)', displayText: 'http.post'},
        {text: 'http.put(url, body)', displayText: 'http.put'},
        {text: 'http.delete(url)', displayText: 'http.delete'},

        // zigbee2mqttEvent
        {text: 'zigbee2mqttEvent = (message) ->', displayText: 'zigbee2mqttEvent'},

        // mqtt
        {text: 'Mqtt.publish(topic, payload, qos, retain)', displayText: 'Mqtt.publish'},
        {text: 'mqttEvent = (ENTITY_ID, actionName) ->', displayText: 'mqttEvent'},
        {text: 'message', displayText: 'message'},
        {text: 'message.payload', displayText: 'message.payload'},
        {text: 'message.topic', displayText: 'message.topic'},
        {text: 'message.qos', displayText: 'message.qos'},
        {text: 'message.duplicate', displayText: 'message.duplicate'},
        {text: 'message.storage', displayText: 'message.storage'},
        {text: 'message.error', displayText: 'message.error'},
        {text: 'message.success', displayText: 'message.success'},
        {text: 'message.new_state', displayText: 'message.new_state'},

        // automation
        {text: 'Action', displayText: 'Action'},
        {text: 'Condition', displayText: 'Condition'},
        {text: 'Trigger', displayText: 'Trigger'},
        {text: 'automationAction = (entityId)->', displayText: 'automationAction'},
        {text: 'automationCondition = (entityId)->', displayText: 'automationCondition'},
        {text: 'automationTriggerAlexa = (msg) ->', displayText: 'automationTriggerAlexa'},
        {text: 'automationTriggerTime = (msg) ->', displayText: 'automationTriggerTime'},
        {text: 'automationTriggerStateChanged = (msg)->', displayText: 'automationTriggerStateChanged'},
        {text: 'automationTriggerSystem = (msg)->', displayText: 'automationTriggerSystem'},

        // entity manager
        {text: 'GetEntity(ENTITY_ID)', displayText: 'GetEntity'},
        {text: 'EntitySetState(ENTITY_ID, state)', displayText: 'EntitySetState'},
        {text: 'EntitySetStateName(ENTITY_ID, name)', displayText: 'EntitySetStateName'},
        {text: 'EntityGetState(ENTITY_ID)', displayText: 'EntityGetState'},
        {text: 'EntitySetAttributes(ENTITY_ID, attr)', displayText: 'EntitySetAttributes'},
        {text: 'EntitySetMetric(ENTITY_ID, name, value)', displayText: 'EntitySetMetric'},
        {text: 'EntityCallAction(ENTITY_ID, action, args)', displayText: 'EntityCallAction'},
        {text: 'EntityCallScene(ENTITY_ID, args)', displayText: 'EntityCallScene'},
        {text: 'EntityGetSettings(ENTITY_ID)', displayText: 'EntityGetSettings'},
        {text: 'EntityGetAttributes(ENTITY_ID)', displayText: 'EntityGetAttributes'},

        // entity
        {text: 'entityAction = (ENTITY_ID, actionName, args)->', displayText: 'entityAction'},

        // telegram
        {text: 'telegramAction = (ENTITY_ID, actionName)->', displayText: 'telegramAction'},

        // camera
        {text: 'Camera.continuousMove(x, y)', displayText: 'Camera.continuousMove'},
        {text: 'Camera.stopContinuousMove(attr)', displayText: 'Camera.stopContinuousMove'},

        // alexa
        {text: 'skillOnLaunch = ()->', displayText: 'skillOnLaunch'},
        {text: 'skillOnSessionEnd = ()->', displayText: 'skillOnSessionEnd'},
        {text: 'skillOnIntent = ()->', displayText: 'skillOnIntent'},
        {text: 'Alexa.slots[\'place\']', displayText: 'Alexa.slots'},
        {text: 'Alexa.sendMessage("#{place}_#{state}")', displayText: 'Alexa.sendMessage'},
        {text: 'Done("#{place}_#{state}")', displayText: 'Done'},

        // miner
        {text: 'Miner.stats()', displayText: 'Miner.stats'},
        {text: 'Miner.devs()', displayText: 'Miner.devs'},
        {text: 'Miner.summary()', displayText: 'Miner.summary'},
        {text: 'Miner.pools()', displayText: 'Miner.pools'},
        {text: 'Miner.addPool(url)', displayText: 'Miner.addPool'},
        {text: 'Miner.version()', displayText: 'Miner.version'},
        {text: 'Miner.enable(poolId)', displayText: 'Miner.enable'},
        {text: 'Miner.disable(poolId)', displayText: 'Miner.disable'},
        {text: 'Miner.delete(poolId)', displayText: 'Miner.delete'},
        {text: 'Miner.switchPool(poolId)', displayText: 'Miner.switchPool'},
        {text: 'Miner.restart()', displayText: 'Miner.restart'},

    ]
};
