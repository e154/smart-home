angular
.module('appDirectives')
.directive 'dashboardWidgetNodes', ['$compile', '$templateCache'
($compile, $templateCache) ->
  restrict: 'A'
  replace: true
  scope:
    widget: '=dashboardWidgetNodes'
  templateUrl: '/core/templates/_widget_nodes.html'
  link: ($scope, $element, $attrs) ->

    $scope.total = 0
    $scope.online = 0
    $scope.disabled = 0
    $scope.error = 0

    $scope.$on 'telemetry_update', (e, data)->
      return if !data.nodes?.status
      $scope.online = 0
      $scope.disabled = 0
      $scope.error = 0
      $scope.total = data.nodes.total if data.nodes.total
      angular.forEach data.nodes.status, (status, node)->
        switch status
          when 'connected'
            $scope.online++
          when 'error'
            $scope.error++

      $scope.disabled = $scope.total - ($scope.online + $scope.error)

    $scope.openSettings =()->
      console.log 'open settings', widget

    $scope.removeWidget =()->
      console.log 'remove widget', widget

    return
]