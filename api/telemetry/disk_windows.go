// +build windows

package telemetry

func NewDisk() (disk *Disk) {

	return
}

type Disk struct {
	DiskStats	string	`json:"disk_stat"`
}

func (d *Disk) Update() {

}