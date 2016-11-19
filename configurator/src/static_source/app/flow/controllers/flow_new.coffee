angular
.module('appControllers')
.controller 'flowNewCtrl', ['$scope', 'Notify', 'Flow', '$state', 'Message', 'Workflow'
($scope, Notify, Flow, $state, Message, Workflow) ->
  vm = this

  $scope.flow = new Flow({
    name: "Процесс"
    status: "enabled"
    description: ""
  })

  $scope.workflows = []
  success = (result)->
    console.log result
    $scope.workflows = result.items
  error = (result)->
    Message result.data.status, result.data.message
  Workflow.all {}, success, error

  $scope.submit =->
    success =(data)->
      $state.go("dashboard.flow.show", {id: data.id})

    error =(result)->
      Message result.data.status, result.data.message

    $scope.flow.$create(success, error)

  vm
]