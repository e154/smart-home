angular
.module('appControllers')
.controller 'flowEditCtrl', ['$scope', 'Message', '$stateParams', 'Flow', '$state', 'Workflow'
($scope, Message, $stateParams, Flow, $state, Workflow) ->
  vm = this

  vm.workflows = []
  success = (result)->
    console.log result
    vm.workflows = result.items
  error = (result)->
    Message result.data.status, result.data.message
  Workflow.all {}, success, error

  Flow.show {id: $stateParams.id}, (flow)->
    vm.flow = flow

  vm.remove =->
    if confirm('точно удалить процесс?')
      remove()

  remove =->
    success =->
      $state.go("dashboard.flow.index")
    error =(result)->
      Message result.data.status, result.data.message
    vm.flow.$delete success, error

  vm.submit =->
    success =(data)->
      $state.go("dashboard.flow.show", {id: data.id})

    error =(result)->
      Message result.data.status, result.data.message

    vm.flow.$update(success, error)

  vm
]