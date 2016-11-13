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

    constructor: ->
      @callbacks = {}
      @socket = socketFactory({
        url: "#{window.server_url}/api/v1/ws"
        socket: new SockJS("#{window.server_url}/api/v1/ws")
      })

      @socket.setHandler "message", @onmessage
      @socket.setHandler "open", @onopen
      @socket.setHandler "close", @onclose

    setHandler: (event, callback)->
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

    onclose: ->

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

  new Stream()
]