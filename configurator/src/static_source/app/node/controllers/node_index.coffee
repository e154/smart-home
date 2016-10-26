angular
.module('appControllers')
.controller 'nodeIndexCtrl', ['$scope', 'Notify', 'Node', '$state'
($scope, Notify, Node, $state) ->
  vm = this

  tableCallback = {}
  vm.options =
    perPage: 20
    resource: Node
    columns: [
      {
        name: 'id'
        field: 'id'
        width: '50px'
      }
      {
        name: 'Node name'
        field: 'name'
        width: '48%'
        clickCallback: ($event, item)->
          $event.preventDefault()
          $state.go('dashboard.node.show', {id: item.id})
          false
      }
      {
        name: 'ip'
        field: 'ip'
        width: '100px'
      }
      {
        name: 'port'
        field: 'port'
        width: '50px'
      }
      {
        name: 'created at'
        field: 'created_at'
        template: '<span>{{item[field] | date:"H:mm dd.MM.yyyy"}}</span>'
      }
      {
        name: 'update at'
        field: 'update_at'
        template: '<span>{{item[field] | date:"H:mm dd.MM.yyyy"}}</span>'
      }
    ]
    menu:
      column: 1
      buttons: []
    callback: tableCallback

  vm
]