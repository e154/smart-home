angular
.module('appControllers')
.controller 'scriptEditCtrl', ['$scope', 'Message', '$stateParams', 'Script', '$state', 'Notify'
($scope, Message, $stateParams, Script, $state, Notify) ->
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
      Notify 'success', 'Скрипт успешно сохранен', 3

    error =(result)->
      Message result.data.status, result.data.message

    vm.script.$update(success, error)

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