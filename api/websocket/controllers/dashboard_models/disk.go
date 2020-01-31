//+build linux windows darwin,!386

package dashboard_models

import (
	"github.com/shirou/gopsutil/disk"
	"sync"
)

func NewDisk() (_disk *Disk) {

	var state_root *disk.UsageStat
	//var state_tmp *disk.UsageStat
	var err error

	_disk = &Disk{}

	if state_root, err = disk.Usage("/"); err == nil {
		_disk.Root = state_root
	}

	//if state_tmp, err = disk.Usage("/tmp"); err == nil {
	//	_disk.Tmp = state_tmp
	//}

	return
}

type Disk struct {
	sync.Mutex
	Root *disk.UsageStat `json:"root"`
	//Tmp		*disk.UsageStat		`json:"tmp"`
}

func (d *Disk) Update() {

	var state_root *disk.UsageStat
	//var state_tmp *disk.UsageStat
	var err error

	if state_root, err = disk.Usage("/"); err == nil {
		d.Lock()
		d.Root = state_root
		d.Unlock()
	}

	//if state_tmp, err = disk.Usage("/tmp"); err == nil {
	//	d.Tmp = state_tmp
	//}
}
