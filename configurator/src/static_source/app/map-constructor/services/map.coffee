angular
.module('angular-map')
.factory 'mapConstructor', ['$rootScope', '$compile', 'Map', 'Message', 'Notify', 'mapEditor'
  ($rootScope, $compile, Map, Message, Notify, mapEditor) ->
    class mapConstructor extends mapEditor

      id: null
      elements: null
      settings: null
      scope: null
      container: null
      wrap_class: 'map-wrapper'
      wrapper: null

      grid: 5
      minHeight: 400
      minWidth: 400
      model: null

      constructor: (@scope, @id)->
        @model = new Map({id: @id})

      update: (cb)->
        @fadeIn()
        success =(data)->
          @fadeOut()
          cb(data)

        error =(result)->
          Message result.data.status, result.data.message

        @model.$update success, error

      load: ()->
        success =()=>
          @init()
        error =(result)->
          Message result.data.status, result.data.message

        @model.$show success, error

      remove: (cb)->
        return if !confirm('Вы точно хотите удалить эту карту?')
        success =(data)=>
          cb(data)
        error =(result)->
          Message result.data.status, result.data.message
        @model.$delete success, error

      init: ()->
      fadeIn: ()->
        @wrapper.find(".page-loader").fadeIn("fast")
      fadeOut: ()->
        @wrapper.find(".page-loader").fadeOut("fast")



    mapConstructor
]