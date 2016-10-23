angular
.module('app')
.config ['$routeProvider', '$locationProvider', '$stateProvider'
($routeProvider, $locationProvider, $stateProvider) ->

  $stateProvider
    .state(
      name: "dashboard"
      url: "/"
      abstract: true
      views:
        '@':
          templateUrl: '/templates/dashboard/dashboard.html'
          controller: 'dashboardCtrl as dashboard'
      )

    .state(
      name: "dashboard.index"
      url: ""
      templateUrl: '/templates/dashboard/dashboard.index.html'
    )

    .state(
      name: "dashboard.node"
      url: "/node"
      templateUrl: '/templates/node/node.html'
      controller: 'nodeIndexCtrl as node'
    )

    .state(
      name: "dashboard.node_show"
      url: "/node/:id"
      templateUrl: '/templates/node/ode.show.html'
      controller: 'nodeShowCtrl as node'
    )

    .state(
      name: "dashboard.node_edit"
      url: "/node/:id"
      templateUrl: '/templates/node/node.edit.html'
      controller: 'nodeEditCtrl as node'
    )

  $locationProvider.html5Mode
    enabled: true
    requireBase: false

  $routeProvider.otherwise
    redirectTo: '/'
]
