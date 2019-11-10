package messagebird

import (
	"github.com/messagebird/go-rest-api"
	"github.com/messagebird/go-rest-api/balance"
	"github.com/messagebird/go-rest-api/sms"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
)

var (
	log = logging.MustGetLogger("message bird")
)

type MBClient struct {
	cfg    *MBClientConfig
	client *messagebird.Client
}

func NewMBClient(cfg *MBClientConfig) (*MBClient, error) {

	if cfg.Name == "" || cfg.AccessKey == "" {
		return nil, errors.New("bad parameters")
	}

	client := &MBClient{
		cfg:    cfg,
		client: messagebird.New(cfg.AccessKey),
	}
	return client, nil
}

func (c *MBClient) SendSMS(phone, body string) (string, error) {

	log.Infof("send sms %v, %v", phone, body)

	msgParams := &sms.Params{
		Type:       "sms",
		DataCoding: "unicode",
	}

	msg, err := sms.Create(c.client, c.cfg.Name, []string{phone}, body, msgParams)
	if err != nil {
		mbErr, ok := err.(messagebird.ErrorResponse)
		if !ok {
			return "", errors.New(err.Error())
		}

		//fmt.Println("Code:", mbErr.Errors[0].Code)
		//fmt.Println("Description:", mbErr.Errors[0].Description)
		//fmt.Println("Parameter:", mbErr.Errors[0].Parameter)

		return "", errors.New(mbErr.Errors[0].Description)
	}

	return msg.ID, nil
}

func (c *MBClient) GetStatus(smsId string) (string, error) {

	msg, err := sms.Read(c.client, smsId)
	if err != nil {
		return "", errors.New(err.Error())
	}

	return msg.Recipients.Items[0].Status, nil
}

func (c *MBClient) Balance() (*balance.Balance, error) {

	b, err := balance.Read(c.client)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return b, nil
}
