angular
.module('angular-bpmn')
.factory 'bpmnPanning', ['log', '$compile', '$timeout'
  (log, $compile, $timeout) ->
    class bpmnPanning
      container: null
      scope: null
      wrapper: null
      constructor: (container, scope, wrapper)->
        @container = container
        @scope = scope
        @wrapper = wrapper

        @init()

      init: ()->
        template = $compile('<div class="cross normal"></div>')(@scope)
        @container.append(template)
        $(template).css({
          position: 'absolute'
          top: -23
          left: -45
        })

        @wrapper.append($compile('<div class="zoom-info" data-help="zoom">x{{zoom}}</div>')(@scope))

        drag =
          x: 0
          y: 0
          state: false

        delta =
          x: 0
          y: 0

        @wrapper.on 'mousedown', (e)=>
          if !@scope.settings.engine.container.movable
            return

          if $(e.target).hasClass('ui-resizable-handle') || $(e.target).hasClass('entry') ||
              $(e.target).hasClass('viewport') || $(e.target).hasClass('minimap')
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
          if !@scope.settings.engine.container.movable
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
          if !@scope.settings.engine.container.zoom
            return

          e.preventDefault()

          @scope.zoom += 0.1 * e.deltaY
          @scope.zoom = parseFloat(@scope.zoom.toFixed(1))
          if @scope.zoom < 0.1
            @scope.zoom = 0.1
          else if @scope.zoom > 2
            @scope.zoom = 2

          @scope.instance.setZoom(@scope.zoom)
          @scope.instance.repaintEverything(@scope.zoom)

          $timeout ()=>
            @scope.$apply(
              @scope.zoom
            )
          , 0

          @container.css({
            '-webkit-transform':@scope.zoom
            '-moz-transform':'scale('+@scope.zoom+')'
            '-ms-transform':'scale('+@scope.zoom+')'
            '-o-transform':'scale('+@scope.zoom+')'
            'transform':'scale('+@scope.zoom+')'
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

    bpmnPanning
]