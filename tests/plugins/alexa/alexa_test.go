// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

package alexa

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	alexa2 "github.com/e154/smart-home/internal/plugins/alexa"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
	"github.com/e154/smart-home/pkg/scripts"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/e154/bus"
	. "github.com/e154/smart-home/tests/plugins"
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
			scriptService scripts.ScriptService,
			supervisor plugins.Supervisor,
			eventBus bus.Bus) {

			// register plugins
			AddPlugin(adaptors, "triggers")
			AddPlugin(adaptors, "alexa")

			var lastVal string
			scriptService.PushFunctions("Done", func(args string) {
				lastVal = args
			})

			// common
			// ------------------------------------------------
			ch := make(chan alexa2.EventAlexaAction)
			defer close(ch)
			fn := func(_ string, m interface{}) {
				switch v := m.(type) {
				case alexa2.EventAlexaAction:
					ch <- v
				}
			}
			eventBus.Subscribe(alexa2.TopicPluginAlexa, fn)
			defer eventBus.Unsubscribe(alexa2.TopicPluginAlexa, fn)

			// add scripts
			// ------------------------------------------------

			alexaSkillScript, err := AddScript("alexa skill script", skillScriptSrc, adaptors, scriptService)
			So(err, ShouldBeNil)

			// add alexa skills
			// ------------------------------------------------

			skill := &models.AlexaSkill{
				SkillId:     "amzn1.ask.skill.1ccc278b-ffbf-440c-87e3-83349761fbab",
				Description: "flat lights",
				Status:      "enabled",
				ScriptId:    common.Int64(alexaSkillScript.Id),
			}
			skill.Id, err = adaptors.AlexaSkill.Add(context.Background(), skill)
			So(err, ShouldBeNil)

			intent := &models.AlexaIntent{
				Name:         "FlatLights",
				AlexaSkillId: skill.Id,
				ScriptId:     alexaSkillScript.Id,
			}
			err = adaptors.AlexaIntent.Add(context.Background(), intent)
			So(err, ShouldBeNil)

			// ------------------------------------------------

			serviceCh := WaitService(eventBus, time.Second*5, "Supervisor")
			supervisor.Start(context.Background())
			defer supervisor.Shutdown(context.Background())
			So(<-serviceCh, ShouldBeTrue)

			// ------------------------------------------------
			plugin, err := supervisor.GetPlugin("alexa")
			So(err, ShouldBeNil)

			alexaPlugin, ok := plugin.(alexa2.AlexaPlugin)
			So(ok, ShouldBeTrue)

			server := alexaPlugin.Server()

			t.Run("on launch", func(t *testing.T) {
				req := &alexa2.Request{}
				err = json.Unmarshal([]byte(launchRequest), req)
				ctx.So(err, ShouldBeNil)

				resp := alexa2.NewResponse()
				server.OnLaunchHandler(nil, req, resp)

				ctx.So(lastVal, ShouldEqual, "skillOnLaunch")
			})

			t.Run("on intent", func(t *testing.T) {

				req := &alexa2.Request{}
				err = json.Unmarshal([]byte(intentRequest), req)
				ctx.So(err, ShouldBeNil)

				resp := alexa2.NewResponse()
				server.OnIntentHandle(nil, req, resp)

				ctx.So(lastVal, ShouldEqual, "kitchen_on")

				// wait message
				msg, ok := WaitT[alexa2.EventAlexaAction](time.Second*2, ch)
				ctx.So(ok, ShouldBeTrue)

				ctx.So(msg.Payload, ShouldEqual, "kitchen_on")
				ctx.So(msg.SkillId, ShouldEqual, 1)
				ctx.So(msg.IntentName, ShouldEqual, "FlatLights")
			})
		})
	})
}
