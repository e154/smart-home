angular
.module('app')
.config ['$routeProvider', '$locationProvider', '$stateProvider'
($routeProvider, $locationProvider, $stateProvider) ->

  $stateProvider

    .state(
      name: "dashboard.node"
      url: "node"
      abstract: true
      views:
        '@dashboard':
          templateUrl: '/node/templates/node.html'
          controller: 'nodeCtrl as node'
    )

    .state(
      name: "dashboard.node.index"
      url: ""
      views:
        '@dashboard.node':
          templateUrl: '/node/templates/node.index.html'
          controller: 'nodeIndexCtrl as node'
    )

    .state(
      name: "dashboard.node.new"
      url: "/new"
      templateUrl: '/node/templates/node.new.html'
      controller: 'nodeNewCtrl as node'
    )

    .state(
      name: "dashboard.node.show"
      url: "/:id"
      templateUrl: '/node/templates/node.show.html'
      controller: 'nodeShowCtrl as node'
    )

    .state(
      name: "dashboard.node.edit"
      url: "/:id/edit"
      templateUrl: '/node/templates/node.edit.html'
      controller: 'nodeEditCtrl as node'
    )

]
