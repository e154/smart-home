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

  .filter 'readableDateTimeMarkup', ->
    (datetime) ->
      if datetime
        data = Date.create(datetime)

        [data.format("DD.MM.YYYY"), '  <span class="badge label-danger">', data.format("HH:mm"), '</span>'].join('')
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

  .filter 'readablePhone', ->
    (input) ->
      input.toString().trim().replace(/(\d{1})(\d{3})(\d{3})(\d{2})(\d{2})/, "+$1\u00a0($2)\u00a0$3-$4-$5") if input && input.length > 0