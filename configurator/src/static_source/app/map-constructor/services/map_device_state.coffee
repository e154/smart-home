angular
.module('angular-map')
.factory 'MapDeviceState', ['$filter',($filter) ->
  class MapDeviceState

    id: null
    device_state: null
    image: null
    style: null

    constructor: (@scope, @device_state)->

    serialize: ()->

      device_state = if @device_state.id then {id: @device_state.id} else null
      {
        id: @id if @id
        image: {id: @image.id} if @image
        device_state: device_state
        style: @style
      }

    deserialize: (m)->
      @id = m.id if m.id
      @device_state = m.device_state if m.device_state
      @image = m.image if m.image
      @style = m.style if m.style

      @

    removeImage: ()->
      @image = null

  MapDeviceState
]