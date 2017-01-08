angular
.module('appDirectives')
.directive 'dashboardWidgetMemory', ['$compile', '$templateCache'
($compile, $templateCache) ->
  restrict: 'A'
  replace: true
  scope:
    widget: '=dashboardWidgetMemory'
  templateUrl: '/core/templates/_widget_memory.html'
  link: ($scope, $element, $attrs) ->

    $scope.percent = 0
    $scope.mem_total = 0
    $scope.mem_free = 0

    $scope.$on 'telemetry_update', (e, data)->
      return if !data.memory
      $scope.mem_total = data.memory?.mem_total || data.memory?.mem_total || 0
      $scope.mem_free = data.memory?.mem_free || data.memory?.mem_free || 0
      $scope.percent = (($scope.mem_total - $scope.mem_free) / ($scope.mem_total/100)).toFixed(2)
      $scope.usage = $scope.mem_total - $scope.mem_free

    $scope.openSettings =()->
      console.log 'open settings', widget

    $scope.removeWidget =()->
      console.log 'remove widget', widget

    return
]