angular
.module('appControllers')
.controller 'scriptNewCtrl', ['$scope', 'Notify', 'Script', '$state', 'Message'
($scope, Notify, Script, $state, Message) ->
  vm = this

  vm.script = new Script({
    name: "Новый скрипт"
    lang: "coffeescript"
    description: ""
    source: ""
  })

  vm.submit =->
    success =(data)->
      $state.go("dashboard.script.show", {id: data.id})

    error =(result)->
      Message result.data.status, result.data.message

    vm.script.$create(success, error)

  vm
]