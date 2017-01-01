angular
.module('angular-map')
.factory 'MapDevice', [ '$http', 'Message', 'Notify', 'DeviceState', 'MapDeviceState'
  ($http, Message, Notify, DeviceState, MapDeviceState) ->
    class MapDevice

      id: null
      scope: null
      device: null
      devices: []
      states: []

      constructor: (@scope)->
        @scope.$watch(()=>
          @device
        , (val, old_val)=>
          return if val == old_val
          @getDeviceStates(val)
        )

      getDeviceStates: (device)->
        success =(states)=>
          @device.states = states
          @states = []
          angular.forEach @device.states, (device_state)=>
            md_state = new MapDeviceState(@scope, device_state)
            @states.push md_state

        error =(result)->
          Message result.data.status, result.data.message
        DeviceState.get_by_device {id: device.id}, success, error

      # get devices (select2)
      refreshDevices: (query)=>
        $http(
          method: 'GET'
          url: window.server_url + "/api/v1/device/search"
          params:
            query: query
            limit: 5
            offset: 0
        ).then (response)=>
          @devices = response.data.devices

      addMapDeviceState: (state)->

      serialize: ()->
        states = []
        angular.forEach @states, (_state)=>
          state = _state.serialize()
          state.map_device = {id: @id}
          states.push state

        {
          id: @id if @id
          device: {id: @device.id} if @device
          states: states
        }

      deserialize: (m)=>
        @id = m.id if m.id
        @device = m.device if m.device
        @status = m.status || 'enabled'

        @states = []
        angular.forEach @device.states, (device_state)=>
          md_state = new MapDeviceState(@scope, device_state)
          angular.forEach m.states, (state)=>
            if state.device_state.id == device_state.id
              md_state.deserialize state
          @states.push md_state

        @

      create: ()->
      update: ()->
      remove: ()->
      update: (cb)->
        @upload(cb)

    MapDevice
]