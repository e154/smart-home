angular
.module('appDirectives')
.directive 'dashboardWidgetDevices', ['$compile', '$templateCache'
($compile, $templateCache) ->
  restrict: 'A'
  replace: true
  scope:
    widget: '=dashboardWidgetDevices'
  templateUrl: '/core/templates/_widget_devices.html'
  link: ($scope, $element, $attrs) ->

    $scope.total = 0
    $scope.online = 0
    $scope.disabled = 0
    $scope.error = 0

    $scope.$on 'telemetry_update', (e, data)->
      return if !data.devices
      $scope.online = 0
      $scope.disabled = 0
      $scope.error = 0
      $scope.total = data.devices.total if data.devices.total
      angular.forEach data.devices.status, (status, device)->
        switch status.system_name
          when 'ENABLED'
            $scope.online++
          when 'DISABLED'
            break
          when 'ERROR'
            $scope.error++

        $scope.disabled = $scope.total - ($scope.online + $scope.error)


    $scope.openSettings =()->
      console.log 'open settings', widget

    $scope.removeWidget =()->
      console.log 'remove widget', widget

    return
]