angular
.module('app')
.config ['$routeProvider', '$locationProvider', '$stateProvider'
($routeProvider, $locationProvider, $stateProvider) ->

  $stateProvider

    .state(
      name: "dashboard.log"
      url: "log"
      abstract: true
      views:
        '@dashboard':
          templateUrl: '/log/templates/log.html'
          controller: 'logCtrl'
    )

    .state(
      name: "dashboard.log.index"
      url: ""
      views:
        '@dashboard.log':
          templateUrl: '/log/templates/log.index.html'
          controller: 'logIndexCtrl'
    )
]
