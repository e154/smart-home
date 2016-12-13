angular
.module('angular-bpmn')
.factory 'bpmnEditor', ['log', 'bpmnUuid', '$compile', 'bpmnObjectFact', '$q', '$timeout'
  (log, bpmnUuid, $compile, bpmnObjectFact, $q, $timeout) ->
    class bpmnEditor
      wrapper: null
      container: null
      scope: null
      pallet: null
      mouseover: null

      constructor: (container, settings)->

      loadEditor: ()->
        log.debug 'load editor'
        if !@pallet
          @wrapper.append($compile('<div class="palette-entries popup" bpmn-editor-palette="settings.editorPallet" settings="settings"></div>')(@scope))
          @pallet = @wrapper.find('.palette-entries')
        @droppableInit()
        @keyboardBindings()

        @wrapper.on 'mousedown', ()=>
          @deselectAll()

      unloadEditor: ()->
        if !@pallet
          return

        log.debug 'unload editor'
        @pallet.remove()
        @pallet = null

      reloadEditor: ()->
        @unloadEditor()
        @loadEditor()

      droppableInit: ()->
        @wrapper.droppable({
          drop: (event, ui)=>
            offset = @wrapper.offset()
            position =
              left: (ui.offset.left - offset.left - @container.position().left) / @scope.zoom
              top: (ui.offset.top - offset.top - @container.position().top) / @scope.zoom

            parent = @selectElementByPoint(position.top, position.left)[0]
            log.debug parent
            if parent
              position =
                left: $(parent.element).offset().left  - position.left
                top: $(parent.element).offset().top  - position.top

              log.debug position

            # type update
            #----------------
            type = $(ui.draggable).attr('entry-type')
            data_group = $(ui.draggable).parent().attr('data-group')
            if !type || type == ''
              return

            id = bpmnUuid.gen()
            objects = []
            if data_group == 'swimlane'
              objects.push($.extend(true, angular.copy(@scope.settings.baseObject), {
                id: id
                type:
                  name: 'swimlane'
                draggable: false
                position: position
              }))
              objects.push($.extend(true, angular.copy(@scope.settings.baseObject), {
                id: bpmnUuid.gen()
                parent: id
                type:
                  name: 'swimlane-row'
                draggable: false
              }))
            else
              objects.push($.extend(true, angular.copy(@scope.settings.baseObject), {
                id: id
                type: JSON.parse(type)
                parent: parent?.data?.id || ''
                draggable: true
                position: position
              }))

            @addObjects(objects)
        })

      addObjects: (objects)->
        if !objects || objects == ''
          return

        promise = []
        newObjects = {}
        angular.forEach objects, (object)=>
          obj = new bpmnObjectFact(object, @scope)
          if !@scope.intScheme.objects[object.id]
            @scope.intScheme.objects[object.id] = obj
            newObjects[object.id] = obj
          promise.push(obj.elementPromise)

        # Ждём когда прогрузятся все шаблоны
        $q.all(promise).then ()=>
          # проходим по массиву ранее созданных объектов,
          # и добавляем в дом
          angular.forEach newObjects, (object)=>
            # добавляем объект в контейнер
            object.appendTo(@container, @scope.settings.point)

      removeObject: (selected)=>
        return if !selected

        switch selected?.type
          when "connector"
            @scope.instance.select().each (c)=>
              if c.id == selected.id
                c.removeOverlay("myLabel")
                @scope.instance.detach(c)
          when "shape"
            #TODO first remove child objects
            index = 0
            for key, object of @scope.intScheme.objects
              if object.data.id == selected.id
#                log.debug 'found', selected
                object.remove()
                delete @scope.intScheme.objects[key]
                break
              index++
          else
            log.error 'unknown object type'

      removeSelected: (scope)=>
        if !scope || !scope.selected
          return

        angular.forEach scope.selected, (selected)=>
          @removeObject(selected)

        scope.selected = []

        @serialise(scope)

      keyboardBindings: ()->
        if !@scope.settings.keyboard
          return

        @wrapper.on 'mouseover', ()=>
          @mouseover = true

        @wrapper.on 'mouseleave', ()=>
          @mouseover = false

        angular.forEach @scope.settings.keyboard, (button, key_id)=>
          log.debug 'bind key:', button.name
          fn = this[button.callback] || window[button.callback]
          if typeof fn != 'function'
            return
          key key_id, (event, handler)=>
            if @mouseover
              event.preventDefault()
              fn.apply(null, [@scope])

      setLabel: (conn, label)->
        label = "" if !label
        overlay = conn.getOverlay('myLabel')
        if !overlay
          overlay = [ "Label", { label: label, cssClass: "aLabel" }, id:"myLabel" ]
          conn.addOverlay(overlay)
        else
          overlay.setLabel(label)


      getLabel: (conn)->
        return if !conn
        return conn.getOverlay('myLabel')?.getLabel() || ''

      selectElementInAabb: (t, l, w, h)->
        @scope.selected = []
        angular.forEach @scope.intScheme.objects, (object)->

          itemOffset = object.elementOffset()
          if itemOffset.top >= t && itemOffset.left >= l &&
              itemOffset.right < l + w && itemOffset.bottom < t + h

            @scope.selected.push(object.data.id)
            object.select(true)
            object.group('select')

          else
            object.select(false)

        @scope.$apply()

      selectElementOutPoint: (t, l, w, h)->
        @scope.selected = []
        for key, object of @scope.intScheme.objects
          itemOffset = object.elementOffset()
          if itemOffset.top <= t && itemOffset.left <= l &&
              itemOffset.right >= l + w && itemOffset.bottom >= t + h
            @scope.selected.push(object.data.id)
            break

      getAllConnections: ()->
        return [] if !@scope.instance
        @scope.instance.getAllConnections()

      selectElementByPoint: (t, l)->
        @scope.selected = []
        log.debug 'check position:',t, l
        for key, object of @scope.intScheme.objects
          if !object.canAParent
            continue

          offset = @wrapper.offset()
          itemOffset = object.elementOffset()
          log.debug 'check:',object.data.type.name
          if itemOffset.top - offset.top <= t && itemOffset.left - offset.left <= l &&
              itemOffset.right - offset.left >= l && itemOffset.bottom - offset.top >= t
            log.debug '===>',object.data.type.name
            log.debug '===>',itemOffset.top, itemOffset.left, itemOffset.right, itemOffset.bottom
            @scope.selected.push(object)
          else
            log.debug 'no:',itemOffset.top, itemOffset.left, itemOffset.right, itemOffset.bottom

        @scope.selected

      serialiseConnection: (connection)->
        params = connection.getParameters()

        id = ""
        if !params['element-id']
          id = bpmnUuid.gen()
          connection.setParameter("element-id", id)
        else
          id = params['element-id']

        id: id
        direction: params['direction']
        start:
          object: $(connection.source).attr('element-id')
          point: connection.endpoints[0].getParameters()['anchor-id'] || 0
        end:
          object: $(connection.target).attr('element-id')
          point: connection.endpoints[1].getParameters()['anchor-id'] || 0
        title: connection.getOverlay('myLabel')?.getLabel() || ''

      serialise: (scope)->

        objects = []
        connectors = []

        # objects
        #----------------
        #log.debug 'objects',scope.intScheme.objects
        angular.forEach scope.intScheme.objects, (obj)->

          draggable = true
          type = obj.data['type']
          if type == 'swimlane' || type == 'swimlane-row'
            draggable = false

          objects.push({
            id: obj.data?.id || bpmnUuid.gen()
            type: obj.data.type || 'default'
            width: obj.data.width || 'auto'
            height: obj.data.height || 'auto'
            draggable: draggable
            position: obj.position
            status: obj.data.status || 'enabled'
            error: ''
            title: obj?.data?.title
            description: obj?.data?.description || ''
          })

        # connectors
        #----------------
        #log.debug 'connectors',scope.intScheme.connectors
        connections = @getAllConnections()
        angular.forEach connections, (connection)=>
          end = false; start = false
          con = @serialiseConnection(connection)
          for id, obj of scope.intScheme.objects
            start = true if id == con.start.object
            end = true if id == con.end.object
            if start && end
              connectors.push angular.copy(con)
              break

        $timeout ()=>
          @scope.$apply()

        @scope.extScheme = {
          objects: objects
          connectors: connectors
        }

      getScheme: ()->
        @serialise(@scope)

    bpmnEditor
]