angular
.module('angular-map')
.factory 'MapDevice', [ '$http', 'Message', 'Notify', 'DeviceState', 'MapDeviceState', 'DeviceAction'
  ($http, Message, Notify, DeviceState, MapDeviceState, DeviceAction) ->
    class MapDevice

      id: null
      scope: null
      device: null
      device_action: null
      device_actions: []
      devices: []
      states: []
      image: null

      constructor: (@scope)->
        @scope.$watch(()=>
          @device
        , (val, old_val)=>
          return if !val || val == old_val
          @getDeviceStates(val)
          @getDeviceActions(val)
        )

      getDeviceActions: (device)->
        success =(actions)=>
          @device.actions = actions
        error =(result)->
          Message result.data.status, result.data.message

        DeviceAction.get_by_device {id: device.id}, success, error

      getDeviceStates: (device)->
        success =(states)=>
          @device.states = states
          @states = []
          def_exist = false
          angular.forEach @device.states, (device_state)=>
            def_exist = true if device_state.system_name == 'DEFAULT'
            md_state = new MapDeviceState(@scope, device_state)
            @states.push md_state

          if !def_exist
            md_state = new MapDeviceState(@scope)
            md_state.getDefault(@device)
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

      serialize: ()->
        return if !@device

        states = []
        angular.forEach @states, (_state)=>
          state = _state.serialize()
          state.map_device = {id: @id} if @id
          states.push state

        {
          id: @id if @id
          device: {id: @device.id} if @device
          states: states
          image: @image
          device_action: {id: @device_action.id} if @device_action
        }

      deserialize: (m)=>
        @id = m.id if m.id
        @device = m.device if m.device
        @device_action = m.device_action if m.device_action
        @status = m.status || 'enabled'
        @image = m.image || null

        @states = []
        angular.forEach @device.states, (device_state)=>
          md_state = new MapDeviceState(@scope, device_state)
          # check default state
          angular.forEach m.states, (state)=>
            if state.device_state?.id && state.device_state.id == device_state.id
              md_state.deserialize state
          @states.push md_state

        @

      removeImage: ()->
        @image = null

      create: ()->
      update: ()->
      remove: ()->
      update: (cb)->
        @upload(cb)

    MapDevice
]