angular
.module('angular-map')
.factory 'MapImage', ['$rootScope', '$compile', 'Message', 'Notify'
  ($rootScope, $compile, Message, Notify) ->
    class MapImage

      scope: null

      constructor: (@scope)->

    MapImage
]