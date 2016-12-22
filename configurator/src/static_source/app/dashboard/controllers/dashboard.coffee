angular
.module('appControllers')
.controller 'dashboardCtrl', ['$scope', 'Notify', 'Stream'
($scope, Notify, Stream) ->
  vm = this

  vm.menu =
    'name': 'Main Menu'
    'items': [
      {
        'label': 'Dashboard'
        'route': 'dashboard.index'
        'icon': 'fa fa-home'
      }
      {
        'label': 'Maps'
        'route': 'dashboard.map'
        'link': 'dashboard.map.index'
        'icon': 'fa fa-map-o'
        'items': [
          {
            'label': 'All devices'
            'route': 'dashboard.map.index'
          }
          {
            'label': 'Add new'
            'route': 'dashboard.map.new'
          }
        ]
      }
      {
        'label': 'Devices'
        'route': 'dashboard.device'
        'link': 'dashboard.device.index'
        'icon': 'fa fa-microchip'
        'items': [
          {
            'label': 'All devices'
            'route': 'dashboard.device.index'
          }
          {
            'label': 'Add new'
            'route': 'dashboard.device.new'
          }
        ]
      }
      {
        'label': 'Nodes'
        'route': 'dashboard.node'
        'link': 'dashboard.node.index'
        'icon': 'fa fa-sitemap'
        'items': [
          {
            'label': 'All nodes'
            'route': 'dashboard.node.index'
          }
          {
            'label': 'Add new'
            'route': 'dashboard.node.new'
          }
        ]
      }
      {
        'label': 'Flows'
        'route': 'dashboard.flow'
        'link': 'dashboard.flow.index'
        'icon': 'bpmn-icon-business-rule'
        'items': [
          {
            'label': 'All flow'
            'route': 'dashboard.flow.index'
          }
          {
            'label': 'Add new'
            'route': 'dashboard.flow.new'
          }
        ]
      }
      {
        'label': 'Workflow'
        'route': 'dashboard.workflow'
        'link': 'dashboard.workflow.index'
        'icon': 'fa fa-cube'
        'items': [
          {
            'label': 'All workflow'
            'route': 'dashboard.workflow.index'
          }
          {
            'label': 'Add new'
            'route': 'dashboard.workflow.new'
          }
        ]
      }
      {
        'label': 'Scripts'
        'route': 'dashboard.script'
        'link': 'dashboard.script.index'
        'icon': 'fa fa-pencil-square-o'
        'items': [
          {
            'label': 'All scripts'
            'route': 'dashboard.script.index'
          }
          {
            'label': 'Add new'
            'route': 'dashboard.script.new'
          }
        ]
      }
      {
        'label': 'Logs'
        'route': 'dashboard.log.index'
        'icon': 'fa fa-file-text-o'

      }
      {
        'label': 'Notifr'
        'route': 'dashboard.notifr.index'
        'icon': 'fa fa-envelope'
        'items': [
          {
            'label': 'All templates'
            'route': 'dashboard.notifr.index'
          }
          {
            'label': 'Add template'
            'route': 'dashboard.notifr.new_template'
          }
          {
            'label': 'All items'
            'route': 'dashboard.notifr.items'
          }
        ]
      }
    ]


  $scope.nodes = {}
  Stream.subscribe "nodes", (nodes)->
    $scope.nodes = nodes

  vm
]