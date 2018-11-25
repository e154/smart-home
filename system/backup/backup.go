package backup

import (
	"fmt"
	"os/exec"
	"os"
	"time"
	"github.com/e154/smart-home-old/api/log"
	"path/filepath"
)

type Backup struct {
	cfg     *BackupConfig
	Options []string
	path string
}

func NewBackup(cfg *BackupConfig) *Backup {
	return &Backup{
		cfg: cfg,
		path: "./snapshots",
	}
}

func (b *Backup) New() {
	fmt.Println("backup")

	options := b.dumpOptions()

	// filename
	filename := fmt.Sprintf("%s/database.psql", b.path)
	options = append(options, "-f", filename)

	//fmt.Println("options", options)

	cmd := exec.Command("pg_dump", options...)
	cmd.Env = append(os.Environ(), "PGPASSWORD=qwe123")

	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}

	err = zipit([]string{"./data/file_storage", filename}, fmt.Sprintf("%s/%s.zip", b.path, time.Now().Format("2006-01-02T15:04:05.999")))
	if err != nil {
		log.Error(err.Error())
	}

	if err = os.Remove(filename); err != nil {
		log.Error(err.Error())
	}
}

func (b *Backup) List() (list []string) {

	filepath.Walk(b.path, func(path string, info os.FileInfo, err error) error {
		if info.Name() == ".gitignore" || info.Name() == b.path || info.IsDir() {
			return nil
		}
		list = append(list, info.Name())
		return nil
	})
	return
}

func (b *Backup) Restore(name string) {
	fmt.Println("restore:", name)


}

func (b Backup) dumpOptions() []string {
	options := b.Options

	// db name
	if b.cfg.PgName != "" {
		options = append(options, "-d")
		options = append(options, b.cfg.PgName)
	}

	// host
	if b.cfg.PgHost != "" {
		options = append(options, "-h")
		options = append(options, b.cfg.PgHost)
	}

	// port
	if b.cfg.PgPort != "" {
		options = append(options, "-p")
		options = append(options, b.cfg.PgPort)
	}

	// user
	if b.cfg.PgUser != "" {
		options = append(options, "-U")
		options = append(options, b.cfg.PgUser)
	}

	// compress level
	options = append(options, "-Z", "9")

	// formats
	//options = append(options, "-F", "t")

	return options
}
