angular
.module('angular-map')
.factory 'mapEditor', ['$rootScope', '$compile', 'Fullscreen'
  ($rootScope, $compile, Fullscreen) ->
    class mapElement

      load_editor: (selector)->
        # container
        # --------------------
        container = angular.element(selector)
        wrapper = container.parent('.' + @wrap_class)
        if wrapper.length == 0
          container.wrap('<div class="' + @wrap_class + '"></div>')
        @wrapper = container.parent('.' + @wrap_class).attr('id', @id)
        @wrapper.append('<div class="page-loader"><div class="spinner">loading...</div></div>')

        # fullscreen
        # --------------------
        @fullscreen = new Fullscreen(@wrapper, @scope)
        @wrapper.append($compile('<div class="fullscreen entry" ng-click="resize()" data-help="resize editor window">full screen</div>')(@scope)) if @fullscreen.available

        # resizable
        # --------------------
        if @wrapper.resizable('instance')
          @wrapper.resizable('destroy')
        @wrapper.resizable
          minHeight: @minHeight
          minWidth: @minWidth
          grid: @grid
          handles: 's'

        @wrapper.find(".page-loader").fadeOut("fast")

        return

    mapElement
]