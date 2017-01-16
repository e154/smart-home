angular
.module('appControllers')
.controller 'nodeIndexCtrl', ['$scope', 'Notify', 'Node', '$state', '$timeout'
($scope, Notify, Node, $state, $timeout) ->
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
      {
        name: 'node.status'
        sort: 'disable'
        width: '50px'
        template: "<span class='label label-success' ng-if='column.getStatus(item[\"id\"]) == \"connected\"'>{{\"node.state.connected\" | translate}}</span>
<span class='label label-warning' ng-if='column.getStatus(item[\"id\"]) == \"wait\"'>{{'node.state.wait' | translate}}</span>
<span class='label label-danger' ng-if='column.getStatus(item[\"id\"]) == \"error\"'>{{'node.state.error' | translate}}</span>
<span class='label label-default' ng-if='item[\"status\"] == \"disabled\"'>{{'node.disabled' | translate}}</span>
"
        getStatus: (id)->
          $scope.nodes[id]
      }
    ]
    menu:null
    callback: tableCallback
    onLoad: (result)->
      $timeout ()->
        $scope.getStatus().then (result)->
          $scope.nodes = result.nodes.status
      , 500

  vm
]