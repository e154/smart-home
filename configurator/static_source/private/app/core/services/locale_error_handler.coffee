angular
.module('appServices')
.factory 'LocaleErrorHandler',[ '$q', '$log', ($q, $log) ->
  (part, lang)->
    $log.error('The " /translates/' + part + '/' + lang + '.json" part was not loaded.')
    $q.when({})
]