// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package onvif

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	wsnt "github.com/eyetowers/gonvif/pkg/generated/onvif/docs_oasisopen_org/wsn/b2"
	deviceWsdl "github.com/eyetowers/gonvif/pkg/generated/onvif/www_onvif_org/ver10/device/wsdl"
	eventsWsdl "github.com/eyetowers/gonvif/pkg/generated/onvif/www_onvif_org/ver10/events/wsdl"
	media1Wsdl "github.com/eyetowers/gonvif/pkg/generated/onvif/www_onvif_org/ver10/media/wsdl"
	"github.com/eyetowers/gonvif/pkg/generated/onvif/www_onvif_org/ver10/schema"
	media2Wsdl "github.com/eyetowers/gonvif/pkg/generated/onvif/www_onvif_org/ver20/media/wsdl"
	ptzWsdl "github.com/eyetowers/gonvif/pkg/generated/onvif/www_onvif_org/ver20/ptz/wsdl"
	"github.com/eyetowers/gonvif/pkg/gonvif"

	"github.com/e154/smart-home/common"
)

const (
	unsubscribeTimeout = 2 * time.Second
	profileIndex       = 0
	pollTimeout        = "PT60S"
)

var (
	subscriptionTimeout wsnt.AbsoluteOrRelativeTimeType = "PT120S"
)

type Client struct {
	username, password, address string
	port                        int64
	requireAuthorization        bool
	cli                         gonvif.Client
	mediaProfiles               []*schema.Profile
	media2Profiles              []*media2Wsdl.MediaProfile
	capabilities                *schema.Capabilities
	pTZConfigurationOptions     *schema.PTZConfigurationOptions
	isStarted                   atomic.Bool
	quit                        chan struct{}
	wg                          sync.WaitGroup
	actorHandler                func(interface{})
}

func NewClient(handler func(interface{})) *Client {
	return &Client{
		actorHandler: handler,
	}
}

func (s *Client) Start(username, password, address string, port int64, requireAuthorization bool) (err error) {
	if s.isStarted.Load() {
		return
	}
	s.isStarted.Store(true)

	s.username = username
	s.password = password
	s.address = address
	s.port = port
	s.requireAuthorization = requireAuthorization

	s.quit = make(chan struct{})

	s.wg.Add(1)

	go func() {
		defer func() {
			s.wg.Done()
		}()

		var counter int

		for {

			counter++
			if counter >= 3 {
				counter = 0
				go s.actorHandler(&ConnectionStatus{false})
			}

			select {
			case <-s.quit:
				return
			default:
			}

			if err != nil {
				time.Sleep(time.Second * 5)
			}

			// Connect to the Onvif device.
			s.cli, err = gonvif.New(context.Background(), fmt.Sprintf("http://%s:%d", address, port), username, password, false)
			if err != nil {
				continue
			}

			if err = s.GetCapabilities(); err != nil {
				continue
			}

			if err = s.getOptions(); err != nil {
				continue
			}

			if len(s.mediaProfiles) == 0 {
				continue
			}

			var streamList []string
			if streamList, err = s.GetStreamList(); err != nil {
				continue
			}

			snapshotURI := s.GetSnapshotURI()

			go s.actorHandler(&StreamList{List: streamList, SnapshotUri: snapshotURI})

			err = s.eventServiceSubscribe()
		}
	}()

	return
}

func (s *Client) Shutdown() (err error) {
	if !s.isStarted.Load() {
		return
	}
	close(s.quit)
	s.wg.Wait()
	s.isStarted.Store(false)
	return
}

func (s *Client) GetCapabilities() error {
	device, err := s.cli.Device()
	if err != nil {
		return err
	}

	resp, err := device.GetCapabilities(&deviceWsdl.GetCapabilities{})
	if err != nil {
		return err
	}
	s.capabilities = resp.Capabilities
	return nil
}

