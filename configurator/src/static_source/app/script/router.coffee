angular
.module('app')
.config ['$routeProvider', '$locationProvider', '$stateProvider'
($routeProvider, $locationProvider, $stateProvider) ->

  $stateProvider

    .state(
      name: "dashboard.script"
      url: "script"
      abstract: true
      views:
        '@dashboard':
          templateUrl: '/script/templates/script.html'
          controller: 'scriptCtrl as script'
    )

    .state(
      name: "dashboard.script.index"
      url: ""
      views:
        '@dashboard.script':
          templateUrl: '/script/templates/script.index.html'
          controller: 'scriptIndexCtrl as script'
    )

    .state(
      name: "dashboard.script.new"
      url: "/new"
      templateUrl: '/script/templates/script.new.html'
      controller: 'scriptNewCtrl as script'
    )

    .state(
      name: "dashboard.script.show"
      url: "/:id"
      templateUrl: '/script/templates/script.show.html'
      controller: 'scriptShowCtrl as script'
    )

    .state(
      name: "dashboard.script.edit"
      url: "/:id/edit"
      templateUrl: '/script/templates/script.edit.html'
      controller: 'scriptEditCtrl as script'
    )

]
