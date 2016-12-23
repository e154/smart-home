angular
.module('appServices')
.factory 'mapConstructor', ['$rootScope', '$compile', 'Map', 'Message', 'Notify'
  ($rootScope, $compile, Map, Message, Notify) ->
    class mapConstructor extends Map

      id: null
      elements: null
      settings: null
      scope: null
      container: null

      constructor: (@scope, @id)->

      update: (cb)=>
        success =(data)->
          cb(data)

        error =(result)->
          Message result.data.status, result.data.message

        @.$update success, error

      load: ()=>
        @scope.isLoad = true
        success =->

        error =(result)->
          Message result.data.status, result.data.message

        @.$show success, error

      remove: (cb)=>
        return if !confirm('Вы точно хотите удалить эту карту?')
        success =(data)=>
          cb(data)
        error =(result)->
          Message result.data.status, result.data.message
        @.$delete success, error


    mapConstructor
]