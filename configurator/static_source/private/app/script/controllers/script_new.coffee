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

  $scope.$watch 'script.script.lang', (lang)->
    return if !lang || lang == ''
    switch lang
      when 'javascript'
        $scope.ace_options.mode = 'javascript'
      when 'coffeescript'
        $scope.ace_options.mode = 'coffee'
      when 'lua'
        $scope.ace_options.mode = 'lua'

  vm.exec =->
    success =(data)->

    error =(result)->
      Message result.data.status, result.data.message

    vm.script.$exec success, error

  vm
]