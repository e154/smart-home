package commands

import (
	. "github.com/e154/smart-home/cmd/server/container"
	. "github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/backup"
	"github.com/e154/smart-home/system/logging"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var (
	log = MustGetLogger("main")
)

var (
	filename string

	restoreCmd = &cobra.Command{
		Use:   "restore",
		Short: "Restore settings from backup archive",
		Run: func(cmd *cobra.Command, args []string) {

			app := BuildContainer(fx.Invoke(func(
				logger *logging.Logging,
				backup *backup.Backup) {

				if err := backup.Restore(filename); err != nil {
					log.Error(err.Error())
				}

			}))
			Start(app)
		},
	}
)

func init() {
	backupCmd.Flags().StringVarP(&filename, "filename", "f", "backup.zip", "backup file name")
}
