angular
.module('angular-map')
.factory 'mapConstructor', ['$rootScope', '$compile', 'Map', 'Message', 'Notify', 'mapEditor', 'mapLayer'
  ($rootScope, $compile, Map, Message, Notify, mapEditor, mapLayer) ->
    class mapConstructor extends mapEditor

      id: null
      settings: null
      scope: null
      panning: null
      container: null
      wrapper: null
      model: null
      _model: null

      constructor: (@scope, @id)->
        super

        @model = {}
        @_model = new Map({
          id: @id
          layers: []
        })
        @scope.zoom = 1.0
        @scope.settings =
          movable: true
          zoom: true
          grid: 5
          minHeight: 400
          minWidth: 400

      update: (cb)->
        @fadeIn()
        success =(data)=>
          @fadeOut()
          cb(data)
          @deserialize()

        error =(result)->
          Message result.data.status, result.data.message

        @serialize()
        @_model.$updateFull success, error

      load: ()->
        success =()=>
          @deserialize()
        error =(result)->
          Message result.data.status, result.data.message

        @_model.$showFull success, error

      remove: (cb)->
        return if !confirm('Вы точно хотите удалить эту карту?')
        success =(data)=>
          cb(data)
        error =(result)->
          Message result.data.status, result.data.message
        @_model.$delete success, error

      fadeIn: ()->
        @wrapper.find(".page-loader").fadeIn("fast")
      fadeOut: ()->
        @wrapper.find(".page-loader").fadeOut("fast")

      serialize: ()=>
        @_model.layers = []
        angular.forEach @model.layers, (layer)=>
          @_model.layers.push layer.serialize()

      deserialize: ()=>
        @model =
          id: @_model.id
          name: @_model.name
          description: @_model.description
          layers: []

        if @_model?.layers && @_model.layers.length != 0
          angular.forEach @_model.layers, (layer)=>
            @model.layers.push new mapLayer(@scope).deserialize(layer)

    mapConstructor
]