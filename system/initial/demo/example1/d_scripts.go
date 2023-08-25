package example1

import (
	"fmt"
	"os"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
	"github.com/e154/smart-home/system/scripts"
)

type ScriptManager struct {
	adaptors      *adaptors.Adaptors
	scriptService scripts.ScriptService
}

func NewScriptManager(adaptors *adaptors.Adaptors,
	scriptService scripts.ScriptService) *ScriptManager {
	return &ScriptManager{
		adaptors:      adaptors,
		scriptService: scriptService,
	}
}

func (s *ScriptManager) addScripts() (scripts []*m.Script, err error) {
	// L3+ script
	var script1 *m.Script
	if script1, err = s.addScript("l3+_script_v1", sourceScript1, "l3+ script v1"); err != nil {
		return
	}
	// api monitor
	var script2 *m.Script
	if script2, err = s.addScript("sensor_script_v1", fmt.Sprintf(sourceScript2, os.Getenv("LC_ADDRESS")), "sensor script v1"); err != nil {
		return
	}
	scripts = []*m.Script{script1, script2}
	return
}

func (s *ScriptManager) addScript(name, source, desc string) (script *m.Script, err error) {

	if script, err = s.adaptors.Script.GetByName(name); err == nil {
		return
	}

	script = &m.Script{
		Lang:        common.ScriptLangCoffee,
		Name:        name,
		Source:      source,
		Description: desc,
	}

	engineScript, err := s.scriptService.NewEngine(script)
	So(err, ShouldBeNil)

	err = engineScript.Compile()
	So(err, ShouldBeNil)

	script.Id, err = s.adaptors.Script.Add(script)
	So(err, ShouldBeNil)
	return
}

const sourceScript1 = `

# entity
# ##################################
ifError =(res)->
    return !res || res.error || res.Error

checkStatus =->
    stats = Miner.stats()
    if ifError(stats)
        Actor.setState
            'new_state': 'ERROR'
        return
    p = JSON.parse(stats.result)
    attrs = {
        heat: false
        chain1_temp_chip: p.temp2_1
        chain2_temp_chip: p.temp2_2
        chain3_temp_chip: p.temp2_3
        chain4_temp_chip: p.temp2_4
        chain1_temp_pcb: p.temp1
        chain2_temp_pcb: p.temp2
        chain3_temp_pcb: p.temp3
        chain4_temp_pcb: p.temp4
        chain_acn1: p.chain_acn1
        chain_acn2: p.chain_acn2
        chain_acn3: p.chain_acn3
        chain_acn4: p.chain_acn4
        ghs_av: p["GHS av"]
        fan1: p.fan1
        fan2: p.fan2
    }
    status = 'ENABLED'
    if p.chain_acn1 != 72 || p.chain_acn2 != 72 || p.chain_acn3 != 72 || p.chain_acn4 != 72 
        status = 'WARNING'
    if p.fan1 == 0 || p.fan2 == 0
        status = 'WARNING'
    if p.temp2_1 >= 60 || p.temp2_2 >= 60 || p.temp2_3 >= 60 || p.temp2_4 >= 60
        status = 'WARNING'
    if p.temp1 >= 60 || p.temp2 >= 60 || p.temp3 >= 60 || p.temp4 >= 60
        status = 'WARNING'
    Actor.setState
        new_state: status
        attribute_values: attrs
        storage_save: false

entityAction = (entityId, actionName)->
    switch actionName
        when 'CHECK' then checkStatus()

# automation
# ##################################
automationTriggerTime = (msg)->
    supervisor.callAction(msg.entity_id, 'CHECK', {})
    return false

automationTriggerStateChanged = (msg)->
    #print '---trigger---'
    p = unmarshal msg.payload
    if !p.old_state.state || !p.old_state.state.name
        return
    newState = p.new_state.state.name
    oldState = p.old_state.state.name
    if newState == oldState
        return false
    return newState == 'WARNING' || newState == 'ERROR'

automationCondition = (entityId)->
    #print '---condition---'
    entity = supervisor.getEntity(entityId)
    if !entity
        return false
    if entity.state && (entity.state.name == 'WARNING' || entity.state.name == 'ERROR')
        return true
    return false

automationAction = (entityId)->
    #print '---action---'
    entity = supervisor.getEntity(entityId)
    attr = entity.getAttributes()
    sendMsg(format(entityId, entity.state.name, attr))

# telegram
# ##################################
telegramSendReport =->
    entities = ['cgminer.l3n1','cgminer.l3n2','cgminer.l3n3','cgminer.l3n4','cgminer.l3n5']
    for entityId, i in entities
        entity = supervisor.getEntity(entityId)
        attr = entity.getAttributes()
        sendMsg(format(entityId, entity.state.name, attr))

format =(entityId, stateName, attr)->
	return entityId + " status: " + stateName + "\\r\\n" +
		"chain_acn1: " +  attr.chain_acn1 + "\\r\\n" +
		"chain_acn2: " +  attr.chain_acn2 + "\\r\\n" +
		"chain_acn3: " +  attr.chain_acn3 + "\\r\\n" +
		"chain_acn4: " +  attr.chain_acn4 + "\\r\\n" +
		"chain1_temp_chip: " +  attr.chain1_temp_chip + "\\r\\n" +
        "chain2_temp_chip: " +  attr.chain2_temp_chip + "\\r\\n" +
		"chain3_temp_chip: " +  attr.chain3_temp_chip + "\\r\\n" +
		"chain4_temp_chip: " +  attr.chain4_temp_chip + "\\r\\n" +
		"chain1_temp_pcb: " +  attr.chain1_temp_pcb + "\\r\\n" +
		"chain2_temp_pcb: " +  attr.chain2_temp_pcb + "\\r\\n" +
		"chain3_temp_pcb: " +  attr.chain3_temp_pcb + "\\r\\n" +
        "chain4_temp_pcb: " +  attr.chain4_temp_pcb + "\\r\\n" +
		"heat: " +  attr.heat + "\\r\\n" +
		"hardware_errors: " +  attr.hardware_errors + "\\r\\n" +
		"GHS av: " +  attr.ghs_av + "\\r\\n" +
		"fan1: " +  attr.fan1 + "\\r\\n" +
		"fan2: " +  attr.fan2 + "\\r\\n"

telegramAction = (entityId, actionName)->
    switch actionName
       when 'CHECK' then telegramSendReport()

sendMsg =(body)->
    msg = notifr.newMessage();
    msg.type = 'telegram';
    msg.attributes = {
        'name': 'clavicus',
        'body': body
    };
    notifr.send(msg);
`

const sourceScript2 = `

# entity
# ##################################
checkStatus =->
    res = http.get("%s")
    if res.error 
        Actor.setState
            'new_state': 'ERROR'
        return
    p = JSON.parse(res.body)
    attrs =
        paid_rewards: p.user.paid_rewards

    Actor.setState
        new_state: 'ENABLED'
        attribute_values: attrs
        storage_save: true

entityAction = (entityId, actionName)->
    switch actionName
        when 'CHECK' then checkStatus()

# automation
# ##################################
automationTriggerTime = (msg)->
    supervisor.callAction(msg.entity_id, 'CHECK', {})
    return false
`
