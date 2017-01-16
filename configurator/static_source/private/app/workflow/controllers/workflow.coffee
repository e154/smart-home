angular
.module('appControllers')
.controller 'workflowCtrl', ['$scope', 'Notify', 'Workflow', 'Stream', '$log'
($scope, Notify, Workflow, Stream, $log) ->
  vm = this

  $scope.workflows = []
  $scope.getStatus = ->
    Stream.sendRequest("get.workflow.status", {})

  vm
]