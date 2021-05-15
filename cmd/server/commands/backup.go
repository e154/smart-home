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
	backupCmd = &cobra.Command{
		Use:   "backup",
		Short: "Backup database settings to file",
		Long:  "Before upgrading your system, it is strongly recommended that you make a full backup, or at least back up any data or configuration information you can't afford to lose.",
		Run: func(cmd *cobra.Command, args []string) {

			app := BuildContainer(fx.Invoke(func(
				logger *logging.Logging,
				backup *backup.Backup) {

				if err := backup.New(); err != nil {
					log.Error(err.Error())
				}
			}))

			Start(app)
		},
	}
)
