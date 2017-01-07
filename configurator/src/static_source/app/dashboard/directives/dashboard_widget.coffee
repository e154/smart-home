angular
.module('appDirectives')
.directive 'dashboardWidget', ['$compile', '$templateCache'
($compile, $templateCache) ->
  restrict: 'A'
  replace: true
  scope:
    widget: '=dashboardWidget'
  link: ($scope, $element, $attrs) ->

    return if !$scope.widget

    compile =->
      template = ''
      switch $scope.widget.type
        when 'memory'
          $scope.percent = 0
          $scope.mem_total = 0
          $scope.mem_free = 0
          template = $templateCache.get('/dashboard/templates/_widget_memory.html')
          $scope.$on 'telemetry_update', (e, data)->
            $scope.mem_total = data.memory?.mem_total || data.memory?.mem_total || 0
            $scope.mem_free = data.memory?.mem_free || data.memory?.mem_free || 0
            $scope.percent = (($scope.mem_total - $scope.mem_free) / ($scope.mem_total/100)).toFixed(2)
            $scope.usage = $scope.mem_total - $scope.mem_free
        when 'swap'
          $scope.percent = 0
          $scope.mem_total = 0
          $scope.mem_free = 0
          template = $templateCache.get('/dashboard/templates/_widget_memory.html')
          $scope.$on 'telemetry_update', (e, data)->
            $scope.mem_total = data.memory?.swap_total || data.memory?.swap_total || 0
            $scope.mem_free = data.memory?.swap_free || data.memory?.swap_free || 0
            $scope.percent = (($scope.mem_total - $scope.mem_free) / ($scope.mem_total/100)).toFixed(2)
            $scope.usage = $scope.mem_total - $scope.mem_free
        when 'cpu_dig'
          $scope.processors = {}
          $scope.usage = 0
          $scope.max = 0
          $scope.mhz = 0
          template = $templateCache.get('/dashboard/templates/_widget_cpu_dig.html')
          $scope.$on 'telemetry_update', (e, data)->
            $scope.usage = (data.cpu?.usage || 0).toFixed(2)
            $scope.max = $scope.usage if $scope.usage > $scope.max
            if data.cpu?.info?.processors
#              console.log $scope.processors
              $scope.processors = data.cpu.info.processors
              $scope.model_name = $scope.processors[0].model_name
              $scope.mhz = $scope.processors[0].mhz.toFixed(0)
              $scope.cores = $scope.processors[0].cores


      previousContent = $compile(template)($scope)
      $element.html(previousContent)

    # settings
    # --------------------
    $scope.openSettings =()->
      console.log 'open settings', widget

    $scope.removeWidget =()->
      console.log 'remove widget', widget

    compile()

    return
]