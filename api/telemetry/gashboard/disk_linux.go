// +build linux

package dasboard

import (
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/e154/smart-home/api/log"
)

func NewDisk() (disk *Disk) {

	stats, err := linuxproc.ReadDiskStats("/proc/diskstats")
	if err != nil {
		log.Error("disk stat read fail")
	}

	disk = &Disk{
		DiskStats: stats,
	}

	return
}

type Disk struct {
	DiskStats	[]linuxproc.DiskStat	`json:"disk_stat"`
}

func (d *Disk) Update() {

}