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
        name: '#'
        field: 'id'
      }
      {
        name: 'node.name'
        field: 'name'
        clickCallback: ($event, item)->
          $event.preventDefault()
          $state.go('dashboard.node.show', {id: item.id})
          false
      }
      {
        name: 'node.ip'
        field: 'ip'
      }
      {
        name: 'node.port'
        field: 'port'
      }
      {
        name: 'node.created_at'
        field: 'created_at'
        template: '<span>{{item[field] | readableDateTime}}</span>'
      }
      {
        name: 'node.update_at'
        field: 'update_at'
        template: '<span>{{item[field] | readableDateTime}}</span>'
      }
    ]
    menu:null
    callback: tableCallback

  vm
]