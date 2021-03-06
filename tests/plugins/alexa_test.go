// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package plugins

import (
	"encoding/json"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/alexa"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestAlexa(t *testing.T) {

	const (
		skillScriptSrc = `
skillOnLaunch = ->
    #print '---action onLaunch---'
    Done('skillOnLaunch')
skillOnSessionEnd = ->
    #print '---action onSessionEnd---'
    Done('skillOnSessionEnd')
skillOnIntent = ->
    #print '---action onIntent---'
    state = 'on'
    if Alexa.slots['state'] == 'off'
        state = 'off'

    place = Alexa.slots['place']

    Done("#{place}_#{state}")
    
    Alexa.sendMessage("#{place}_#{state}")
`

		launchRequest = `{
    "version": "1.0",
    "session": {
        "new": true,
        "sessionId": "...",
        "application": {
            "applicationId": "amzn1.ask.skill.1ccc278b-ffbf-440c-87e3-83349761fbab"
        },
        "user": {
            "userId": "..."
        }
    },
    "context": {
        
        "Extensions": {
            "available": {
                "aplext:backstack:10": {}
            }
        },
        "System": {
            "application": {
                "applicationId": "amzn1.ask.skill.1ccc278b-ffbf-440c-87e3-83349761fbab"
            },
            "user": {
                "userId": "..."
            },
            "device": {
                "deviceId": "...",
                "supportedInterfaces": {}
            },
            "apiEndpoint": "https://api.amazonalexa.com",
            "apiAccessToken": "..."
        }
    },
    "request": {
        "type": "LaunchRequest",
        "requestId": "amzn1.echo-api.request.122e6887-0ddb-4781-ba88-67e15e928209",
        "locale": "en-US",
        "timestamp": "2021-05-13T16:35:16Z",
        "shouldLinkResultBeReturned": false
    }
}`

		intentRequest = `{
	"version": "1.0",
	"session": {
		"new": false,
		"sessionId": "amzn1.echo-api.session.b7bcc77c-0165-47d0-a9b3-dfbc4f21b65b",
		"application": {
			"applicationId": "amzn1.ask.skill.1ccc278b-ffbf-440c-87e3-83349761fbab"
		},
		"user": {
			"userId": "..."
		}
	},
	"context": {
		"Extensions": {
			"available": {
				"aplext:backstack:10": {}
			}
		},
		"System": {
			"application": {
				"applicationId": "amzn1.ask.skill.1ccc278b-ffbf-440c-87e3-83349761fbab"
			},
			"user": {
				"userId": "..."
			},
			"device": {
				"deviceId": "...",
				"supportedInterfaces": {}
			},
			"apiEndpoint": "https://api.amazonalexa.com",
			"apiAccessToken": "..."
		}
	},
	"request": {
		"type": "IntentRequest",
		"requestId": "amzn1.echo-api.request.3c295be0-3b79-49a6-a274-b41162e17b52",
		"locale": "en-US",
		"timestamp": "2021-05-15T05:39:46Z",
		"intent": {
			"name": "FlatLights",
			"confirmationStatus": "NONE",
			"slots": {
				"state": {
					"name": "state",
					"value": "on",
					"resolutions": {
						"resolutionsPerAuthority": [
							{
								"authority": "...",
								"status": {
									"code": "ER_SUCCESS_MATCH"
								},
								"values": [
									{
										"value": {
											"name": "on",
											"id": "ed2b5c0139cec8ad2873829dc1117d50"
										}
									}
								]
							}
						]
					},
					"confirmationStatus": "NONE",
					"source": "USER",
					"slotValue": {
						"type": "Simple",
						"value": "on",
						"resolutions": {
							"resolutionsPerAuthority": [
								{
									"authority": "...",
									"status": {
										"code": "ER_SUCCESS_MATCH"
									},
									"values": [
										{
											"value": {
												"name": "on",
												"id": "ed2b5c0139cec8ad2873829dc1117d50"
											}
										}
									]
								}
							]
						}
					}
				},
				"place": {
					"name": "place",
					"value": "kitchen",
					"resolutions": {
						"resolutionsPerAuthority": [
							{
								"authority": "...",
								"status": {
									"code": "ER_SUCCESS_MATCH"
								},
								"values": [
									{
										"value": {
											"name": "kitchen",
											"id": "09228dac155633b13780552bc01dc2e0"
										}
									}
								]
							}
						]
					},
					"confirmationStatus": "NONE",
					"source": "USER",
					"slotValue": {
						"type": "Simple",
						"value": "kitchen",
						"resolutions": {
							"resolutionsPerAuthority": [
								{
									"authority": "...",
									"status": {
										"code": "ER_SUCCESS_MATCH"
									},
									"values": [
										{
											"value": {
												"name": "kitchen",
												"id": "09228dac155633b13780552bc01dc2e0"
											}
										}
									]
								}
							]
						}
					}
				}
			}
		}
	}
}`
	)

	Convey("alexa", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService,
			entityManager entity_manager.EntityManager,
			automation automation.Automation,
			eventBus event_bus.EventBus,
			pluginManager common.PluginManager) {

			err := migrations.Purge()
			So(err, ShouldBeNil)

			// register plugins
			err = AddPlugin(adaptors, "triggers")
			err = AddPlugin(adaptors, "alexa")
			ctx.So(err, ShouldBeNil)

			// add scripts
			// ------------------------------------------------
			alexaSkillScript := &m.Script{
				Lang:        common.ScriptLangCoffee,
				Name:        "alexa skill script",
				Source:      skillScriptSrc,
				Description: "",
			}

			engineAlexaSkill, err := scriptService.NewEngine(alexaSkillScript)
			So(err, ShouldBeNil)
			err = engineAlexaSkill.Compile()
			So(err, ShouldBeNil)

			alexaSkillScript.Id, err = adaptors.Script.Add(alexaSkillScript)
			So(err, ShouldBeNil)

			// add alexa skills
			// ------------------------------------------------

			skill := &m.AlexaSkill{
				SkillId:     "amzn1.ask.skill.1ccc278b-ffbf-440c-87e3-83349761fbab",
				Description: "flat lights",
				Status:      "enabled",
				ScriptId:    common.Int64(alexaSkillScript.Id),
			}
			skill.Id, err = adaptors.AlexaSkill.Add(skill)
			So(err, ShouldBeNil)

			intent := &m.AlexaIntent{
				Name:         "FlatLights",
				AlexaSkillId: skill.Id,
				ScriptId:     alexaSkillScript.Id,
			}
			err = adaptors.AlexaIntent.Add(intent)
			So(err, ShouldBeNil)

			var lastVal string
			scriptService.PushFunctions("Done", func(args string) {
				lastVal = args
			})

			// start system
			// ------------------------------------------------

			pluginManager.Start()
			automation.Reload()
			entityManager.LoadEntities(pluginManager)

			defer func() {
				entityManager.Shutdown()
				automation.Shutdown()
				pluginManager.Shutdown()
			}()

			time.Sleep(time.Millisecond * 500)

			// ...
			plugin, err := pluginManager.GetPlugin("alexa")
			So(err, ShouldBeNil)

			alexaPlugin, ok := plugin.(alexa.AlexaPlugin)
			So(ok, ShouldBeTrue)

			server := alexaPlugin.Server()

			t.Run("on launch", func(t *testing.T) {
				req := &alexa.Request{}
				err = json.Unmarshal([]byte(launchRequest), req)
				ctx.So(err, ShouldBeNil)

				resp := alexa.NewResponse()
				server.OnLaunchHandler(nil, req, resp)

				ctx.So(lastVal, ShouldEqual, "skillOnLaunch")
			})

			t.Run("on intent", func(t *testing.T) {

				ch := make(chan alexa.EventAlexaAction)
				eventBus.Subscribe(alexa.TopicPluginAlexa, func(_ string, msg alexa.EventAlexaAction) {
					ch <- msg
				})

				req := &alexa.Request{}
				err = json.Unmarshal([]byte(intentRequest), req)
				ctx.So(err, ShouldBeNil)

				resp := alexa.NewResponse()
				server.OnIntentHandle(nil, req, resp)

				ctx.So(lastVal, ShouldEqual, "kitchen_on")

				ticker := time.NewTimer(time.Second * 2)
				defer ticker.Stop()

				var msg alexa.EventAlexaAction
				select {
				case v := <-ch:
					msg = v
					break
				case <-ticker.C:
					break
				}

				ctx.So(msg.Payload, ShouldEqual, "kitchen_on")
				ctx.So(msg.SkillId, ShouldEqual, 1)
				ctx.So(msg.IntentName, ShouldEqual, "FlatLights")
			})
		})
	})
}
