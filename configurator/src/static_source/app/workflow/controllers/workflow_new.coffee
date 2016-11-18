angular
.module('appControllers')
.controller 'workflowNewCtrl', ['$scope', 'Notify', 'Workflow', '$state', 'Message'
($scope, Notify, Workflow, $state, Message) ->
  vm = this

  vm.workflow = new Workflow({
    name: "Новый процесс"
    ip: "127.0.0.1"
    port: 3000
    status: "enabled"
    description: ""
  })

  vm.submit =->
    success =(data)->
      $state.go("dashboard.workflow.show", {id: data.id})

    error =(result)->
      Message result.data.status, result.data.message

    vm.workflow.$create(success, error)

  vm
]