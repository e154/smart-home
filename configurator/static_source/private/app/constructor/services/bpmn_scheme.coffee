LEFT_MB = 1
MIDDLE_MB = 2
RIGHT_MB = 3

angular
.module('angular-bpmn')
.factory 'bpmnScheme', [
  '$rootScope', 'log', 'bpmnUuid', '$compile', 'bpmnSettings', '$templateCache', '$templateRequest', '$q', '$timeout'
  'bpmnObjectFact', 'bpmnPanning', 'bpmnFullscreen', 'bpmnEditor'
  ($rootScope, log, bpmnUuid, $compile, bpmnSettings, $templateCache, $templateRequest, $q, $timeout, bpmnObjectFact
   bpmnPanning, bpmnFullscreen, bpmnEditor) ->
    class bpmnScheme extends bpmnEditor

      isDebug: true
      isStarted: false
      id: null
      scope: null
      container: null
      wrapper: null
      cache: null
      style_id: 'bpmn-style-theme'
      wrap_class: 'bpmn-wrapper'
      schemeWatch: null
      stopListen: null
      panning: null
      fullscreen: null

      constructor: (container, settings)->
        super

#        set unique id
        @id = bpmnUuid.gen()

#        set main container
        @container = container

#        set parent wrap
        wrapper = container.parent('.' + @wrap_class)
        if wrapper.length == 0
          container.wrap('<div class="' + @wrap_class + '"></div>')
        @wrapper = container.parent('.' + @wrap_class).attr('id', @id)
#        preventSelection(document)

#        create scope
        @scope = $rootScope.$new()
        @scope.extScheme = {}          # scheme
        @scope.intScheme =          # real scheme objects
          objects: {}
          connectors: []
        @scope.settings = {}
        @scope.selected = []
        @scope.zoom = 1
#        @setSettings(settings)
        @scope.changeTheme = @changeTheme

        @wrapper.append('<div class="page-loader"><div class="spinner">loading...</div></div>')
        @wrapper.append($compile('<div ng-if="settings.engine.container.theme_selector" class="theme-selector entry">
<select class="form-control" ng-model="settings.engine.theme" ng-change="changeTheme()" style="width: auto;" data-help="select theme">
<option ng-repeat="them in settings.theme.list" value="{{them}}">{{them}}</option></select>
</div>')(@scope))

        @fullscreen = new bpmnFullscreen(@wrapper, @scope)
        @wrapper.append($compile('<div ng-if="settings.engine.container.fullscreen" class="fullscreen entry" ng-click="resize()" data-help="resize editor window">full screen</div>')(@scope)) if @fullscreen.available

      setStatus: ()->
        if @scope.settings.engine.status == 'editor'
          @reloadEditor()
          @container.addClass('editor')
          @container.removeClass('viewer')
        else
          @container.addClass('viewer')
          @container.removeClass('editor')

      setScheme: (scheme)->
        if !scheme?
          scheme = @scope.extScheme || {}

        @scope.extScheme = scheme

      # настройки, тема оформления, объекты
      #------------------------------------------------------------------------------
      setSettings: (settings)->
        if !settings?
          settings = {}

        @scope.settings = $.extend(true, angular.copy(bpmnSettings), angular.copy(settings))
        if settings?.editorPallet
          @scope.settings.editorPallet = angular.copy(settings.editorPallet)
        @scope.settings = $.extend(true, @scope.settings,
          point:
            isSource: true
            isTarget: true
            endpoint: ["Rectangle", {width:15, height:15}]
            paintStyle:
#              strokeStyle: "blue"
              outlineWidth: 1
#            hoverPaintStyle:
#              outlineColor: "blue"
            maxConnections: -1
        ) if @scope.settings.engine.status == 'editor'

        @loadStyle()
        @cacheTemplates()
        @objectsUpdate()
        @setStatus()

      loadStyle: ()->
        theme_name = @scope.settings.engine.theme
        file = @scope.settings.theme.root_path + '/' + theme_name + '/style.css'
        # id="bpmn-style-theme-default"
        themeStyle = $('link#' + @style_id + '-' + theme_name)
        if themeStyle.length == 0
          $("<link/>", {
            rel: "stylesheet"
            href: file
            id: @style_id + '-' + theme_name
          }).appendTo("head")
        else
          themeStyle.attr('href', file)

        log.debug 'load style file:', file

        @wrapper.removeClass()
        @wrapper.addClass(@wrap_class + ' ' + theme_name)

      cacheTemplates: ()->
        if !@cache
          @cache = []

        angular.forEach @scope.settings.templates, (type)=>
          if !type.templateUrl? || type.templateUrl == ''
            return

          templateUrl = @scope.settings.theme.root_path + '/' + @scope.settings.engine.theme + '/' + type.templateUrl
          template = $templateCache.get(templateUrl)
          if !template?
            templatePromise = $templateRequest(templateUrl)
            @cache.push(templatePromise)
            templatePromise.then (result)=>
              log.debug 'load template file:', templateUrl
              $templateCache.put(templateUrl, result)

      # пакетная обработка объектов
      #------------------------------------------------------------------------------
      makePackageObjects: (resolve)->
        @isStarted = false
        log.debug 'make package objects'
        # Создадим все объекты, сохраним указатели в массиве
        # потому как возможны перекрёстные ссылки
        promise = []
        angular.forEach @scope.extScheme.objects, (object)=>
          obj = new bpmnObjectFact(object, @scope)
          @scope.intScheme.objects[object.id] = obj
          promise.push(obj.elementPromise)

        # Ждём когда прогрузятся все шаблоны
        $q.all(promise).then ()=>
          # проходим по массиву ранее созданных объектов,
          # и добавляем в дом
          angular.forEach @scope.intScheme.objects, (object)=>
            # добавляем объект в контейнер
            object.appendTo(@container, @scope.settings.point)

          @connectPackageObjects()
          @scope.instance.repaintEverything()

          @isStarted = true
          @wrapper.find(".page-loader").fadeOut("slow")
          resolve()

      connectPackageObjects: ()->
        if !@scope.extScheme?.connectors?
          return

        log.debug 'connect package objects'

        angular.forEach @scope.extScheme.connectors, (connector)=>
          if (!connector.start || !connector.end) ||
              (!@scope.intScheme.objects[connector.start.object] || !@scope.intScheme.objects[connector.end.object]) ||
              (!@scope.intScheme.objects[connector.start.object].points || !@scope.intScheme.objects[connector.end.object].points)
            return

          # связь создаётся по точкам созданным ранее
          source_obj_points = {}
          target_obj_points = {}

          angular.forEach @scope.intScheme.objects, (object)->
            if object.data.id == connector.start.object
              source_obj_points = object.points
            if object.data.id == connector.end.object
              target_obj_points = object.points

          if !source_obj_points[connector.start.point]?
            log.error 'connect: source not found', this
            return

          if !target_obj_points[connector.end.point]?
            log.error 'connect: target not found', this
            return

          # связь создаётся по точкам
          points = {
            sourceEndpoint: source_obj_points[connector.start.point]

          # Привязка к объекту, якорь выбирается автоматически
#          target: @scope.intScheme.objects[connector.end.object]['object']

          # Привязка к конкретному якорю объекта
            targetEndpoint: target_obj_points[connector.end.point]

            # параметры соединения: id ...
            parameters:
              'element-id': connector.id || bpmnUuid.gen()
              'direction': connector.direction || ''
          }

          # подпись для связи
          if connector.title && connector.title != ""
            points['overlays'] = [
              [ "Label", { label:connector.title, cssClass: "aLabel" }, id:"myLabel" ]
            ]

          @scope.intScheme.connectors.push(@scope.instance.connect(points, @scope.settings.connector))

      deselectAll: ()->
        angular.forEach @scope.intScheme.objects, (object)->
          object.select(false)

        @scope.instance.select().each (c)->
          c.removeClass('selected')

        @scope.selected = []

      objectsUpdate: ()->
        angular.forEach @scope.intScheme.objects, (object)->
          object.templateUpdate()

      start: ()->
        log.debug 'start'
        return $q (resolve)=>
          @instart(resolve)

      instart: (resolve)->

        if !@panning
          @panning = new bpmnPanning(@container, @scope, @wrapper)

        @loadStyle()

        if !@scope.instance
          @scope.instance = jsPlumb.getInstance($.extend(true, @scope.settings.instance, {Container: @container}))

        @cache = []
        @cacheTemplates()
        @container.addClass('bpmn')
#        @setStatus()

        # watchers
#        if @schemeWatch
#          @schemeWatch()
#        @schemeWatch = @scope.$watch 'extScheme', (val, old_val)=>
#          if val == old_val
#            return
#          @restart()

        # make objects
        $q.all(@cache).then ()=>
          @makePackageObjects(resolve)

        if @scope.settings.engine.container?.resizable?
          if @wrapper.resizable('instance')
            @wrapper.resizable('destroy')
          @wrapper.resizable
            minHeight: 200
            minWidth: 400
            grid: @scope.settings.draggable.grid
            handles: 's'

        # interceptors
        #-------------------------------
        @scope.instance.bind 'click', (e)=>

          shift = key.getPressedKeyCodes().indexOf(16) > -1
          @deselectAll() if !shift

          e.addClass('selected')

          #TODO fix, circular dependencies
          @scope.selected.push({
            id: e.id
            #object: e
            type: 'connector'
          })

          @scope.$apply()

        # disable loopback
        @scope.instance.bind 'beforeDrop', (e)=>
          e.sourceId != e.targetId

        @stopListen = @scope.$on '$stateChangeSuccess', ()=>
          @destroy()

        @wrapper.on 'mousedown', (e)=>
#          log.debug '@mousedown'
          @deselectAll()

      destroy: ()->
        log.debug 'destroy'
#        log.debug 'total objects:', @scope.intScheme.objects
        log.debug 'total connectors:', @scope.intScheme.connectors.length

        @wrapper.find(".page-loader").fadeIn("fast")

#        if @scope.settings.engine.container?.resizable?
#          @wrapper.resizable('destroy')

        angular.forEach @scope.intScheme.objects, (obj)->
          obj.remove()

        @scope.intScheme.objects = []
        @scope.intScheme.connectors = []
        @scope.instance.empty(@container)

        if @schemeWatch
          @schemeWatch()

        if @panning
          @panning.destroy()
        @panning = null

        @wrapper.off 'mousedown'

      #TODO fix restart
      restart: ()->
        log.debug 'restart'
        $timeout ()=>
          if @isStarted?
            @destroy()
          @start()

      changeTheme: ()=>
        @loadStyle()
        @cacheTemplates()
        @objectsUpdate()

    bpmnScheme
  ]