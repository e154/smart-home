angular
.module('angular-map')
.factory 'MapText', ['$rootScope', '$compile', 'Message', 'Notify', 'Image'
  ($rootScope, $compile, Message, Notify, Image) ->
    class MapImage

      id: null
      scope: null
      text: ''
      style: ''

      constructor: (@scope)->

      serialize: ()->
        id: @id if @id
        text: @text
        style: @style

      deserialize: (m)->
        @id = m.id if m.id
        @text = m.text if m.text
        @style = m.style if m.style

        @

    MapImage
]