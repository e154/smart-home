angular
.module('appControllers')
.controller 'flowNewCtrl', ['$scope', 'Notify', 'Flow', '$state', 'Message', 'Workflow'
($scope, Notify, Flow, $state, Message, Workflow) ->
  vm = this

  vm.flow = new Flow({
    name: "Процесс"
    status: "enabled"
    description: ""
  })

  vm.workflows = []
  success = (result)->
    console.log result
    vm.workflows = result.items
  error = (result)->
    Message result.data.status, result.data.message
  Workflow.all {}, success, error

  vm.submit =->
    success =(data)->
      $state.go("dashboard.flow.show", {id: data.id})

    error =(result)->
      Message result.data.status, result.data.message

    vm.flow.$create(success, error)

  vm
]