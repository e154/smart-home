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

  /**
   * Data type for the attribute.
   */
  interface Attribute {
    /**
     * The name of the attribute.
     */
    name: string;

    /**
     * The type of the attribute (e.g., string, number, object).
     */
    type: AttributeType;

    /**
     * The value of the attribute, which can be any primitive type or another attribute.
     */
    value: any | Attribute;
  }

  // Data type for the attribute collection
  interface Attributes {
    [key: string]: Attribute;
  }

  /**
   * Interface for a point on the map.
   */
  interface Point {
    /**
     * The longitude of the point, represented as a big integer.
     */
    lon: number;

    /**
     * The latitude of the point, represented as a big integer.
     */
    lat: number;
  }

  /**
   * Data type for the area.
   */
  interface Area {
    /**
     * The polygon defining the boundary of the area, represented as an array of points.
     */
    polygon: Point[];

    /**
     * The date and time when the area was created.
     */
    created_at: Date;

    /**
     * The date and time when the area was last updated.
     */
    updated_at: Date;

    /**
     * The name of the area.
     */
    name: string;

    /**
     * A brief description of the area.
     */
    description: string;

    /**
     * The center point of the area, represented as an object with `lat` and `lon` properties.
     */
    center: Point;

    /**
     * The unique identifier of the area.
     */
    id: number;

    /**
     * The zoom level associated with the area.
     */
    zoom: number;

    /**
     * The resolution associated with the area.
     */
    resolution: number;
  }

  /**
   * Interface for the metric options element.
   */
  interface MetricOptionsItem {
    /**
     * The name of the option.
     */
    name: string;

    /**
     * A brief description of the option.
     */
    description: string;

    /**
     * The color associated with the option.
     */
    color: string;

    /**
     * A translated value for the option.
     */
    translate: string;

    /**
     * A label for the option.
     */
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

  /**
   * Data type for the metric.
   */
  interface Metric {
    /**
     * The unique identifier of the metric.
     */
    id: number;

    /**
     * The name of the metric.
     */
    name: string;

    /**
     * A brief description of the metric.
     */
    description: string;

    /**
     * Additional options related to the metric.
     */
    options: MetricOptions;

    /**
     * An array of data items for the metric.
     */
    data: MetricDataItem[];

    /**
     * The type of the metric chart ('line', 'bar', 'doughnut', 'radar', 'pie', or 'horizontal').
     */
    type: 'line' | 'bar' | 'doughnut' | 'radar' | 'pie' | 'horizontal';

    /**
     * An array of ranges for the metric.
     */
    ranges: string[];

    /**
     * The date and time when the metric was last updated.
     */
    updated_at: Date;

    /**
     * The date and time when the metric was created.
     */
    created_at: Date;
  }

  // Data type for the attribute value
  interface AttributeValue {
    [key: string]: any;
  }

  /**
   * Interface for the parameters related to an entity's state.
   */
  interface EntityStateParams {
    /**
     * The new state of the entity. This is optional.
     */
    new_state?: string;

    /**
     * Attribute values associated with the entity's state. This is optional.
     */
    attribute_values?: AttributeValue;

    /**
     * The settings value associated with the entity's state. This is optional.
     */
    settings_value?: AttributeValue;

    /**
     * A boolean indicating whether to save the state in storage. This is optional.
     */
    storage_save?: boolean;
  }

  /**
   * Interface for briefly representing the entity's action.
   */
  interface EntityActionShort {
    /**
     * The name of the action.
     */
    name: string;

    /**
     * A brief description of the action.
     */
    description: string;

    /**
     * An optional URL to an image representing the action.
     */
    image_url?: string;

    /**
     * An optional icon representing the action.
     */
    icon?: string;
  }

  /**
   * Interface for briefly representing the state of an entity.
   */
  interface EntityStateShort {
    /**
     * The name of the state.
     */
    name: string;

    /**
     * A brief description of the state.
     */
    description: string;

    /**
     * An optional image URL for the state.
     */
    image_url?: string;

    /**
     * An optional icon URL for the state.
     */
    icon?: string;
  }

  /**
   * Interface for the entity.
   */
  interface Entity {
    /**
     * The unique identifier of the entity.
     */
    id: string;

    /**
     * The type of the entity.
     */
    type: string;

    /**
     * The name of the entity.
     */
    name: string;

    /**
     * A brief description of the entity.
     */
    description: string;

    /**
     * An optional icon URL for the entity.
     */
    icon?: string;

    /**
     * An optional image URL for the entity.
     */
    image_url?: string;

    /**
     * A list of short actions associated with the entity.
     */
    actions: EntityActionShort[];

    /**
     * A list of short states associated with the entity.
     */
    states: EntityStateShort[];

    /**
     * The current state of the entity. This is optional.
     */
    state?: EntityStateShort;

    /**
     * Additional attributes related to the entity.
     */
    attributes: Attributes;

    /**
     * Additional settings related to the entity.
     */
    settings: Attributes;

    /**
     * An optional area associated with the entity.
     */
    area?: Area;

    /**
     * A list of metrics associated with the entity.
     */
    metrics: Metric[];

    /**
     * A boolean indicating whether the entity is hidden.
     */
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
   *
   * @example
   * ```ts
   * Storage.push(settingsKey, marshal(settings))
   * ```
   */
  function marshal(value: any, reviver?: (this: any, key: string, value: any) => any): string;

  /**
   * Converts a JSON string back to a JavaScript value.
   * @param {any} value - JSON string to convert.
   * @param {(key: string, value: any) => any} replacer - Optional conversion function for each key-value pair.
   * @param {string | number} space - Optional parameter to control the indentation in the resulting line.
   * @returns {any[]} - Array of values.
   *
   * @example
   * ```ts
   * const _settings = Storage.getByName(settingsKey)
   * const _key = this.name + '.' + level
   * if (_settings) {
   *   const settings: Settings = unmarshal(_settings);
   *   if (!settings.notify[_key]) {
   *     //...
   *   }
   * }
   * ```
   */
  function unmarshal(value: any, replacer?: (this: any, key: string, value: any) => any, space?: string | number): any[];

  /**
   * Converts a string to an array of hexadecimal values.
   * @param {string} data - Input string in hex format.
   * @returns {number[]} - An array of numbers.
   *
   * @example
   * ```ts
   * var hexString = "48656C6C6F20576F726C64"; // Hexadecimal representation of "Hello World"
   * var byteArr = hex2arr(hexString);
   * console.log(byteArr); // [ 72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100 ]
   * ```
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
   *
   * @example
   * ```ts
   * r = ExecuteSync("data/scripts/ping.sh", "google.com")
   * if (r.out == 'ok') {
   *   print "site is available ^^"
   * }
   * ```
   */
  function ExecuteSync(command: string, ...optionalParams: any[]): ExecuteResponse;

  /**
   * Executes a command asynchronously and returns a response.
   * @param {string} command - The command to execute.
   * @param {...any} optionalParams - Additional command parameters.
   * @returns {ExecuteResponse} - Response to command execution.
   *
   * @example
   * ```ts
   * const file = 'script.js';
   * const args = { param1: 'value1', param2: 'value2' };
   *
   * ExecuteAsync(file, args);
   * ```
   */
  function ExecuteAsync(command: string, ...optionalParams: any[]): ExecuteResponse;

  /**
   * Encrypts data.
   * @param {string} data - Data to be encrypted.
   * @returns {string} - Encrypted data.
   *
   * @example
   * ```ts
   * const originalData = "Secret Information"; // Original data
   * const encryptedData = Encrypt(originalData); // Encryption
   *
   * console.log("Original Data:", originalData);
   * console.log("Encrypted Data:", encryptedData);
   * ```
   */
  function Encrypt(data: string): string;

  /**
   * Decrypts data.
   * @param {string} data - Encrypted data.
   * @returns {string} - Decrypted data.
   *
   * @example
   * ```ts
   * // Data Decryption
   * const decryptedData = Decrypt(encryptedData); // Decryption
   *
   * console.log("Decrypted Data:", decryptedData);
   * ```
   */
  function Decrypt(data: string): string;

  /////////////////////////////
  /// supervisor
  /////////////////////////////

  /**
   * Retrieves information about an entity by identifier.
   * @param {string} entityId - Entity identifier.
   * @returns {Entity} - Information about the entity.
   *
   * @example
   * ```ts
   * const entity = GetEntity(ENTITY_ID);
   * ```
   */
  function GetEntity(entityId: string): Entity;

  /**
   * Sets the state of an entity.
   * @param {string} entityId - Entity identifier.
   * @param {EntityStateParams} params - Parameters for setting the state.
   *
   * @example
   * ```ts
   * const attrs = {
   *   foo: bar,
   * }
   * const stateName = 'connected'
   * EntitySetState(ENTITY_ID, {
   *   new_state: stateName,
   *   attribute_values: attrs,
   *   storage_save: true
   * });
   * ```
   */
  function EntitySetState(entityId: string, params: EntityStateParams): void;

  /**
   * Sets the entity state name.
   * @param {string} entityId - Entity identifier.
   * @param {string} state - State name.
   *
   * @example
   * ```ts
   * const stateName = 'connected'
   * EntitySetStateName(ENTITY_ID, stateName);
   * ```
   */
  function EntitySetStateName(entityId: string, state: string): void;

  /**
   * Retrieves brief information about the state of the entity.
   * @param {string} entityId - Entity identifier.
   * @returns {EntityStateShort} - Brief information about the state of the entity.
   *
   * @example
   * ```ts
   * const currentState = EntityGetState(ENTITY_ID);
   * print(marshal(homeState))
   * ```
   */
  function EntityGetState(entityId: string): EntityStateShort;

  /**
   * Sets the attributes of an entity.
   * @param {string} entityId - Entity identifier.
   * @param {AttributeValue} params - Attributes to set.
   *
   * @example
   * ```ts
   * const attrs = {
   *     foo: bar,
   * }
   * EntitySetAttributes(ENTITY_ID, attrs);
   * ```
   */
  function EntitySetAttributes(entityId: string, params: AttributeValue): void;

  /**
   * Gets the attributes of an entity.
   * @param {string} entityId - Entity identifier.
   * @returns {AttributeValue} - Entity attributes.
   *
   * @example
   * ```ts
   * const checkerEntity = 'sensor.battery_checker'
   *
   * const attr = EntityGetAttributes(checkerEntity);
   * ```
   */
  function EntityGetAttributes(entityId: string): AttributeValue;

  /**
   * Retrieves entity settings.
   * @param {string} entityId - Entity identifier.
   * @returns {AttributeValue} - Entity settings.
   *
   * @example
   * ```ts
   * const checkerEntity = 'sensor.battery_checker'
   *
   * const settings = EntityGetSettings(checkerEntity)
   * if (!settings[msg.entity_id]) {
   *   return false;
   * }
   * const telegram = settings['telegram'];
   * ```
   */
  function EntityGetSettings(entityId: string): AttributeValue;

  /**
   * Sets the entity metric.
   * @param {string} entityId - Entity identifier.
   * @param {{ [key: string]: any }} params - Metric parameters.
   *
   * @example
   * ```ts
   * const attrs = {
   *     foo: bar,
   * }
   * const metricName = 'counter'
   * EntitySetMetric(ENTITY_ID, name, attrs);
   * ```
   */
  function EntitySetMetric(entityId: string, params: { [key: string]: any }): void;

  /**
   * Causes the action of the entity.
   * @param {string} entityId - Entity identifier.
   * @param {string} action - The name of the action.
   * @param {{ [key: string]: any }} params - Parameters for the action.
   *
   * @example
   * ```ts
   * EntityCallAction('sensor.internet_checker', 'PING', {});
   * ```
   */
  function EntityCallAction(entityId: string, action: string, params: { [key: string]: any }): void;

  /**
   * Calls a function or method of an entity's script.
   * @param {string} entityId - Entity identifier.
   * @param {string} fn - The name of the function.
   * @param {any} payload - Parameters for the action.
   *
   * @example
   * ```ts
   * const automationTriggerBle = (msg: TriggerBleMessage) => {
   *   // console.log(marshal(msg))
   *   EntityCallScript('ble.term1', 'handler', msg.payload)
   * }
   * ```
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
   *
   * @example
   * ```ts
   * EntitiesCallAction({
   *   action_name: 'ON',
   *   tags: ['room'],
   * })
   * ```
   */
  function EntitiesCallAction(params1: CallAction, params2?: { [key: string]: any }): void;

  /**
   * Calls the entity script.
   * @param {string} entityId - Entity identifier.
   * @param {{ [key: string]: any }} params - Parameters for the script.
   *
   * @example
   * ```ts
   * const attrs = {
   *   foo: bar,
   * }
   * EntityCallScene(ENTITY_ID, attrs);
   * ```
   */
  function EntityCallScene(entityId: string, params: { [key: string]: any }): void;

  /**
   * Function to perform an entity action.
   * @param {string} entityId - Entity identifier.
   * @param {string} actionName - Action name.
   * @param {{ [key: string]: any }} params - Parameters for performing the action.
   *
   * @example
   * ```ts
   * // Function to perform entity action
   * const entityAction = (entityId: string, actionName: string): void => {
   *   switch (actionName) {
   *     case 'CHECK':
   *       return checkStatus();
   *   }
   * };
   * ```
   */
  type entityAction = (entityId: string, actionName: string, params: { [key: string]: any }) => void;


  /////////////////////////////
  /// geo
  /////////////////////////////
  /**
   * Calculates the distance from a point to the center of an area.
   * @param {number} areaId - Area ID.
   * @param {Point} point - Point to calculate the distance.
   * @returns {number} - Distance from the point to the center of the area.
   *
   * @example
   * ```ts
   * distance = GeoDistanceToArea("my_area", { lat: 55.7558, lon: 37.6176 })
   *
   * if distance < 1000
   *   console.log("The point is close to the geographic area.")
   * else
   *   console.log("The point is located far from the geographic area.")
   * ```
   */
  function GeoDistanceToArea(areaId: number, point: Point): number;

  /**
   * Calculates the distance between two points.
   * @param {Point} point1 - The first point.
   * @param {Point} point2 - Second point.
   * @returns {number} - Distance between two points.
   *
   * @example
   * ```ts
   * point1 = { lat: 34.0522, lon: -118.2437 }
   * point2 = { lat: 37.7749, lon: -122.4194 }
   *
   * distance = GeoDistanceBetweenPoints(point1, point2)
   *
   * console.log("Distance between points:", distance, "km")
   * ```
   */
  function GeoDistanceBetweenPoints(point1: Point, point2: Point): number;

  /**
   * Checks if a point is inside an area.
   * @param {number} areaId - Area ID.
   * @param {Point} point - Point to check.
   * @returns {boolean} - true if the point is inside the region, false otherwise.
   *
   * @example
   * ```ts
   * isInside = GeoPointInsideArea(123, { lat: 40.7128, lon: -74.0060 })
   *
   * if isInside
   *   console.log("The point is located inside the Area")
   * else
   *   console.log("The point is located outside the Area.")
   * ```
   */
  function GeoPointInsideArea(areaId: number, point: Point): boolean;

  /////////////////////////////
  /// http
  /////////////////////////////
  /**
   * Interface for responding to an HTTP request.
   */
  interface HttpResponse {
    /**
     * The body content of the response.
     */
    body: string;

    /**
     * A boolean indicating whether an error occurred. This is optional.
     */
    error: boolean;

    /**
     * An error message if an error occurred. This is optional.
     */
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
     *
     * @example
     * ```ts
     * res = http.get("%s")
     * if res.error
     *   return
     * p = JSON.parse(res.body)
     * ```
     */
    get(url: string): HttpResponse;

    /**
     * Performs an HTTP POST request.
     * @param {string} url - URL for the POST request.
     * @param {string} data - Data for the POST request.
     * @returns {HttpResponse} - Response to a POST request.
     *
     * @example
     * ```ts
     * res = http.post("%s", {'foo': 'bar'})
     * if res.error
     *   return
     * p = JSON.parse(res.body)
     * ```
     */
    post(url: string, data: string): HttpResponse;

    /**
     * Performs an HTTP PUT request.
     * @param {string} url - URL for the PUT request.
     * @param {string} data - Data for the PUT request.
     * @returns {HttpResponse} - Response to the PUT request.
     *
     * @example
     * ```ts
     * res = http.put("%s", {'foo': 'bar'})
     * if res.error
     *   return
     * p = JSON.parse(res.body)
     * ```
     */
    put(url: string, data: string): HttpResponse;

    /**
     * Performs an HTTP HEAD request.
     * @param {string} url - URL for the PUT request.
     * @param {string} data - Data for the PUT request.
     * @returns {HttpResponse} - Response to the PUT request.
     *
     * @example
     * ```ts
     * res = http.head("%s", {'foo': 'bar'})
     * if res.error
     *   return
     * p = JSON.parse(res.body)
     * ```
     */
    put(url: string, data: string): HttpResponse;

    /**
     * Performs an HTTP DELETE request.
     * @param {string} url - URL for the DELETE request.
     * @returns {HttpResponse} - Response to a DELETE request.
     *
     * @example
     * ```ts
     * res = http.delete("%s")
     * if res.error
     *   return
     * p = JSON.parse(res.body)
     * ```
     */
    delete(url: string): HttpResponse;

    /**
     * Sets headers for an HTTP request.
     * @param {{ [key: string]: any }} params - Request headers.
     *
     * @example
     * ```ts
     * res = http.headers([{'apikey': 'some text'}]).get("%s")
     * if res.error
     * return
     * p = JSON.parse(res.body)
     * ```
     * @returns {http} - HTTP client with set headers.
     */
    headers(params: { [key: string]: any }): http;

    /**
     * Sets the basic authorization for the HTTP request.
     * @param {string} username - Username.
     * @param {string} password - User password.
     * @returns {http} - HTTP client with basic authorization set.
     *
     * @example
     * const res = http.basicAuth('user', 'password').download(uri);
     * ```
     */
    basicAuth(username: string, password: string): http;

    /**
     * Sets Digest protocol authentication for an HTTP request.
     * @param {string} username - Username.
     * @param {string} password - User password.
     * @returns {http} - HTTP client with Digest protocol authentication installed.
     *
     *
     * @example
     * ```ts
     * const res = http.digestAuth('user', 'password').download(uri);
     * ```
     */
    digestAuth(username: string, password: string): http;

    /**
     * Performs an HTTP request to download a file.
     * @param {string} post - URL for downloading the file.
     * @returns {HttpResponse} - Response to a file download request.
     *
     * @example
     * ```ts
     * const res = http.basicAuth('user', 'password').download(uri);
     * const res = http.digestAuth('user', 'password').download(uri);
     * const res = http.download(uri);
     *
     * ```
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
    /**
     * The user IDs associated with the notification. This is optional.
     */
    userIDS?: string;

    /**
     * The title of the HTML5 notification.
     */
    title: string;

    /**
     * The body content of the HTML5 notification.
     */
    body: string;
  }

  /**
   * Message interface for Telegram.
   */
  interface MessageTelegram {
    /**
     * The chat ID to which the message will be sent. This is optional.
     */
    chat_id?: number;

    /**
     * The URI of a file to send via the message. This is optional.
     */
    file_uri?: string;

    /**
     * The local path of a file to send via the message. This is optional.
     */
    file_path?: string;

    /**
     * The URI of a photo to send via the message. This is optional.
     */
    photo_uri?: string;

    /**
     * The local path of a photo to send via the message. This is optional.
     */
    photo_path?: string;

    /**
     * An array of keys associated with the message. This is optional.
     */
    keys?: string[];

    /**
     * The body content of the Telegram message.
     */
    body: string;
  }

  /**
   * Message interface for SMS notifications.
   */
  interface MessageSMS {
    /**
     * The phone number to which the SMS will be sent.
     */
    phone: string;

    /**
     * The body content of the SMS message.
     */
    body: string;
  }

  /**
   * Message interface for Web Push notifications.
   */
  interface MessageWebpush {
    /**
     * The user IDs associated with the notification. This is optional.
     */
    userIDS?: string;

    /**
     * The title of the web push notification.
     */
    title: string;

    /**
     * The body content of the web push notification.
     */
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
    /**
     * The type of the notification. This can be one of predefined types such as `MessageType`.
     */
    type?: MessageType;

    /**
     * The entity ID associated with the notification, if any.
     */
    entity_id?: string;

    /**
     * Attributes specific to the message format (HTML5, Telegram, SMS, WebPush).
     * This is a union type that can be one of `MessageHtml5`, `MessageTelegram`, `MessageSMS`, or `MessageWebpush`.
     */
    attributes: MessageHtml5 | MessageTelegram | MessageSMS | MessageWebpush;
  }

  /**
   * Interface for sending notifications.
   *
   * @example
   * ```ts
   * const push = (title: string, body: string) => {
   *   const msg = notifr.newMessage();
   *   msg.type = 'webpush';
   *   msg.attributes = {
   *     'title': title,
   *     'body': body
   *   };
   *   notifr.send(msg);
   * }
   * ```
   *
   * @example
   * ```ts
   * const sendMsg = function (body: string) {
   *   const settings = EntityGetSettings(checkerEntity)
   *   if (!settings[msg.entity_id]) {
   *     return false;
   *   }
   *   const telegram = settings['telegram'];
   *   var msg = notifr.newMessage();
   *   msg.entity_id = telegram;
   *   msg.attributes = {
   *     body: body,
   *     // keys: keyboard
   *   } as MessageTelegram;
   *   notifr.send(msg);
   * };
   * ```
   */
  interface notifrInterface {
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
  interface StorageInterface {
    /**
     * Adds a value to the store using the specified key.
     * @param {string} key - The key to add the value.
     * @param {string} value - The value to add.
     * @returns {string} - Added value.
     *
     *
     * @example
     * ```ts
     * Storage.push('temperature', 25.5);
     * ```
     */
    push(key: string, value: string): string;

    /**
     * Retrieves a value from storage for the specified key.
     * @param {string} key - The key to get the value.
     * @returns {string} - The returned value.
     *
     * @example
     * ```ts
     * const temperature = Storage.getByName('temperature');
     * console.log(temperature);
     * ```
     */
    getByName(key: string): string;

    /**
     * Searches the storage for values using the specified key.
     * @param {string} key - Key to search for values.
     * @returns {{ [key: string]: string }} - Found values as an object.
     *
     * @example
     * ```ts
     * const result = Storage.search('temperature');
     * console.log(result);
     * ```
     */
    search(key: string): { [key: string]: string };

    /**
     * Removes and returns a value from storage for the specified key.
     * @param {string} key - The key to delete the value.
     * @returns {string} - The removed value.
     *
     * @example
     * ```ts
     * const removedValue = Storage.pop('temperature');
     * console.log(removedValue);
     * ```
     */
    pop(key: string): string;
  }

  interface Tag {
    /**
     * The unique ID of the tag.
     */
    id: number;

    /**
     * The name of the tag.
     */
    name: string;
  }


  /**
   * Interface for representing a variable.
   */
  interface Variable {
    /**
     * The date and time when the variable was created.
     */
    created_at: string;

    /**
     * The date and time when the variable was last updated.
     */
    updated_at: string;

    /**
     * The name of the variable.
     */
    name: string;

    /**
     * The value of the variable, represented as a string.
     */
    value: string;

    /**
     * The ID of the entity associated with this variable. Optional.
     */
    entity_id?: string;

    /**
     * The system or context in which the variable is defined.
     */
    system: string;

    /**
     * An array of tags associated with the variable.
     */
    tags: string[];
  }

  /**
   * Interface for representing the response when listing variables.
   */
  interface VariableListResponse {
    /**
     * An array of variable objects.
     */
    items: Variable[];

    /**
     * The total number of variables returned in the list.
     */
    total: number;

    /**
     * An optional error object, which can contain details about any errors that occurred.
     */
    error?: any;
  }

  /**
   * Interface for representing the options when listing variables.
   */
  interface ListVariableOptions {
    /**
     * The maximum number of items to return in the list.
     */
    limit: number;

    /**
     * The offset from which to start fetching items in the list.
     */
    offset: number;

    /**
     * The field by which to order the results (e.g., `created_at`, `name`).
     */
    orderBy: string;

    /**
     * The sort direction of the results (`asc` for ascending, `desc` for descending).
     */
    sort: string;

    /**
     * A boolean indicating whether to filter by system variables.
     */
    system?: boolean;

    /**
     * An array of names to filter the variables by.
     */
    names?: string[];

    /**
     * A query string to search for specific variable names or values.
     */
    query?: string;

    /**
     * An array of tags to filter the variables by.
     */
    tags?: string[];

    /**
     * An array of entity IDs to filter the variables by.
     */
    entityIds: string[];
  }

  /**
   * Interface for representing a request to push a new variable.
   */
  interface VariablePushRequest {
    /**
     * The name of the variable.
     */
    name: string;

    /**
     * The value of the variable, represented as a string.
     */
    value: string;

    /**
     * An array of tags associated with the variable.
     */
    tags: string[];
  }

  /**
   * Interface for representing the response when getting a variable by name.
   */
  interface GetByNameResponse {
    /**
     * The variable object retrieved by name.
     */
    variable: Variable;

    /**
     * An optional error object, which can contain details about any errors that occurred.
     */
    error?: any;
  }

  /**
   * Interface for representing a variable and its associated methods.
   *
   * @example
   * ```ts
   * const variable:VariablePushRequest = {
   *   name: "v2",
   *   value: "Lorem ipsum dolore",
   *   tags: ["v2", "v3", "v4"]
   * }
   * ```
   */
  interface VariablesInterface {
    /**
     * Lists variables based on the given options.
     * @param options - The list options including limit, offset, order by, sort direction, etc.
     * @returns A response containing an array of variables and possibly an error object.
     *
     * @example
     * ```ts
     * const options: ListVariableOptions = {
     *   limit: 10,
     *   offset: 0,
     *   orderBy: "",
     *   sort: "",
     *   tags: ["v4"]
     * }
     * const list:VariableListResponse = Variables.list(options)
     * print(list.error)
     *
     * print(marshal(list.items[0]))
     * ```
     */
    list(options: ListVariableOptions): VariableListResponse;

    /**
     * Pushes a new variable with the given name and value.
     * @param options - The push request including the variable's name and value, tags, etc.
     * @returns A string indicating the success or failure of the operation.
     *
     * @example
     * ```ts
     * err = Variables.push(variable)
     * if (err) {
     *   print(err)
     * }
     * ```
     */
    push(options: VariablePushRequest): string;

    /**
     * Gets a variable by its name.
     * @param name - The name of the variable to retrieve.
     * @returns A response containing the retrieved variable and possibly an error object.
     *
     * @example
     * ```ts
     * const response:GetByNameResponse = Variables.getByName("v2")
     * if (response.error) {
     *   print(response.error)
     * }
     * ```
     */
    getByName(name: string): GetByNameResponse;

    /**
     * Deletes a variable by its name.
     * @param name - The name of the variable to delete.
     * @returns A string indicating the success or failure of the operation.
     *
     * @example
     * ```ts
     * err = Variables.delete("v2")
     * if (err) {
     *   print(err)
     * }
     * ```
     */
    delete(name: string): string;
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
   *
   * @example
   * ```ts
   * const entityAction = (ENTITY_ID, actionName, args) => {
   *   if (actionName == 'ON') {
   *     PushSystemEvent('command_enable_task', {id: 51});
   *   }
   *   if (actionName == 'OFF') {
   *     PushSystemEvent('command_disable_task', {id: 51});
   *   }
   * }
   * ```
   */
  function PushSystemEvent(command: SystemEventCommand, params: { [key: string]: any }): void;

  /**
   * Function to perform an action in Telegram.
   * @param {string} entityId - Entity identifier.
   * @param {string} actionName - Action name.
   * @param {{ [key: string]: any }} params - Parameters for performing the action.
   *
   * @example
   * ```ts
   * telegramAction = (entityId, actionName)->
   *  switch actionName
   *   when 'CHECK' then telegramSendReport()
   *
   *  sendMsg = (body)->
   *   msg = notifr.newMessage();
   *   msg.entity_id = 'telegram.testbot';
   *   msg.attributes = {
   *     'body': body
   *   };
   *   notifr.send(msg);
   * ```
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
     *
     * @example
     * ```ts
     * const continuousMove = (args) => {
     *   var X, Y;
     *   X = args['X'] || 0;
     *   Y = args['Y'] || 0;
     *   if (Math.abs(X) > Math.abs(Y)) {
     *     Y = 0;
     *   } else {
     *     X = 0;
     *   }
     *   Camera.continuousMove(X, Y);
     * };
     * ```
     */
    continuousMove(x: number, y: number): void;

    /**
     * Stops continuous camera movement.
     *
     * @example
     * ```ts
     * const stopStop = (args) => {
     *   Camera.stopContinuousMove();
     * };
     * ```
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
    /**
     * Indicates whether an error occurred during the operation. `true` if there was an error, `false` otherwise.
     */
    error: boolean;

    /**
     * Error message containing details about the error if one occurred. This will be present only if `error` is `true`.
     */
    errMessage: string;

    /**
     * The result of the operation, which can include various types of data depending on the interaction.
     */
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
    /**
     * The content of the message payload.
     */
    payload: string;

    /**
     * The topic on which the MQTT message was received or published.
     */
    topic: string;

    /**
     * The Quality of Service (QoS) level associated with the MQTT message.
     * This indicates the level of assurance for the delivery of the message.
     */
    qos: number;

    /**
     * Indicates whether the message is a duplicate of a previously published message at the same QoS level.
     */
    duplicate: boolean;

    /**
     * Error message if the message was not delivered successfully.
     * This property will be present only if `success` is `false`.
     */
    error: string;

    /**
     * Indicates whether the message delivery was successful.
     */
    success: boolean;

    /**
     * The new state parameters of the entity, which may include additional details about the state change.
     */
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
   * @param {Z2MMessage} msg - Zigbee2mqtt event.
   *
   * @example
   * ```ts
   * const zigbee2mqttEvent = (msg: Z2MMessage) => {
   *   console.log("Received MQTT message:");
   *   console.log("Payload:", msg.payload);
   *   console.log("Topic:", msg.topic);
   *   console.log("QoS:", msg.qos);
   *   console.log("Duplicate:", msg.duplicate);
   *
   *   if (msg.error) {
   *     console.error("Error:", msg.error);
   *   } else if (msg.success) {
   *     console.log("Operation successful!");
   *     console.log("New state:", msg.new_state);
   *   }
   *
   *   const value = msg.storage.getByName("key");
   *   console.log("Value from storage:", value);
   * }
   * ```
   */
  function zigbee2mqttEvent(msg: Z2MMessage): void;


  /////////////////////////////
  /// mqtt
  /////////////////////////////
  /**
   * Interface for interacting with MQTT.
   *
   * @example
   * ```ts
   * topic = entityId.split(".")[1] + '/cmnd/POWER'
   * Mqtt.publish(topic, 'ON', 0, false)
   * ```
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
   * Interface for MQTT messages.
   */
  interface MqttMessage {
    /**
     * The content of the message payload.
     */
    payload: string;

    /**
     * The topic on which the MQTT message was received or published.
     */
    topic: string;

    /**
     * The Quality of Service (QoS) level associated with the MQTT message.
     * This indicates the level of assurance for the delivery of the message.
     */
    qos: number;

    /**
     * Indicates whether the message is a duplicate of a previously published message at the same QoS level.
     */
    duplicate: boolean;

    /**
     * Error message if the message was not delivered successfully.
     * This property will be present only if `success` is `false`.
     */
    error: string;

    /**
     * Indicates whether the message delivery was successful.
     */
    success: boolean;

    /**
     * The new state parameters of the entity, which may include additional details about the state change.
     */
    new_state: EntityStateParams;
  }

  /**
   * Function called on MQTT event.
   *
   * @example
   * ```ts
   * const mqttEvent = (msg: MqttMessage) => {
   *   if (!msg || !msg.payload) {
   *     return;
   *   }
   *
   *   print(marshal(msg))
   * }
   * ```
   * @example
   * ```ts
   * mqttEvent = ->
   * #print '---mqtt new event---'
   *   arrayOfStrings = message.topic.split('/')
   *   if arrayOfStrings[1] == 'stat' && arrayOfStrings[2] == 'POWER'
   *     setState()
   *   if arrayOfStrings[1] == 'tele' && arrayOfStrings[2] == 'STATE'
   *     setTele()
   * ```
   */
  function mqttEvent(msg: MqttMessage): void;

  /////////////////////////////
  /// automation
  /////////////////////////////
  interface Action {
  }

  interface Condition {
  }

  /**
   * Trigger description.
   *
   * @example
   * ```ts
   * const trigger: Trigger = {
   *     id: 1,
   *     name: "every minute",
   *     plugin_name: "time",
   *     description: "This is an updated trigger with new settings.",
   *     entity_ids: ["entity1", "entity2"],
   *     script_id: 1,
   *     payload: {"cron": {"name": "cron", "type": "string", "value": "0 * * * * *"}},
   *     area_id: 1,
   *     enabled: false
   * };
   * ```
   */
  interface Trigger {
    /**
     * Primary key.
     */
    id: number;
    /**
     * Trigger name.
     */
    name: string;
    /**
     * Name of the plugin that uses this trigger.
     */
    plugin_name: string;
    /**
     * Description of the trigger function or behavior. (optional)
     */
    description?: string;
    /**
     * Array of entity or object IDs that the trigger acts on. (optional)
     */
    entity_ids?: string[];
    /**
     * ID of the script or program code that the trigger calls when executed. (optional)
     */
    script_id?: number;
    /**
     * An object with attributes that are passed along with the trigger. (optional)
     */
    payload?: Attributes;
    /**
     * ID of the area or zone in which the trigger acts. (optional)
     */
    area_id?: number;
    /**
     * A boolean value indicating whether the trigger is enabled or not. (optional)
     */
    enabled?: boolean;
  }

  /**
   * Function called when the automation condition is met.
   * @param {string} entityId - The ID of the entity associated with the automation condition.
   * @returns {boolean} - The result of the condition being met.
   *
   * @example
   * ```ts
   * const automationCondition = (entityId) => {
   *   const currentState = EntityGetState('mqtt.iphone_se');
   *   print('currentState: ' + currentState.name);
   *   return currentState.name == 'is at home';
   * };
   * ```
   */
  type automationCondition = (entityId: string) => boolean;

  /**
   * A function called when an automation action is executed.
   * @param {string} entityId - The ID of the entity associated with the automation action.
   *
   * @example
   * ```ts
   * const automationAction = (entityId: string): void => {
   *   checkStatus();
   * };
   * ```
   */
  type automationAction = (entityId: string) => void;

  /**
   * Interface for Alexa Time trigger messages.
   */
  interface TriggerAlexaTimeMessage {
    /**
     * The payload containing the specific details or data associated with the Alexa time event. This can vary depending on the event type.
     */
    payload: any;

    /**
     * A unique name identifying the type of trigger event or action that initiated the message, specifically related to time-based triggers in Alexa.
     */
    trigger_name: string;

    /**
     * The entity identifier associated with this specific Alexa time trigger instance.
     */
    entity_id: string;
  }

  /**
   * Function called when an Alexa Time trigger event occurs.
   * @param {TriggerAlexaTimeMessage} msg - Alexa Time trigger message.
   * @returns {boolean} - The result of the trigger execution.
   */
  type automationTriggerAlexa = (msg: TriggerAlexaTimeMessage) => boolean;

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
   *
   * @example
   * ```ts
   * const automationTriggerTime = function (msg) => {
   *   EntityCallAction('sensor.internet_checker', 'CHECK', {});
   *   return true;
   * };
   * ```
   */
  type automationTriggerTime = (msg: TriggerTimeMessage) => boolean;

  /**
   * Interface for the entity state change event.
   */
  interface EventEntityState {
    /**
     * The identifier of the entity whose state has changed.
     */
    entity_id: string;

    /**
     * The current value of the entity's state.
     */
    value: any;

    /**
     * A brief representation or code representing the new state of the entity.
     */
    state: EntityStateShort;

    /**
     * Additional attributes associated with the entity, such as metadata or other properties.
     */
    attributes: AttributeValue;

    /**
     * Settings or configuration parameters for the entity.
     */
    settings: AttributeValue;

    /**
     * The date and time when the entity's state was last changed.
     */
    last_changed: Date;

    /**
     * The date and time when the entity's state was last updated.
     */
    last_updated: Date;
  }

  /**
   * Interface for trigger state change event.
   */
  interface EventStateChanged {
    /**
     * Indicates whether the storage should save the new state. `true` if it should be saved, otherwise `false`.
     */
    storage_save: boolean;

    /**
     * Indicates whether the metric associated with this entity should not be saved. `true` if the metric should not be saved, otherwise `false`.
     */
    do_not_save_metric: boolean;

    /**
     * The name of the plugin responsible for the state change.
     */
    plugin_name: string;

    /**
     * The identifier of the entity whose state has changed.
     */
    entity_id: string;

    /**
     * Optional. The old state of the entity before the state change occurred.
     */
    old_state?: EventEntityState;

    /**
     * The new state of the entity after the state change occurred.
     */
    new_state: EventEntityState;
  }

  /**
   * Interface for trigger state change event.
   */
  interface TriggerStateChangedMessage {
    /**
     * The payload containing details about the state change event. Optional.
     */
    payload?: EventStateChanged;

    /**
     * A unique name identifying the type of trigger that has changed its state.
     */
    trigger_name: string;

    /**
     * The entity identifier associated with this specific trigger instance.
     */
    entity_id: string;
  }

  /**
   * Function called on a state change trigger event.
   * @param {TriggerStateChangedMessage} msg - State change trigger message.
   * @returns {boolean} - The result of the trigger execution.
   *
   * @example
   * ```ts
   * const automationTriggerStateChanged = (msg: TriggerStateChangedMessage): boolean => {
   *   if (!msg?.entity_id) {
   *     return false
   *   }
   *
   *   const level = msg.payload.new_state.attributes.battery;
   *   // exit if battery level has not changed
   *   if (level == msg.payload.old_state.attributes.battery) {
   *     return false
   *   }
   *   console.log(msg.entity_id, ' -> ', checkerEntity, ' -> get settings')
   *
   *   return true;
   * }
   * ```
   */
  type automationTriggerStateChanged = (msg: TriggerStateChangedMessage) => boolean;

  /**
   * Interface for system trigger message.
   */
  interface SystemTriggerMessage {
    /**
     * The topic or category associated with the trigger message.
     */
    topic: string;

    /**
     * A unique name identifying the specific event or action that triggered the message.
     */
    event_name: string;

    /**
     * The actual event data or payload related to the trigger event.
     */
    event: any;
  }

  /**
   * Interface for system trigger message.
   */
  interface TriggerSystemMessage {
    /**
     * The payload containing the specific details or data associated with the system event.
     */
    payload: SystemTriggerMessage;

    /**
     * A unique name identifying the type of trigger event or action that initiated the message.
     */
    trigger_name: string;

    /**
     * The entity identifier associated with this specific system trigger instance.
     */
    entity_id: string;
  }

  /**
   * Function called when a system trigger event occurs.
   * @param {TriggerSystemMessage} msg - System trigger message.
   * @returns {boolean} - The result of the trigger execution.
   *
   * @example
   * ```ts
   * const automationTriggerSystem = (msg: TriggerSystemMessage) => {
   *   print('---trigger---', ' ', msg.payload.event_name, ' ', msg.payload.topic, ' ', msg.payload.event)
   *   return false
   * }
   * ```
   */
  type automationTriggerSystem = (msg: TriggerSystemMessage) => boolean;

  /**
   * Interface for Ble trigger messages.
   */
  interface TriggerBleMessage {
    /**
     * The binary data payload associated with the Bluetooth Low Energy (BLE) event.
     */
    payload: Uint8Array;

    /**
     * A unique name identifying the trigger event or action that initiated the message.
     */
    trigger_name: string;

    /**
     * The entity identifier associated with this specific BLE trigger instance.
     */
    entity_id: string;
  }

  /**
   * Function called when the Ble trigger event occurs.
   * @param {TriggerBleMessage} msg - Ble trigger message.
   * @returns {boolean} - The result of the trigger execution.
   *
   * @example
   * ```ts
   * const automationTriggerBle = (msg: TriggerBleMessage) => {
   *   // console.log(marshal(msg))
   *   EntityCallScript('ble.term1', 'handler', msg.payload)
   * }
   * ```
   */
  type automationTriggerBle = (msg: TriggerBleMessage) => boolean;

  /**
   * Interface for Stt trigger messages.
   */
  interface TriggerSttMessage {
    /**
     * The content of the message to be processed by the speech-to-text system.
     */
    payload: string;

    /**
     * A unique name identifying the trigger event or action that initiated the message.
     */
    trigger_name: string;

    /**
     * The entity identifier associated with this specific trigger instance.
     */
    entity_id: string;
  }

  /**
   * Function called when the "Speech to text" event occurs.
   * @param {TriggerSttMessage} msg - stt trigger message.
   * @returns {boolean} - The result of the trigger execution.
   *
   * @example
   * ```ts
   * const automationTriggerStt = (data: TriggerSttMessage): boolean => {
   *   console.log(`stt command: ${data.payload}`);
   *   return true
   * };
   * ```
   */
  type automationTriggerStt = (msg: TriggerSttMessage) => boolean;

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
   *
   * @example
   * ```ts
   *  const UUID_UNITS = 'EBE0CCBE-7A0A-4B0C-8A1A-6FF2997DA3A6'        // 0x00 - F, 0x01 - C    READ WRITE
   *
   *  const u: number = unit == UNITS.C ? 0xFF : 0x01;
   *
   *  const result = BleWrite(UUID_UNITS, new Uint8Array([u]), true)
   *  if (result.error) {
   *     this.error(result.error)
   *     return
   *  }
   * ```
   */
  function BleWrite(char: string, payload: Uint8Array, withResponse: boolean): BleResponse;

  /**
   * Read reads the current characteristic value.
   * @param {string} char - UUID Device characteristic.
   *
   * @example
   * ```ts
   *  const UUID_UNITS = 'EBE0CCBE-7A0A-4B0C-8A1A-6FF2997DA3A6'        // 0x00 - F, 0x01 - C    READ WRITE
   *
   *  const result = BleRead(UUID_UNITS)
   *  if (result.error) {
   *     this.error(result.error)
   *     return
   *  }
   *
   *  if (result.response) {
   *     this._unit = result.response[0] == 1 ? UNITS.C : UNITS.F
   *  }
   * ```
   */
  function BleRead(char: string): BleResponse;

  /////////////////////////////
  /// triggers
  /////////////////////////////

  /**
   * Function called if a trigger needs to be deleted.
   * @param {number} triggerId - The ID of the trigger.
   *
   * @example
   * ```ts
   * TriggerDelete(1)
   * ```
   */
  function TriggerDelete(triggerId: number): void;

  /**
   * Trigger description.
   *
   * @example
   * ```ts
   * const trigger: NewTrigger = {
   *     name: "every minute",
   *     plugin_name: "time",
   *     description: "This is an updated trigger with new settings.",
   *     entity_ids: ["entity1", "entity2"],
   *     script_id: 1,
   *     payload: {"cron": {"name": "cron", "type": "string", "value": "0 * * * * *"}},
   *     area_id: 1,
   *     enabled: false
   * };
   * ```
   */
  interface NewTrigger {
    /**
     * Trigger name.
     */
    name: string;

    /**
     * Name of the plugin that uses this trigger.
     */
    plugin_name: string;

    /**
     * Description of the trigger function or behavior. (optional)
     */
    description?: string;

    /**
     * Array of entity or object IDs that the trigger acts on. (optional)
     */
    entity_ids?: string[];

    /**
     * ID of the script or program code that the trigger calls when executed. (optional)
     */
    script_id?: number;

    /**
     * An object with attributes that are passed along with the trigger. (optional)
     */
    payload?: Attributes;

    /**
     * ID of the area or zone in which the trigger acts. (optional)
     */
    area_id?: number;

    /**
     * A boolean value indicating whether the trigger is enabled or not. (optional)
     */
    enabled?: boolean;
  }

  /**
   * Trigger result.
   */
  interface TriggerResult {
    /**
     * Error occurred during trigger execution. (optional)
     */
    error?: any;

    /**
     * ID of the trigger result. (optional)
     */
    id?: number;
  }

  /**
   * The function is called when a trigger needs to be created programmatically.
   * @param {NewTrigger} params - Request parameter.
   *
   * @example
   * ```ts
   * const trigger: NewTrigger = {
   *     name: "every minute",
   *     plugin_name: "time",
   *     description: "This is an updated trigger with new settings.",
   *     entity_ids: ["entity1", "entity2"],
   *     script_id: 1,
   *     payload: {"cron": {"name": "cron", "type": "string", "value": "0 * * * * *"}},
   *     area_id: 1,
   *     enabled: false
   * };
   *
   * const result = TriggerAdd(trigger);
   * console.log(result); // Output the new trigger with new parameters.
   * ```
   */
  function TriggerAdd(params: NewTrigger): TriggerResult;

  /**
   * The function is called when it is necessary to programmatically update the trigger parameters.
   * @param {Trigger} params - Request parameter.
   *
   * @example
   * ```ts
   * const updatedParams: Trigger = {
   *     id: 1,
   *     name: "every minute",
   *     plugin_name: "time",
   *     description: "This is an updated trigger with new settings.",
   *     entity_ids: ["entity1", "entity2"],
   *     script_id: 1,
   *     payload: {"cron": {"name": "cron", "type": "string", "value": "0 * * * * *"}},
   *     area_id: 1,
   *     enabled: false
   * };
   *
   * const result = TriggerUpdate(updatedParams);
   * console.log(result); // Output the updated trigger with new parameters.
   * ```
   */
  function TriggerUpdate(params: Trigger): TriggerResult;

type init = () => void
type main = () => void
}
