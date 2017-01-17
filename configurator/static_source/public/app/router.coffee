angular
.module('app')
.config ['$locationProvider', '$stateProvider'
($locationProvider, $stateProvider) ->

  $stateProvider
    .state(
      name: "signin"
      url: "/signin"
      controller: 'signinCtrl'
      templateUrl: '/templates/signin.html'
    )

    .state(
      name: "recovery"
      url: "/recovery"
      controller: 'recoveryCtrl'
      templateUrl: '/templates/recovery.html'
    )

    .state(
      name: "reset"
      url: "/reset"
      controller: 'resetCtrl'
      templateUrl: '/templates/reset.html'
    )
]
