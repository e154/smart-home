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
          controller: 'flowCtrl as flow'
    )

    .state(
      name: "dashboard.flow.index"
      url: ""
      views:
        '@dashboard.flow':
          templateUrl: '/flow/templates/flow.index.html'
          controller: 'flowIndexCtrl as flow'
    )

    .state(
      name: "dashboard.flow.new"
      url: "/new"
      templateUrl: '/flow/templates/flow.new.html'
      controller: 'flowNewCtrl as flow'
    )

    .state(
      name: "dashboard.flow.show"
      url: "/:id"
      templateUrl: '/flow/templates/flow.show.html'
      controller: 'flowShowCtrl as flow'
    )

    .state(
      name: "dashboard.flow.edit"
      url: "/:id/edit"
      templateUrl: '/flow/templates/flow.edit.html'
      controller: 'flowEditCtrl as flow'
    )

]
