angular
.module('appControllers')
.controller 'flowShowCtrl', ['$scope', 'Notify', 'Flow', '$stateParams', '$state', '$timeout'
($scope, Notify, Flow, $stateParams, $state, $timeout) ->
  vm = this

  # get flow
  $scope.flow = {}
  success = (flow) ->
    $scope.flow = flow
    $timeout ()->
      $scope.getStatus().then (result)->
        $scope.flows = result.flows

        angular.forEach $scope.flows, (value, id)->
          if flow.id == parseInt(id, 10)
            $scope.flow.state = value
    , 500

  error = ->
    $state.go 'dashboard.flow.index'

  Flow.get_redactor {id: $stateParams.id}, success, error

  vm
]