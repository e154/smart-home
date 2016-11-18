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
    ]


  $scope.nodes = {}
  Stream.subscribe "nodes", (nodes)->
    $scope.nodes = nodes

  vm
]