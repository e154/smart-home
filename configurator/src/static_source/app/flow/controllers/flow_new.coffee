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
      Notify 'success', 'Скрипт успешно создан', 3
      $timeout ()->
        $state.go("dashboard.flow.show", {id: data.id})
      , 2000

    error =(result)->
      Message result.data.status, result.data.message

    $scope.flow.$create(success, error)

  vm
]