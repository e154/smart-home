angular
.module('app')
.config ['$routeProvider', '$locationProvider', '$stateProvider'
($routeProvider, $locationProvider, $stateProvider) ->

  $stateProvider

    .state(
      name: "dashboard.notifr"
      url: "notifr"
      abstract: true
      views:
        '@dashboard':
          templateUrl: '/notifr/templates/notifr.html'
    )

    .state(
      name: "dashboard.notifr.index"
      url: ""
      views:
        '@dashboard.notifr':
          templateUrl: '/notifr/templates/email.index.html'
          controller: 'emailIndexCtrl as ctrl'
    )

    .state(
      name: "dashboard.notifr.template"
      url: "/template/:name"
      templateUrl: '/notifr/templates/email.template.html'
      controller: 'emailTemplateCtrl as ctrl'
    )

    .state(
      name: "dashboard.notifr.new_template"
      url: "/template"
      templateUrl: '/notifr/templates/email.template.html'
      controller: 'emailTemplateCtrl as ctrl'
    )

    .state(
      name: "dashboard.notifr.items"
      url: "/items"
      templateUrl: '/notifr/templates/email.item.html'
      controller: 'emailItemCtrl as ctrl'
    )

]
