angular
.module('appControllers')
.controller 'deviceIndexCtrl', ['$scope', 'Notify', 'Device', '$state', 'Node'
($scope, Notify, Device, $state, Node) ->
  vm = this

  Node.get {
    limit:99
    offset: 0
    order: 'desc'
    query: {}
    sortby: 'created_at'
  }, (data)->
    vm.nodes = data.nodes

  tableCallback = {}
  vm.options =
    perPage: 20
    resource: Device
    columns: [
      {
        name: '#'
        field: 'id'
      }
      {
        name: 'device.name'
        field: 'name'
        clickCallback: ($event, item)->
          $event.preventDefault()
          $state.go('dashboard.device.show', {id: item.id})
          false
      }
      {
        name: 'device.node'
        field: 'node_id'
      }
      {
        name: 'device.created_at'
        field: 'created_at'
        template: '<span>{{item[field] | readableDateTime}}</span>'
      }
      {
        name: 'device.update_at'
        field: 'update_at'
        template: '<span>{{item[field] | readableDateTime}}</span>'
      }
      {
        name: 'device.status'
        field: 'status'
        width: '50px'
        template: "
<span class='label label-success' ng-if='item[\"status\"] == \"enabled\"'>{{'device.enabled' | translate}}</span>
<span class='label label-default' ng-if='item[\"status\"] == \"disabled\"'>{{'device.disabled' | translate}}</span>
"
      }
    ]
    menu:null
    callback: tableCallback

  vm
]