###
# Created by delta54 on 06.12.14.
###

angular
.module('appDirectives')
.directive 'wpmenu', ['$window', '$log', 'storage', '$timeout', '$state'
  ($window, $log, storage, $timeout, $state) ->
    {
      restrict: 'A'
      scope:
        wpmenu: '='
      template: '<!--- Sidebar navigation -->
        <script type="text/ng-template" id="categoryTree">
          <ul>
            <li ng-repeat="item in item.items" ng-class="{open: (item.route | includedByState), has_sub: item.items.length}">
              <a ui-sref="{{item.route}}" ui-sref-active="active"><i ng-if="item.icon" ng-class=[item.icon]></i><span translate="{{item.label}}"></span>
                <span ng-if="item.items.length" class="pull-right"><i class="fa" ng-class="(item.route | includedByState) ? \'fa-chevron-down\' : \'fa-chevron-left\'"></i></span>
              </a>
              <div ng-if="item.items" ng-include src="\'categoryTree\'" include-replace></div>
            </li>
          </ul>
        </script>
        <!-- If the main navigation has sub navigation, then add the class "has_sub" to "li" of main navigation. -->
        <ul id="nav">
          <!-- Main menu with font awesome icon -->
          <li ng-repeat="item in wpmenu.items" ng-class="{active: (item.route | includedByState), has_sub: item.items.length}">
              <span class="glow"></span>
              <a ui-sref="{{item.link || item.route}}" ng-class="{has_sub: item.items}"><i ng-class=[item.icon]></i><span translate="{{item.label}}"></span>
                <span ng-if="item.items.length" class="pull-right"><i class="fa" ng-class="(item.route | includedByState) ? \'fa-chevron-down\' : \'fa-chevron-left\'"></i></span>
              </a>
              <div ng-if="item.items" ng-include src="\'categoryTree\'" include-replace></div>
          </li>
        </ul>'
      link: ($scope, $element, $attrs) ->

        menu_storage = new storage('main_menu')


        MainMenu =
          minimized: true
          init: ->

            # востановление состояния min/max
            if localStorage
              minimized = menu_storage.getBool('minimized')
            if minimized
              MainMenu.min()

            MainMenu.update()


            $timeout ()->
              MainMenu.restore()
            , 1000

          update: ->
            $('body')
            .off 'click', '#main-content #nav a'
            .on 'click', '#main-content #nav a', (e) ->
              if $(this).parents('#main-content:first').hasClass('enlarged')
                e.preventDefault()
                $state.go $(this).attr('ui-sref')
                return false
              if $(this).parent().hasClass('has_sub')
                e.preventDefault()
              # open <--> hide elements menu
              if !$(this).hasClass('subdrop')
                MainMenu.open this
              else if $(this).hasClass('subdrop')
                MainMenu.close this

            .off 'click', '.menubutton'
            .on 'click', '.menubutton', (e) ->
              e.preventDefault()
              MainMenu.toggle()

            .off 'click', '.sidebar-dropdown a'
            .on 'click', '.sidebar-dropdown a', (e) ->
              e.preventDefault()
              if !$(this).hasClass('open')
                # hide any open menus and remove all other classes
                $('.sidebar #nav').slideUp 350
                $('.sidebar-dropdown a').removeClass 'open'
                # open our new menu and add the open class
                $('.sidebar #nav').slideDown 350
                $(this).addClass 'open'
              else if $(this).hasClass('open')
                $(this).removeClass 'open'
                $('.sidebar #nav').slideUp 350

            $('#nav > li.has_sub > a.open').addClass('subdrop').next('ul').show()

          close: (e) ->
            $(e).removeClass 'subdrop'
            $(e).next('ul').slideUp 350
            $('.pull-right i', $(e).parent()).removeClass('fa-chevron-down').addClass 'fa-chevron-left'
            #$(".pull-right i",$(this).parents("ul:eq(1)")).removeClass("fa-chevron-down").addClass("fa-chevron-left");

          open: (e) ->
            # hide any open menus and remove all other classes
            $('ul', $(e).parents('ul:first')).slideUp 350
            $('a', $(e).parents('ul:first')).removeClass 'subdrop'
            $('#nav .pull-right i').removeClass('fa-chevron-down').addClass 'fa-chevron-left'
            # open our new menu and add the open class
            $(e).next('ul').slideDown 350
            $(e).addClass 'subdrop'
            $('.pull-right i', $(e).parents('.has_sub:last')).removeClass('fa-chevron-left').addClass 'fa-chevron-down'
            $('.pull-right i', $(e).siblings('ul')).removeClass('fa-chevron-down').addClass 'fa-chevron-left'

          toggle: ->
            if !MainMenu.minimized
              MainMenu.min()
            else
              MainMenu.max()

            if localStorage
              menu_storage.setItem('minimized', MainMenu.minimized)

          max: ->
            if $('#main-content').hasClass('enlarged')
#              $('#nav .has_sub .pull-right i').addClass('fa-chevron-left').removeClass('fa-chevron-down').removeClass 'fa-chevron-right'
              $('#main-content').removeClass 'enlarged'
              MainMenu.minimized = false
              MainMenu.restore()

          min: ->
            if !$('#main-content').hasClass('enlarged')
              $('#nav .has_sub ul').removeAttr 'style'
#              $('#nav .has_sub .pull-right i').removeClass('fa-chevron-left').addClass 'fa-chevron-down'
              $('#nav ul .has_sub .pull-right i').removeClass('fa-chevron-down').addClass 'fa-chevron-right'
              $('#main-content').addClass 'enlarged'
              MainMenu.minimized = true

          restore: ->
            # развернуть активный пункт меню
            if !MainMenu.minimized
              $('#nav a.active').each ->
                activeParentLink = $(this).parents('.has_sub').find('a:first')
                MainMenu.open activeParentLink
                return

        $timeout ()->
          MainMenu.init()
        , 0
    }
  ]
