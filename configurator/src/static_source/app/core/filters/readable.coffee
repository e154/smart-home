###*
# Created by delta54 on 20.12.14.
###

angular
.module 'appServices'

  .filter 'readableDate', ->
    (date) ->
      if date
        moment(date).format("DD MM YYYY")
      else
        ''

  .filter 'readableDateTime', ->
    (datetime) ->
      if datetime
        moment(datetime).format("DD.MM.YYYY HH:mm")
      else
        ''

  .filter 'readableTime', ->
    (time) ->
      if time
        Date.create(time.replace('.000Z', '')).format('HH:mm')
      else
        ''

  .filter 'readableTimeTZ', ->
    (time) ->
      if time
        Date.create(time.replace('.000Z', '')).format('{HH}:{mm}')
        Date.create(time.slice(0, '1970-01-01T00:00:00.000'.length)).format('HH:mm')
      else
        ''