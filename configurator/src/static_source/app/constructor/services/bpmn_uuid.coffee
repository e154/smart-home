
#
# http://jsfiddle.net/briguy37/2mvfd/
#

angular
.module('angular-bpmn')
.factory 'bpmnUuid', () ->
  class uuid

    @generator: ()->
      'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c)->
        d = new Date().getTime()
        r = (d + Math.random()*16)%16 | 0
        d = Math.floor(d/16)

        q = null
        if c == 'x'
          q = r
        else
          q = (r&0x3|0x8)
        return (q).toString(16)
      )

    @gen: ()->
      return @generator()

  uuid