func (s *Client) GetStreamList() ([]string, error) {

	var list = make([]string, 0)

	var protocol schema.TransportProtocol
	protocol = schema.TransportProtocolTCP
	if s.capabilities.Media.StreamingCapabilities.RTP_RTSP_TCP {
		protocol = schema.TransportProtocolRTSP
	}
	var stream = schema.StreamTypeRTPUnicast
	if media, err := s.cli.Media(); err == nil {
		for _, profile := range s.mediaProfiles {
			resp, _ := media.GetStreamUri(&media1Wsdl.GetStreamUri{
				StreamSetup: &schema.StreamSetup{
					Transport: &schema.Transport{
						Protocol: &protocol,
					},
					Stream: &stream,
				},
				ProfileToken: profile.Token,
			})
			if resp != nil && resp.MediaUri != nil {
				list = append(list, s.prepareUri(resp.MediaUri.Uri))
			}
		}
	}

	if media, err := s.cli.Media2(); err == nil {
		for _, profile := range s.media2Profiles {
			resp, _ := media.GetStreamUri(&media2Wsdl.GetStreamUri{
				Protocol:     string(protocol),
				ProfileToken: profile.Token,
			})
			if resp != nil {
				list = append(list, s.prepareUri(resp.Uri))
			}
		}
	}

	return list, nil
}

func (s *Client) GetSnapshotURI() *string {
	var uri string
	if media, err := s.cli.Media(); err == nil {
		resp, _ := media.GetSnapshotUri(&media1Wsdl.GetSnapshotUri{
			ProfileToken: s.mediaProfiles[profileIndex].Token,
		})
		if resp != nil && resp.MediaUri != nil {
			uri = resp.MediaUri.Uri
		}
	}
	if media, err := s.cli.Media2(); err == nil {
		resp, _ := media.GetSnapshotUri(&media2Wsdl.GetSnapshotUri{
			ProfileToken: s.mediaProfiles[profileIndex].Token,
		})
		if resp != nil {
			uri = resp.Uri
		}
	}
	if uri == "" {
		return nil
	}
	return common.String(s.prepareUri(uri))
}

func (s *Client) ContinuousMove(X, Y float32) error {

	if X == 0 && Y == 0 {
		return nil
	}

	ptz, err := s.cli.PTZ()
	if err != nil {
		return err
	}

	options := s.pTZConfigurationOptions.Spaces.ContinuousPanTiltVelocitySpace[profileIndex]
	if Y > options.YRange.Max {
		Y = options.YRange.Max
	}
	if Y < options.YRange.Min {
		Y = options.YRange.Min
	}

	if X > options.XRange.Max {
		X = options.XRange.Max
	}
	if X < options.XRange.Min {
		X = options.XRange.Min
	}

	var profileToken *schema.ReferenceToken
	if s.mediaProfiles != nil {
		profileToken = s.mediaProfiles[profileIndex].Token
	}
	if s.media2Profiles != nil {
		profileToken = s.media2Profiles[profileIndex].Token
	}
	_, err = ptz.ContinuousMove(&ptzWsdl.ContinuousMove{
		ProfileToken: profileToken,
		Velocity: &schema.PTZSpeed{
			PanTilt: &schema.Vector2D{
				X: X,
				Y: Y,
			},
		},
	})
	if err != nil {
		log.Warn(err.Error())
	}
	return err
}

func (s *Client) StopContinuousMove() error {

	ptz, err := s.cli.PTZ()
	if err != nil {
		return err
	}

	var profileToken *schema.ReferenceToken
	if s.mediaProfiles != nil {
		profileToken = s.mediaProfiles[profileIndex].Token
	}
	if s.media2Profiles != nil {
		profileToken = s.media2Profiles[profileIndex].Token
	}
	_, err = ptz.Stop(&ptzWsdl.Stop{
		ProfileToken: profileToken,
	})
	if err != nil {
		log.Warn(err.Error())
	}
	return err
}

