package twilio

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/sfreiberg/gotwilio"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	log = logging.MustGetLogger("twilio")
)

type TWClient struct {
	cfg    *TWConfig
	client *gotwilio.Twilio
}

func NewTWClient(cfg *TWConfig) (*TWClient, error) {
	if cfg.sid == "" || cfg.authToken == "" {
		return nil, errors.New("bad parameters")
	}

	tw := &TWClient{
		cfg:    cfg,
		client: gotwilio.NewTwilioClient(cfg.sid, cfg.authToken),
	}
	return tw, nil
}

func (t *TWClient) SendSMS(phone, body string) (string, error) {

	log.Infof("send sms %v, %v", phone, body)

	var resp *gotwilio.SmsResponse
	var ex *gotwilio.Exception

	if !strings.Contains(phone, "+") {
		phone = fmt.Sprintf("+%s", phone)
	}

	var err error
	resp, ex, err = t.client.SendSMS(t.cfg.from, phone, body, "", "")
	if err != nil {
		return "", errors.New(err.Error())
	}

	if ex != nil {
		return "", errors.New(ex.Message)
	}

	return resp.Sid, nil
}

func (t *TWClient) GetStatus(smsId string) (string, error) {

	var resp *gotwilio.SmsResponse
	var ex *gotwilio.Exception
	var err error

	resp, ex, err = t.client.GetSMS(smsId)
	if err != nil {
		return "", errors.New(err.Error())
	}

	if ex != nil {
		return "", errors.New(ex.Message)
	}

	return resp.Status, nil
}

func (t *TWClient) Balance() (float32, error) {

	uri, err := url.Parse(fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Balance.json", t.cfg.sid))
	if err != nil {
		return 0, errors.New(err.Error())
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		return 0, errors.New(err.Error())
	}

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", t.cfg.sid, t.cfg.authToken)))
	req.Header.Add("Authorization", "Basic "+auth)

	resp, err := client.Do(req)
	if err != nil {
		return 0, errors.New(err.Error())
	}
	defer resp.Body.Close()

	balance := &Balance{}
	if err = json.NewDecoder(resp.Body).Decode(balance); err != nil {
		return 0, errors.New(err.Error())
	}
	amount, _ := strconv.ParseFloat(balance.Balance, 64)

	return float32(amount), nil
}
