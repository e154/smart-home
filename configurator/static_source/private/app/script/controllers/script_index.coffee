angular
.module('appControllers')
.controller 'scriptIndexCtrl', ['$scope', 'Notify', 'Script', '$state', '$timeout'
($scope, Notify, Script, $state, $timeout) ->
  vm = this

  tableCallback = {}
  vm.options =
    perPage: 20
    resource: Script
    columns: [
      {
        name: '#'
        field: 'id'
      }
      {
        name: 'script.name'
        field: 'name'
        clickCallback: ($event, item)->
          $event.preventDefault()
          $state.go('dashboard.script.show', {id: item.id})
          false
      }
      {
        name: 'script.lang'
        field: 'lang'
      }
      {
        name: 'script.created_at'
        field: 'created_at'
        template: '<span>{{item[field] | readableDateTime}}</span>'
      }
      {
        name: 'script.update_at'
        field: 'update_at'
        template: '<span>{{item[field] | readableDateTime}}</span>'
      }
    ]
    menu:null
    callback: tableCallback
    onLoad: (result)->

  vm
]