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
        'dashboard@':
          templateUrl: '/dashboard/templates/dashboard.html'
          controller: 'dashboardCtrl as dashboard'
      )

    .state(
      name: "dashboard.index"
      url: ""
      controller: 'dashboardIndexCtrl'
      templateUrl: '/dashboard/templates/dashboard.index.html'
      onExit: ()->
        angular.element(document).find('body').removeClass('dashboard')
    )

    .state(
      name: "dashboard.account"
      url: "account"
      controller: 'accountCtrl'
      templateUrl: '/dashboard/templates/dashboard.account.html'
    )
]
