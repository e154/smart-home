angular
.module('angular-map')
.factory 'mapPanning', ['log', '$compile', '$timeout'
  (log, $compile, $timeout) ->
    class mapPanning

      container: null
      scope: null
      wrapper: null

      constructor: (@container, @scope, @wrapper)->
        @init()

      init: ()->
        @scope.zoom = 1.0 if !@scope.zoom
        @setZoom(@scope.zoom)

        if @scope.settings.cross?
          template = $compile('<div class="cross normal"></div>')(@scope)
          @container.append(template)

        $(template).css({
          position: 'absolute'
          top: -23
          left: -45
        })

        if @scope.settings.zoom
          @wrapper.append($compile('<div class="zoom-info" data-help="zoom">x{{zoom}}</div>')(@scope))

        drag =
          x: 0
          y: 0
          state: false

        delta =
          x: 0
          y: 0

        @wrapper.on 'mousedown', (e)=>
          if !@scope.settings.movable
            return

          if $(e.target).is('.ui-resizable-handle, .entry, .viewport, .minimap, .draggable, .draggable-entity')
            return

          if !drag.state && e.which == LEFT_MB
            drag.x = e.pageX
            drag.y = e.pageY
            drag.state = true

        @wrapper.on 'mousemove', (e)=>
          if drag.state
            delta.x = e.pageX - drag.x
            delta.y = e.pageY - drag.y

            cur_offset = @container.offset()

            @container.offset({
              left: (cur_offset.left + delta.x)
              top: (cur_offset.top + delta.y)
            })

            drag.x = e.pageX
            drag.y = e.pageY

        @wrapper.on 'mouseup', (e)=>
          if !@scope.settings.movable
            return

          if drag.state
            drag.state = false

        @wrapper.on 'contextmenu', (e)->
          return false

        @wrapper.on 'mouseleave', (e)->
          if drag.state
            drag.state = false

        # zoom
        #-----------------------------------
        @wrapper.on 'mousewheel', (e)=>
#          shift = key.getPressedKeyCodes().indexOf(16) > -1
          if !@scope.settings.zoom #|| !shift
            return

          e.preventDefault()

          @scope.zoom += 0.1 * e.deltaY
          @scope.zoom = parseFloat(@scope.zoom.toFixed(1))
          if @scope.zoom < 0.1
            @scope.zoom = 0.1
          else if @scope.zoom > 4
            @scope.zoom = 4

          $timeout ()=>
            @scope.$apply(
              @scope.zoom
            )
          , 0

          @setZoom(@scope.zoom)

      setZoom: (zoom)->
        @container.css({
          '-webkit-transform':zoom
          '-moz-transform':'scale('+zoom+')'
          '-ms-transform':'scale('+zoom+')'
          '-o-transform':'scale('+zoom+')'
          'transform':'scale('+zoom+')'
        })

      destroy: ()->
        @wrapper.off 'mousedown'
        @wrapper.off 'mousemove'
        @wrapper.off 'mouseup'
        @wrapper.off 'contextmenu'
        @wrapper.off 'mousewheel'

        @container = null
        @scope = null
        @wrapper = null

    mapPanning
]