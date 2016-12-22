angular
.module('app')
.config ['$routeProvider', '$locationProvider', '$stateProvider'
($routeProvider, $locationProvider, $stateProvider) ->

  $stateProvider

    .state(
      name: "dashboard.device"
      url: "device"
      abstract: true
      views:
        '@dashboard':
          templateUrl: '/device/templates/device.html'
          controller: 'deviceCtrl as device'
    )

    .state(
      name: "dashboard.device.index"
      url: ""
      views:
        '@dashboard.device':
          templateUrl: '/device/templates/device.index.html'
          controller: 'deviceIndexCtrl as device'
    )

    .state(
      name: "dashboard.device.new"
      url: "/new"
      templateUrl: '/device/templates/device.new.html'
      controller: 'deviceNewCtrl as device'
    )

    .state(
      name: "dashboard.device.show"
      url: "/:id"
      views:
        "":
          templateUrl: '/device/templates/device.show.html'
          controller: 'deviceShowCtrl as device'
        "tab-view@dashboard.device.show":
          templateUrl: '/device/templates/device.action.html'
          controller: 'deviceActionCtrl as action'
    )

    .state(
      name: "dashboard.device.edit"
      url: "/:id/edit"
      templateUrl: '/device/templates/device.edit.html'
      controller: 'deviceEditCtrl as device'
    )

    .state(
      name: "dashboard.device.show.statuses"
      url: "/statuses"
      views:
        "tab-view@dashboard.device.show":
          templateUrl: '/device/templates/device.statuses.html'
          controller: 'deviceStatusCtrl as status'
    )
]
