angular
.module('appControllers')
.controller 'scriptEditCtrl', ['$scope', 'Message', '$stateParams', 'Script', '$state'
($scope, Message, $stateParams, Script, $state) ->
  vm = this

  Script.show {id: $stateParams.id}, (script)->
    vm.script = script

  vm.remove =->
    if confirm('точно удалить узел?')
      remove()

  remove =->
    success =->
      $state.go("dashboard.script.index")
    error =(result)->
      Message result.data.status, result.data.message
    vm.script.$delete success, error

  vm.submit =->
    success =(data)->

    error =(result)->
      Message result.data.status, result.data.message

    vm.script.$update(success, error)

  vm
]