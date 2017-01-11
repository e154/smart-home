angular
.module('app')
.config ['$routeProvider', '$locationProvider', '$stateProvider'
($routeProvider, $locationProvider, $stateProvider) ->

  $stateProvider

    .state(
      name: "dashboard.user"
      url: "user"
      abstract: true
      views:
        '@dashboard':
          templateUrl: '/user/templates/user.html'
          controller: 'userCtrl'
    )

    .state(
      name: "dashboard.user.index"
      url: ""
      views:
        '@dashboard.user':
          templateUrl: '/user/templates/user.index.html'
          controller: 'userIndexCtrl'
    )

    .state(
      name: "dashboard.user.new"
      url: "/new"
      templateUrl: '/user/templates/user.new.html'
      controller: 'userNewCtrl'
    )

    .state(
      name: "dashboard.user.show"
      url: "/:id"
      templateUrl: '/user/templates/user.show.html'
      controller: 'userShowCtrl'
    )

    .state(
      name: "dashboard.user.edit"
      url: "/:id/edit"
      templateUrl: '/user/templates/user.edit.html'
      controller: 'userEditCtrl'
    )

]
