// +build linux

package telemetry

func NewDisk() (disk *Disk) {

	return
}

type Disk struct {
	DiskStats	[]linuxproc.DiskStat	`json:"disk_stat"`
}

func (d *Disk) Update() {

}