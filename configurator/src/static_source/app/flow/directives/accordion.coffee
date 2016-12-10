angular
.module('appDirectives')
.directive 'accordionV1', ["log", (log)->
  restrict: 'A'
  scope:
    accordionV1: '=accordionV1'
    title: '@'
    expand: '@'
  transclude: true
  replace: true
  template: '<div class="panel panel-default" ng-class="{collapsed : !expand}">
  <div class="panel-heading" ng-click="callback($event)">{{title}} <span class="icon-triangle">◀</span></div>
  <div class="panel-body" ng-transclude></div>
</div>'
  link: ($scope, element, attrs) ->
    $scope.callback =->
      if element.hasClass('collapsed')
        element.removeClass('collapsed')
      else
        element.addClass('collapsed')

      return
]
