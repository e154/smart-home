angular
.module('angular-map')
.factory 'MapDeviceState', ['$filter',($filter) ->
  class MapDeviceState

    id: null
    device_state: null
    image: null

    constructor: (@scope, @device_state)->

    serialize: ()->
      return null if !@device_state

      {
        id: @id if @id
        image: {id: @image.id} if @image
        device_state: {id: @device_state.id}
      }

    deserialize: (m)->
      @id = m.id if m.id
      @device_state = m.device_state if m.device_state
      @image = m.image if m.image

      @

    removeImage: ()->
      @image = null

    getDefault: (device)->
      @device_state =
        system_name: 'DEFAULT'
        device: device
        description: $filter('translate')('Default state')

  MapDeviceState
]