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
          templateUrl: '/dashboard/templates/dashboard.html'
          controller: 'dashboardCtrl as dashboard'
      )

    .state(
      name: "dashboard.index"
      url: ""
      templateUrl: '/dashboard/templates/dashboard.index.html'
    )

    .state(
      name: "dashboard.node"
      url: "node"
      templateUrl: '/node/templates/node.index.html'
      controller: 'nodeIndexCtrl as node'
    )

    .state(
      name: "dashboard.node_new"
      url: "node/new"
      templateUrl: '/node/templates/node.new.html'
      controller: 'nodeNewCtrl as node'
    )

    .state(
      name: "dashboard.node_show"
      url: "node/:id"
      templateUrl: '/node/templates/node.show.html'
      controller: 'nodeShowCtrl as node'
    )

    .state(
      name: "dashboard.node_edit"
      url: "node/:id"
      templateUrl: '/node/templates/node.edit.html'
      controller: 'nodeEditCtrl as node'
    )

  $locationProvider.html5Mode
    enabled: true
    requireBase: false

  $routeProvider.otherwise
    redirectTo: '/'
]
