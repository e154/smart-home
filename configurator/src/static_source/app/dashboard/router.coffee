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
]
