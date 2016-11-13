angular
.module('appServices')
.factory 'uuid', [() ->
  class uuid
    stack: []

    gen: ->
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

    getId: ->
      id = @gen()
      if @stack.indexOf(id) == -1
        @stack.push(id)
        return id
      else
      return @getId()

  new uuid
]