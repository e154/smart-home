package commands

import (
	"fmt"
	. "github.com/e154/smart-home/cmd/server/container"
	. "github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/initial"
	"github.com/e154/smart-home/system/logging"
	"github.com/e154/smart-home/version"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var (
	Server = &cobra.Command{
		Use:   "server",
		Short: fmt.Sprintf(version.ShortVersionBanner, ""),
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Printf(version.ShortVersionBanner, "")

			app := BuildContainer(fx.Invoke(func(
				logger *logging.Logging,
				dbSaver logging.ISaver,
				initialService *initial.Initial,
			) {
				logger.SetSaver(dbSaver)
			}))

			Start(app)

			Work()

			Stop(app)
		},
	}
)

func init() {
	Server.AddCommand(backupCmd)
	Server.AddCommand(demoCmd)
	Server.AddCommand(resetCmd)
	Server.AddCommand(restoreCmd)
	Server.AddCommand(versionCmd)
}
