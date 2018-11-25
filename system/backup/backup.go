package backup

import (
	"fmt"
	"os/exec"
	"os"
	"path/filepath"
	"path"
	"github.com/jinzhu/gorm"
	"time"
	"errors"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("backup")
)

type Backup struct {
	cfg     *BackupConfig
	Options []string
	db      *gorm.DB
}

func NewBackup(cfg *BackupConfig,
	db *gorm.DB) *Backup {
	return &Backup{
		cfg: cfg,
		db:  db,
	}
}

func (b *Backup) New() (err error) {
	log.Info("backup")

	options := b.dumpOptions()

	tmpDir := path.Join(os.TempDir(), "smart_home")
	if err = os.MkdirAll(tmpDir, 0755); err != nil {
		return
	}

	// filename
	filename := path.Join(tmpDir, "database.tar")
	options = append(options, "-f", filename)

	//log.Info()("options", options)

	cmd := exec.Command("pg_dump", options...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", b.cfg.PgPass))

	_, err = cmd.CombinedOutput()
	if err != nil {
		return
	}

	err = zipit([]string{path.Join("data", "file_storage"), filename}, path.Join(b.cfg.Path, fmt.Sprintf("%s.zip", time.Now().Format("2006-01-02T15:04:05.999")), ))
	if err != nil {
		return
	}

	os.Remove(tmpDir)

	log.Info("complete")

	return
}

func (b *Backup) List() (list []string) {

	filepath.Walk(b.cfg.Path, func(path string, info os.FileInfo, err error) error {
		if info.Name() == ".gitignore" || info.Name() == b.cfg.Path || info.IsDir() {
			return nil
		}
		list = append(list, info.Name())
		return nil
	})
	return
}

func (b *Backup) Restore(name string) (err error) {
	log.Infof("restore: %s", name)

	file := path.Join(b.cfg.Path, name)

	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		err = errors.New("file not found")
		return
	}

	tmpDir := path.Join(os.TempDir(), "smart_home")
	if err = unzip(file, tmpDir); err != nil {
		return
	}

	//log.Info()("tmpDir", tmpDir)

	log.Info("Purge database")

	if err = b.db.Exec(`DROP SCHEMA IF EXISTS "public" CASCADE;`).Error; err != nil {
		return
	}
	if err = b.db.Exec(`CREATE SCHEMA "public";`).Error; err != nil {
		return
	}

	options := b.restoreOptions()

	options = append(options, "-f", path.Join(tmpDir, "database.tar"))

	//log.Info()("options", options)

	cmd := exec.Command("psql", options...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", b.cfg.PgPass))

	if _, err = cmd.CombinedOutput(); err != nil {
		return
	}

	os.Remove(path.Join("data", "file_storage"))

	if err = Copy(path.Join(tmpDir, "file_storage"), path.Join("data", "file_storage")); err != nil {
		return
	}

	os.Remove(tmpDir)

	log.Info("complete")

	return
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
	//options = append(options, "-Z", "9")

	// formats
	options = append(options, "-F", "t")

	return options
}


func (b Backup) restoreOptions() []string {
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

	return options
}
