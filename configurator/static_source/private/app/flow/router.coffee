angular
.module('app')
.config ['$routeProvider', '$locationProvider', '$stateProvider'
($routeProvider, $locationProvider, $stateProvider) ->

  $stateProvider

    .state(
      name: "dashboard.flow"
      url: "flow"
      abstract: true
      views:
        '@dashboard':
          templateUrl: '/flow/templates/flow.html'
          controller: 'flowCtrl'
    )

    .state(
      name: "dashboard.flow.index"
      url: ""
      views:
        '@dashboard.flow':
          templateUrl: '/flow/templates/flow.index.html'
          controller: 'flowIndexCtrl'
    )

    .state(
      name: "dashboard.flow.new"
      url: "/new"
      templateUrl: '/flow/templates/flow.new.html'
      controller: 'flowNewCtrl'
    )

    .state(
      name: "dashboard.flow.show"
      url: "/:id"
      templateUrl: '/flow/templates/flow.show.html'
      controller: 'flowShowCtrl'
    )

    .state(
      name: "dashboard.flow.edit"
      url: "/:id/edit"
      templateUrl: '/flow/templates/flow.edit.html'
      controller: 'flowEditCtrl'
    )

]
