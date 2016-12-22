angular
.module('app')
.config ['$routeProvider', '$locationProvider', '$stateProvider'
($routeProvider, $locationProvider, $stateProvider) ->

  $stateProvider

    .state(
      name: "dashboard.map"
      url: "map"
      abstract: true
      views:
        '@dashboard':
          templateUrl: '/map/templates/map.html'
          controller: 'mapCtrl'
    )

    .state(
      name: "dashboard.map.index"
      url: ""
      views:
        '@dashboard.map':
          templateUrl: '/map/templates/map.index.html'
          controller: 'mapIndexCtrl'
    )

    .state(
      name: "dashboard.map.new"
      url: "/new"
      views:
        '@dashboard.map':
          templateUrl: '/map/templates/map.new.html'
          controller: 'mapNewCtrl'
    )

    .state(
      name: "dashboard.map.edit"
      url: "/:id"
      views:
        '@dashboard.map':
          templateUrl: '/map/templates/map.edit.html'
          controller: 'mapEditCtrl'
    )

    .state(
      name: "dashboard.map.edit.settings"
      url: "/settings"
      views:
        '@dashboard.map':
          templateUrl: '/map/templates/map.settings.html'
          controller: 'mapEditCtrl'
    )
]