func (s *Client) getOptions() error {

	// MEDIA PROFILES
	if media, err := s.cli.Media(); err == nil {
		var resp *media1Wsdl.GetProfilesResponse
		resp, err = media.GetProfiles(&media1Wsdl.GetProfiles{})
		if err == nil {
			s.mediaProfiles = resp.Profiles
		}
	}

	if media, err := s.cli.Media2(); err == nil {
		var resp *media2Wsdl.GetProfilesResponse
		resp, err = media.GetProfiles(&media2Wsdl.GetProfiles{
			Type: []string{"All"},
		})
		if err == nil {
			s.media2Profiles = resp.Profiles
		}
	}

	// PTZ
	ptz, err := s.cli.PTZ()
	if err == nil {
		var configurationToken *schema.ReferenceToken
		if s.mediaProfiles != nil {
			configurationToken = s.mediaProfiles[profileIndex].PTZConfiguration.Token
		}
		if s.media2Profiles != nil {
			configurationToken = s.media2Profiles[profileIndex].Configurations.PTZ.Token
		}
		configurationOptions, err := ptz.GetConfigurationOptions(&ptzWsdl.GetConfigurationOptions{
			ConfigurationToken: configurationToken,
		})
		if err == nil {
			s.pTZConfigurationOptions = configurationOptions.PTZConfigurationOptions
		}
	}

	return nil
}

func (s *Client) eventServiceSubscribe() error {
	events, err := s.cli.Events()
	if err != nil {
		return err
	}
	resp, err := events.CreatePullPointSubscription(&eventsWsdl.CreatePullPointSubscription{
		InitialTerminationTime: &subscriptionTimeout,
	})
	if err != nil {
		return err
	}
	headers := gonvif.ComposeHeaders(resp.SubscriptionReference)
	subscription, err := s.cli.Subscription(string(*resp.SubscriptionReference.Address), headers...)
	if err != nil {
		return err
	}
	return s.processEvents(subscription)
}

func (s *Client) processEvents(subscription eventsWsdl.PullPointSubscription) error {
	defer func() { _ = s.unsubscribe(subscription) }()
	ch := make(chan *eventsWsdl.PullMessagesResponse)
	chErr := make(chan error)
	defer func() {
		close(ch)
		close(chErr)
	}()

	for {

		go func() {
			resp, err := subscription.PullMessages(&eventsWsdl.PullMessages{MessageLimit: 100, Timeout: pollTimeout})
			select {
			case <-s.quit:
				return
			default:
			}
			if err != nil {
				chErr <- err
				return
			}
			ch <- resp
		}()

		select {
		case <-s.quit:
			return nil
		case v := <-ch:
			s.eventHandler(v.NotificationMessage)
			if _, err := subscription.Renew(&wsnt.Renew{TerminationTime: &subscriptionTimeout}); err != nil {
				return err
			}
		case err := <-chErr:
			return err
		}
	}
}

func (s *Client) unsubscribe(subscription eventsWsdl.PullPointSubscription) error {
	ctx, cancel := context.WithTimeout(context.Background(), unsubscribeTimeout)
	defer cancel()

	var empty eventsWsdl.EmptyString
	_, err := subscription.UnsubscribeContext(ctx, &empty)
	return err
}

func (s *Client) eventHandler(messages []*wsnt.NotificationMessage) {
	for _, msg := range messages {
		switch msg.Topic.Value {
		case "tns1:VideoSource/MotionAlarm":
			s.prepareMotionAlarm(msg)
		case "tns1:VideoSource/GlobalSceneChange/ImagingService":
			s.prepareImagingService(msg)
		default:
			log.Debugf("unknown message topic: \"%s\"", msg.Topic.Value)
		}
	}
}

func (s *Client) prepareMotionAlarm(msg *wsnt.NotificationMessage) {
	if msg.Message.Message == nil || msg.Message.Message.PropertyOperation != "Changed" {
		return
	}
	var state = false
	var t time.Time
	if msg.Message.Message != nil && msg.Message.Message.Data != nil &&
		msg.Message.Message.Data.SimpleItem != nil && len(msg.Message.Message.Data.SimpleItem) > 0 {
		state = msg.Message.Message.Data.SimpleItem[profileIndex].Value == "true"
	}
	if msg.Message.Message != nil && msg.Message.Message.UTCTime != nil {
		t = msg.Message.Message.UTCTime.Time
	}
	go s.actorHandler(&MotionAlarm{State: state, Time: t})
}

func (s *Client) prepareImagingService(msg *wsnt.NotificationMessage) {

}

var re = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

func (s *Client) prepareUri(uri string) string {
	if !s.requireAuthorization || !re.MatchString(uri) {
		return uri
	}
	ip := re.FindString(uri)
	return strings.ReplaceAll(uri, ip, fmt.Sprintf("%s:%s@%s", s.username, s.password, ip))
}
