angular
.module('appControllers')
.controller 'workflowShowCtrl', ['$scope', 'Notify', 'Workflow', '$stateParams', '$state', '$timeout'
($scope, Notify, Workflow, $stateParams, $state, $timeout) ->
  vm = this

  success = (workflow) ->
    vm.workflow = workflow
    $timeout ()->
      $scope.getStatus().then (result)->
        $scope.workflows = result.workflows

        angular.forEach $scope.workflows, (value, id)->
          if workflow.id == parseInt(id, 10)
            vm.workflow.state = value
    , 500

  error = ->
    $state.go 'dashboard.workflow.index'

  Workflow.show {id: $stateParams.id}, success, error

  vm
]