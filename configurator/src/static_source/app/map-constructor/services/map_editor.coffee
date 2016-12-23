angular
.module('angular-map')
.factory 'mapEditor', ['$rootScope', '$compile', 'mapFullscreen', 'mapPanning'
  ($rootScope, $compile, mapFullscreen, mapPanning) ->
    class mapElement

      scope: null

      load_editor: (selector)->
        # container
        # --------------------
        @container = angular.element(selector)
        wrapper = @container.parent('.' + @wrap_class)
        if wrapper.length == 0
          @container.wrap('<div class="' + @wrap_class + '"></div>')
        @wrapper = @container.parent('.' + @wrap_class).attr('id', @id)
        @wrapper.append('<div class="page-loader"><div class="spinner">{{"loading..."| translate}}</div></div>')

        # fullscreen
        # --------------------
        @fullscreen = new mapFullscreen(@wrapper, @scope)
        @wrapper.append($compile('<div class="fullscreen entry" ng-click="resize()" data-help="resize editor window">{{"full screen" | translate}}</div>')(@scope)) if @fullscreen.available

        # resizable
        # --------------------
        if @wrapper.resizable('instance')
          @wrapper.resizable('destroy')
        @wrapper.resizable
          minHeight: @scope.settings.minHeight
          minWidth: @scope.settings.minWidth
          grid: @scope.settings.grid
          handles: 's'

        @panning = new mapPanning(@container, @scope, @wrapper)
        @wrapper.find(".page-loader").fadeOut("fast")

        return

    mapElement
]