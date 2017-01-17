angular
.module('angular-map')
.factory 'MapDevice', [ '$http', 'Message', 'Notify', 'DeviceState', 'MapDeviceState', 'DeviceAction', 'MapDeviceAction'
  ($http, Message, Notify, DeviceState, MapDeviceState, DeviceAction, MapDeviceAction) ->
    class MapDevice

      id: null
      scope: null
      device: null
      devices: []
      states: []
      actions: []
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
          @actions = []
          angular.forEach @device.actions, (device_action)=>
            md_action = new MapDeviceAction(@scope, device_action)
            @actions.push md_action
        error =(result)->
          Message result.data.status, result.data.message

        DeviceAction.get_by_device {id: device.device?.id || device.id}, success, error

      getDeviceStates: (device)->
        success =(states)=>
          @device.states = states
          @states = []
          angular.forEach @device.states, (device_state)=>
            md_state = new MapDeviceState(@scope, device_state)
            @states.push md_state

        error =(result)->
          Message result.data.status, result.data.message
        DeviceState.get_by_device {id: device.device?.id || device.id}, success, error

      # get devices (select2)
      refreshDevices: (query)=>
        $http(
          method: 'GET'
          url: window.app_settings.server_url + "/api/v1/device/search"
          params:
            query: query
            limit: 15
            offset: 0
        ).then (response)=>
          devices = angular.copy(response.data.devices)
          angular.forEach devices, (device, index)->
            if !device.device_id? && !device.address?
              devices.splice(index, 1)
          @devices = devices

      serialize: ()->
        return if !@device

        states = []
        angular.forEach @states, (_state)=>
          state = _state.serialize()
          state.map_device = {id: @id} if @id
          states.push state

        actions = []
        angular.forEach @actions, (_action)=>
          action = _action.serialize()
          action.map_device = {id: @id} if @id
          actions.push action

        {
          id: @id if @id
          device: {id: @device.id} if @device
          states: states
          actions: actions
          image: @image
        }

      deserialize: (m)=>
        @id = m.id if m.id
        @device = m.device if m.device
        @status = m.status || 'enabled'
        @image = m.image || null

        @states = []
        angular.forEach @device.states, (device_state)=>
          md_state = new MapDeviceState(@scope, device_state)
          angular.forEach m.states, (state)=>
            if state.device_state?.id && state.device_state.id == device_state.id
              md_state.deserialize state
          @states.push md_state

        @actions = []
        angular.forEach @device.actions, (device_action)=>
          md_action = new MapDeviceAction(@scope, device_action)
          angular.forEach m.actions, (action)=>
            if action.device_action?.id && action.device_action.id == device_action.id
              md_action.deserialize action
          @actions.push md_action

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