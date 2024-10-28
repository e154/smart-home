export {};

declare global {
  // Entity ID
  const ENTITY_ID: string;

  // Enumerate attribute triggers
  enum AttributeType {
    INT = 'int',
    STRING = 'string',
    FLOAT = 'float',
    BOOL = 'bool',
    ARRAY = 'array',
    TIME = 'time',
    MAP = 'map',
    IMAGE = 'image',
    POINT = 'point',
  }

  // Data type for the attribute
  type Attribute = {
    name: string;
    type: AttributeType;
    value: any | Attribute;
  };

  // Data type for the attribute collection
  type Attributes = {
    [key: string]: Attribute;
  }

  // Interface of a point on the map
  interface Point {
    lon: bigint;
    lat: bigint;
  }

  // Data type for the area
  type Area = {
    polygon: Point[];
    created_at: Date;
    updated_at: Date;
    name: string;
    description: string;
    center: Point;
    id: number;
    zoom: number;
    resolution: number;
  };

  // Interface for the metric options element
  interface MetricOptionsItem {
    name: string;
    description: string;
    color: string;
    translate: string;
    label: string;
  }

  // Metric options interface
  interface MetricOptions {
    items: MetricOptionsItem[];
  }

  // Metric data element interface
  interface MetricDataItem {
    value: {
      [key: string]: any;
    };
    metric_id: number;
    time: Date;
  }

  // Data type for the metric
  type Metric = {
    id: number;
    name: string;
    description: string;
    options: MetricOptions;
    data: MetricDataItem[];
    type: 'line' | 'bar' | 'doughnut' | 'radar' | 'pie' | 'horizontal'
    ranges: string[];
    updated_at: Date;
    created_at: Date;
  };

  // Data type for the attribute value
  type AttributeValue = {
    [key: string]: any;
  }

  // Entity state parameters interface
  interface EntityStateParams {
    new_state?: string;
    attribute_values?: AttributeValue;
    settings_value?: AttributeValue;
    storage_save?: boolean;
  }

  // Interface for briefly representing the entity's action
  interface EntityActionShort {
    name: string;
    description: string;
    image_url?: string;
    icon?: string;
  }

  // Interface for briefly representing the state of an entity
  interface EntityStateShort {
    name: string;
    description: string;
    image_url?: string;
    icon?: string;
  }

  // Interface for the entity
  interface Entity {
    id: string;
    type: string;
    name: string;
    description: string;
    icon?: string;
    image_url?: string;
    actions: EntityActionShort[];
    states: EntityStateShort[];
    state?: EntityStateShort;
    attributes: Attributes;
    settings: Attributes;
    area?: Area;
    metrics: Metric[];
    hidden: boolean;
  }

  /////////////////////////////
  /// common
  /////////////////////////////

  /**
   * Converts an object or value to a JSON string.
   * @param {any} value - The value or object to convert.
   * @param {(key: string, value: any) => any} reviver - An optional conversion function for each key-value pair.
   * @returns {string} - JSON string.
   */
  function marshal(value: any, reviver?: (this: any, key: string, value: any) => any): string;

  /**
   * Converts a JSON string back to a JavaScript value.
   * @param {any} value - JSON string to convert.
   * @param {(key: string, value: any) => any} replacer - Optional conversion function for each key-value pair.
   * @param {string | number} space - Optional parameter to control the indentation in the resulting line.
   * @returns {any[]} - Array of values.
   */
  function unmarshal(value: any, replacer?: (this: any, key: string, value: any) => any, space?: string | number): any[];

  /**
   * Converts a string to an array of hexadecimal values.
   * @param {string} data - Input string in hex format.
   * @returns {number[]} - An array of numbers.
   */
  function hex2arr(data: string): Array<number>;

  /**
   * Response interface for command execution.
   */
  interface ExecuteResponse {
    out: string;
    err: string;
  }

  /**
   * Executes the command synchronously and returns a response.
   * @param {string} command - The command to execute.
   * @param {...any} optionalParams - Additional command parameters.
   * @returns {ExecuteResponse} - Response to command execution.
   */
  function ExecuteSync(command: string, ...optionalParams: any[]): ExecuteResponse;

  /**
   * Executes a command asynchronously and returns a response.
   * @param {string} command - The command to execute.
   * @param {...any} optionalParams - Additional command parameters.
   * @returns {ExecuteResponse} - Response to command execution.
   */
  function ExecuteAsync(command: string, ...optionalParams: any[]): ExecuteResponse;

  /**
   * Encrypts data.
   * @param {string} data - Data to be encrypted.
   * @returns {string} - Encrypted data.
   */
  function Encrypt(data: string): string;

  /**
   * Decrypts data.
   * @param {string} data - Encrypted data.
   * @returns {string} - Decrypted data.
   */
  function Decrypt(data: string): string;

  /////////////////////////////
  /// supervisor
  /////////////////////////////

  /**
   * Retrieves information about an entity by identifier.
   * @param {string} entityId - Entity identifier.
   * @returns {Entity} - Information about the entity.
   */
  function GetEntity(entityId: string): Entity;

  /**
   * Sets the state of an entity.
   * @param {string} entityId - Entity identifier.
   * @param {EntityStateParams} params - Parameters for setting the state.
   */
  function EntitySetState(entityId: string, params: EntityStateParams): void;

  /**
   * Sets the entity state name.
   * @param {string} entityId - Entity identifier.
   * @param {string} state - State name.
   */
  function EntitySetStateName(entityId: string, state: string): void;

  /**
   * Retrieves brief information about the state of the entity.
   * @param {string} entityId - Entity identifier.
   * @returns {EntityStateShort} - Brief information about the state of the entity.
   */
  function EntityGetState(entityId: string): EntityStateShort;

  /**
   * Sets the attributes of an entity.
   * @param {string} entityId - Entity identifier.
   * @param {AttributeValue} params - Attributes to set.
   */
  function EntitySetAttributes(entityId: string, params: AttributeValue): void;

  /**
   * Gets the attributes of an entity.
   * @param {string} entityId - Entity identifier.
   * @returns {AttributeValue} - Entity attributes.
   */
  function EntityGetAttributes(entityId: string): AttributeValue;

  /**
   * Retrieves entity settings.
   * @param {string} entityId - Entity identifier.
   * @returns {AttributeValue} - Entity settings.
   */
  function EntityGetSettings(entityId: string): AttributeValue;

  /**
   * Sets the entity metric.
   * @param {string} entityId - Entity identifier.
   * @param {{ [key: string]: any }} params - Metric parameters.
   */
  function EntitySetMetric(entityId: string, params: { [key: string]: any }): void;

  /**
   * Causes the action of the entity.
   * @param {string} entityId - Entity identifier.
   * @param {string} action - The name of the action.
   * @param {{ [key: string]: any }} params - Parameters for the action.
   */
  function EntityCallAction(entityId: string, action: string, params: { [key: string]: any }): void;

  /**
   * Calls a function or method of an entity's script.
   * @param {string} entityId - Entity identifier.
   * @param {string} fn - The name of the function.
   * @param {any} payload - Parameters for the action.
   */
  function EntityCallScript(entityId: string, fn: string, payload: any): void;

  /**
   * Interface representing parameters for calling an action on an entity
   */
  interface CallAction {
    entity_id?: string;
    action_name: string;
    tags?: string[];
    area_id?: number;
  }

  /**
   * Function to call an action on an entity (version 2).
   * @param {CallAction} params1 - Parameters for calling the action.
   * @param {Object} params2 - Additional parameters.
   */
  function EntitiesCallAction(params1: CallAction, params2?: { [key: string]: any }): void;

  /**
   * Calls the entity script.
   * @param {string} entityId - Entity identifier.
   * @param {{ [key: string]: any }} params - Parameters for the script.
   */
  function EntityCallScene(entityId: string, params: { [key: string]: any }): void;

  /**
   * Function to perform an entity action.
   * @param {string} entityId - Entity identifier.
   * @param {string} actionName - Action name.
   * @param {{ [key: string]: any }} params - Parameters for performing the action.
   */
  function entityAction(entityId: string, actionName: string, params: { [key: string]: any }): void;


  /////////////////////////////
  /// geo
  /////////////////////////////
  /**
   * Calculates the distance from a point to the center of an area.
   * @param {number} areaId - Area ID.
   * @param {Point} point - Point to calculate the distance.
   * @returns {number} - Distance from the point to the center of the area.
   */
  function GeoDistanceToArea(areaId: number, point: Point): number;

  /**
   * Calculates the distance between two points.
   * @param {Point} point1 - The first point.
   * @param {Point} point2 - Second point.
   * @returns {number} - Distance between two points.
   */
  function GeoDistanceBetweenPoints(point1: Point, point2: Point): number;

  /**
   * Checks if a point is inside an area.
   * @param {number} areaId - Area ID.
   * @param {Point} point - Point to check.
   * @returns {boolean} - true if the point is inside the region, false otherwise.
   */
  function GeoPointInsideArea(areaId: number, point: Point): boolean;

  /////////////////////////////
  /// http
  /////////////////////////////
  /**
   * Interface for responding to an HTTP request.
   */
  interface HttpResponse {
    body: string;
    error: boolean;
    errorMessage: string;
  }

  /**
   * HTTP client interface with basic request methods.
   */
  interface http {
    /**
     * Performs an HTTP GET request.
     * @param {string} url - URL for the GET request.
     * @returns {HttpResponse} - Response to a GET request.
     */
    get(url: string): HttpResponse;

    /**
     * Performs an HTTP POST request.
     * @param {string} url - URL for the POST request.
     * @param {string} data - Data for the POST request.
     * @returns {HttpResponse} - Response to a POST request.
     */
    post(url: string, data: string): HttpResponse;

    /**
     * Performs an HTTP PUT request.
     * @param {string} url - URL for the PUT request.
     * @param {string} data - Data for the PUT request.
     * @returns {HttpResponse} - Response to the PUT request.
     */
    put(url: string, data: string): HttpResponse;

    /**
     * Performs an HTTP DELETE request.
     * @param {string} url - URL for the DELETE request.
     * @returns {HttpResponse} - Response to a DELETE request.
     */
    delete(url: string): HttpResponse;

    /**
     * Sets headers for an HTTP request.
     * @param {{ [key: string]: any }} params - Request headers.
     * @returns {http} - HTTP client with set headers.
     */
    headers(params: { [key: string]: any }): http;

    /**
     * Sets the basic authorization for the HTTP request.
     * @param {string} username - Username.
     * @param {string} password - User password.
     * @returns {http} - HTTP client with basic authorization set.
     */
    basicAuth(username: string, password: string): http;

    /**
     * Sets Digest protocol authentication for an HTTP request.
     * @param {string} username - Username.
     * @param {string} password - User password.
     * @returns {http} - HTTP client with Digest protocol authentication installed.
     */
    digestAuth(username: string, password: string): http;

    /**
     * Performs an HTTP request to download a file.
     * @param {string} post - URL for downloading the file.
     * @returns {HttpResponse} - Response to a file download request.
     */
    download(post: string): HttpResponse;
  }

  /////////////////////////////
  /// notifr
  /////////////////////////////
  /**
   * Message interface for HTML5 notifications.
   */
  interface MessageHtml5 {
    userIDS?: string;
    title: string;
    body: string;
  }

  /**
   * Message interface for Telegram.
   */
  interface MessageTelegram {
    chat_id?: number;
    file_uri?: string;
    file_path?: string;
    photo_uri?: string;
    photo_path?: string;
    keys?: string[];
    body: string;
  }

  /**
   * Message interface for SMS notifications.
   */
  interface MessageSMS {
    phone: string;
    body: string;
  }

  /**
   * Message interface for Web Push notifications.
   */
  interface MessageWebpush {
    userIDS?: string;
    title: string;
    body: string;
  }

  /**
   * Message type, which determines its format.
   */
  enum MessageType {
    HTML5_NOTIFY = 'html5_notify',
    TELEGRAM = 'telegram',
    WEBPUSH = 'webpush'
  }

  /**
   * Common message interface for different triggers of notifications.
   */
  interface Message {
    type?: MessageType;
    entity_id?: string;
    attributes: MessageHtml5 | MessageTelegram | MessageSMS | MessageWebpush;
  }

  /**
   * Interface for sending notifications.
   */
  interface notifr {
    /**
     * Creates a new message.
     * @returns {Message} - New message.
     */
    newMessage(): Message;

    /**
     * Sends a notification.
     * @param {Message} msg - Message to send.
     */
    send(msg: Message): void;
  }

  /**
   * Function for creating a message template.
   * @param {string} templateName - The name of the template.
   * @param {{ [key: string]: any }} params - Parameters for filling out the template.
   * @returns {any} - The result of applying the template.
   */
  function template(templateName: string, params: { [key: string]: any }): any;

  /////////////////////////////
  /// logging
  /////////////////////////////
  /**
   * Console output interface.
   */
  interface Console {
    /**
     * Prints a message to the console with the level "log".
     * @param {any} message - Message to display.
     * @param {...any} optionalParams - Additional parameters for output.
     */
    log(message?: any, ...optionalParams: any[]): void;

    /**
     * Displays an informational message to the console.
     * @param {any} message - Message to display.
     * @param {...any} optionalParams - Additional parameters for output.
     */
    info(message?: any, ...optionalParams: any[]): void;

    /**
     * Prints a debug message to the console.
     * @param {any} message - Message to display.
     * @param {...any} optionalParams - Additional parameters for output.
     */
    debug(message?: any, ...optionalParams: any[]): void;

    /**
     * Displays a warning to the console.
     * @param {any} message - Message to display.
     * @param {...any} optionalParams - Additional parameters for output.
     */
    warn(message?: any, ...optionalParams: any[]): void;

    /**
     * Displays an error message to the console.
     * @param {any} message - Message to display.
     * @param {...any} optionalParams - Additional parameters for output.
     */
    error(message?: any, ...optionalParams: any[]): void;
  }


  /**
   * Object for console output.
   */
  let console: Console;

  /**
   * Function to output a message to the console.
   * @param {any} message - Message to display.
   * @param {...any} optionalParams - Additional parameters for output.
   */
  function print(message?: any, ...optionalParams: any[]): void;

  /////////////////////////////
  /// storage
  /////////////////////////////
  /**
   * Interface for storage management.
   */
  interface Storage {
    /**
     * Adds a value to the store using the specified key.
     * @param {string} key - The key to add the value.
     * @param {string} value - The value to add.
     * @returns {string} - Added value.
     */
    push(key: string, value: string): string;

    /**
     * Retrieves a value from storage for the specified key.
     * @param {string} key - The key to get the value.
     * @returns {string} - The returned value.
     */
    getByName(key: string): string;

    /**
     * Searches the storage for values using the specified key.
     * @param {string} key - Key to search for values.
     * @returns {{ [key: string]: string }} - Found values as an object.
     */
    search(key: string): { [key: string]: string };

    /**
     * Removes and returns a value from storage for the specified key.
     * @param {string} key - The key to delete the value.
     * @returns {string} - The removed value.
     */
    pop(key: string): string;
  }

  /////////////////////////////
  /// system events
  /////////////////////////////
  /**
   * System event command type.
   */
  type SystemEventCommand =
    | 'command_enable_task'
    | 'command_disable_task'
    | 'command_enable_trigger'
    | 'command_disable_trigger'
    | 'event_call_trigger'
    | 'event_call_action'
    | 'command_load_entity'
    | 'command_unload_entity';

  /**
   * Function to send a system event.
   * @param {SystemEventCommand} command - System event command.
   * @param {{ [key: string]: any }} params - Parameters for the system event.
   */
  function PushSystemEvent(command: SystemEventCommand, params: { [key: string]: any }): void;

  /**
   * Function to perform an action in Telegram.
   * @param {string} entityId - Entity identifier.
   * @param {string} actionName - Action name.
   * @param {{ [key: string]: any }} params - Parameters for performing the action.
   */
  function telegramAction(entityId: string, actionName: string, params: { [key: string]: any }): void;

  /**
   * Camera interface with control methods.
   */
  interface Camera {
    /**
     * Starts continuous camera movement along specified coordinates.
     * @param {number} x - X coordinate.
     * @param {number} y - Y coordinate.
     */
    continuousMove(x: number, y: number): void;

    /**
     * Stops continuous camera movement.
     */
    stopContinuousMove(): void;
  }

  /**
   * Interface for working with the Alexa skill.
   */
  interface Alexa {
    slots: { [key: string]: string };

    /**
     * Sets the output speech for the Alexa skill.
     * @param {string} text - Speech text.
     * @returns {Alexa} - The Alexa object for the call chain.
     */
    outputSpeech(text: string): Alexa;

    /**
     * Adds a card for the Alexa skill.
     * @param {string} title - Title of the card.
     * @param {string} content - Contents of the card.
     * @returns {Alexa} - The Alexa object for the call chain.
     */
    card(title: string, content: string): Alexa;

    /**
     * Sets the session end flag for the Alexa skill.
     * @param {boolean} flag - Session end flag.
     * @returns {Alexa} - The Alexa object for the call chain.
     */
    endSession(flag: boolean): Alexa;

    /**
     * Returns the current session status of the Alexa skill.
     * @returns {string} - Session status.
     */
    session(): string;

    /**
     * Sends a message from the Alexa skill.
     * @param {any} msg - Message to send.
     */
    sendMessage(msg: any): void;
  }

  /**
   * Function called when the Alexa skill is launched.
   */
  function skillOnLaunch(): void;

  /**
   * Function called when the Alexa skill session ends.
   */
  function skillOnSessionEnd(): void;

  /**
   * Function called when an Alexa skill intent is received.
   */
  function skillOnIntent(): void;

  /////////////////////////////
  /// miner
  /////////////////////////////
  /**
   * Interface for interaction with the miner.
   */
  interface MinerResult {
    error: boolean;
    errMessage: string;
    result: string;
  }

  /**
   * Miner interface with control methods.
   */
  interface Miner {
    /**
     * Gets miner statistics.
     * @returns {MinerResult} - The result of the statistics request.
     */
    stats(): MinerResult;

    /**
     * Receives information about miner devices.
     * @returns {MinerResult} - The result of executing a request for information about devices.
     */
    devs(): MinerResult;

    /**
     * Retrieves summary information about the miner.
     * @returns {MinerResult} - The result of the summary information request.
     */
    summary(): MinerResult;

    /**
     * Receives information about the pools to which the miner is connected.
     * @returns {MinerResult} - The result of executing a request for information about pools.
     */
    pools(): MinerResult;

    /**
     * Adds a new pool to the miner.
     * @param {string} url - URL of the new pool.
     * @returns {MinerResult} - The result of the add pool operation.
     */
    addPool(url: string): MinerResult;

    /**
     * Gets the miner version.
     * @returns {MinerResult} - The result of the miner version request.
     */
    version(): MinerResult;

    /**
     * Enables the pool with the specified ID.
     * @param {number} poolId - Pool ID to enable.
     * @returns {MinerResult} - The result of the pool enable operation.
     */
    enable(poolId: number): MinerResult;

    /**
     * Disables the pool with the specified ID.
     * @param {number} poolId - The pool ID to disable.
     * @returns {MinerResult} - The result of the pool shutdown operation.
     */
    disable(poolId: number): MinerResult;

    /**
     * Deletes the pool with the specified ID.
     * @param {number} poolId - ID of the pool to delete.
     * @returns {MinerResult} - The result of the pool deletion operation.
     */
    delete(poolId: number): MinerResult;

    /**
     * Switches to the pool with the specified ID.
     * @param {number} poolId - ID of the pool to switch.
     * @returns {MinerResult} - The result of the switch to the pool operation.
     */
    switchPool(poolId: number): MinerResult;

    /**
     * Restarts the miner.
     * @returns {MinerResult} - The result of the miner restart operation.
     */
    restart(): MinerResult;
  }


  /////////////////////////////
  /// zigbee2mqttEvent
  /////////////////////////////
  /**
   * Interface for Zigbee2mqtt events.
   */
  interface Z2MMessage {
    payload: string;
    topic: string;
    qos: number;
    duplicate: boolean;
    error: string;
    success: boolean;
    new_state: EntityStateParams;

    /**
     * Checks whether the Zigbee2mqtt event was successful.
     * @returns {boolean} - The result of the success check.
     */
    ok(): boolean;

    /**
     * Clears Zigbee2mqtt event data.
     * @returns {boolean} - Result of data clearing.
     */
    clear(): boolean;

    /**
     * Creates a copy of the Zigbee2mqtt event.
     * @returns {Z2MMessage} - A copy of the Zigbee2mqtt event.
     */
    copy(): Z2MMessage;

    /**
     * Gets the value of a variable by key from the Zigbee2mqtt event.
     * @param {string} key - The key of the variable.
     * @returns {string} - The value of the variable.
     */
    getVar(key: string): string;

    /**
     * Sets the value of a variable by key in the Zigbee2mqtt event.
     * @param {string} key - The key of the variable.
     * @param {string} value - The value of the variable.
     */
    setVar(key: string, value: string): void;

    /**
     * Updates Zigbee2mqtt event data based on another event.
     * @param {Z2MMessage} msg - Zigbee2mqtt event to update.
     */
    Update(msg: Z2MMessage): void;
  }

  /**
   * Function called on Zigbee2mqtt event.
   * @param {Z2MMessage} message - Zigbee2mqtt event.
   */
  function zigbee2mqttEvent(message: Z2MMessage): void;


  /////////////////////////////
  /// mqtt
  /////////////////////////////
  /**
   * Interface for interacting with MQTT.
   */
  interface Mqtt {
    /**
     * Publishes a message to MQTT.
     * @param {string} topic - The subject of the message.
     * @param {string} payload - Message body.
     * @param {number} qos - Message service level (0, 1, 2).
     * @param {boolean} retain - Message retention flag.
     */
    publish(topic: string, payload: string, qos: number, retain: boolean): void;
  }

  /**
   * Function called on MQTT event.
   * @param {string} entityId - The ID of the entity associated with the MQTT event.
   * @param {string} actionName - The name of the action associated with the MQTT event.
   */
  function mqttEvent(entityId: string, actionName: string): void;

  /////////////////////////////
  /// automation
  /////////////////////////////
  interface Action {
  }

  interface Condition {
  }

  interface Trigger {
  }

  /**
   * Function called when the automation condition is met.
   * @param {string} entityId - The ID of the entity associated with the automation condition.
   * @returns {boolean} - The result of the condition being met.
   */
  function automationCondition(entityId: string): boolean;

  /**
   * A function called when an automation action is executed.
   * @param {string} entityId - The ID of the entity associated with the automation action.
   */
  function automationAction(entityId: string): void;

  /**
   * Interface for Alexa Time trigger messages.
   */
  interface TriggerAlexaTimeMessage {
    payload: any;
    trigger_name: string;
    entity_id: string;
  }

  /**
   * Function called when an Alexa Time trigger event occurs.
   * @param {TriggerAlexaTimeMessage} msg - Alexa Time trigger message.
   * @returns {boolean} - The result of the trigger execution.
   */
  function automationTriggerAlexa(msg: TriggerAlexaTimeMessage): boolean;

  /**
   * Interface for Time trigger messages.
   */
  interface TriggerTimeMessage {
    payload: Date;
    trigger_name: string;
    entity_id: string;
  }

  /**
   * Function called when the Time trigger event occurs.
   * @param {TriggerTimeMessage} msg - Time trigger message.
   * @returns {boolean} - The result of the trigger execution.
   */
  function automationTriggerTime(msg: TriggerTimeMessage): boolean;

  /**
   * Interface for the entity state change event.
   */
  interface EventEntityState {
    entity_id: string;
    value: any;
    state: EntityStateShort;
    attributes: AttributeValue;
    settings: AttributeValue;
    last_changed: Date;
    last_updated: Date;
  }

  /**
   * Interface for trigger state change event.
   */
  interface EventStateChanged {
    storage_save: boolean;
    do_not_save_metric: boolean;
    plugin_name: string;
    entity_id: string;
    old_state?: EventEntityState;
    new_state: EventEntityState;
  }

  /**
   * Interface for trigger state change event.
   */
  interface TriggerStateChangedMessage {
    payload?: EventStateChanged;
    trigger_name: string;
    entity_id: string;
  }

  /**
   * Function called on a state change trigger event.
   * @param {TriggerStateChangedMessage} msg - State change trigger message.
   * @returns {boolean} - The result of the trigger execution.
   */
  function automationTriggerStateChanged(msg: TriggerStateChangedMessage): boolean;

  /**
   * Interface for system trigger message.
   */
  interface SystemTriggerMessage {
    topic: string;
    event_name: string;
    event: any;
  }

  /**
   * Interface for system trigger message.
   */
  interface TriggerSystemMessage {
    payload: SystemTriggerMessage;
    trigger_name: string;
    entity_id: string;
  }

  /**
   * Function called when a system trigger event occurs.
   * @param {TriggerSystemMessage} msg - System trigger message.
   * @returns {boolean} - The result of the trigger execution.
   */
  function automationTriggerSystem(msg: TriggerSystemMessage): boolean;

  /**
   * Interface for Ble trigger messages.
   */
  interface TriggerBleMessage {
    payload: Uint8Array;
    trigger_name: string;
    entity_id: string;
  }

  /**
   * Function called when the Ble trigger event occurs.
   * @param {TriggerBleMessage} msg - Ble trigger message.
   * @returns {boolean} - The result of the trigger execution.
   */
  function automationTriggerBle(msg: TriggerBleMessage): boolean;

  /**
   * Interface for Stt trigger messages.
   */
  interface TriggerSttMessage {
    payload: string;
    trigger_name: string;
    entity_id: string;
  }

  /**
   * Function called when the "Speech to text" event occurs.
   * @param {TriggerSttMessage} msg - stt trigger message.
   * @returns {boolean} - The result of the trigger execution.
   */
  function automationTriggerStt(msg: TriggerSttMessage): boolean;

  /**
   * Interface for responding to bluetooth commands.
   */
  interface BleResponse {
    response?: Uint8Array;
    error?: string;
  }

  /**
   * Write replaces the characteristic value with a new value. The call will return after all data has been written.
   * @param {string} char - UUID Device characteristic.
   * @param {Uint8Array} payload - Command.
   * @param {boolean} withResponse - if a response is expected.
   */
  function BleWrite(char: string, payload: Uint8Array, withResponse: boolean): BleResponse;

  /**
   * Read reads the current characteristic value.
   * @param {string} char - UUID Device characteristic.
   */
  function BleRead(char: string): BleResponse;

}
