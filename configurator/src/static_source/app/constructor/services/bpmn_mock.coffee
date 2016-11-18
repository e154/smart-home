angular
.module('angular-bpmn')
.service 'bpmnMock', () ->
  {
    scheme1: {
      name: 'Simply bpmn scheme'
      description: ''
      objects: [
        {
          id: 1
          type:
            name: 'event'
            start:
              0:
                0: true
          position:
            top: 80
            left: 50
          status: ''
          error: ''
          title: 'start event'
          description: ''
        }
        {
          id: 2
          type:
            name: 'task'
            status: ''
            action: ''
          position:
            top: 60
            left: 260
          status: ''
          error: ''
          title: 'task'
          description: ''
        }
        {
          id: 3
          type:
            name: 'event'
            end:
              simply:
                top_level: true
          position:
            top: 80
            left: 530
          status: ''
          error: ''
          title: 'end event'
          description: ''
        }
      ]
      connectors: [
        {
          id: 1
          start:
            object: 1
            point: 1
          end:
            object: 2
            point: 10
          flow_type: "default"
          title: "connector №1"
        }
        {
          id: 2
          start:
            object: 2
            point: 4
          end:
            object: 3
            point: 3
          flow_type: "default"
          title: "connector №2"
        }
      ]
    }
    scheme2: {
      name: 'Parallel Event-Based Gateway'
      description: ''
      objects: [
        {
          id: 1
          type:
            name: 'event'
            start:
              0:
                0:
                  true
          position:
            top: 80
            left: 50
          status: ''
          error: ''
          title: 'message 1'
          description: ''
        }
        {
          id: 2
          type:
            name: 'event'
            start:
              0:
                0:
                  true
          position:
            top: 240
            left: 50
          status: ''
          error: ''
          title: 'message 2'
          description: ''
        }
        {
          id: 3
          type:
            name: 'gateway'
            base: 'data'
            status: 'xor'
          position:
            top: 160
            left: 190
          status: ''
          error: ''
          title: ''
          description: ''
        }
        {
          id: 4
          type:
            name: 'task'
          position:
            top: 140
            left: 370
          status: ''
          error: ''
          title: 'task'
          description: ''
        }
      ]
      connectors: [
        {
          id: 1
          start:
            object: 1
            point:1
          end:
            object: 3
            point: 0
          flow_type: "default"
          title: "connector №1"
        }
        {
          id: 2
          start:
            object: 2
            point: 1
          end:
            object: 3
            point: 2
          flow_type: "default"
          title: "connector №2"
        }
        {
          id: 2
          start:
            object: 3
            point: 1
          end:
            object: 4
            point: 10
          flow_type: "default"
          title: "connector №3"
        }
      ]
    }
    scheme3: {
      name: 'Base scheme'
      description: ''
      objects: [
        {
          id: 1
          type:
            name: 'event'
            start:
              simply:
                top_level: true
          parent: 0
          position:
            top: 210
            left: 120
          status: ''
          error: ''
          title: ''
          description: ''
        }
        {
          id: 2
          type:
            name: 'task'
          parent: 0
          position:
            top: 190
            left: 220
          status: ''
          error: ''
          title: 'Оформить заявку'
          description: ''
        }
        {
          id: 3
          type:
            name: 'task'
          parent: 0
          position:
            top: 190
            left: 370
          status: ''
          error: ''
          title: 'Рассмотреть заявку'
          description: ''
        }
        {
          id: 4
          type:
            name: 'gateway'
          parent: 0
          position:
            top: 210
            left: 540
          status: ''
          error: ''
          title: 'Одобрена?'
          description: ''
        }
        {
          id: 5
          type:
            name: 'task'
          parent: 0
          position:
            top: 190
            left: 650
          status: ''
          error: ''
          title: 'Выделить машину'
          description: ''
        }
        {
          id: 6
          type:
            name: 'task'
          parent: 0
          position:
            top: 390
            left: 510
          status: ''
          error: ''
          title: 'Заявка отклоннена'
          description: ''
        }
        {
          id: 7
          type:
            name: 'event'
            end:
              simply:
                top_level: true
          parent: 0
          position:
            top: 550
            left: 540
          status: ''
          error: ''
          title: 'Заявка отклонена'
          description: ''
        }
        {
          id: 8
          type:
            name: 'gateway'
          parent: 0
          position:
            top: 210
            left: 810
          status: ''
          error: ''
          title: 'Выделена?'
          description: ''
        }
        {
          id: 9
          type:
            name: 'gateway'
            status: 'parallel'
          parent: 0
          position:
            top: 210
            left: 930
          status: ''
          error: ''
          title: ''
          description: ''
        }
        {
          id: 10
          type:
            name: 'task'
          parent: 0
          position:
            top: 100
            left: 1050
          status: ''
          error: ''
          title: 'Машина выделена'
          description: ''
        }
        {
          id: 11
          type:
            name: 'task'
          parent: 0
          position:
            top: 290
            left: 1050
          status: ''
          error: ''
          title: 'Выполнить рейс'
          description: ''
        }
        {
          id: 12
          type:
            name: 'gateway'
            status: 'parallel'
          parent: 0
          position:
            top: 210
            left: 1200
          status: ''
          error: ''
          title: ''
          description: ''
        }
        {
          id: 13
          type:
            name: 'event'
            end:
              simply:
                top_level: true
          parent: 0
          position:
            top: 210
            left: 1340
          status: ''
          error: ''
          title: 'Заявка выполнена'
          description: ''
        }
      ]
      connectors: [
        {
          id: 1
          start:
            object: 1
            point:1
          end:
            object: 2
            point: 10
          flow_type: "default"
          title: ""
        }
        {
          id: 2
          start:
            object: 2
            point:4
          end:
            object: 3
            point: 10
          flow_type: "default"
          title: ""
        }
        {
          id: 3
          start:
            object: 3
            point:4
          end:
            object: 4
            point: 3
          flow_type: "default"
          title: ""
        }
        {
          id: 4
          start:
            object: 4
            point:2
          end:
            object: 6
            point: 1
          flow_type: "default"
          title: "Нет"
        }
        {
          id: 4
          start:
            object: 4
            point:1
          end:
            object: 5
            point: 10
          flow_type: "default"
          title: "Да"
        }
        {
          id: 5
          start:
            object: 6
            point: 7
          end:
            object: 7
            point: 0
          flow_type: "default"
          title: ""
        }
        {
          id: 6
          start:
            object: 5
            point: 4
          end:
            object: 8
            point: 3
          flow_type: "default"
          title: ""
        }
        {
          id: 7
          start:
            object: 8
            point: 0
          end:
            object: 3
            point: 1
          flow_type: "default"
          title: "Нет"
        }
        {
          id: 8
          start:
            object: 8
            point: 1
          end:
            object: 9
            point: 3
          flow_type: "default"
          title: "Да"
        }
        {
          id: 9
          start:
            object: 9
            point: 1
          end:
            object: 10
            point: 10
          flow_type: "default"
          title: ""
        }
        {
          id: 10
          start:
            object: 9
            point: 1
          end:
            object: 11
            point: 10
          flow_type: "default"
          title: ""
        }
        {
          id: 11
          start:
            object: 10
            point: 4
          end:
            object: 12
            point: 0
          flow_type: "default"
          title: ""
        }
        {
          id: 12
          start:
            object: 11
            point: 4
          end:
            object: 12
            point: 2
          flow_type: "default"
          title: ""
        }
        {
          id: 13
          start:
            object: 12
            point: 1
          end:
            object: 13
            point: 3
          flow_type: "default"
          title: ""
        }
      ]
    }
    scheme4: {
      name: 'Base scheme with grouping'
      description: ''
      objects: [
        {
          id: 1
          type:
            name: 'event'
            start:
              simply:
                top_level: true
          position:
            top: 80
            left: 50
          status: ''
          error: ''
          title: ''
          description: ''
        }
        {
          id: 2
          type:
            name: 'task'
          parent: 14
          position:
            top: 40
            left: 40
          status: ''
          error: ''
          title: 'Оформить заявку'
          description: ''
        }
        {
          id: 3
          type:
            name: 'task'
          parent: 15
          position:
            top: 40
            left: 40
          status: ''
          error: ''
          title: 'Рассмотреть заявку'
          description: ''
        }
        {
          id: 4
          type:
            name: 'gateway'
          position:
            top: 290
            left: 370
          status: ''
          error: ''
          title: 'Одобрена?'
          description: ''
        }
        {
          id: 5
          type:
            name: 'task'
          parent: 14
          position:
            top: 40
            left: 200
          status: ''
          error: ''
          title: 'Заявка отклонена'
          description: ''
        }
        {
          id: 6
          type:
            name: 'event'
            end:
              simply:
                top_level: true
          parent: 14
          position:
            top: 60
            left: 380
          status: ''
          error: ''
          title: 'Заявка отклонена'
          description: ''
        }
        {
          id: 7
          type:
            name: 'task'
          parent: 16
          position:
            top: 30
            left: 30
          status: ''
          error: ''
          title: 'Выделить машину'
          description: ''
        }
        {
          id: 8
          type:
            name: 'gateway'
          parent: 16
          position:
            top: 50
            left: 210
          status: ''
          error: ''
          title: 'Выделена?'
          description: ''
        }
        {
          id: 9
          type:
            name: 'gateway'
            status: 'parallel'
          parent: 16
          position:
            top: 50
            left: 340
          status: ''
          error: ''
          title: ''
          description: ''
        }
        {
          id: 10
          type:
            name: 'task'
          parent: 14
          position:
            top: 40
            left: 480
          status: ''
          error: ''
          title: 'Машина выделена'
          description: ''
        }
        {
          id: 11
          type:
            name: 'gateway'
            status: 'parallel'
          position:
            top: 80
            left: 780
          status: ''
          error: ''
          title: ''
          description: ''
        }
        {
          id: 12
          type:
            name: 'event'
            end:
              simply:
                top_level: true
          position:
            top: 80
            left: 890
          status: ''
          error: ''
          title: 'Заявка выполнена'
          description: ''
        }
        {
          id: 13
          type:
            name: 'task'
          parent: 16
          position:
            top: 30
            left: 440
          status: ''
          error: ''
          title: 'Выполнить рейс'
          description: ''
        }
        {
          id: 14
          type:
            name: 'group'
          position:
            top: 20
            left: 140
          width: '614px'
          height: '170px'
          status: ''
          error: ''
          title: 'Заказчик'
          description: ''
          style: 'solid'
        }
        {
          id: 15
          type:
            name: 'group'
          position:
            top: 230
            left: 140
          width: '186px'
          height: '176px'
          status: ''
          error: ''
          title: 'Канцелярия'
          description: ''
          style: 'solid'
        }
        {
          id: 16
          type:
            name: 'group'
          position:
            top: 430
            left: 310
          width: '588px'
          height: '168px'
          status: ''
          error: ''
          title: 'Гараж'
          description: ''
          style: 'solid'
        }
      ]
      connectors: [
        {
          id: 1
          start:
            object: 1
            point:1
          end:
            object: 2
            point: 10
          flow_type: "default"
          title: ""
        }
        {
          id: 2
          start:
            object: 2
            point: 7
          end:
            object: 3
            point: 1
          flow_type: "default"
          title: ""
        }
        {
          id: 3
          start:
            object: 3
            point: 4
          end:
            object: 4
            point: 3
          flow_type: "default"
          title: ""
        }
        {
          id: 4
          start:
            object: 4
            point: 0
          end:
            object: 5
            point: 7
          flow_type: "default"
          title: "нет"
        }
        {
          id: 5
          start:
            object: 5
            point: 4
          end:
            object: 6
            point: 3
          flow_type: "default"
          title: ""
        }
        {
          id: 6
          start:
            object: 4
            point: 2
          end:
            object: 7
            point: 1
          flow_type: "default"
          title: "да"
        }
        {
          id: 7
          start:
            object: 7
            point: 4
          end:
            object: 8
            point: 3
          flow_type: "default"
          title: ""
        }
        {
          id: 7
          start:
            object: 8
            point: 2
          end:
            object: 3
            point: 7
          flow_type: "default"
          title: "нет"
        }
        {
          id: 8
          start:
            object: 8
            point: 1
          end:
            object: 9
            point: 3
          flow_type: "default"
          title: "да"
        }
        {
          id: 9
          start:
            object: 9
            point: 0
          end:
            object: 10
            point: 7
          flow_type: "default"
          title: ""
        }
        {
          id: 10
          start:
            object: 10
            point: 4
          end:
            object: 11
            point: 3
          flow_type: "default"
          title: ""
        }
        {
          id: 11
          start:
            object: 11
            point: 1
          end:
            object: 12
            point: 3
          flow_type: "default"
          title: ""
        }
        {
          id: 12
          start:
            object: 13
            point: 1
          end:
            object: 11
            point: 2
          flow_type: "default"
          title: ""
        }
        {
          id: 13
          start:
            object: 9
            point: 1
          end:
            object: 13
            point: 10
          flow_type: "default"
          title: ""
        }
      ]
    }
    scheme5: {
      name: 'Base scheme with swimlane'
      description: ''
      objects: [
        {
          id: 1
          type:
            name: 'event'
            start:
              simply:
                top_level: true
          parent: 15
          position:
            top: 70
            left: 120
          status: ''
          error: ''
          title: ''
          description: ''
        }
        {
          id: 2
          type:
            name: 'task'
          parent: 15
          position:
            top: 50
            left: 240
          status: ''
          error: ''
          title: 'Оформить заявку'
          description: ''
        }
        {
          id: 3
          type:
            name: 'task'
          parent: 16
          position:
            top: 60
            left: 240
          status: ''
          error: ''
          title: 'Рассмотреть заявку'
          description: ''
        }
        {
          id: 4
          type:
            name: 'gateway'
          parent: 16
          position:
            top: 80
            left: 420
          status: ''
          error: ''
          title: 'Одобрена?'
          description: ''
        }
        {
          id: 5
          type:
            name: 'task'
          parent: 15
          position:
            top: 50
            left: 390
          status: ''
          error: ''
          title: 'Заявка отклонена'
          description: ''
        }
        {
          id: 6
          type:
            name: 'event'
            end:
              simply:
                top_level: true
          parent: 15
          position:
            top: 70
            left: 570
          status: ''
          error: ''
          title: 'Заявка отклонена'
          description: ''
        }
        {
          id: 7
          type:
            name: 'task'
          parent: 17
          position:
            top: 60
            left: 390
          status: ''
          error: ''
          title: 'Выделить машину'
          description: ''
        }
        {
          id: 8
          type:
            name: 'gateway'
          parent: 17
          position:
            top: 80
            left: 570
          status: ''
          error: ''
          title: 'Выделена?'
          description: ''
        }
        {
          id: 9
          type:
            name: 'gateway'
            status: 'parallel'
          parent: 17
          position:
            top: 80
            left: 690
          status: ''
          error: ''
          title: ''
          description: ''
        }
        {
          id: 10
          type:
            name: 'task'
          parent: 15
          position:
            top: 50
            left: 660
          status: ''
          error: ''
          title: 'Машина выделена'
          description: ''
        }
        {
          id: 11
          type:
            name: 'gateway'
            status: 'parallel'
          parent: 15
          position:
            top: 70
            left: 830
          status: ''
          error: ''
          title: ''
          description: ''
        }
        {
          id: 12
          type:
            name: 'event'
            end:
              simply:
                top_level: true
          parent: 15
          position:
            top: 70
            left: 940
          status: ''
          error: ''
          title: 'Заявка выполнена'
          description: ''
        }
        {
          id: 13
          type:
            name: 'task'
          parent: 17
          position:
            top: 60
            left: 800
          status: ''
          error: ''
          title: 'Выполнить рейс'
          description: ''
        }
        {
          id: 14
          type:
            name: 'swimlane'
          position:
            top: 10
            left: 10
          width: '1040px'
          height: ''
          status: ''
          error: ''
          title: ''
          description: ''
          style: ''
        }
        {
          id: 15
          type:
            name: 'swimlane-row'
          parent: 14
          position:
            top: 0
            left: 0
          width: '100%'
          height: '200px'
          status: ''
          error: ''
          title: 'Заказчик'
          description: ''
          style: ''
        }
        {
          id: 16
          type:
            name: 'swimlane-row'
          parent: 14
          position:
            top: 0
            left: 0
          width: '100%'
          height: '200px'
          status: ''
          error: ''
          title: 'Канцелярия'
          description: ''
          style: ''
        }
        {
          id: 17
          type:
            name: 'swimlane-row'
          parent: 14
          position:
            top: 0
            left: 0
          width: '100%'
          height: '200px'
          status: ''
          error: ''
          title: 'Гараж'
          description: ''
          style: ''
        }
      ]
      connectors: [
        {
          id: 1
          start:
            object: 1
            point:1
          end:
            object: 2
            point: 10
          flow_type: "default"
          title: ""
        }
        {
          id: 2
          start:
            object: 2
            point: 7
          end:
            object: 3
            point: 1
          flow_type: "default"
          title: ""
        }
        {
          id: 3
          start:
            object: 3
            point: 4
          end:
            object: 4
            point: 3
          flow_type: "default"
          title: ""
        }
        {
          id: 4
          start:
            object: 4
            point: 0
          end:
            object: 5
            point: 7
          flow_type: "default"
          title: "нет"
        }
        {
          id: 5
          start:
            object: 5
            point: 4
          end:
            object: 6
            point: 3
          flow_type: "default"
          title: ""
        }
        {
          id: 6
          start:
            object: 4
            point: 2
          end:
            object: 7
            point: 1
          flow_type: "default"
          title: "да"
        }
        {
          id: 7
          start:
            object: 7
            point: 4
          end:
            object: 8
            point: 3
          flow_type: "default"
          title: ""
        }
        {
          id: 7
          start:
            object: 8
            point: 2
          end:
            object: 3
            point: 7
          flow_type: "default"
          title: "нет"
        }
        {
          id: 8
          start:
            object: 8
            point: 1
          end:
            object: 9
            point: 3
          flow_type: "default"
          title: "да"
        }
        {
          id: 9
          start:
            object: 9
            point: 0
          end:
            object: 10
            point: 7
          flow_type: "default"
          title: ""
        }
        {
          id: 10
          start:
            object: 10
            point: 4
          end:
            object: 11
            point: 3
          flow_type: "default"
          title: ""
        }
        {
          id: 11
          start:
            object: 11
            point: 1
          end:
            object: 12
            point: 3
          flow_type: "default"
          title: ""
        }
        {
          id: 12
          start:
            object: 13
            point: 1
          end:
            object: 11
            point: 2
          flow_type: "default"
          title: ""
        }
        {
          id: 13
          start:
            object: 9
            point: 1
          end:
            object: 13
            point: 10
          flow_type: "default"
          title: ""
        }
      ]
    }
  }

