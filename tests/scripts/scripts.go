package scripts

const coffeeScript1 = `
"use strict";

r = IC.ExecuteSync "data/scripts/ping.sh", "google.com"

if r.out == 'ok'
    print "site is available ^^"

IC.store(r.out)
`

const coffeeScript2 = `
"use strict";

require('data/scripts/util.js')
require('data/scripts/vars.js')

# var foo
print 'foo value:', foo

# var bar
print 'bar value:', bar

# func from external lib
print 'run external function: ', area(4)

console.log months[0]

IC.store(foo + '-' + bar + '-' + months[0])
`