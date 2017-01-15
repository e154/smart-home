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

    .state(
      name: "dashboard.role"
      url: "role"
      abstract: true
      views:
        '@dashboard':
          templateUrl: '/user/templates/role.html'
          controller: 'roleCtrl'
    )

    .state(
      name: "dashboard.role.index"
      url: ""
      views:
        '@dashboard.role':
            templateUrl: '/user/templates/role.index.html'
            controller: 'roleIndexCtrl'
    )

    .state(
      name: "dashboard.role.new"
      url: "/new"
      templateUrl: '/user/templates/role.new.html'
      controller: 'roleNewCtrl'
    )

    .state(
      name: "dashboard.role.show"
      url: "/:name"
      templateUrl: '/user/templates/role.show.html'
      controller: 'roleShowCtrl'
    )

    .state(
      name: "dashboard.role.edit"
      url: "/:name/edit"
      templateUrl: '/user/templates/role.edit.html'
      controller: 'roleEditCtrl'
    )

    .state(
      name: "dashboard.permissions"
      url: "permissions"
      views:
        '@dashboard':
          templateUrl: '/user/templates/permissions.index.html'
          controller: 'permissionsIndexCtrl'
    )

    .state(
      name: "dashboard.permissions.show"
      url: "/permissions/:name"
      templateUrl: '/user/templates/permissions.show.html'
      controller: 'permissionsShowCtrl'
    )
]
