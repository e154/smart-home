angular
.module('app')
.config ['$translatePartialLoaderProvider', '$translateProvider', '$locationProvider', '$routeProvider', 'pikadayConfigProvider'
($translatePartialLoaderProvider, $translateProvider, $locationProvider, $routeProvider, pikadayConfigProvider) ->

  $translatePartialLoaderProvider.addPart('dashboard');

  $translateProvider.useLoader('$translatePartialLoader', {
    urlTemplate: '/static/translates/{part}/{lang}.json'
    loadFailureHandler: 'LocaleErrorHandler'
  })

  $translateProvider.preferredLanguage 'ru'
  $translateProvider.useSanitizeValueStrategy null

  $locationProvider.html5Mode
    enabled: true
    requireBase: false

  $routeProvider.otherwise
    redirectTo: '/'

#   Pikaday
    locales =
      ru:
        previousMonth : 'Назад',
        nextMonth     : 'Следующий',
        months        : ["Январь", "Февраль", "Март", "Апрель", "Май", "Июнь", "Июль", "Август", "Сентабрь", "Октябрь", "Ноябрь", "Декабрь"],
        weekdays      : ["Понедельник", "Вторник", "Среда","Четверг", "Пятница", "Суббота", "Воскресенье"],
        weekdaysShort : ["Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"]


  pikadayConfigProvider.setConfig
    i18n: locales.ru
    locales: locales
    theme: 'smart-theme'

]

angular
.module('app')
.run ['$rootScope', '$state',($rootScope, $state) ->
  $rootScope.$on '$stateChangeSuccess', (event, toState, toParams, fromState, fromParams) ->
    document.getElementsByTagName('body')[0].classList.remove('loading')

  # http://stackoverflow.com/questions/13547007/fancybox2-fancybox-causes-page-to-to-jump-to-the-top
  $(document).ready ()->
    $('.fancybox').fancybox helpers:
      overlay:
        locked: false
]