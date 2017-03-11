package scripts

import "github.com/e154/smart-home/api/models"

type Model struct {

}

func (m *Model) NewDevice() *models.Device {
	return &models.Device{}
}

func (m *Model) NewFlow() *models.Flow {
	return &models.Flow{}
}

func (m *Model) NewNode() *models.Node {
	return &models.Node{}
}
