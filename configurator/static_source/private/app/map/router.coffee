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
      name: "dashboard.map.show"
      url: "/:id"
      views:
        '@dashboard.map':
          templateUrl: '/map/templates/map.show.html'
          controller: 'mapShowCtrl'
    )

    .state(
      name: "dashboard.map.edit"
      url: "/:id/edit"
      abstract: true
      views:
        '@dashboard.map':
          templateUrl: '/map/templates/map.edit.html'
          controller: 'mapEditCtrl'
    )

    .state(
      name: "dashboard.map.edit.main"
      url: ""
      views:
        'editortabs@dashboard.map.edit':
          templateUrl: '/map/templates/map.editor.main_window.html'
          controller: 'mapEditMainWindowCtrl'
    )

    .state(
      name: "dashboard.map.edit.settings"
      url: "/settings"
      views:
        'editortabs@dashboard.map.edit':
          templateUrl: '/map/templates/map.editor.settings.html'
    )

    .state(
      name: "dashboard.map.edit.callback"
      url: "/callback"
      views:
        'editortabs@dashboard.map.edit':
          templateUrl: '/map/templates/map.editor.callback.html'
    )
]
