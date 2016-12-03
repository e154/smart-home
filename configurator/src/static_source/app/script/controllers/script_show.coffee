angular
.module('appControllers')
.controller 'scriptShowCtrl', ['$scope', 'Notify', 'Script', '$stateParams', '$state', '$timeout'
($scope, Notify, Script, $stateParams, $state, $timeout) ->
  vm = this

  success = (script) ->
    vm.script = script
    $timeout ()->
      $scope.getStatus().then (result)->
        $scope.scripts = result.scripts

        angular.forEach $scope.scripts, (value, id)->
          if script.id == parseInt(id, 10)
            vm.script.state = value
    , 500

  error = ->
    $state.go 'dashboard.script.index'

  Script.show {id: $stateParams.id}, success, error

  vm
]