#
#
# signature:
# Stream.send(object)
#
# promise = Stream.sendRequest(command, value)
# promise.then (result)->
#   $log.debug result
#
# example:
# Stream.send "request", {value: "qw"}
#
# Stream.sendRequest("request", {value: "qw"}).then (result)->
#   $log.debug result
#
#

angular
.module('appServices')
.service 'Stream', [ 'socketFactory', '$q', '$log', 'Notify', 'uuid', '$rootScope', '$timeout'
(socketFactory, $q, $log, Notify, uuid, $rootScope, $timeout) ->
  class Stream
    socket: null
    callbacks: null
    subscribers: null

    constructor: ->
      @subscribers = {}
      @callbacks = {}
      @connect()

    connect: ->
      url = "#{window.app_settings.server_url}/api/v1/ws?token=#{window.app_settings.token}"
      @socket = socketFactory({
        socket: new SockJS(url)
      })

      @setHandler "message", @onmessage
      @setHandler "open", @onopen
      @setHandler "close", @onclose

    setHandler: (event, callback)->
      @connect() if !@socket
      @socket.setHandler event, callback

    onmessage: (e)=>
      m = {type: ""}
      try
        m = angular.fromJson(e.data)
      catch
        $log.debug "from the stream came a string value"
        return

      switch m.type
        when "notify"
          Notify m.value.type, m.value.body
        when "broadcast"
          if (m.value.type in Object.keys(@subscribers))
            @subscribers[m.value.type](m.value.body)

      if !m.id || m.id == ""
        return

      if @callbacks[m.id]?.defer
        $rootScope.$apply @callbacks[m.id].defer.resolve(m)
        delete @callbacks[m.id]

      return

    onopen: ->
      clientInfo =->
        n = navigator

        width: window.screen.width
        height: window.screen.height
        cookie: n.cookieEnabled
        language: n.language
        platform: n.platform
#        plugins: n.plugins
        location: window.location.pathname
        href: window.location.href

      @send angular.toJson({client_info: {value: clientInfo()}})

    onclose: =>
      $timeout ()=>
        @connect()
      , 1000

    send: (m)->
      if typeof m != 'object'
        return

      @socket.send angular.toJson(m)

    sendRequest: (c, m)=>
      if typeof m != 'object'
        return

      defer = $q.defer()
      id = uuid.getId()
      m.id = id
      @callbacks[id] =
        time: new Date()
        defer: defer

      $timeout ()=>
        q = {}
        q[c] = m
        @send q

      defer.promise

    subscribe: (name, func)=>
      if !(name in @subscribers)
        @subscribers[name] = func
      else
        Notify "Ошибка", "Подписчик #{name} уже есть в списке"

    unsubscribe: (name)->
      if (name in @subscribers)
        delete @subscribers[name]


  new Stream()
]