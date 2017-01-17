angular
.module('app')
.config ['$locationProvider'
( $locationProvider) ->

  $locationProvider.html5Mode
    enabled: true
    requireBase: false

]

angular
.module('app')
.run ['$rootScope', '$state'
($rootScope, $state) ->

#    http://stackoverflow.com/questions/24764764/conditionally-set-angulars-ng-class-based-on-state
  $rootScope.$state = $state;

  return
]