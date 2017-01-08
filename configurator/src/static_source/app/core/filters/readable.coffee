###*
# Created by delta54 on 20.12.14.
###

angular
.module 'appServices'

  .filter 'readableDate', ->
    (date) ->
      return '' if !date
      moment(date).format("DD MM YYYY")

  .filter 'readableDateTime', ->
    (datetime) ->
      return '' if !datetime
      moment(datetime).format("DD.MM.YYYY HH:mm")

  .filter 'readableTime', ->
    (time) ->
      return '' if !time
      Date.create(time.replace('.000Z', '')).format('HH:mm')

  .filter 'readableTimeTZ', ->
    (time) ->
      return '' if !time
      Date.create(time.replace('.000Z', '')).format('{HH}:{mm}')
      Date.create(time.slice(0, '1970-01-01T00:00:00.000'.length)).format('HH:mm')

  .filter 'readableBytes', ->
    (bytes, si, decimals)->
      return '0' if !bytes
      k = if si then 1000 else 1024
      if Math.abs(bytes) < k
        return bytes + ' B'
      dm = decimals + 1 || 3
      units = if !si then [
        'kB','MB','GB','TB','PB','EB','ZB','YB'
      ] else [
        'KiB','MiB','GiB','TiB','PiB','EiB','ZiB','YiB'
      ]
      i = Math.floor(Math.log(bytes) / Math.log(k))
      parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + units[i]

  .filter 'toHHMMSS', ->
    (sec_num) ->
      days = Math.floor(sec_num / (3600 * 24))
      hours = Math.floor((sec_num - days) / 3600)
      minutes = Math.floor((sec_num - (hours * 3600)) / 60)
      seconds = Math.floor(sec_num - (hours * 3600) - (minutes * 60))

      days = "0" + days if days < 10
      hours = "0" + hours if hours < 10
      minutes = "0" + minutes if minutes < 10
      seconds = "0" + seconds if seconds < 10

      "#{days}:#{hours}:#{minutes}:#{seconds}"

  .filter 'uptime', ['$filter', ($filter)->
    (sec_num) ->
      days = Math.floor(sec_num / (3600 * 24))
      hours = Math.floor((sec_num - days) / 3600)
      minutes = Math.floor((sec_num - (hours * 3600)) / 60)
      seconds = Math.floor(sec_num - (hours * 3600) - (minutes * 60))

      translate = $filter('translate')
      _days = translate('uptime.days')
      _hours = translate('uptime.hours')
      _minutes = translate('uptime.minutes')
      _seconds = translate('uptime.seconds')
      _online = translate('uptime.online_for')

      time = "#{_online} "
      time += " 0" if days > 0 && days < 10
      time += "#{days} #{_days}, " if days > 0
      time += "0" if hours < 10
      time += "#{hours} #{_hours}, "
      time += "0" if minutes < 10
      time += "#{minutes} #{_minutes}, "
      time += "0" if seconds < 10
      time += "#{seconds} #{_seconds}"
      time
  ]